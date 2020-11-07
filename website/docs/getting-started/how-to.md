---
id: how-to
title: How does it work ?
sidebar_label: How does it work ?
---


TRASA is a relatively young project with rapidly evolving API and workflow. But we thrive to keep the best and up to date documentation to help you succeed with TRASA. 

If you are having trouble operating a feature, ping us in our community discussion forum.

## Working concepts

TRASA controls `users` (DevOps, Software, Marketing teams, 3rd party vendors, etc.) `access to services` (SSH, RDP, Web, Database). The decision to control access depends on `access policy` and `access privilege` defined by administrators. 

Once TRASA is installed and configured, [users](./glossary.md#user) will  access server and services with through TRASA access proxy.

### How will user access the services?
 Currently supported methods of connection:
- **RDP :** only possible from trasa RDP client (web application, access from the browser)
- **SSH :** trasa ssh client (web application, access from browser ), *nix terminal ssh client,  putty, bitvise, moba xterm. Basically any ssh client available.
- **HTTP :** normal as you would browse web page (session is managed with trasa browser extension
- **Database(MySQL only) :** using mysql client you already use.


##  Basic Workflow
1. Create/import users profiles
2. Create/import service profiles
3. Create access policy
4. Assign user or user groups to the service or service group (Access Mapping)
5. Test Access.
6. Monitor Access.
