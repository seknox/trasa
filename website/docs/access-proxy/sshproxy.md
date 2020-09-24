---
id: sshproxy
title: SSH Proxy
sidebar_label: SSH Proxy
draft: true
keywords: 
    - ssh
    - certificate
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';



TRASA can serve as a gateway between the client and the upstream service.
TRASA will
There are different methods for accessing services through the TRASA ssh proxy.





## Upstream Service Authentication
TRASA supports password, private key, and certificate to authenticate the upstream server.


### 1. Saving Credentials in Vault
You can [store private key of a service or password in vault](https://trasa.io/docs/secret-vault/index)



### 2. Using SSH Certificates
You can use TRASA as a SSH CA

#### 2.1 Initialize CA
To use SSH certificates you must first  [initialise CA](/trasa/docs/guides/ca) (if you haven't already) from TRASA dashboard

#### 2.2 Using User Certificates
TRASA ssh proxy will generate and use a temporary certificate signed by itself while connecting any upstream service.
To make use of that user certificate, you must tell each upstream server to trust any certificate signed by our CA.
To do that,
##### 2.2.1 Distribute CA public keys
* Download client CA  from dashboard (Providers->Certificate Authority->Download SSH client CA)
* Copy the downloaded ssh keys into upstream servers
* Edit /etc/ssh/sshd_config of upstream server and add the following
`TrustedUserCAKeys <path to ca public key>`
* Restart ssh daemon
`sudo systemctl restart sshd`




#### 2.3 Host Certificates

##### 2.3.1 Configure Client Device
* Download host CA  from dashboard (manage->Certificate Authority->Download SSH host CA)
* Copy its contents to /etc/ssh/ssh_known_hosts

##### 2.3.2 Configure Upstream Server
* Go to service page in dashboard
* In App Config tab, download "Generate and Download" button
* Copy the downloaded zip file to upstream server
* Extract the files into /etc/ssh
* Edit /etc/ssh/sshd_config and add the following
`HostKey /etc/ssh/id_rsa`
`HostCertificate /etc/ssh/id_rsa-cert.pub`
* Restart sshd daemon
`sudo systemctl restart sshd`






