---
id: index
title: Secret Vault
sidebar_label: Introduction
---


TRASA uses Vault to store secret credentials and keys. 


First, you need to initialize Vault (one time only)
* Go to Providers -> Secret Storage in dashboard
* Click the Initialise button 
* Copy the keys and keep them safely


>If the trasa-server restarts, you need to decrypt the Vault using any three of these five keys.



Storing Service Credentials 

If the Vault is in a decrypted state, you can use it to store service credentials like passwords and keys.

* Go to services and click on the service you want to configure credentials
* Go to “Manage Credentials” tab
* Fill in username and password/key
* Click on + sign to save
>From now on, users won’t be asked for a password while logging into this service
