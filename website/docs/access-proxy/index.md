---
id: introduction
title: Access Proxy
sidebar_label: Introduction
---
Access proxies are core server components in TRASA. Users typically access protected services via access proxies. 
Access proxies give us full visibility at the application layer (Layer 7 OSI) and enforce more granular policies. 

<img alt="generate-ca-btn" src={('/img/docs/access-proxy/trasa-access-proxy.png')} />  


Currently, TRASA supports four types of protocol.
1. HTTPs
2. SSH
3. RDP
4. Database (beta, Mysql only)

You need to [configure firewall rules](../../install/initial-setup.md#3-firewall-configuration-optional) to enforce remote access to servers only from TRASA.
