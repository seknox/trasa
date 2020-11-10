---
id: concepts
title: Concepts
sidebar_label: Concepts
description: Basic concepts for using TRASA
---





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










