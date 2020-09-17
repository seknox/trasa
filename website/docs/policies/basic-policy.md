---
id: basic-policy
title: Basic Policy
sidebar_label: Basic Policy
---

Basic Policy supports time and location-based policy along with control over two-factor authentication, option for session recording, and permission for file transfers. 

## Creating Policy

1. Second Factor Auth 
   + Enforce mandatory 2FA
  
2. Session Recording:
   + Enable or disable session recording. Only supported in SSH and RDP.
  
3. File Transfers:
   + Allow or disallow file transfers in SSH or RDP session.
  
4. IP Source:
   + Whitelist IP source. Default is `0.0.0.0/0` which allows from all sources.
   + You can whitelist multiple sources as comma seperated values. `192.168.0.1/24,192.168.0.10`

5. Day and Time:
   + Allows to resttrict access based on day and time.
   + E.g. Sunday 11AM-4PM, Monday 1AM-9PM. Lets you create multiple day and time based policies.
  
6. Expiry
   + Set time when this policy expires. Important when you have to allow access to 3rd party support for example 1 day. Policy automatically revokes access after expiry.