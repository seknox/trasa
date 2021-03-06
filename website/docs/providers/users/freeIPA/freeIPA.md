---
id: free-ipa
title: FreeIPA
sidebar_label: FreeIPA
---



FreeIPA is an integrated Identity and Authentication solution for Linux/UNIX networked environments. A FreeIPA server provides centralized authentication, authorization, and account information by storing data about users, groups, hosts, and other objects necessary to manage the security aspects of a network of computers. 
 
TRASA has pre-built integration with freeIPA using Light Weight Directory Access Protocol or LDAP.


## Create FreeIPA IDP

![create freeIPA idp](./create-freeipa.png 'create freeIPA idp')

## Configure FreeIPA integration
As freeIPA integration is done via standard LDAP integration, it is required to have a general understanding of LDAP.  [Digital Ocean](https://www.digitalocean.com/community/tutorials/understanding-the-ldap-protocol-data-hierarchy-and-entry-components) and [ldap.com](https://ldap.com/basic-ldap-concepts/) has an excellent article to understand basic LDAP. Make sure to read them first.

As minimum requirements to configure the integration between TRASA and FreeIPA,
you will need the server address (domain name) of freeIPA server, root DN for LDAP, search base for LDAP, service account name and password for LDAP binding. Details in the image below shows samples value that should be similar to your requirements. It is easy if you test with ldapsearch first and configure TRASA with similar parameters.
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

## Transfering users to or from TRASA and freeIPA.
The chances are that you would want to authenticate users, which are already created by TRASA by freeIPA. To achieve this, you will need to transfer users from trasaIDP to freeIPA. 
Under **C
onvert trasa users to/from freeipa** menu of Idp configuration, you can click on **click to transfer user** button to open transfer list and perform user transfer. 
