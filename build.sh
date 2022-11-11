docker build -t echo-server:latest .
docker build -t echo-server-gateway:latest . -f Dockerfile.gateway

docker push apiteamdevops/echo-server:latest
docker push apiteamdevops/echo-server-gateway:latest
