clear
git pull

export PATH=$PATH:/usr/local/go/bin
go build -o build/app

sudo docker build -t cufee/am-api .

sudo docker stop am-api
sudo docker rm am-api

sudo docker run -d --restart unless-stopped --name am-api -p 4000:4000 cufee/am-api:latest
