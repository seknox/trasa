---
id: proxy-vs-direct
title: Access Proxy vs Direct Access
sidebar_label: Access Proxy vs Direct Access
draft: true
keywords: 
    - access proxy
---


### Access Proxy
<img alt="generate-ca-btn" src={('/img/docs/access-proxy/trasa-access-proxy.png')} />  

TRASA has a built-in proxy for HTTP, SSH, RDP, and SQL protocol.
You don't need to enable it, but you may want to configure a firewall policy to enforce access from TRASA proxy only.
If users directly access the service, it won't be protected with TRASA.


* Easier to implement. No need to install and configure agents in each server.
* Session can be recorded (text or video log).


### Direct Access with 2FA agent
<img alt="generate-ca-btn" src={('/img/docs/access-proxy/trasa-direct-access.png')} />  

TRASA agent will make a request to TRASA server after password authentication.
Then TRASA will take care of second factor authentication.
  
* Can implement 2FA in local login.







