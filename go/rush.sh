sudo docker buildx build --platform linux/amd64 -t registry.cn-shanghai.aliyuncs.com/2024-happy-hack/scaler:v1.1.4 . 
kind load docker-image registry.cn-shanghai.aliyuncs.com/2024-happy-hack/scaler:v1.1.4
# kind load docker-image registry.cn-shanghai.aliyuncs.com/2024-happy-hack/my_simulator:v1
docker image prune
kubectl delete -f ../manifest//serverless-simulation.yaml
kubectl apply -f ../manifest/serverless-simulation.yaml 
