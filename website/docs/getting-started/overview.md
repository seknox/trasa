---
id: overview
title: Welcome
sidebar_label: Overview
slug: /
---



Welcome to TRASA docs. Before anything, let's answer your first question.

## What is TRASA?

TRASA is an open-source zero trust service access platform built by [Seknox](https://www.seknox.com). It is unique in the sense that it bundles **lots of access control features** into a single platform and, as such, allows you to achieve zero trust access control strategies in your infrastructure. 


It essentially:
+ Is a Layer 7 protocol, user identity, and privilege aware access proxy.
+ Can enforce security policies (time, file transfers, location, context, 2FA) to SSH, RDP, web, database access. 
+ Can enforce access policy based on the security hygiene of user devices.
+ Add two-factor authentication agent (native integration) to protect console access to SSH, RDP, and hardware appliance. 

If you have used Bastian server to jump access or centralized access to internal infrastructure, you can also think of TRASA as a **Bastian server on steroids!**


### How is it different from what we already have implemented to control access?
To distinguish how TRASA and zero trust systems differ from legacy access control products, see how legacy vs. zero trust access control system decides to allow access in the below image:

<img alt="trasa vs legacy" src={('/img/docs/getting-started/zero-notzero.svg')} />

<br />

TRASA Access Proxy is a drop in upgrade for your homegrown bastian/jump server.
Whether you are using a linux server (configured as Jump server) or Microsoft Remote Desktop Gateway, TRASA offers all those features along with best practices enabled and configurable by default.





## Immediate use cases and benefits

### TRASA server as an access point (Bastian/Jump server) for your internal infrastructure

Access Proxy is a reverse proxy that also understands RDP, SSH, HTTP, and Database protocols and makes a forwarding decision based on the access policy. 


<!-- <img alt="enrol device" src={('/img/docs/tutorial/all-users.png')} /> -->

<img alt="trasa as bastian for internal services" src={('/img/docs/getting-started/zero-trust-service-access.svg')} />


### Follow best practice security for remote (SSH, RDP, WEB, Database).

- Centralized authorization for remote access.
- Enforce security policy for remote access.
- Add two factor authentication to remote access.
- Easily add and implement SSH Certificate authentication to every SSH access.

