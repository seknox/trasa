---
id: index
title: Secret Vault
sidebar_label: Introduction
---
import useBaseUrl from '@docusaurus/useBaseUrl';

TRASA uses Vault to store secret credentials and keys. 


### Initialize Vault (one time only)
* Open Menu Drawer and click on Providers
<img alt="providers-menu" src={useBaseUrl('img/docs/providers/providers-menu.png')} />
* Go to "Secret Storage" tab
<img alt="secret-storage-tab" src={useBaseUrl('img/docs/providers/secret-vault/secret-storage-tab.png')} />
* Click the Initialise button 
<img alt="initialise" src={useBaseUrl('img/docs/providers/secret-vault/initialise.png')} />
* Copy the keys and keep them safely
<img alt="keys" src={useBaseUrl('img/docs/providers/secret-vault/keys.png')} />


:::tip
If the trasa-server restarts, you need to decrypt the Vault using any three of these five keys.
:::


### Storing Service Credentials 

If the Vault is in a decrypted state, you can use it to store service credentials like passwords and keys.

* Go to services and click on the service you want to configure credentials
* Go to “Manage Credentials” tab
<img alt="manage-creds-tab" src={useBaseUrl('img/docs/providers/secret-vault/manage-creds-tab.png')} />
* Fill in username and password/key
* Click on + sign to save
<img alt="save-creds" src={useBaseUrl('img/docs/providers/secret-vault/save-creds.png')} />  

>From now on, users won’t be asked for a password while logging into this service with this privilege
