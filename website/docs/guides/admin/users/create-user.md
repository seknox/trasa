---
id: create-user
title: Create User
sidebar_label: Create User
---



 **Users** in TRASA have dual identities., 
 
 1) **TRASA user** - user identity used to uniquely identify a user in TRASA platform.  
 2) **privileges** - user identity mapped to authorized application or service that is local to particular authapp. 


## 1. TRASA users
TRASA users are uniquely identified by enrolled email address or assigned username. They are unique to throughout specefic organization and are created by administrators.
TRASA users have roles of ethier `orgAdmin` or `selfUser`. 


## 2. Privilege
Authapp users are TRASA users mapped to authorized application or service that is local to particular authapp.

For example, TRASA user `james@nepsec.org` can be assigned to `winserver001` authapp authorized to login as `Administrator`.
Authapp users are often mapped to authapps with [TRASA Policy](https://seknox.com/trasa/docs/concepts/permissions-policies)
