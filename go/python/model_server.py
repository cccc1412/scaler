import sys
import grpc
import model_pb2
import model_pb2_grpc
from concurrent import futures
import predict_lgbm as model

SEQUENCE_SIZE = 10

class Model(model_pb2_grpc.ModelServicer):
    def Predict(self, request, context):
        if len(request.sequences) != SEQUENCE_SIZE:
            return model_pb2.PredictResponse(result=-1)
        result = model.Predict(request.meta_key, request.sequences)
        return model_pb2.PredictResponse(result=result)

def serve():
    sequence_size = 10
    window_size = 60
    if len(sys.argv) > 1:
        global SEQUENCE_SIZE
        SEQUENCE_SIZE = int(sys.argv[1])
        sequence_size = int(sys.argv[1])
    if len(sys.argv) > 2:
        window_size = int(sys.argv[2])
    model.Run(sequence_size, window_size)
    server = grpc.server(futures.ThreadPoolExecutor())
    model_pb2_grpc.add_ModelServicer_to_server(Model(), server)
    server.add_insecure_port('[::]:9030')
    server.start()
    print('Model server start at 127.0.0.1:9030')
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
