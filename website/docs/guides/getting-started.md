---
id: getting-started
title: TRASA User Guides
sidebar_label: Getting Started
---
import useBaseUrl from '@docusaurus/useBaseUrl';

> If you are looking to install and configure TRASA, look at our [Administrative  Docs](https://www.trasa.io/docs/ "Administrative Docs")

## Hi !

If your security team has deployed TRASA in your infrastructure and all of a sudden, you are required to use TRASA for remote access, this guide is for you to get started.



* You probably got a link from your administrator. It will take you to initial password setup page.
<img  alt="password-setup" src={useBaseUrl('img/docs/user-guides/account/password-setup.png')} />


* Then you will be redirected to TRASA login page. Before you login, install TRASA mobile app from [Play Store](https://play.google.com/store/apps/details?id=com.trasa&hl=en) or [App Store](https://apps.apple.com/np/app/trasa/id1411267389).
* Now, enter your email and password you just set to login.

Since this is your first time logging into TRASA, you have  not yet added your 2FA device yet.
QR code will appear on screen.
<img  alt="qr-code" src={useBaseUrl('img/docs/user-guides/device/qr-code.png')} />


* Open TRASA mobile app and press + button on bottom right and then press QR icon

<img width="20%" alt="mobile-app-add-qr" src={useBaseUrl('img/docs/quickstart/mobile-app-add-qr.png')} />  

* Scan the QR code on the browser
* If everything goes well, you will see the following icon on your app

<img width="20%" alt="mobile-app-added-totp" src={useBaseUrl('img/docs/quickstart/mobile-app-added-totp.png')} />  

* Press the icon to get TOTP codes


Now 2FA device is added.

* Try logging in again
* Now you need to choose TOTP and enter TOTP code from the mobile app


When you log into the TRASA dashboard, you will be redirected to your account page. 

### Services 
It displays all the services you are assigned to. You can access them from here.

### Account
It displays your account details. You can change your password and download the SSH key from here.

### Device
It Displays your current devices and lets you delete them or enroll new devices.


### Access Stats
It shows your authentication logs.

