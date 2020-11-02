---
id: radius-server
title: Radius Server
sidebar_label: Radius Server
---


For software and hardware appliance to which you do not wish users to access through TRASA access proxy, or to protect protocols currently not supported in TRASA access proxy, you can integrate with TRASA radius server to enforce two-factor authentication and policies.
 

:::important
To add native two-factor authentication and enforce TRASA policy to software and hardware appliance that do not support agent integration, radius authentication is preferred to protect those systems. These systems may include database servers, network devices such as firewalls, router, and switch. 
:::


## How does it work?

TRASA has a built-in radius server. To add radius authentication, you will

1. Need to create a radius service profile in TRASA (radius server) and
2. Configure radius authentication in software or appliance (radius client).
3. 
This guide will show you how you can add a radius server profile in TRASA.

Radius client configuration depends specifically on software or appliance in context but usually requires creating a radius authentication group with a radius configuration.


## Creating radius service profile in TRASA
<iframe width="100%" height='600' src="https://www.youtube.com/embed/TUdPuBdiIAM" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Adding TRASA as radius server in radius client (system you want to protect)

On radius client, you will be required to supply:

1. **RADIUS server address:** This will be `Hostname/IP` of TRASA server
2. **Radius server port:** `1812`
3. **Shared Secret:** Service key of radius service. Copy this value from service profile (reference image below).

<img alt="copy radius service ServiceKey" src={('/img/docs/services/radius-service-integration.png')} />