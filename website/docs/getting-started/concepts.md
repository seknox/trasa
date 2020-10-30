---
id: concepts
title: Concepts
sidebar_label: Concepts
description: Basic concepts for using TRASA
---






TRASA Access Proxy is a drop in upgrade for your homegrown bastian/jump server.
Whether you are using a linux server (configured as Jump server) or Microsoft Remote Desktop Gateway, TRASA offers all those features along with best practices enabled and configurable by default.

## TRASA server as an access point (Bastian/Jump server) for your internal infrastructure

Access Proxy is just a reverse proxy that also understans RDP, SSH, HTTP, and Database protocols and makes a forwarding decision based on the access policy. 


<!-- <img alt="enrol device" src={('/img/docs/tutorial/all-users.png')} /> -->

<img alt="trasa as bastian for internal services" src={('/img/docs/getting-started/zero-trust-service-access.svg')} />


## Working concepts

TRASA controls `users` (DevOps, Software, Marketing teams, 3rd party vendors, etc.) `access to services` (SSH, RDP, Web, Database). The decision to control access depends on `access policy` and `access privilege` defined by administrators. 


###  Basic Workflow
1. Create/import users profiles
2. Create/import service profiles
3. Create access policy
4. Assign user or user groups to the service or service group (Access Mapping)
5. Test Access.
6. Monitor Access.




