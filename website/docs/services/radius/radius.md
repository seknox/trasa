---
id: radius-server
title: Radius Server
sidebar_label: Radius Server
---

TRASA has a built-in radius server.

So if you configure a device to use TRASA as a RADIUS server, TRASA will handle authentications and policies. 

TRASA listens on port 1812 for RADIUS requests. 

To configure a server/device to use TRASA as RADIUS, 
* Create a service of type RADIUS and correct hostname/IP
* Enable RADIUS authentication on the device/server with
```
RADIUS server address: hostname/IP of TRASA
Port: 1812
Shared Secret: Service secret
```

:::tip
You can get the service secret in the service settings.
:::

