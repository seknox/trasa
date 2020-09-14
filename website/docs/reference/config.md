---
id: config
title: Config reference
sidebar_label: Config
---

### Database 

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

    Database password


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


#### `sslenabled`
* type: `bool`  
* default: `false`

    Database ssl status

---
### Logging

#### `level`
* type: `string`  
* default: `TRACE`  
    log level
    
 ---
### Minio

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
### Redis

#### `server`
* type: `string`  
* default: `127.0.0.1`  
    Redis server hostname/IP
 



---
### Security

#### `insecureSkipVerify`
* type: `bool`  
* default: `false`  
    Skip ssl verify while making http requests.
    
### Trasa

#### `autoCert`
* type: `bool`  
* default: `false`  
    Use autoCert to generate signed certificates from Let's Encrypt.


 
#### `listenAddr`
* type: `string`  
* default: `localhost`  


 
#### `cloudServer`
* type: `string`  
* default: `https://sg.cpxy.trasa.io`  
    Address of proxy server to forward U2F requests

---

### Proxy

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

### Vault

#### `tsxvault`
* type: `bool`  
* default: `true`  
