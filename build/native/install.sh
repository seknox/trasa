sudo apt-get -y update && \

#for parsing json
sudo apt-get -y install jq && \

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






INIT_CONFIG=`cat initconfig.json ` && \
TRASA_SIGNUP_DATA=`echo $INIT_CONFIG | jq -r '.signup'` && \
TRASA_HOSTNAME=`echo $INIT_CONFIG | jq -r '.hostname'` && \


#echo INIT_CONFIG
echo $TRASA_SIGNUP_DATA
echo $TRASA_HOSTNAME

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
sudo cp -r  dashboard /etc/trasa/build


sudo cp -r config /etc/trasa/config && \
sudo cp -r static /etc/trasa/static && \


sudo cp service-files-single-binary/* /etc/systemd/system && \
sudo systemctl daemon-reload && \


#sudo docker load -i bins/guacd.tar && \






#mkdir certs my-safe-directory
#cockroach cert create-ca --certs-dir=certs --ca-key=my-safe-directory/ca.key
#cockroach cert create-node localhost beta.trasa.io $(hostname) --certs-dir=certs --ca-key=my-safe-directory/ca.key
#



#Start services
sudo systemctl start cockroach && \
sudo systemctl start minio && \
sudo systemctl start redis && \
sudo systemctl start guacd && \




#create database and user

sudo systemctl start trasa



sleep 10


#signup and get orgID
RESPONSE=`curl  --insecure 'http://localhost/idp/signup' -H 'content-type: application/json' --data "$TRASA_SIGNUP_DATA"` && \

echo $RESPONSE

STATUS=`echo $RESPONSE | jq '.status'`


#if [ "$STATUS"  !=  "success"]
#then
#  echo "Failed to sign up"
#  exit
#endif

ORG_ID=`echo $RESPONSE | jq '.data' | jq -r '.[0]'`


echo $ORG_ID


#write orgID in config
sudo sed -i -e 's|<ORG_ID>|'$ORG_ID'|g' /etc/trasa/config/config.toml && \


sudo systemctl restart trasa
#'{"email":"bhargab@seknox.com","timezone":"Asia/Kathmandu","phoneNumber":9802084613,"orgName":"Seknox_Staging","userName":"bhrg3se","country":"Nepal"}'



echo ""
echo ""
echo ""
echo "______________________________________"
echo "DONE"

