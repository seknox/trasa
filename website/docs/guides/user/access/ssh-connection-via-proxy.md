---
id: ssh-connection-via-proxy
title: Access SSH Service
sidebar_label: SSH
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Access methods
You can access SSH service either via Browser or SSH client.

### Using Browser

* Login into your TRASA account
* Search for the service you want to connect
* Click connect and choose service username
* Enter password and TOTP 
> You may be asked to save the new host key





### Using SSH clients

#### Download TRASA User Key (Optional)
* Go to dashboard and go to "Account" tab
* Click the menu to get dropdown menu items
* Click the "get ssh private key" button" to download key

>This key is used to authenticate to TRASA server NOT upstream SSH server.
>So, you might still be asked for upstream password.

:::tip 
If you're using Putty, use PuTTYgen to convert the downloaded key `id_rsa` into `id_rsa.ppk`.
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
* Enter the IP address of service you want to connect to   
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


### Using private key instead of password

#### Save private keys in vault (Recommended)
Ask your administrator to save the private key in vault.

#### Using agent forwarding
>SSH Agent forwarding is not recommended since it allows user with root priviglege in server to use your SSH keys.
* Add the private key to ssh agent `ssh-add <private_key_path>`
* Use -A flag `ssh -A -i <private_key_path> root@trasa.hostname -p 8022`



:::tip
You can use trasa private key and service private key at the same time.   
`ssh -A -i <trasa_private_key_path> -i <service_private_key_path> root@trasa.hostname -p 8022`
:::




