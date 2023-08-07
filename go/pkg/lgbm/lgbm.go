package lgbm


import (
    "context"
    "log"

    "google.golang.org/grpc"
    model_pb "github.com/AliyunContainerService/scaler/model_proto"
)

var Model model_pb.ModelClient

func Init() {
	conn, err := grpc.Dial("[::]:9030", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
        log.Fatalf("Connect python server failed: %+v", err)
  }
  Model = model_pb.NewModelClient(conn)
	log.Println("Connect python server success")
}

func Predict(meta_key string, sequences []int32) int32  {
	if Model == nil {
    log.Printf("Model is nil")
		return -1
	}
  resp, err := Model.Predict(context.Background(), &model_pb.PredictRequest{MetaKey: meta_key, Sequences: sequences})
  if err != nil {
      log.Fatalf("Predict service failed: %+v", err)
  }
  log.Printf("Predict : %d", resp.Result)

  return resp.Result
}
