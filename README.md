#### 不用代理:
docker build -t vps-mall:latest -f ./Dockerfile .

#### 使用代理:
配置docker使用proxy
mkdir /etc/systemd/system/docker.service.d
vi /etc/systemd/system/docker.service.d/http-proxy.conf
[Service]
Environment="HTTP_PROXY=http://192.168.0.132:1081/"
Environment="HTTPS_PROXY=http://192.168.0.132:1081/"

docker build --build-arg HTTP_PROXY=http://192.168.0.132:1081 --build-arg HTTPS_PROXY=http://192.168.0.132:1081 -t vps-mall:latest -f ./Dockerfile .

### RUN:
docker run -d --name vps-mall -p 5577:5577 vps-mall:latest