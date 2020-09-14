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
    
 
