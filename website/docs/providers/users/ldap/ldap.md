---
id: ldap
title: Connecting TRASA with LDAP server
sidebar_label: LDAP
---

Follow this document as a reference for all LDAP based Identity Provider, including
+ Active Directory
+ FreeIPA
+ LDAP server

> Active Directory and FreeIPA have a prebuilt configuration in TRASA and are available in the Identity Provider menu. 

To test for LDAP binding outside of TRASA, you can use [ldapsearch](https://linux.die.net/man/1/ldapsearch 'ldapsearch') in Linux systems or `ldap.exe` tool found in Windows server (under LDAP server configuration panel).

[Digital Ocean](https://www.digitalocean.com/community/tutorials/understanding-the-ldap-protocol-data-hierarchy-and-entry-components) and [ldap.com](https://ldap.com/basic-ldap-concepts/) has an excellent article to understand the basics of LDAP. If you are new to LDAP, make sure to read them first.


### Prerequisite

There are two things needed when integrating TRASA with the LDAP server. 
1. LDAP service account for account Binding (authentication)
2. User group in the LDAP server. TRASA will import users from this group.

**Note that TsxVault must be already initialized and must be in an unsealed state.**




## 1. Create New Identity Provider
![create idp](/img/docs/users/ldap/create-idp.png 'create idp')


## 2. Configure Identity Provider
![Configuring LDAP IDP](/img/docs/users/ldap/configuring-ldap.png 'Configuring LDAP IDP')



+ **Server Domain** - IP or the domain name where the LDAP server is hosted. We have used IP `34.87.105.20` here
+ **LDAP DN** - usually a LDAP base where users can be queried. We have used `CN=Users,DC=trasatest,DC=internal` as base user DN here.
+ **Service Account Name** - Service account name used to authenticate (aka bind) LDAP. Any user account with access rights to query the LDAP server can be used here. But as a better security option, always create and use a separate service account for similar use cases. We have used `serviceaccount` as a LDAP user here
+ **Service Account Password** - Password for the above service account.


## 3. Import Users from LDAP server
![Importing Users from LDAP](/img/docs/users/ldap/import-ldap-users.png 'Importing Users from LDAP')

Provide the full path to the LDAP user group. As an example, here we have used user group name `ldapgroup`


## Finishing
If all went well, users from the LDAP group would be imported in TRASA. Users can use the same LDAP credentials to authenticate themself in TRASA
