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
	"log"
	"sync"
	"time"
  "math"
	"github.com/AliyunContainerService/scaler/go/pkg/config"
	model2 "github.com/AliyunContainerService/scaler/go/pkg/model"
  my_model "github.com/AliyunContainerService/scaler/go/pkg/lgbm"
	platform_client2 "github.com/AliyunContainerService/scaler/go/pkg/platform_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
  preload_mu     sync.Mutex
	wg             sync.WaitGroup
	instances      map[string]*model2.Instance //不包括idle
	idleInstance   *list.List
  preloadList    *list.List
  collected      bool
  running_cnt    int // 并发量
  max_running    int // 最大并发量
  max_running_interval int
  max_running_mov int
  instance_cnt   int // 当前所有instacne = running + idle
  exe_time       int64
  assign_time    int64
  dt_time        int64
  cold_start_cnt int
  assign_cnt     int
  reuse_cnt      int
  preloading_cnt int
  start_time     map[string]int64
  waiting_cnt    int
  base           float64

  pred_running   int
  running_history RingBuffer 
  p              int
  d              int
  q              int
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
    preload_mu:     sync.Mutex{},
		wg:             sync.WaitGroup{},
		instances:      make(map[string]*model2.Instance),
    start_time: make(map[string]int64),
		idleInstance:   list.New(),
    preloadList: list.New(),

    running_history:  *NewRingBuffer(10), 
    max_running: 1,
    base: 0,
    p : 1,
    d : 0,
    q : 1,
	}
  log.Printf("New scaler : key %s, mem %d, timeout %d, runtime: %s", metaData.GetKey(), metaData.GetMemoryInMb(), metaData.GetTimeoutInSecs(), metaData.GetRuntime())
	scheduler.wg.Add(3)
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

  // go func() {
  //   defer scheduler.wg.Done()
  //   for i :=0; i < 20; i++ {
  //     scheduler.Preload(context.Background())
  //     log.Printf("preload loop for app: %s is stoped", metaData.Key)
  //   }
  // }()


	return scheduler
}

func (s *Simple) delayAssignLoop(delay_req *pb.AssignRequest) (*pb.AssignReply, error) {
  for {
    s.mu.Lock()
    if element := s.idleInstance.Front(); element != nil { //idel列表不为空,取队头
        instance := element.Value.(*model2.Instance) //返回这个示例
		    instance.Busy = true
        s.hist.num_reuse++
    		s.idleInstance.Remove(element)
        s.start_time[instance.Id] = time.Now().UnixMilli() 
        s.running_cnt++
        if(s.running_cnt > s.max_running) {
          s.max_running = s.running_cnt
        }
        if(s.running_cnt > s.max_running_interval) {
          s.max_running_interval = s.running_cnt
        }
    
		    s.mu.Unlock()
		    log.Printf("delay assign succeed, request id: %s, instance %s reused", delay_req.RequestId, instance.Id)
        instance.ExeStartTime = time.Now()
        s.reuse_cnt++
        s.waiting_cnt--
	  	  return &pb.AssignReply{
		    	Status: pb.Status_Ok,
		    	Assigment: &pb.Assignment{
		    		RequestId:  delay_req.RequestId,
		    		MetaKey:    instance.Meta.Key,
		    		InstanceId: instance.Id,
		    	},
			    ErrorMessage: nil,
		    }, nil
    
    } else {
      s.mu.Unlock()
      time.Sleep(20 * time.Millisecond)
    }
  }
}

