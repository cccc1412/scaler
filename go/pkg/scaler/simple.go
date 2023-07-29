/*
Copyright 2023 The Alibaba Cloud Serverless Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scaler

import (
	"container/list"
	"context"
	"fmt"
	"github.com/AliyunContainerService/scaler/go/pkg/config"
	model2 "github.com/AliyunContainerService/scaler/go/pkg/model"
	platform_client2 "github.com/AliyunContainerService/scaler/go/pkg/platform_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"

	pb "github.com/AliyunContainerService/scaler/proto"
	"github.com/google/uuid"

)

type InstanceConfig struct {
  InsertTime time.Time
  Meta  *model2.Meta
}

type Simple struct {
  hist           Distribution
  lastCallEndTs  int64
	config         *config.Config
	metaData       *model2.Meta
	platformClient platform_client2.Client
	mu             sync.Mutex
	wg             sync.WaitGroup
	instances      map[string]*model2.Instance
	idleInstance   *list.List
  preloadList    *list.List
  collected      bool
}


func New(metaData *model2.Meta, config *config.Config) Scaler {
	client, err := platform_client2.New(config.ClientAddr)
	if err != nil {
		log.Fatalf("client init with error: %s", err.Error())
	}
	scheduler := &Simple{
    hist: *NewDistribution(10, 200),
    lastCallEndTs: time.Now().UnixMilli(),
		config:         config,
		metaData:       metaData,
		platformClient: client,
		mu:             sync.Mutex{},
		wg:             sync.WaitGroup{},
		instances:      make(map[string]*model2.Instance),
		idleInstance:   list.New(),
    preloadList: list.New(),
	}
	log.Printf("New scaler for app: %s is created", metaData.Key)
	scheduler.wg.Add(1)
	go func() {
		defer scheduler.wg.Done()
		scheduler.gcLoop()
		log.Printf("gc loop for app: %s is stoped", metaData.Key)
	}()

  go func() {
    scheduler.preloadLoop()
  }()

	return scheduler
}

//返回的是实例,每类app都会有一个独立的scaler
func (s *Simple) Assign(ctx context.Context, request *pb.AssignRequest) (*pb.AssignReply, error) {
  start := time.Now()
	startTs := start.UnixMilli()
  dt := startTs - s.lastCallEndTs;
  s.hist.Add(int(dt))
  log.Printf("hist mean : %f", s.hist.mean)
  if(s.hist.collected) {
    s.config.PreloadInterval = s.hist.GetQuantiles(0.1)
    s.config.KeepAliveInterval = s.hist.GetQuantiles(0.9)
    log.Printf("window update : %d , %d", s.config.PreloadInterval, s.config.KeepAliveInterval)
  }
	instanceId := uuid.New().String()
	defer func() {
		log.Printf("Assign, request id: %s, instance id: %s, cost %dms", request.RequestId, instanceId, time.Since(start).Milliseconds())
	}()
	log.Printf("Assign, request id: %s", request.RequestId)
	s.mu.Lock()
	if element := s.idleInstance.Front(); element != nil { //idel列表不为空,取队头
		instance := element.Value.(*model2.Instance) //返回这个示例
		instance.Busy = true
		s.idleInstance.Remove(element)
		s.mu.Unlock()
		log.Printf("Assign, request id: %s, instance %s reused", request.RequestId, instance.Id)
		instanceId = instance.Id
		return &pb.AssignReply{
			Status: pb.Status_Ok,
			Assigment: &pb.Assignment{
				RequestId:  request.RequestId,
				MetaKey:    instance.Meta.Key,
				InstanceId: instance.Id,
			},
			ErrorMessage: nil,
		}, nil
	}
	s.mu.Unlock()

	//Create new Instance
	resourceConfig := model2.SlotResourceConfig{
		ResourceConfig: pb.ResourceConfig{
			MemoryInMegabytes: request.MetaData.MemoryInMb,
		},
	}
	slot, err := s.platformClient.CreateSlot(ctx, request.RequestId, &resourceConfig) //创建slot
	if err != nil {
		errorMessage := fmt.Sprintf("create slot failed with: %s", err.Error())
		log.Printf(errorMessage)
		return nil, status.Errorf(codes.Internal, errorMessage)
	}

	meta := &model2.Meta{
		Meta: pb.Meta{
			Key:           request.MetaData.Key,
			Runtime:       request.MetaData.Runtime,
			TimeoutInSecs: request.MetaData.TimeoutInSecs,
		},
	}
	instance, err := s.platformClient.Init(ctx, request.RequestId, instanceId, slot, meta) //初始化instance, 一个slot只能分配给一个实例,一个实例就是一个app
	if err != nil {
		errorMessage := fmt.Sprintf("create instance failed with: %s", err.Error())
		log.Printf(errorMessage)
		return nil, status.Errorf(codes.Internal, errorMessage)
	}

	//add new instance
	s.mu.Lock()
	instance.Busy = true
	s.instances[instance.Id] = instance
	s.mu.Unlock()
	log.Printf("request id: %s, instance %s for app %s is created, init latency: %dms", request.RequestId, instance.Id, instance.Meta.Key, instance.InitDurationInMs)

	return &pb.AssignReply{
		Status: pb.Status_Ok,
		Assigment: &pb.Assignment{
			RequestId:  request.RequestId,
			MetaKey:    instance.Meta.Key,
			InstanceId: instance.Id,
		},
		ErrorMessage: nil,
	}, nil
}

//idle 调用
func (s *Simple) Preload(ctx context.Context, meta *model2.Meta) error{
  //
  requestId := uuid.New().String()
  instanceId := uuid.New().String()
  resourceConfig := model2.SlotResourceConfig{
		ResourceConfig: pb.ResourceConfig{
			MemoryInMegabytes: meta.MemoryInMb,
		},
	}
  slot, err := s.platformClient.CreateSlot(ctx, requestId, &resourceConfig) //创建slot
  if err != nil {
		errorMessage := fmt.Sprintf("AddPreloadList create slot failed with: %s", err.Error())
		log.Printf(errorMessage)
		return status.Errorf(codes.Internal, errorMessage)
	}
	instance, err := s.platformClient.Init(ctx, requestId, instanceId, slot, meta) //初始化instance, 一个slot只能分配给一个实例,一个实例就是一个app
  if err != nil {
    errorMessage := fmt.Sprintf("Preload instance init failed with: %s", err.Error())
    log.Printf(errorMessage)
    return status.Errorf(codes.Internal, errorMessage)
  }
  instance.Busy = false
	instance.LastIdleTime = time.Now()
  s.mu.Lock()
  s.instances[instance.Id] = instance
  s.idleInstance.PushFront(instance)
  s.mu.Unlock()
  log.Printf("Preload instance Push")

  return err
}

func (s *Simple) AddPreloadList(meta *model2.Meta) {
  InstanceConfig := InstanceConfig{
    InsertTime: time.Now(),
    Meta: meta,
  }
  log.Printf("waiting for insert preloadList")
  // s.mu.Lock()
  s.preloadList.PushBack(InstanceConfig)
  log.Printf("AddpreloadList, preloadList len : %d", s.preloadList.Len())
  // s.mu.Unlock()
}


func (s *Simple) Idle(ctx context.Context, request *pb.IdleRequest) (*pb.IdleReply, error) {
  now := time.Now()
  s.lastCallEndTs = now.UnixMilli()	
  if request.Assigment == nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("assignment is nil"))
	}
	// reply := &pb.IdleReply{
	// 	Status:       pb.Status_Ok,
	// 	ErrorMessage: nil,
	// }
	start := time.Now()
	instanceId := request.Assigment.InstanceId
	defer func() {
		log.Printf("Idle, request id: %s, instance: %s, cost %dus", request.Assigment.RequestId, instanceId, time.Since(start).Microseconds())
	}()

  slotId := ""
  needDestroy := true
  s.mu.Lock()
	defer s.mu.Unlock()
  if instance := s.instances[instanceId]; instance != nil {
    slotId = instance.Slot.Id
    if s.config.PreloadInterval == 0 {
      needDestroy = false
    }
    if needDestroy {
      log.Printf("preload meta,key : %s", instance.Meta.Key)
      s.AddPreloadList(instance.Meta)
      s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
    } else {
      instance.LastIdleTime = time.Now()
      instance.Busy = false
    	s.idleInstance.PushFront(instance)
    }
  }
  



	// //log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	// needDestroy := false
 //  // needDestroy  := true
	// if request.Result != nil && request.Result.NeedDestroy != nil && *request.Result.NeedDestroy {
	// 	needDestroy = true
	// }
	// defer func() {
	// 	if needDestroy {
	// 		s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
	// 	}
	// }()
	// log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	// s.mu.Lock()
	// defer s.mu.Unlock()
	// if instance := s.instances[instanceId]; instance != nil {
	// 	slotId = instance.Slot.Id
	// 	instance.LastIdleTime = time.Now()
	// 	if needDestroy {
	// 		log.Printf("request id %s, instance %s need be destroy", request.Assigment.RequestId, instanceId)
	// 		return reply, nil
	// 	}

	// 	if instance.Busy == false {
	// 		log.Printf("request id %s, instance %s already freed", request.Assigment.RequestId, instanceId)
	// 		return reply, nil
	// 	}
	// 	instance.Busy = false
	// 	s.idleInstance.PushFront(instance)
	// } else {
	// 	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("request id %s, instance %s not found", request.Assigment.RequestId, instanceId))
	// }
	return &pb.IdleReply{
		Status:       pb.Status_Ok,
		ErrorMessage: nil,
	}, nil
}

func (s *Simple) deleteSlot(ctx context.Context, requestId, slotId, instanceId, metaKey, reason string) {
	log.Printf("start delete Instance %s (Slot: %s) of app: %s", instanceId, slotId, metaKey)
	if err := s.platformClient.DestroySLot(ctx, requestId, slotId, reason); err != nil {
		log.Printf("delete Instance %s (Slot: %s) of app: %s failed with: %s", instanceId, slotId, metaKey, err.Error())
	}
}

func (s *Simple) gcLoop() {
	log.Printf("gc loop for app: %s is started", s.metaData.Key)
	ticker := time.NewTicker(s.config.GcDuration)
	for range ticker.C {
		for {
			s.mu.Lock()
			if element := s.idleInstance.Back(); element != nil {
				instance := element.Value.(*model2.Instance)
				idleDuration := time.Now().Sub(instance.LastIdleTime) / 1000000
        log.Printf("check gc, idleDuration : %d, s.config.KeepAliveIntervalInterval: %d", idleDuration, time.Duration(s.config.KeepAliveInterval))
				if idleDuration > time.Duration(s.config.KeepAliveInterval){
					//need GC
					s.idleInstance.Remove(element)
					delete(s.instances, instance.Id)
					s.mu.Unlock()
					go func() {
						reason := fmt.Sprintf("Idle duration: %fs, excceed configured duration: %fs", idleDuration.Seconds(), s.config.IdleDurationBeforeGC.Seconds())
						ctx := context.Background()
						ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
						defer cancel()
						s.deleteSlot(ctx, uuid.NewString(), instance.Slot.Id, instance.Id, instance.Meta.Key, reason)
					}()

					continue
				}
			}
			s.mu.Unlock()
			break
		}
	}
}

func (s *Simple) preloadLoop() {
  log.Printf("preloadLoop start")
  log.Printf("PreloadInterval: %d", s.config.PreloadInterval)
  ticker := time.NewTicker(s.config.PreloadDuration)
  for range ticker.C {
    for {
      s.mu.Lock()
      if element := s.preloadList.Front(); element != nil {
        instanceConfig := element.Value.(InstanceConfig)
        s.preloadList.Remove(element)
        s.mu.Unlock()
        preloadDuration := time.Now().Sub(instanceConfig.InsertTime) / 1000000 //ms
        log.Printf("check preload, preloadDuration : %d, s.config.PreloadInterval: %d", preloadDuration, time.Duration(s.config.PreloadInterval))
        if preloadDuration > time.Duration(s.config.PreloadInterval) {
          ctx := context.Background()
          s.Preload(ctx, instanceConfig.Meta)
          // ticker.Reset(time.Duration(s.config.PreloadIterval * 1000000))
          log.Printf("Preload execute, preloadDuration : %d", preloadDuration)
        }
        continue
      }
      s.mu.Unlock()
      break
    }
  }
}

func (s *Simple) Stats() Stats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Stats{
		TotalInstance:     len(s.instances),
		TotalIdleInstance: s.idleInstance.Len(),
	}
}
