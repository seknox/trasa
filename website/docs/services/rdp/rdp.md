---
id: rdp-service
title: Microsoft Remote Desktop
sidebar_label: RDP
---

## Setup Methods

There are two ways to protect RDP servers:
1. Native 2FA agents 
2. Via TRASA access proxy

Either way, you need to [create a service](../index.md#creating-a-new-service) first.


### 1. Native 2FA agents
You need to install and configure 2fa agents in all RDP servers you want to protect.
[This guide](../../native-tfa/windows/windows-two-factor-authentication.md) will help you configure native agents in windows.

:::note Important
If you install a native 2fa agent, the 2FA prompt will be shown in local login as well as RDP login.
:::

### 2. RDP Access Proxy
To use TRASA as an RDP proxy, you need to [configure firewall rules](../../install/initial-setup.md#3-firewall-configuration-optional) to enforce RDP access from TRASA only.

We only support access through RDP proxy from a browser. So, users need to log into the TRASA dashboard(web app) to access RDP. 
Learn more about accessing RDP proxy [here](../../guides/user/access/rdp-connection-via-proxy.md).

TRASA uses guacamole to connect to the RDP server.
To enable RDP, guacd (guacamole server daemon) must be running. By default, TRASA will look for guacd on 127.0.0.1:4822, but you can change that in [config](../../system/config-reference.md#guacdaddr).

