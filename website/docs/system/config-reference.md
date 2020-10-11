---
id: config-reference
title: Config reference
sidebar_label: Config
---


Config file is located at `/etc/trasa/config/config.toml`

:::tip
You can override config file values with environment variables.
e.g. To override logging.level, you need to set env variable LOGGING.LEVEL
:::



### database 

#### `dbname`
* type: `string`  
* default: `trasadb`

    Database name


#### `dbuser`
* type: `string`  
* default: `trasauser`

    Database username

#### `dbpass`
* type: `string`  
* default: `trasauser`

    Database password


#### `dbtype`
* type: `enum string {"postgres","cockroachdb"}`  
* default: `cockroachdb`  
    Database type


#### `server`
* type: `string`  
* default: `127.0.0.1`

    Database hostname/ip


#### `port`
* type: `string`  
* default: `5432`

    Database port


#### `sslenabled`
* type: `bool`  
* default: `false`

    Database ssl status


#### `usercert`
* type: `string`  
* default: ``

    Database ssl user certificate file path (only applicable if `sslenabled=true`). 


#### `userkey`
* type: `string`  
* default: ``

    Database ssl user private key file path (only applicable if `sslenabled=true`). 


#### `cacert`
* type: `string`  
* default: ``

    Database ssl ca certificate file path (only applicable if `sslenabled=true`). 



---
### logging

#### `level`
* type: `string`  
* default: `TRACE`  
    log level
    
 ---
### minio

#### `status`
* type: `bool`  
* default: `false`  
    Enable minio. If disabled local file system is used to store session recordings.

#### `server`
* type: `string`  
* default: `127.0.0.1`  
    Minio server hostname/IP    
 
#### `key`
* type: `string`  
* default: `minioadmin`  
    Minio access key
    
 
#### `secret`
* type: `string`  
* default: `minioadmin`  
    Minio access secret
    
 
#### `sslenabled`
* type: `bool`  
* default: `false`  
    Enable ssl in minio 
    
 
---
### redis

#### `server`
* type: `string`  
* default: `127.0.0.1`  
    Redis server hostname/IP
 



---
### security

#### `insecureSkipVerify`
* type: `bool`  
* default: `false`  
    Skip ssl verify while making http requests.
---     
### trasa

#### `autoCert`
* type: `bool`  
* default: `true`  
    Use autoCert to generate signed certificates from Let's Encrypt.


 
#### `listenAddr`
* type: `string`  
* default: `localhost`  
    Listen address of TRASA server
 
#### `cloudServer`
* type: `string`  
* default: `https://sg.cpxy.trasa.io`  
    Address of proxy server to forward U2F requests


#### `proxyDashboard`
* type: `bool`  
* default: `false`  
    If enabled TRASA server will proxy dashboard instead of serving. It is useful during development when dashboard is served using `npm start` command.
    
#### `dashboardAddr`
* type: `bool`  
* default: `localhost:3000`  
    If proxyDashboard is enabled TRASA server will proxy dashboard from this address.
    



---

### proxy

#### `sshlistenAddr`
* type: `string`  
* default: `:8022`  
     Listen address of ssh proxy


#### `guacdAddr`
* type: `string`  
* default: `127.0.0.1:4822`  
     Address of guacd


#### `guacdEnabled`
* type: `bool`  
* default: `false`  
     
     
---

### vault

#### `tsxvault`
* type: `bool`  
* default: `true`  
