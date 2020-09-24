---
id: rdp-connection-via-proxy
title: Windows logon or RDP Service
---

> When your administrator configures access proxy, you can access RDP service only via the TRASA dashboard.

> If access proxy is not configured and only the TRASAWIN 2fa agent is configured, you can access remote computers with any regular RDP client you have been using.

### Windows logon TFA Prompt
When the windows server or workstation is protected with TRASA 2FA agent, just after the username(privilege) and password is validated in the windows logon screen, a prompt will appear on your screen. You will have to enter your trasaID(email or username) and select the TFA method in order to perform second step verification.


### RDP via TRASA access proxy
* Login into your TRASA account
* Search for the service you want to connect
* Click connect and choose service username
* Enter password and TOTP 

