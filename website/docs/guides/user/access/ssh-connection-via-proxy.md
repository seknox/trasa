---
id: ssh-connection-via-proxy
title: Access SSH Service
sidebar_label: SSH
---

You can access SSH service either via Browser or SSH client.

### Using SSH clients

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


##### Agent forwarding.

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




### Using Browser

* Login into your TRASA account
* Search for the service you want to connect
* Click connect and choose service username
* Enter password and TOTP 
> You may be asked to save new host key