//返回的是实例,每类app都会有一个独立的scaler
func (s *Simple) Assign(ctx context.Context, request *pb.AssignRequest, req_total int) (*pb.AssignReply, error) {
  s.mu.Lock()
  s.assign_cnt++
  s.mu.Unlock()
  log.Printf("assign for metaKey : %s, cur state, running_cnt : %d, idle_cnt : %d, instance_cnt : %d, max_running : %d, keepalive window: %d, preloading_cnt : %d", s.metaData.Key, s.running_cnt, s.idleInstance.Len(), s.instance_cnt, s.max_running, s.config.KeepAliveInterval, s.preloading_cnt)
  log.Printf("moving average cur state, exe_time :%d, assign_time: %d, dt :%d", s.exe_time, s.assign_time, s.dt_time)
  log.Printf("cur state , reuse rate : %f", float64(s.reuse_cnt) / float64(s.assign_cnt))
  log.Printf("cur state, pred_running : %d", s.pred_running)
  start := time.Now()
	startTs := start.UnixMilli()
  if s.lastCallStartTs != 0 {
    dt := startTs - s.lastCallStartTs;
    s.dt_time = int64(0.2*float64(s.dt_time) + 0.8 * float64(dt))
    s.mu.Lock()
    s.hist.Add(int(dt))
    s.hist.PrintReuseRate()
    log.Printf("hist mean : %f, cv : %f", s.hist.mean, s.hist.CV())
    // if(s.hist.collected && s.hist.CV() < 1.5) {
    //   // s.config.PreloadInterval = s.hist.GetQuantiles(0.05)
    //   s.config.KeepAliveInterval = s.hist.GetQuantiles(0.95) - int(s.exe_time)
    //   log.Printf("window update : %d , %d, cv : %f", s.config.PreloadInterval, s.config.KeepAliveInterval, s.hist.CV())
    // }
    s.mu.Unlock()
  }     
  // if(s.assign_time != 0) {
  //   // s.config.KeepAliveInterval = int(s.exe_time*s.exe_time / s.assign_time)
  //   s.config.KeepAliveInterval = int(s.assign_time) + s.hist.GetQuantiles(0.95)
  //   s.config.GcDuration = time.Duration(s.config.KeepAliveInterval * 1e6 / 10)
  //   log.Printf("keepalive window update : %d, GcDuration : %d", s.config.KeepAliveInterval, s.config.GcDuration / 1e6)
  // }
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
	// defer func() {
 //    data_for_arima := s.running_history.ToArray()
 //    if(len(data_for_arima) > 10) {
 //      for i := range(data_for_arima) {
 //        log.Printf("%d", data_for_arima[i])
 //      }
 //      predictions := predictARIMA(intToFloatArray(data_for_arima), s.p, s.d, s.q, 10)
 //      tmp := 0.0
 //      end := data_for_arima[len(data_for_arima) - 1]
 //      for _, pred := range predictions {
 //        tmp += pred + float64(end)
 //      }
 //      s.pred_running = int(tmp / float64(len(predictions)))
 //      go func() {
 //        if(s.idleInstance.Len() <  s.pred_running && s.pred_running > s.running_cnt) {
 //          log.Printf("preload, reason : s.idleInstance.Len() <  s.pred_running && s.pred_running > s.running_cnt , idle_len : %d, pre_running : %d, running_cnt : %d", s.idleInstance.Len(), s.pred_running, s.running_cnt)
 //          s.Preload(context.Background())
 //        }
 //      }()
 //    }
	// 	log.Printf("Assign, request id: %s, instance id: %s, cost %dms", request.RequestId, instanceId, time.Since(start).Milliseconds())
	// }()
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
    s.start_time[instance.Id] = startTs 
    s.running_cnt++
    if(s.running_cnt > s.max_running) {
      s.max_running = s.running_cnt
    }
    if(s.running_cnt > s.max_running_interval) {
      s.max_running_interval = s.running_cnt
    }
    
		s.mu.Unlock()
		log.Printf("Assign, request id: %s, instance %s reused", request.RequestId, instance.Id)
		instanceId = instance.Id
    instance.ExeStartTime = time.Now()
    s.reuse_cnt++
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

  s.mu.Lock()
  s.metaData.Meta.MemoryInMb = request.MetaData.MemoryInMb
  pred_wait_time := math.Max(float64(s.waiting_cnt/(s.running_cnt + 1)), 1) * float64(s.exe_time)
  // if(s.preloading_cnt > s.waiting_cnt || s.assign_time >= s.exe_time && s.metaData.GetTimeoutInSecs() * 1000 > uint32(s.exe_time) && s.waiting_cnt < s.running_cnt) {
  if(s.preloading_cnt > s.waiting_cnt || s.assign_time >= int64(pred_wait_time) && s.metaData.GetTimeoutInSecs()*1000 > uint32(pred_wait_time)){
    s.waiting_cnt++
    log.Printf("delay assign,metaKey : %s ,waiting_cnt : %d, running_cnt : %d, preloading_cnt : %d", s.metaData.Key, s.waiting_cnt, s.running_cnt, s.preloading_cnt)
    s.mu.Unlock()
    // s.preload_mu.Lock()
    if(s.preloading_cnt + s.running_cnt < s.waiting_cnt) {
      for i:=0; i < int(5 * math.Pow(2, s.base)); i++ {
        s.base++
        go func() {
          s.Preload(context.Background())
        }()
      }
    }
    // s.preload_mu.Unlock()
    return s.delayAssignLoop(request)    
  } else {
    s.mu.Unlock()
  }

  // if(s.preloading_cnt == 0) {
    for i:=0; i < int(5 * math.Pow(2, s.base)); i++ {
        s.base++
        go func() {
          s.Preload(context.Background())
        }()
    }
  // }
	//Create new Instance
  s.mu.Lock()
  s.cold_start_cnt++
  s.mu.Unlock()
  log.Printf("cold start, metaKey : %s, running_cnt : %d, idle_cnt : %d, waiting_cnt: %d, pred_running : %d ,cold_start_cnt : %d, assign_cnt : %d, preloading_cnt : %d", s.metaData.Key, s.running_cnt, s.idleInstance.Len(), s.waiting_cnt, s.pred_running, s.cold_start_cnt, s.assign_cnt, s.preloading_cnt)
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
  s.running_cnt++
  s.instance_cnt++
	instance.Busy = true
	s.instances[instance.Id] = instance
  start_time := time.Now()
  s.start_time[instance.Id] = start_time.UnixMilli()  
  // if(s.running_cnt > 0) {
  //   s.Preload(ctx)
  // }
  // log.Printf("AddPreloadList : meta mem : %d", instance.Meta.MemoryInMb)
  if(s.running_cnt > s.max_running) {
    s.max_running = s.running_cnt
  }
  if(s.running_cnt > s.max_running_interval) {
    s.max_running_interval = s.running_cnt
  }
	s.mu.Unlock()
	log.Printf("request id: %s, instance %s for app %s is created, init latency: %dms", request.RequestId, instance.Id, instance.Meta.Key, instance.InitDurationInMs)
  s.assign_time = int64(0.2 * float64(s.assign_time) + 0.8 * float64(instance.InitDurationInMs))
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
  s.mu.Lock()
  s.preloading_cnt++
  s.mu.Unlock()
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
  s.mu.Lock()
  s.instance_cnt++
  s.instances[instance.Id] = instance
  s.idleInstance.PushFront(instance)
  s.preloading_cnt--
  s.mu.Unlock()
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
  s.preloadList.PushFront(InstanceConfig)
  log.Printf("AddpreloadList, preloadList len : %d", s.preloadList.Len())
}


