---
id: rdp-service
title: Microsoft Remote Desktop
sidebar_label: RDP
---
There are two ways to protect RDP servers.
1. Through TRASA access proxy
2. Native 2FA agents 

#### 1. Native 2FA agents
You need to install and configure 2fa agents in all RDP servers you want to protect.
[This guide](../../native-tfa/windows/windows-two-factor-authentication.md) will help you configure native agents.


#### 2. RDP Access Proxy
TRASA uses guacamole to connect to the RDP server.
To enable RDP, guacd (guacamole server daemon) must be running. By default TRASA will look for guacd on 127.0.0.1:4822, but you can change that in [config](../../system/config-reference.md#guacdaddr).

