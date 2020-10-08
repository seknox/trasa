---
id: proxy-vs-direct
title: Access Proxy vs Direct Access
sidebar_label: Access Proxy vs Direct Access
draft: true
keywords: 
    - access proxy
---


:::note
You can use both the access proxy and the native agent at the same time.
You can access the native agent installed server through TRASA access proxy.
:::


### Access Proxy
<img alt="generate-ca-btn" src={('/img/docs/access-proxy/trasa-access-proxy.png')} />  

TRASA has a built-in proxy for HTTP, SSH, RDP, and SQL protocol.
You don't need to enable it, but you may want to configure a firewall policy to enforce access from TRASA proxy only.
If users directly access the service, it won't be protected with TRASA.

Pros of using access proxy:
* Easier to implement. No need to install and configure agents in each server.
* Session can be recorded (text or video log).
* Can autofill the saved passwords and keys from the vault.


### Direct Access with 2FA agent
<img alt="generate-ca-btn" src={('/img/docs/access-proxy/trasa-direct-access.png')} />
  
Native agents are installed on each upstream servers.  

The agent will make a request to TRASA server after password authentication to verify policy and second factor.

Pros of using native agents:
* Can implement 2FA in local login (not just remote login).
* Last-mile security protection. 






