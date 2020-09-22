---
id: concepts
title: Concepts
sidebar_label: Concepts
description: Basic concepts for using TRASA
---

TRASA is unique in the sense that it bundles lots of access control features into a single platform and, as such, allows you to achieve zero trust access control strategies in your infrastructure. 

## Baseline security concept

TRASA is build upon the baseline security concept, which adheres to zero trust access control paradigm.

### 1. Whitelisted entry point for remote access

You may have heard in the wild that zero-trust is all about allowing access from anywhere, anytime. While not delusive, many of us have interpreted it as allowing direct access from users to remote services. It is dangerous, and access via jump/Bastian servers is still the safest way to enable remote access.

Zero trust does not promote to remove your Bastian server but rather upgrade your Bastian server to allow remote access based on risks. It means that users can access remote service from anywhere as long as their connection is not deemed risky. 


### 2. Security hygiene of user device is crucial

Although it is users who access your remote service, their devices (mobile device, workstations) handles every aspect of remote access. Once connected, all sensitive data are also processed and stored in users' devices. If the user's device is already compromised, the firewall, Intrusion detections, SIEMs, malware detection on the server-side are of little help, if any. Users may be choosing strong passwords, handling API keys safely. Still, a compromised device will probably have snooped all those secrets already.

So granting access to servers and services based on user devices' security hygiene is of utmost importance and a fundamental step to achieve zero trust.


### 3. Monitoring trusted access

Almost every security compromises involve the misuse of trusted credentials, trusted networks, and trusted devices. There must be complete visibility to an active authorized session so that any malicious intent hidden in trusted access can be audited in realtime or in the future.


### 4. Realtime view of remote access

+ Administrators must have a realtime view of all authorized users and authorized devices for remote access. 
+ Administrators must have a realtime view of all remote entry points to your infrastructure.
+ Administrators must have a realtime view of all services that have remote access enabled.


## TRASA implementation

TRASA Access Proxy is a drop in upgrade for your homegrown bastian/jump server.
Whether you are using a linux server (configured as Jump server) or Microsoft Remote Desktop Gateway, TRASA offers all those features along with best practices enabled and configurable by default.

### TRASA server as an access point for your internal infrastructure
Access Proxy is just a reverse proxy that happens to know RDP, SSH, HTTP, and Database protocols and makes a forwarding decision based on the administrator's policy. For TRASA access proxy to work, the network must be configured in a way such that every remote access to your server and services are only allowed from the server IP address in which TRASA is installed.

### Configuring your network
+ Make changes in your network firewall such that ingress traffic to your server and services listening for SSH, RDP, HTTPs, and DB traffic are only allowed from TRASA server IP address.
 
