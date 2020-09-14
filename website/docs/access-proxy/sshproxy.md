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



TRASA can serve as gateway between client and upstream service.
TRASA will
There are different methods for accessing services through TRASA ssh proxy.


### Web Access (Dashboard)

* Login into your TRASA account
* Click connect and choose service username
* Enter password and TOTP 
> You may be asked to save new host key


### SSH clients

#### Download User Key (Optional)
* Go to dashboard and go to "Account" tab
* Click the menu to get dropdown menu items
* Click the "get ssh private key" button" to download key

:::tip 
If you're using Putty, use PuTTYgen to convert downloaded key `id_rsa` into `id_rsa.ppk`
:::

<Tabs
    defaultValue="openssh"
    values={[
  
            {label: 'OpenSSH Client', value: 'openssh'},
            {label: 'TRASA cli', value: 'trasacli'},
            {label: 'Putty', value: 'putty'},
            {label: 'Bitvise', value: 'bitvise'},
        ]}
>

<TabItem value="openssh">

* `ssh -i <private_key_path> root@trasa.hostname -p 8022`     
* Enter TRASA email and password   
* Enter IP address of service you want to connect to   
* Enter TOTP code or leave it blank for U2F   
* Enter Service password   

</TabItem>

<TabItem value="trasacli">

* [Setup device agent](#)
* `trasacli -u username`
* Enter TRASA URL if asked     
* Enter TRASA email and password   
* Enter IP address of service you want to connect to   
* Enter TOTP code or leave it blank for U2F   
* Enter Service password   

</TabItem>


<TabItem value="putty">

* Enter the TRASA hostname or IP address under Session.
* Navigate to Connection > SSH > Auth.
* Click Browse... under Authentication parameters / Private key file for authentication.
* Locate the id_rsa.ppk private key and click Open.
* Click Open again to log into the remote server with key pair authentication.

</TabItem>


  <TabItem value="bitvise">
  todo
  </TabItem>


</Tabs>



## Upstream Service Authentication
TRASA supports password, private key and certificate to authenticate upstream server.


### 1. Using Private Key
If you have private key of a service instead of password, you should [save the key in vault](https://trasa.io/docs/secret-vault/index)

##### Or, use agent forwarding.

Since TRASA ssh gateway doesn't have access to user private key you need to enable agent forwarding to use private keys/certificate while logging into upstream server

>SSH Agent forwarding is not recommended since it allows proxy server to use your SSH keys and also eavesdrop on your ongoing session.
>But in our case TRASA ssh proxy should and is able to record all the traffic. So you need to trust TRASA server.


* Add the private key to ssh agent `ssh-add <private_key_path>`
* Use -A flag `ssh -A -i <private_key_path> root@trasa.hostname -p 8022`



:::tip
You can use trasa private key and service private key at the same time. `ssh -A -i <trasa_private_key_path> -i <service_private_key_path> root@trasa.hostname -p 8022`
:::

:::tip
* If you have been using google cloud, gcloud saves private key in `~/.ssh/google_compute_engine`. You can use that to login through TRASA.
:::



### 2. Using SSH Certificates
You can use TRASA as a ssh CA

#### 2.1 Initialize CA
To use SSH certificates you must first  [initialise CA](/trasa/docs/guides/ca) (if you haven't already) from TRASA dashboard

#### 2.2 Using User Certificates
TRASA ssh proxy will generate and use a temporary certificate signed by itself while connecting any upstream service.
So to make use of that user certificate, you must tell each upstream server to trust any certificate signed our CA.
To do that,
##### 2.2.1 Distribute CA public keys
* Download client CA  from dashboard (Providers->Certificate Authority->Download SSH client CA)
* Copy the downloaded ssh keys into upstream servers
* Edit /etc/ssh/sshd_config of upstream server and add the folowing
`TrustedUserCAKeys <path to ca public key>`
* Restart ssh daemon
`sudo systemctl restart sshd`




#### 2.3 Host Certificates

##### 2.3.1 Configure Client Device
* Download host CA  from dashboard (manage->Certificate Authority->Download SSH host CA)
* Copy its contents to /etc/ssh/ssh_known_hosts

##### 2.3.2 Configure Upstream Server
* Go to authapp page in dashboard
* In App Config tab, download "Generate and Download" button
* Copy the downloaded zip file to upstream server
* Extract the files into /etc/ssh
* Edit /etc/ssh/sshd_config and add the following
`HostKey /etc/ssh/id_rsa`
`HostCertificate /etc/ssh/id_rsa-cert.pub`
* Restart sshd daemon
`sudo systemctl restart sshd`






