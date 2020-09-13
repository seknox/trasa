---
id: installation
title: Installation
sidebar_label: Install
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs
  defaultValue="linux"
  values={[
    {label: 'docker', value: 'docker'},
    {label: 'ubuntu', value: 'linux'},
    {label: 'kubernetes', value: 'kubernetes'},
  ]}>
  
<TabItem value="linux">


* Download [trasa-server](https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/trasa-server) binary
```shell script
wget https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/trasa-server
chmod +x trasa-server
```


* Download [dashboard](https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/dashboard.tar) and extract into /var/trasa/dashboard 
```shell script
wget https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/dashboard.tar
mkdir -p /var/trasa/dashboard
tar -C /var/trasa -xf dashboard.tar
```


* Download [GeoLiteCity db](https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb) and move to /etc/trasa/static 
```shell script
wget https://https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb
mkdir -p /etc/trasa/static
mv GeoLite2-City.mmdb /etc/trasa/static/
```

* Run [Postgres](https://www.postgresql.org/) or [cockroachdb](https://cockroachlabs.com) on port 5432
```shell script
docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres
```
* Run [Redis](https://redis.io/download) on port 6379 
```shell script
docker run -d -p 6379:6379 --name redis redis
```

* Run trasa-server binary
```shell script
./trasa-server
```
>Add -f while running trasa-server to enable logging to file /var/log/trasa.log


* Edit `/etc/trasa/config/config.toml` if needed and restart trasa-server

* Run guacamole proxy if you use rdp
```shell script
docker run  --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy:/tmp/trasa/accessproxy -v /tmp/trasa/accessproxy/guac/shared/:/tmp/trasa/accessproxy/guac/shared/  seknox/guacd:v0.0.1
```


   
   </TabItem>
  <TabItem value="docker"> 

* Run [Postgres](https://www.postgresql.org/) or [cockroachdb](https://cockroachlabs.com) on port 5432
```shell script
docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres
```
* Run [Redis](https://redis.io/download) on port 6379 
```shell script
docker run -d -p 6379:6379 --name redis redis
```

* Run trasa-server 
```shell script
docker run --link db:db --link guacd:guacd --link redis:redis -p 443:443 -p 80:80 -p 8022:8022 seknox/trasa:v0.0.1
```



* Run guacamole proxy if you use rdp
```shell script
docker run  --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy:/tmp/trasa/accessproxy -v /tmp/trasa/accessproxy/guac/shared/:/tmp/trasa/accessproxy/guac/shared/  seknox/guacd:v0.0.1
```



Map `/etc/trasa/config/` to custom config file if needed 
```shell script
docker run --link db:db --link guacd:guacd --link redis:redis -p 443:443 -p 80:80 -p 8022:8022 -v my:/etc/trasa/config/ seknox/trasa:v0.0.1
```


   
   </TabItem>
  <TabItem value="kubernetes"> coming soon... </TabItem>


    

</Tabs>
