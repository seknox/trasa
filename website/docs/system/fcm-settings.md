---
id: fcm-setup
title: FCM Setup
sidebar_label: FCM Setup
---

import useBaseUrl from '@docusaurus/useBaseUrl';


FCM is a mobile notification service offered by Google. In order to perform TRASA U2F based 2FA. 
By default, TRASA proxies U2F request via Seknox's Cloud server. It means that by default when users choose to perform U2F via TRASA mobile app, the U2F request is forwarded to Seknox's Cloud server and proxied back when the user confirms response (authorize or deny).

If you do not want to route U2F requests via Seknox server, you can configure your own FCM service.
You need to [build a mobile app from the source](https://github.com/seknox/trasa/blob/master/app/Readme.md). 
 


To configure Seknox's Cloud server as FCM proxy,
* Initialise and decrypt the [vault](../providers/secret-vault/index.md) if you haven't done it already.
* Open the main menu and click in "System".
* Click the "settings" tab.
* Click the "FCM Setting" section.
<img alt="fcm-settings" src={('/img/docs/system/fcm-settings.png')} />  

* On cloud proxy address, enter `https://sg.cpxy.trasa.io`. 
* Enter your email and click "Request Access".
* Open your email inbox and click the verification link.
* Click the "Obtain Key" button in TRASA dashboard.

