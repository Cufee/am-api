clear
git pull

sudo docker stop am-api
sudo docker rm am-api

export PATH=$PATH:/usr/local/go/bin
go build -o build/app

sudo docker build -t cufee/am-api .
sudo docker run -d --name am-api -p 4000:4000 cufee/am-api:latest
