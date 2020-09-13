#!/usr/bin/env sh

TRASA_VERSION=0.0.1

#Install postgres
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install postgresql

#Install docker
sudo apt-get -y remove docker docker-engine docker.io containerd runc
sudo apt-get -y update && \
sudo apt-get -y install  \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common && \
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - && \
sudo apt-key fingerprint 0EBFCD88 && \
sudo add-apt-repository -y \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable" && \
sudo apt-get -y update && \
sudo apt-get -y install docker-ce docker-ce-cli containerd.io && \

sudo apt-get install make
sudo apt-get -y install ffmpeg && \



wget http://download.redis.io/releases/redis-6.0.8.tar.gz
tar xzf redis-6.0.8.tar.gz
cd redis-6.0.8 && make && cp src/redis-server /usr/local/bin/ && cd ..


sudo mkdir -p /etc/trasa/config && \
sudo mkdir -p /etc/trasa/certs && \
sudo mkdir -p /etc/trasa/static && \

#chown $USER /etc/trasa && \
sudo mkdir -p /var/trasa && \
sudo mkdir -p /var/trasa/crdb
sudo mkdir -p /var/trasa/minio && \




mkdir bins

wget https://storage.googleapis.com/trasa-public-download-assets/release/v$TRASA_VERSION/trasa-server -O bins/trasa-server
wget https://storage.googleapis.com/trasa-public-download-assets/release/v$TRASA_VERSION/dashboard.tar -O dashboard.tar

tar xzf dashboard.tar


#copy binaries to binchow
chmod +x bins/* && \
sudo cp bins/* /usr/local/bin/ && \
sudo cp -r  dashboard/dashboard /var/trasa/dashboard

wget https://raw.githubusercontent.com/seknox/trasa/master/build/etc/trasa/config/config.toml

sudo mv config.toml /etc/trasa/config/config.toml

wget https://storage.googleapis.com/trasa-public-download-assets/GeoLite2-City.mmdb
sudo mv GeoLite2-City.mmdb /etc/trasa/static/GeoLite2-City.mmdb



mkdir service-files
wget https://raw.githubusercontent.com/seknox/trasa/master/build/native/trasa.service -O service-files/trasa.service
wget https://raw.githubusercontent.com/seknox/trasa/master/build/native/trasa.service -O service-files/redis.service

sudo cp service-files/* /etc/systemd/system
sudo systemctl daemon-reload && \


psql <<- EOSQL
    CREATE USER docker;
    CREATE DATABASE docker;
    GRANT ALL PRIVILEGES ON DATABASE docker TO docker;
EOSQL

#Start services
sudo systemctl start postgresql && \
sudo systemctl start redis && \
sudo systemctl start guacd && \




#create database and user

sudo systemctl start trasa



sleep 10



sudo systemctl restart trasa


echo ""
echo ""
echo ""
echo "______________________________________"
echo "DONE"

