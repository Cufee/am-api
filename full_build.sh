clear
git pull

sudo docker stop am-api-dev
sudo docker rm am-api-dev

export PATH=$PATH:/usr/local/go/bin
go build -o build/app

sudo docker build -t cufee/am-api-dev .
sudo docker run -d --restart unless-stopped --name am-api-dev -p 4001:4000 cufee/am-api-dev:latest
