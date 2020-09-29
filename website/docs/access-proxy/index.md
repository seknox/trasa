---
id: introduction
title: Access Proxies
sidebar_label: Introduction
---
Access proxies are core server components in TRASA. Users typically access protected services via access proxies. 
Access proxies give us full visibility at the application layer (Layer 7 OSI) and enforce more granular policies. 

Currently, TRASA supports four types of protocol.
1. HTTPs
2. SSH
3. RDP
4. Database (beta, Mysql only)

You need to configure firewall rules to make remote access to servers only from TRASA.
We have guides to configure that in [AWS](../cloud/aws.md), [GCP](../cloud/gcp.md) and [Digital Ocean](../cloud/digital-ocean.md).
