---
id: index
title: Secret Vault
sidebar_label: Introduction
---


TRASA uses Vault to  store secret credentials and keys. 


First you need to initialise vault (one time only)
* Go to Providers -> Secret Storage in dashboard
* Click the Initialise button 
* Copy the keys and keep them safely


>If the trasa-server restarts, you need to decrypt the vault using any 3 of these 5 keys.



Storing Service Credentials 

If the vault is in a decrypted state, you can use it to store service credentials like password and keys.

* Go to services and click on the service you want to configure credentials
* Go to “Manage Credentials” tab
* Fill in username and password/key
* Click on + sign to save
>From now on users won’t be asked for password while logging into this service
