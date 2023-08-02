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
  InitTime  int //ms
  Imm       bool
}

type Simple struct {
  hist           Distribution
  lastCallStartTs  int64
	config         *config.Config
	metaData       *model2.Meta
	platformClient platform_client2.Client
	mu             sync.Mutex
	wg             sync.WaitGroup
	instances      map[string]*model2.Instance
	idleInstance   *list.List
  preloadList    *list.List
  collected      bool
  running_cnt    int // 并发量
  max_running    int
  instance_cnt   int
  exe_time       int
  start_time     map[string]time.Time
}


func New(metaData *model2.Meta, config *config.Config) Scaler {
	client, err := platform_client2.New(config.ClientAddr)
	if err != nil {
		log.Fatalf("client init with error: %s", err.Error())
	}
	scheduler := &Simple{
    hist: *NewDistribution(100, 200),
    lastCallStartTs: 0,
		config:         config,
		metaData:       metaData,
		platformClient: client,
		mu:             sync.Mutex{},
		wg:             sync.WaitGroup{},
		instances:      make(map[string]*model2.Instance),
    start_time: make(map[string]time.Time),
		idleInstance:   list.New(),
    preloadList: list.New(),
	}
  log.Printf("New scaler : key %s, mem %d, timeout %d, runtime: %s", metaData.GetKey(), metaData.GetMemoryInMb(), metaData.GetTimeoutInSecs(), metaData.GetRuntime())
	scheduler.wg.Add(2)
	go func() {
		defer scheduler.wg.Done()
		scheduler.gcLoop()
		log.Printf("gc loop for app: %s is stoped", metaData.Key)
	}()

  go func() {
    defer scheduler.wg.Done()
    scheduler.preloadLoop()
    log.Printf("preload loop for app: %s is stoped", metaData.Key)
  }()

	return scheduler
}

