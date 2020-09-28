---
id: ssh-service
title: Secure Shell (SSH)
sidebar_label: SSH
---

import useBaseUrl from '@docusaurus/useBaseUrl';




There are two ways to protect SSH servers.
1. Through TRASA access proxy
2. Native 2FA agents 

#### 1. Native 2FA agents
You need to install and configure 2fa agents in all SSH servers you want to protect.
[This guide](../../native-tfa/linux-two-factor-authentication.md) will help you configure native agents.


#### 2. SSH Access Proxy
By default TRASA will listen at port 8022. You can change that in [config](../system/config-reference.md#sshlistenaddr) if you want.





### Using SSH Certificates

You can use TRASA as a SSH CA. TRASA access proxy injects a temporary signed certificate with expiry of few minutes. 
If you configure upstream servers to trust TRASA CA, they will accept any ssh key(certificate) signed by TRASA CA.
This makes remote access very easy and secure since user doesn't need ko know password or store keys.





#### Initialize CA
To use SSH certificates you must first  [initialise CA](/trasa/docs/guides/ca) (if you haven't already) from TRASA dashboard

* Go to Providers -> Certificate Authority page
<img alt="download-user-ca" src={useBaseUrl('img/docs/providers/providers-menu.png')} />  
<img alt="ca-tab" src={useBaseUrl('img/docs/providers/ca/ca-tab.png')} />  

* Click the "Generate certs" button
<img alt="generate-ca-btn" src={useBaseUrl('img/docs/providers/ca/generate-ca-btn.png')} />  
* Generate both "SSH User CA" and "SSH Host CA"
<img alt="generate-ca-dialog" src={useBaseUrl('img/docs/providers/ca/generate-ca-dialog.png')} />  


#### Using User Certificates
TRASA ssh proxy will generate and use a temporary certificate signed by itself while connecting any upstream service.
To make use of that user certificate, you must tell each upstream server to trust any certificate signed by our CA.
To do that,


* Download client CA  from dashboard (Providers->Certificate Authority->Download SSH client CA)

<img alt="download-user-ca" src={useBaseUrl('img/docs/providers/ca/download-user-ca.png')} />  

* Download client CA

* Copy the downloaded ssh keys into upstream servers
* Edit /etc/ssh/sshd_config of upstream server and add the following
`TrustedUserCAKeys <path to ca public key>`
* Restart ssh daemon
`sudo systemctl restart sshd`




#### Host Certificates

##### Configure Client Device
* Download host CA  from dashboard (manage->Certificate Authority->Download SSH host CA)
<img alt="download-host-ca" src={useBaseUrl('img/docs/providers/ca/download-host-ca.png')} />  

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


