# Use a multi-stage build to keep the final image size small.
# In the first stage, we'll use the official Go image to build the binary.
# FROM golang:1.17.1 as builder
FROM registry.cn-beijing.aliyuncs.com/cloudnative-challenge/golang:1.17.1 as builder

# Set the working directory to /build.
WORKDIR /build

# Copy the source code excluding the 'data' directory
COPY .. ./
RUN go build -o scaler cmd/scaler/main.go

# In the second stage, we'll use a small, lightweight base image.
FROM registry.cn-beijing.aliyuncs.com/cloudnative-challenge/ubuntu:latest
# FROM python:3.8.5-slim-buster

# FROM python:3.9.1
FROM registry.cn-shanghai.aliyuncs.com/2024-happy-hack/python3.9:v1.0

# Update the package list and install required packages.
# RUN apt-get update 
# RUN apt-get install -y netcat curl
# RUN pip install pmdarima

# Set the working directory to /app.
WORKDIR /app


COPY --from=builder /build/scaler /app/scaler
RUN chmod +x /app/scaler

# Copy the source code into the container, excluding the 'data' directory
COPY --from=builder /build/python /app
COPY --from=builder /build /app/source


# Copy the startup script.
COPY run.sh run.sh
RUN chmod +x run.sh