//返回的是实例,每类app都会有一个独立的scaler
func (s *Simple) Assign(ctx context.Context, request *pb.AssignRequest, req_total int) (*pb.AssignReply, error) {
  log.Printf("assign, cur state, running_cnt : %d, idle_cnt : %d, instance_cnt : %d, max_running : %d", s.running_cnt, s.idleInstance.Len(), s.instance_cnt, s.max_running)
  start := time.Now()
	startTs := start.UnixMilli()
  if s.lastCallStartTs != 0 {
    dt := startTs - s.lastCallStartTs;
    s.mu.Lock()
    s.hist.Add(int(dt))
    s.hist.PrintReuseRate()
    log.Printf("hist mean : %f, cv : %f", s.hist.mean, s.hist.CV())
    if(s.hist.collected && s.hist.CV() < 1.5) {
      // s.config.PreloadInterval = s.hist.GetQuantiles(0.05)
      s.config.KeepAliveInterval = s.hist.GetQuantiles(0.95) - s.exe_time
      log.Printf("window update : %d , %d, cv : %f", s.config.PreloadInterval, s.config.KeepAliveInterval, s.hist.CV())
    }
    s.mu.Unlock()
  }
  s.lastCallStartTs = startTs
  meta := &model2.Meta{
		Meta: pb.Meta{
			Key:           request.MetaData.Key,
			Runtime:       request.MetaData.Runtime,
			TimeoutInSecs: request.MetaData.TimeoutInSecs,
      MemoryInMb: request.MetaData.MemoryInMb,
		},
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
    s.hist.num_reuse++
    // if(s.running_cnt > 0 && s.idleInstance.Len() == 1) { //说明当前是并发请求，进入pipeline
    //   s.Preload(ctx)
    //   log.Printf("AddPreloadList : meta mem : %d", instance.Meta.MemoryInMb)
    // }
		s.idleInstance.Remove(element)
    s.start_time[instanceId] = start
    s.running_cnt++
    if(s.running_cnt > s.max_running) {
      s.max_running = s.running_cnt
    }
		s.mu.Unlock()
		log.Printf("Assign, request id: %s, instance %s reused", request.RequestId, instance.Id)
		instanceId = instance.Id
    instance.ExeStartTime = time.Now()
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

	
	instance, err := s.platformClient.Init(ctx, request.RequestId, instanceId, slot, meta) //初始化instance, 一个slot只能分配给一个实例,一个实例就是一个app
  s.metaData.Meta.MemoryInMb = instance.Meta.MemoryInMb
	if err != nil {
		errorMessage := fmt.Sprintf("create instance failed with: %s", err.Error())
		log.Printf(errorMessage)
		return nil, status.Errorf(codes.Internal, errorMessage)
	}

	//add new instance
	s.mu.Lock()
  s.instance_cnt++
	instance.Busy = true
	s.instances[instance.Id] = instance
  start_time := time.Now()
  s.start_time[instanceId] = start_time  
  // if(s.running_cnt > 0) {
  //   s.Preload(ctx)
  // }
  // log.Printf("AddPreloadList : meta mem : %d", instance.Meta.MemoryInMb)
  s.running_cnt++
  if(s.running_cnt > s.max_running) {
    s.max_running = s.running_cnt
  }
	s.mu.Unlock()
	log.Printf("request id: %s, instance %s for app %s is created, init latency: %dms", request.RequestId, instance.Id, instance.Meta.Key, instance.InitDurationInMs)
  instance.ExeStartTime = time.Now()
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
func (s *Simple) Preload(ctx context.Context) error{
  requestId := "xxxxxxx-preload-xxxxxxx" 
  instanceId := uuid.New().String()
  resourceConfig := model2.SlotResourceConfig{
		ResourceConfig: pb.ResourceConfig{
			MemoryInMegabytes: s.metaData.GetMemoryInMb(),
    },
	}
  log.Printf("preload, mem: %d", s.metaData.GetMemoryInMb())
  slot, err := s.platformClient.CreateSlot(ctx, requestId, &resourceConfig) //创建slot
  if err != nil {
		errorMessage := fmt.Sprintf("AddPreloadList create slot failed with: %s", err.Error())
		log.Printf(errorMessage)
		return status.Errorf(codes.Internal, errorMessage)
	}
  meta := &model2.Meta{
		Meta: pb.Meta{
			Key:           s.metaData.GetKey(),
			Runtime:       s.metaData.GetRuntime(),
			TimeoutInSecs: s.metaData.GetTimeoutInSecs(),
		},
	}
  log.Printf("preload, meta key: %s, runtime : %s, TimeoutInSecs : %d", s.metaData.GetKey(), s.metaData.GetRuntime(), s.metaData.GetTimeoutInSecs())
	instance, err := s.platformClient.Init(ctx, requestId, instanceId, slot, meta) //初始化instance, 一个slot只能分配给一个实例,一个实例就是一个app
  if err != nil {
    errorMessage := fmt.Sprintf("Preload instance init failed with: %s", err.Error())
    log.Printf(errorMessage)
    return status.Errorf(codes.Internal, errorMessage)
  }
  instance.Busy = false
	instance.LastIdleTime = time.Now()
  // s.mu.Lock()
  s.instances[instance.Id] = instance
  s.idleInstance.PushFront(instance)
  // s.mu.Unlock()
  log.Printf("Preload instance Push")

  return err
}

func (s *Simple) AddPreloadList(t time.Time, meta *model2.Meta, initTime int, imm bool) {
  InstanceConfig := InstanceConfig{
    InsertTime: t,
    Meta: meta,
    InitTime: initTime,
    Imm: imm,
  }
  // s.mu.Lock()
  s.preloadList.PushFront(InstanceConfig)
  log.Printf("AddpreloadList, preloadList len : %d", s.preloadList.Len())
  // s.mu.Unlock()
}


func (s *Simple) Idle(ctx context.Context, request *pb.IdleRequest) (*pb.IdleReply, error) {
	if request.Assigment == nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("assignment is nil"))
	}
	reply := &pb.IdleReply{
		Status:       pb.Status_Ok,
		ErrorMessage: nil,
	}
	start := time.Now()
	instanceId := request.Assigment.InstanceId

	//log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	needDestroy := false
	slotId := ""
	// if request.Result != nil && request.Result.NeedDestroy != nil && *request.Result.NeedDestroy {
	// 	needDestroy = true
	// }
  if s.config.PreloadInterval > 0 {
     needDestroy = true 
  }
	defer func() {
	if needDestroy {
			s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
	}
	}()
	log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	s.mu.Lock()
  exe_time := time.Since(s.start_time[instanceId]).Milliseconds()
  s.exe_time = int(float64(s.exe_time) * 0.2 + float64(exe_time) * 0.8)
  log.Printf("exe_time in ms : %d", exe_time)
  s.running_cnt--
  // if(s.idleInstance.Len() > 0 && s.running_cnt == 0) {
  //   log.Printf("need destroy, reason s.idleInstance.Len() > 0 && s.running_cnt == 0, %d, %d", s.idleInstance.Len(), s.running_cnt)
  //   needDestroy = true
  //   s.instance_cnt--
  // }
	defer s.mu.Unlock()
	if instance := s.instances[instanceId]; instance != nil {
		slotId = instance.Slot.Id

  if  instance.InitDurationInMs < int64(s.hist.moving_ave) && instance.InitDurationInMs < int64(s.hist.GetQuantiles(0.05)) {
      needDestroy = true
      delete(s.instances, instanceId)
      s.instance_cnt--
      log.Printf("instance %s destroy, reason: instance.InitDurationInMs < moving_ave(dt), %d, %d", request.Assigment.RequestId, instance.InitDurationInMs, int64(s.hist.moving_ave))
      return reply ,nil
    }

		if needDestroy {
      delete(s.instances, instanceId)
			log.Printf("request id %s, instance %s need be destroy", request.Assigment.RequestId, instanceId)
			return reply, nil
		}

		if instance.Busy == false {
			log.Printf("request id %s, instance %s already freed", request.Assigment.RequestId, instanceId)
			return reply, nil
		}
		instance.Busy = false
    instance.LastIdleTime = start //keepalive start time
		s.idleInstance.PushFront(instance)
	} else {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("request id %s, instance %s not found", request.Assigment.RequestId, instanceId))
	}
	return &pb.IdleReply{
		Status:       pb.Status_Ok,
		ErrorMessage: nil,
	}, nil
}




 //  now := time.Now()
 //  s.lastCallEndTs = now.UnixMilli()	
 //  if request.Assigment == nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("assignment is nil"))
	// }
	// 
	// start := time.Now()
	// instanceId := request.Assigment.InstanceId
	// defer func() {
	// 	log.Printf("Idle, request id: %s, instance: %s, cost %dus", request.Assigment.RequestId, instanceId, time.Since(start).Microseconds())
	// }()

 //  slotId := ""
 //  needDestroy := false
 //  s.mu.Lock()
	// defer s.mu.Unlock()
 //  if instance := s.instances[instanceId]; instance != nil {
 //    slotId = instance.Slot.Id
 //    // if s.config.PreloadInterval != 0 {
 //    //   needDestroy = false
 //    // }
 //    if needDestroy {
 //      log.Printf("preload meta,key : %s", instance.Meta.Key)
 //      s.AddPreloadList(instance.Meta)
 //      s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
 //    } else {
 //      instance.LastIdleTime = time.Now() //keepalive start time
 //      instance.Busy = false
 //    	s.idleInstance.PushFront(instance)
 //    }
 //  }
 //  
	// return &pb.IdleReply{
	// 	Status:       pb.Status_Ok,
	// 	ErrorMessage: nil,
	// }, nil
