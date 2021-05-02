---
id: freeIPA
sidebar_label: FreeIPA Users
title: FreeIPA Users
---

<br /><br /><br />

# **FreeIPA**

FreeIPA is an integrated Identity and Authentication solution for Linux/UNIX networked environments. A FreeIPA server provides centralized authentication, authorization and account information by storing data about user, groups, hosts and other objects necessary to manage the security aspects of a network of computers. 
 
TRASA has pre-built integration with freeIPA using Light Weight Directory Access protocol or LDAP.


## Create FreeIPA IDP

![create freeIPA idp](./create-freeipa.png 'create freeIPA idp')

## Configure FreeIPA integration
As freeIPA integration is done via standard LDAP integration, it is required to have general understanding of LDAP.  [Digital ocean](https://www.digitalocean.com/community/tutorials/understanding-the-ldap-protocol-data-hierarchy-and-entry-components) and [ldap.com](https://ldap.com/basic-ldap-concepts/) has excellent artical to understand basic LDAP. Make sure to read them first.

As minimum requirements to configure integration between TRASA and FreeIPA,
you will need server address (domain name) of freeIPA server, root DN for LDAP, search base for LDAP, service account name and password for LDAP binding. Details in image below shows samples value that should be similar to your requirements. It is easy if you test with ldapsearch first and configure trasa with similer parameters.
![configure-freeipa-idp](./configure-freeipa.png 'configure-freeipa-idp')

The update action will first try to validate your values and only update Idp data.
If all values are correct, LDAP details will be saved and will be used to perform user imports and user authentication.



<br /><br />  

## Importing users from freeIPA
Once freeIPA configuration is updated(see above), importing users is easy.

* ### Importing multiple users
As a reference from configuration image above, to import all image, we would enter search query as <br />
`cn=users,cn=accounts,dc=freeipa,dc=trasa,dc=io`

* ### Importing single user
As a reference from configuration image above, to import all image, we would enter search query as <br /> `cn=testuser, cn=users,cn=accounts,dc=freeipa,dc=trasa,dc=io`


<br /><br />  

## Transfering users to or from trasa and freeIPA.
Chances are that you would want to authenticate users which are already created by trasa by freeIPA. To achieve this, you will need to transfer users from trasaIDP to freeIPA. 
Under **Convert trasa users to/from freeipa** menu of Idp configuration, you can click on **click to transfer user** button to open transfer list and perform user transfer. 
