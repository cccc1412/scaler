kubectl delete --all pods
kubectl delete --all jobs
sudo docker buildx build --platform linux/amd64 -t registry.cn-shanghai.aliyuncs.com/2024-happy-hack/scaler:v1.9.44 . --push
kubectl apply -f ../manifest/serverless-simulation.yaml 
