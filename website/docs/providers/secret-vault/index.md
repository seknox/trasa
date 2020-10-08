---
id: index
title: Secret Vault
sidebar_label: Introduction
---

import useBaseUrl from '@docusaurus/useBaseUrl';

TRASA stores all keys and secrets in secure vault known as TsxVault.

  _Passwords_, _Secret keys_, _API tokens_ etc. are needed by TRASA to integrate with 3rd party services. For example, FCM tokens, Email config settings, IDP integration keys.

### Vault States
* Uninitialised
* Initialised
    - Encrypted
    - Decrypted

Initially, after installation, the vault is in the "Uninitialised" state. You need to initialize the vault.

Now vault decryption keys are generated, and the vault will be in the "Decrypted" state. 
The decryption keys are stored in memory. So if the TRASA service restarts, the vault will be in the "Decrypted" state.
You need to decrypt the vault using the decryption keys to start using it again.



### Initialize Vault (one time only)
* Open Menu Drawer and click on Providers
<img alt="providers-menu" src={('/img/docs/providers/providers-menu.png')} />
* Go to "Secret Storage" tab
* Click the Initialise button 
<img alt="initialise" src={('/img/docs/providers/secret-vault/initialise.png')} />
* Copy the keys and keep them safely
<img alt="keys" src={('/img/docs/providers/secret-vault/keys.png')} />


### Decrypt the Vault
If TRASA service restarts, you need to decrypt the vault to start using it again.

To do that,
* Go to the Providers page.
* Click the "Secret Storage" tab.
* Click the "Enter Decryption Key" dropdown.
* Enter a decryption key and click submit.
* Submit two more decryption keys.

<img alt="keys" src={('/img/docs/providers/secret-vault/decrypt-vault.png')} />


### Storing Service Credentials 

If the Vault is in a decrypted state, you can use it to store service credentials like passwords and keys.

* Go to services and click on the service you want to configure credentials
* Go to “Manage Credentials” tab
* Fill in username and password/key
* Click on + sign to save
<img alt="manage-creds-tab" src={('/img/docs/providers/secret-vault/manage-creds.png')} />

>From now on, users won’t be asked for a password while logging into this service with this privilege
