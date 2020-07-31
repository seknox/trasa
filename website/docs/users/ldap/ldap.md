---
id: ldap
title: Connecting TRASA with LDAP server
sidebar_label: LDAP
---

Follow this document as a reference for all LDAP based Identity Provider including
+ Active Directory
+ FreeIPA
+ LDAP server

> Active Directory and FreeIPA have prebuilt configuration in TRASA and are available in Identity Provider menu. 

To test for LDAP binding outside of TRASA, you can use [ldapsearch](https://linux.die.net/man/1/ldapsearch 'ldapsearch') in Linux systems or `ldap.exe` tool found in Windows server (under ldap server configuration panel).

[Digital ocean](https://www.digitalocean.com/community/tutorials/understanding-the-ldap-protocol-data-hierarchy-and-entry-components) and [ldap.com](https://ldap.com/basic-ldap-concepts/) has excellent article to understand basic of LDAP. If you are new to LDAP, make sure to read them first.


### Prerequisite

There are 2 things needed when integrating TRASA with LDAP server. 
1. LDAP service account for account Binding (authentication)
2. User group in LDAP server. TRASA will import users from this group.

**Note that TsxVault must be already initialized and must be in unsealed state.**




## 1. Create New Identity Provider
![create idp](./create-idp.png 'create idp')


## 2. Configure Identity Provider
![Configuring LDAP IDP](./configuring-ldap.png 'Configuring LDAP IDP')



+ **Server Domain** - IP or domain name where LDAP server is hosted. We have used IP `34.87.105.20` here
+ **LDAP DN** - usually LDAP base where users can be queried. We have used `CN=Users,DC=trasatest,DC=internal` as base user DN here.
+ **Service Account Name** - Service account name used to authenticate (aka bind) LDAP. Any user account with access rights to query ldap server can be used here. But as a better security option, always create and use separate service account for similar use case. We have used `serviceaccount` as ldap user here
+ **Service Account Password** - Password for above service account.


## 3. Import Users from LDAP server
![Importing Users from LDAP](./import-ldap-users.png 'Importing Users from LDAP')

Provide full path to LDAP user group. As an example, here we have used user group name `ldapgroup`


## Finishing
If all wen well, users from LDAP group will be imported in TRASA. Users can use same LDAP credentials to authenticate themself in TRASA