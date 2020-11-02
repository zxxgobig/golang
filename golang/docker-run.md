
docker build -t docker_go:0.1.1 .

docker run --name e_go1 --rm -it docker_go:0.1.0 bash

docker run -p 63333:63333 -v /opt/go/:/opt/go/:rw --name dhr --rm -it docker_go:0.1.0 bash

docker run -p 63333:63333 -v /opt/go/:/opt/go/:rw --name dhr --rm -it docker_go:0.1.0 bash



git config --global url."git@gitlab.ifchange.com:".insteadOf "https://gitlab.ifchange.com"
git config --global user.name "xiaoxiao"
git config --global user.email "xiaoxiao.zheng@ifchange.com"

docker run -d -p 63333:63333 -p 8002:8002 -v /opt/go:/opt/go/:rw --name dhr --rm -it docker_go:0.1.0 bash


ssh-keygen -t ed25519 -C "xiaoxiao.zheng@ifchange.com"