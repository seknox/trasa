---
id: ssh-service
title: Secure Shell (SSH)
sidebar_label: SSH
---

import useBaseUrl from '@docusaurus/useBaseUrl';




There are two ways to protect SSH servers.
1. Through TRASA access proxy
2. Native 2FA agents 

### 1. Native 2FA agents
You need to install and configure 2fa agents in all SSH servers you want to protect.
[This guide](../../native-tfa/linux-two-factor-authentication.md) will help you configure native agents.


### 2. SSH Access Proxy
To use TRASA as an SSH proxy, you need to [configure firewall rules](../../install/initial-setup.md#3-firewall-configuration-optional) to enforce ssh access from TRASA only.

Now users need to SSH into TRASA proxy instead of the upstream server.
```shell script
ssh user@TRASA_HOSTNAME -p 8022
``` 
Here the port 8022 is the default port of TRASA proxy.
:::tip 
You can change the default port in [config](../../system/config-reference.md#sshlistenaddr) if you want.
:::

Learn more about accessing SSH proxy [here](../../guides/user/access/ssh-connection-via-proxy.md)




## Configuring SSH Authentication (Optional)
If you configure SSH authentication, users don't need to enter upstream server password while accessing through TRASA proxy.

You can 
* Store password in vault
* Store private key in vault
* Configure SSH User CA

### Store Password/Keys in vault
Follow [this guide](../../providers/secret-vault/index.md#storing-service-credentials) to configure and store credentials in vault

### Using SSH Certificates

You can use TRASA as a SSH CA. TRASA access proxy injects a temporary signed certificate with expiry of few minutes. 
If you configure upstream servers to trust TRASA CA, they will accept any ssh key(certificate) signed by TRASA CA.
This makes remote access very easy and secure since user doesn't need ko know password or store keys.





#### Initialize CA
To use SSH certificates you must first  [initialise CA](/trasa/docs/guides/ca) (if you haven't already) from TRASA dashboard

* Go to Providers -> Certificate Authority page
<img alt="download-user-ca" src={('/img/docs/providers/providers-menu.png')} />  
<img alt="ca-tab" src={('/img/docs/providers/ca/ca-tab.png')} />  

* Click the "Generate certs" button
<img alt="generate-ca-btn" src={('/img/docs/providers/ca/generate-ca-btn.png')} />  
* Generate both "SSH User CA" and "SSH Host CA"
<img alt="generate-ca-dialog" src={('/img/docs/providers/ca/generate-ca-dialog.png')} />  


#### Using User Certificates
TRASA ssh proxy will generate and use a temporary certificate signed by itself while connecting any upstream service.
To make use of that user certificate, you must tell each upstream server to trust any certificate signed by our CA.
To do that,


* Download client CA  from dashboard (Providers->Certificate Authority->Download SSH client CA)

<img alt="download-user-ca" src={('/img/docs/providers/ca/download-user-ca.png')} />  

* Download client CA

* Copy the downloaded ssh keys into upstream servers
* Edit /etc/ssh/sshd_config of upstream server and add the following
`TrustedUserCAKeys <path to ca public key>`
* Restart ssh daemon
`sudo systemctl restart sshd`




#### Host Certificates
TRASA proxy will automatically validate host keys and certificates when accessing through TRASA proxy.

##### Configure Client Device
Configuring client device is applicable when accessing SSH server directly instead through TRASA proxy. 

* Download host CA  from dashboard (manage->Certificate Authority->Download SSH host CA)
<img alt="download-host-ca" src={('/img/docs/providers/ca/download-host-ca.png')} />  

* Copy its contents to /etc/ssh/ssh_known_hosts

##### Configure Upstream Server
* Go to service page in dashboard
* In App Config tab, download "Generate and Download" button
* Copy the downloaded zip file to upstream server
* Extract the files into /etc/ssh
* Edit /etc/ssh/sshd_config and add the following
`HostKey /etc/ssh/id_rsa`
`HostCertificate /etc/ssh/id_rsa-cert.pub`
* Restart sshd daemon
`sudo systemctl restart sshd`


