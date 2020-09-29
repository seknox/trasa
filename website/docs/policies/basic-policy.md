---
id: basic-policy
title: Basic Policy
sidebar_label: Basic Policy
---
"Basic Policy" supports time and location-based policy along with control over two-factor authentication, option for session recording, and permission for file transfers. 


1. Second Factor Auth 
   + Enforce mandatory 2FA
  
2. Session Recording:
   + Enable or disable session recording. Only supported in SSH and RDP.
  
3. File Transfers:
   + Allow or disallow file transfers in SSH or RDP session.
  
4. IP Source:
   + Whitelist IP source. The default is `0.0.0.0/0`, which allows from all sources.
   + You can whitelist multiple sources as comma-separated values. `192.168.0.1/24,192.168.0.10`

5. Day and Time:
   + Allows restricting access based on day and time.
   + E.g. Sunday 11 AM-4 PM, Monday 1 AM-9 PM. Lets you create multiple day and time-based policies.
  
6. Expiry
   + Set time when this policy expires. Important when you have to allow access to 3rd party support, for example, one day. The policy automatically revokes access after expiry.
