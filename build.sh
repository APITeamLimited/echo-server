docker build -t apiteamdevops/echo-server:latest .
docker build -t apiteamdevops/echo-server-gateway:latest . -f Dockerfile.gateway

docker push apiteamdevops/echo-server:latest
docker push apiteamdevops/echo-server-gateway:latest
