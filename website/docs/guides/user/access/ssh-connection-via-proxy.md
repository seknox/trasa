---
id: ssh-connection-via-proxy
title: Access SSH Service
sidebar_label: SSH
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Access methods


### Direct Access
When the linux server or workstation is protected with TRASA 2FA agent, just after the username(privilege) and password validation in the login screen, a prompt will appear on your screen. 
```shell script
$ ssh root@host
$ password:
$ Enter your trasaID: 
$ Choose TFA method (enter blank for U2F):
```
You will have to enter your trasaID(email or username) and enter the TOTP code to perform second step verification.


### Via TRASA access proxy
You can access the SSH service either via Browser or SSH client.

#### Using Browser

* Login into your TRASA account.
* Search for the service you want to connect to.
* Click connect and choose the service username.
* Enter password and TOTP. 
> You may be asked to save the new host key.





#### Using SSH clients

* `ssh -i <private_key_path> root@TRASA_HOST -p 8022`     
<img alt="ssh-proxy-email" src={('/img/docs/user-guides/access/ssh-proxy-email.png')} />  

* Enter the TRASA email and password   
* Enter the IP address of the service you want to connect to.   
* Enter TOTP code or leave it blank for U2F.   
* Enter the Service password.   
> You may be asked to save the new host key.



:::note
Download a TRASA user key to save you from entering TRASA email and password every time you use SSH proxy.

* Go to the dashboard and go to the "Account" tab.
* Click the menu to get dropdown menu items.
* Click the "get ssh private key" button" to download the SSH key.
* If you're using PuTTY, use PuTTYgen to convert the downloaded key `id_rsa` into `id_rsa.ppk`.


>This key is used to authenticate to TRASA server, NOT the upstream SSH server.
>So, you might still be asked for an upstream password.

:::


<!---

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
* Enter the IP address of service you want to connect to   
* Enter TOTP code or leave it blank for U2F   
* Enter Service password   

</TabItem>


<TabItem value="putty">

* Enter the TRASA hostname or IP address under Session.

:::note 
If you have dowloaded the TRASA user key
* Navigate to Connection > SSH > Auth.
* Use PuTTYgen to convert the downloaded key `id_rsa` into `id_rsa.ppk`.
* Click Browse... under Authentication parameters / Private key file for authentication.
* Locate the id_rsa.ppk private key and click Open.
* Click Open again to log into the remote server with key pair authentication.
::: 

* Enter TRASA email and password.
* Enter IP address of service you want to connect to.   
* Enter TOTP code or leave it blank for U2F.
* Enter Service password.

</TabItem>


  <TabItem value="bitvise">
  todo
  </TabItem>


</Tabs>



--->



#### Using private key instead of password

##### Save private keys in vault (Recommended)
Ask your administrator to save the private key in the vault.

##### Using agent forwarding
>SSH Agent forwarding is not recommended since it allows users with root privilege in the server to use your SSH keys.

* Add the private key to ssh agent `ssh-add <private_key_path>`
* Use -A flag `ssh -A -i <private_key_path> root@trasa.hostname -p 8022`


You can use TRASA private key and service private key at the same time.   
`ssh -A -i <trasa_private_key_path> -i <service_private_key_path> root@trasa.hostname -p 8022`




