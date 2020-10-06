---
id: radius-server
title: Radius Server
sidebar_label: Radius Server
---

import useBaseUrl from '@docusaurus/useBaseUrl';


TRASA has a built-in radius server.

So if you configure a device to use TRASA as a RADIUS server, TRASA will handle authentications and policies. 

TRASA listens on port 1812 for RADIUS requests. 

To configure a server/device to use TRASA as RADIUS, 
* Create a service of type RADIUS with correct hostname/IP
* Enable RADIUS authentication on the device/server with following config
```
RADIUS server address: hostname/IP of TRASA
Port: 1812
Shared Secret: Service key
```

:::tip
You can get the service key in the service settings.
<img  alt="service-secret" src={('/img/docs/services/service-secret.png')} />

:::

