---
id: installation
title: Installation
sidebar_label: Install
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<Tabs
defaultValue="docker"
values={[
{label: 'docker', value: 'docker'},
{label: 'ubuntu', value: 'linux'},
{label: 'kubernetes', value: 'kubernetes'},
]}>

<TabItem value="linux">

- Download [trasa](https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/trasa.tar.gz) binary

```shell script
wget https://storage.googleapis.com/trasa-public-download-assets/release/v0.0.1/trasa.tar.gz
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
sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres
```

- Run [Redis](https://redis.io/download) on port 6379

```shell script
sudo docker run -d -p 6379:6379 --name redis redis
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

- Edit `/etc/trasa/config/config.toml` if needed and restart trasa-server


   </TabItem>
  <TabItem value="docker">

- Run [Postgres](https://www.postgresql.org/) or [CockroachDB](https://cockroachlabs.com) on port 5432

```shell script
sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres
```

- Run [Redis](https://redis.io/download) on port 6379

```shell script
sudo docker run -d -p 6379:6379 --name redis redis
```

- Run guacamole proxy if you use RDP

```shell script
sudo docker run -d --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac --user root  seknox/guacd:v0.0.1
```

- Run trasa-server

```shell script
sudo docker run --link db:db \
--link guacd:guacd \
--link redis:redis \
-p 443:443 \
-p 80:80 \
-p 8022:8022 \
-e TRASA.LISTENADDR=app.trasa \
-v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac \
seknox/trasa:v0.0.1
```


:::tip
Replace app.trasa with hostname/IP of TRASA server.
:::

   </TabItem>
  <TabItem value="kubernetes"> coming soon... </TabItem>

</Tabs>

:::info
Go through [config reference](../system/config-reference) to run TRASA in environment according to your need
:::
