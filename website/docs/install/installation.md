---
id: installation
title: Installation Guide
sidebar_label: Installation
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


## Requirements
To run TRASA with a bare minimal setup, you will need :
1. Database (PostgreSQL), 
2. Redis,  
3. Guacamole Guacd server (Only required If you need to protect Remote Desktop Access),
4. TRASA server itself. 


The minimum server requirement to run TRASA is:
- 1 core CPU,
- 1 GB ram,
- 20 GB storage.


:::important
Planning for extra storage space is critical since session recording of remote access can quickly cover up available spaces. It depends on frequency of access, policy that enables session recording. 
:::


<Tabs
defaultValue="docker"
values={[
{label: 'docker', value: 'docker'},
{label: 'ubuntu', value: 'linux'},
{label: 'kubernetes', value: 'kubernetes'},
]}>

<TabItem value="linux">

- Download [trasa](https://storage.googleapis.com/trasa-public-download-assets/release/v1.1.1/trasa.tar.gz) binary

```shell script
wget https://storage.googleapis.com/trasa-public-download-assets/release/v1.1.1/trasa.tar.gz
```

- Extract and place static files into respective dirs

```shell script
mkdir -p /var/trasa/dashboard
mkdir -p /etc/trasa/static

tar -xzf trasa.tar.gz
mv  trasa/dashboard /var/trasa/
mv trasa/GeoLite2-City.mmdb /etc/trasa/static/
```

- Run [Postgres](https://www.postgresql.org/) or [CockroachDB](https://cockroachlabs.com) on port 5432

```shell script
sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres:13.0
```

- Run [Redis](https://redis.io/download) on port 6379

```shell script
sudo docker run -d -p 6379:6379 --name redis redis:6.0.8
```


- Run guacamole proxy if you use rdp

```shell script
sudo docker run -d --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac --user root seknox/guacd:v0.0.1
```


- Run trasa-server binary

```shell script
sudo ./trasa/trasa-server
```

> Add -f argument while running trasa-server to enable logging to file /var/log/trasa.log

- Edit `/etc/trasa/config/config.toml`
    - Change `trasa.listenAddr` to the domain/IP of server.
    - If `trasa.listenAddr` is not a public domain, turn off the autocert in config by setting `trasa.autoCert=false`

- Restart the trasa-server


   </TabItem>
  <TabItem value="docker">

- Run [Postgres](https://www.postgresql.org/) or [CockroachDB](https://cockroachlabs.com) on port 5432

```shell script
sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres:13.0
```

- Run [Redis](https://redis.io/download) on port 6379

```shell script
sudo docker run -d -p 6379:6379 --name redis redis:6.0.8
```

- Run guacd (Apache Guacamole RDP proxy server). This is only required if you need to protect RDP service

```shell script
sudo docker run -d --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac --user root  seknox/guacd:v0.0.1
```

- Run trasa-server

:::tip
- Replace app.trasa in `TRASA.LISTENADDR` with hostname/IP you want TRASA server to listen on.
- If `TRASA.LISTENADDR` is not a public domain, turn off the autocert by passing `-e TRASA.AUTOCERT="false"`
:::

```shell script
sudo docker run -d --link db:db \
--link guacd:guacd \
--link redis:redis \
-p 443:443 \
-p 80:80 \
-p 8022:8022 \
-e TRASA.LISTENADDR=app.trasa \
-v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac \
seknox/trasa:v1.1.2
```




   </TabItem>
  <TabItem value="kubernetes"> coming soon... </TabItem>

</Tabs>

:::info
Go through [config reference](../system/config-reference) to run TRASA in environment according to your need
:::