// }

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
        // idleDuration := time.Since(instance.LastIdleTime).Milliseconds()
        // log.Printf("check gc, idleDuration : %d, s.config.KeepAliveIntervalInterval: %d", idleDuration, time.Duration(s.config.KeepAliveInterval))
				if s.idleInstance.Len() > s.max_running && (s.config.KeepAliveInterval > 0 && idleDuration > time.Duration(s.config.KeepAliveInterval)){
					//need GC
          log.Printf("gc, idleDuration : %d, s.config.KeepAliveIntervalInterval: %d", idleDuration, time.Duration(s.config.KeepAliveInterval))
					s.idleInstance.Remove(element)
          s.instance_cnt--
					delete(s.instances, instance.Id)
					s.mu.Unlock()
					go func() {
						reason := fmt.Sprintf("Idle duration: %d, excceed configured duration: %d", idleDuration, s.config.KeepAliveInterval)
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
      if(len(s.instances) > s.idleInstance.Len() && s.idleInstance.Len() == 0 && s.hist.moving_ave < float64(s.exe_time)) {
        log.Printf("prelod, reason : len(s.instances) > s.idleInstance.Len() && s.idleInstance.Len() == 0")
        s.Preload(context.Background())
      }
      // if element := s.preloadList.Back(); element != nil {
      //   instanceConfig := element.Value.(InstanceConfig)
      //   preloadDuration := time.Now().Sub(instanceConfig.InsertTime) / 1000000 //ms
      //   log.Printf("check preload, preloadDuration : %d, s.config.PreloadInterval: %d", preloadDuration, time.Duration(s.config.PreloadInterval))
      //   if instanceConfig.Imm == true || preloadDuration > time.Duration(s.config.PreloadInterval - instanceConfig.InitTime) {
      //     ctx := context.Background()
      //     s.Preload(ctx)
      //     s.preloadList.Remove(element)
      //     s.mu.Unlock()
      //     log.Printf("Preload execute, preloadDuration : %d", preloadDuration)
      //   } else {
      //     s.mu.Unlock()
      //   } 
      //   continue
      // }
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