func (s *Simple) Idle(ctx context.Context, request *pb.IdleRequest) (*pb.IdleReply, error) {
  // go func(ctx context.Context, request *pb.IdleRequest) (*pb.IdleReply, error){
    if request.Assigment == nil {
	  	return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("assignment is nil"))
	  }
	  reply := &pb.IdleReply{
	  	Status:       pb.Status_Ok,
	  	ErrorMessage: nil,
	  }
	  start := time.Now()
	  instanceId := request.Assigment.InstanceId

	  needDestroy := false
	  slotId := ""
	  // if request.Result != nil && request.Result.NeedDestroy != nil && *request.Result.NeedDestroy {
	  // 	needDestroy = true
	  // }
    // if s.config.PreloadInterval > 0 {
    //    needDestroy = true 
    // }
	  defer func() {
	  if needDestroy {
	  		s.deleteSlot(ctx, request.Assigment.RequestId, slotId, instanceId, request.Assigment.MetaKey, "bad instance")
	  }
	  }()
	  log.Printf("Idle, request id: %s", request.Assigment.RequestId)
	  s.mu.Lock()
    exe_time := time.Now().UnixMilli() - s.start_time[instanceId]
    s.exe_time = int64(float64(s.exe_time) * 0.2 + float64(exe_time) * 0.8)
    log.Printf("metaKey : %s, moving_ave exe_time in ms : %d", s.metaData.Key, exe_time)
    s.running_cnt--
    // if(s.idleInstance.Len() > 0 && s.running_cnt == 0) {
    //   log.Printf("need destroy, reason s.idleInstance.Len() > 0 && s.running_cnt == 0, %d, %d", s.idleInstance.Len(), s.running_cnt)
    //   needDestroy = true
    //   s.instance_cnt--
    // }
	  defer s.mu.Unlock()
	  if instance := s.instances[instanceId]; instance != nil {
	  	slotId = instance.Slot.Id

      // if(s.waiting_cnt > 0) {
      //   log.Printf("do not destroy, reason: s.waiting_cnt > 0")
      //   instance.Busy = false
      //   instance.LastIdleTime = start //keepalive start time
	  	  // s.idleInstance.PushFront(instance)
      //   return reply, nil
      // }

      // if(float32(s.reuse_cnt)/float32(s.assign_cnt) < 0.8) {
      //   log.Printf("do not destroy, reason: float32(s.reuse_cnt)/float32(s.assign_cnt) < 0.8")
      //   instance.Busy = false
      //   instance.LastIdleTime = start //keepalive start time
	  	  // s.idleInstance.PushFront(instance)
      //   return reply, nil
      // }

      // if  instance.InitDurationInMs < int64(s.dt_time) && instance.InitDurationInMs < int64(s.hist.GetQuantiles(0.05)) { //初始化时间 < moving_ave dt && < 5%分位数
      //   needDestroy = true
      //   delete(s.instances, instanceId)
      //   s.instance_cnt--
      //   log.Printf("instance %s destroy, reason: instance.InitDurationInMs < moving_ave(dt), %d, %d", request.Assigment.RequestId, instance.InitDurationInMs, int64(s.dt_time))
      //   return reply ,nil
      // }

      // if s.idleInstance.Len() > s.max_running { //当前idle个数大于最大并发量
      //   needDestroy = true
      //   delete(s.instances, instanceId)
      //   s.instance_cnt--
      //   log.Printf("instance %s destroy, reason: s.idleInstance.Len() > s.max_running, %d, %d", request.Assigment.RequestId, s.idleInstance.Len(), s.max_running)
      //   return reply ,nil

      // }

      // if s.idleInstance.Len() > (s.running_cnt + s.max_running) / 2 { //当前idle个数大于最大running_cnt
      //   needDestroy = true
      //   delete(s.instances, instanceId)
      //   s.instance_cnt--
      //   log.Printf("instance %s destroy, reason: s.idleInstance.Len() > (s.max_running + s.running_cnt) / 2, %d, %d", request.Assigment.RequestId, s.idleInstance.Len(), (s.running_cnt + s.max_running) / 2)
      //   return reply ,nil

      // }

      // if s.pred_running < s.idleInstance.Len() && s.pred_running != 0 {
      //   needDestroy = true
      //   delete(s.instances, instanceId)
      //   s.instance_cnt--
      //   log.Printf("instance %s destroy, reason: s.pred_running < s.idleInstance.Len() , %d, %d", request.Assigment.RequestId, s.pred_running,s.idleInstance.Len())
      //   return reply ,nil
      // }

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

 //  }(ctx, request)


 //  reply := &pb.IdleReply{
	// 	Status:       pb.Status_Ok,
	// 	ErrorMessage: nil,
	// }
 //  return reply ,nil 

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
        // idleDuration := time.Since(instance.LastIdleTime).Milliseconds()
        // log.Printf("check gc, idleDuration : %d, s.config.KeepAliveIntervalInterval: %d", idleDuration, time.Duration(s.config.KeepAliveInterval))

        needgc := false
        // if(s.idleInstance.Len() > s.max_running && s.max_running !=0) {
        //   needgc = true
        // }

        if(s.config.KeepAliveInterval > 0 && idleDuration > time.Duration(s.config.KeepAliveInterval)) {
          needgc = true
        }
        
        // if(s.idleInstance.Len() < s.max_running / 2) {
        //   log.Printf("do not gc, reason : s.idleInstance.Len() < s.max_running / 5,  %d, %d", s.idleInstance.Len(), s.max_running / 5)
        //   needgc = false
        // }

        if(s.pred_running > s.idleInstance.Len() + s.preloading_cnt + s.running_cnt) {
          needgc = false
        }

        if(s.max_running_interval > s.max_running / 2) {
          needgc = false
        }

        if(s.waiting_cnt > 0) {
          needgc = false
        }

				if (needgc){
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
    s.running_history.Enqueue(s.max_running_interval) 
    needPreload := false
    preload_num := 1
    s.base = 0
    data_for_arima := s.running_history.ToArray()
    if(len(data_for_arima) > 10) {
      for i := range(data_for_arima) {
        log.Printf("%d", data_for_arima[i])
      }
      // predictions := predictARIMA(intToFloatArray(data_for_arima), s.p, s.d, s.q, 1)
      // tmp := 0.0
      // end := data_for_arima[len(data_for_arima) - 1]
      // for _, pred := range predictions {
      //   tmp += pred + float64(end)
      // }
      // s.pred_running = int(tmp / float64(len(predictions)))
      s.pred_running = int(my_model.Predict(s.metaData.Key ,intToInt32Array(data_for_arima)))
      log.Printf("lgbm predict : %d", s.pred_running)
      if(s.preloading_cnt + s.idleInstance.Len() + s.running_cnt<  s.pred_running ) {
        needPreload = true
        preload_num = s.pred_running - s.idleInstance.Len() - s.running_cnt - s.preloading_cnt
      }
    } else {
      s.pred_running = s.max_running_interval

      if(s.preloading_cnt + s.idleInstance.Len() + s.running_cnt<  s.pred_running ) {
        needPreload = true
        preload_num = s.pred_running - s.idleInstance.Len() - s.running_cnt - s.preloading_cnt
      }
    }
    s.max_running_interval = 0
    // preload_num = int(math.Max(float64(preload_num), float64(2 * s.max_running - s.idleInstance.Len() - s.running_cnt - s.preloading_cnt)))
    // if(s.running_cnt > s.idleInstance.Len() && s.idleInstance.Len() == 0 && float64(s.dt_time) < float64(s.exe_time)) {
    //   log.Printf("prelod, reason :  s.running_cnt > s.idleInstance.Len() && s.idleInstance.Len() == 0 && float64(s.dt_time) < float64(s.exe_time)")
    //   needPreload = true
    // }
    if(needPreload) {
      for i := 0; i < preload_num; i++{
        log.Printf("preload, history_len : %d, pred_running : %d, idle_len : %d, running : %d", len(data_for_arima) ,  s.pred_running, s.idleInstance.Len(), s.running_cnt)
        go func() {
          s.Preload(context.Background())
        }()
      }
      // for i := 0; i < preload_num; i++ {
      //   log.Printf("preload, preload_num : %d, idle_len : %d", preload_num, s.idleInstance.Len())
      //   s.Preload(context.Background())
      // }
    } else {
      log.Printf("not preload, history_len : %d, pred_running : %d, idle_len : %d, running : %d", len(data_for_arima), s.pred_running, s.idleInstance.Len(), s.running_cnt)
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
