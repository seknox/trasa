---
id: create-user
title: Create User
sidebar_label: Create User
---



 **Users** in TRASA have dual identities., 
 
 1) **TRASA user** - user identity used to uniquely identify a user in TRASA platform.  
 2) **privileges** - user identity mapped to authorized application or service that is local to a particular service. 


## 1. TRASA users
TRASA users are uniquely identified by an enrolled email address or assigned username. They are unique throughout specific organizations and are created by administrators.
TRASA users have roles of either `orgAdmin` or `selfUser`. 


## 2. Privilege
Privilege are TRASA users mapped to authorized application or service that is local to particular service.

For example, TRASA user `james@nepsec.org` can be assigned to `winserver001` service authorized to log in as `Administrator`.
Privilege are often mapped to services with [TRASA Policy](https://seknox.com/trasa/docs/concepts/permissions-policies)
