#!/usr/bin/env sh


#Install postgres
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get install postgresql

#Install docker
sudo apt-get -y remove docker docker-engine docker.io containerd runc && \
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


sudo apt-get -y install ffmpeg && \


sudo mkdir -p /etc/trasa && \
sudo mkdir -p /etc/trasa/certs && \

#chown $USER /etc/trasa && \
sudo mkdir -p /var/trasa && \

sudo mkdir -p /var/trasa/crdb

sudo mkdir -p /var/trasa/minio && \



#generate rsa keys for trasagw
sudo ssh-keygen -t rsa -b 4096 -f /etc/trasa/certs/id_rsa -q -N ""


#copy binaries to binchow
chmod +x bins/* && \
sudo cp bins/* /usr/local/bin/ && \
sudo cp -r  dashboard /var/trasa/dashboard


sudo cp -r config /etc/trasa/config && \
sudo cp -r static /etc/trasa/static && \


sudo cp service-files-single-binary/* /etc/systemd/system && \
sudo systemctl daemon-reload && \



#Start services
sudo systemctl start cockroach && \
sudo systemctl start minio && \
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

