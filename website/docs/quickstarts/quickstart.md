---
id: quickstart
title: TRASA Quickstart guide
sidebar_label: Quickstart
---
import useBaseUrl from '@docusaurus/useBaseUrl';


This tutorial will show you the basic usage of TRASA.


### STEP 1 Install TRASA

Follow [installation guide](../install/installation.md) to install TRASA. 

### STEP 2 Login

* Download TRASA mobile app from [App Store](https://apps.apple.com/us/app/trasa/id1411267389) or [Play Store](https://play.google.com/store/apps/details?id=com.trasa)
* Open TRASA_URL in your browser.
* Use default credentials to login (root:changeme)

Since this is your first time logging into TRASA, you have not yet added your 2FA device yet.
QR code will appear on the screen.

* Open TRASA mobile app and press + button on the bottom right and then press the QR icon.

<img width="20%" alt="mobile-app-add-qr" src={useBaseUrl('img/docs/quickstart/mobile-app-add-qr.png')} />  

* Scan the QR code on the browser.
* If everything goes well, you will see the following icon on your app.

<img  width="20%"  alt="mobile-app-added-totp" src={useBaseUrl('img/docs/quickstart/mobile-app-added-totp.png')} />
  
* Press the icon to get TOTP codes.


Now 2FA device is added.

* Try logging in again.
* Now you need to choose TOTP and enter TOTP code from the mobile app.


### STEP 3 Create Policy
* Click the menu button on top left to open main menu drawer.
<img  alt="main-menu" src={useBaseUrl('img/docs/main-menu.png')} />  
* Click the Control menu to open policy page.
* Click the "Create new policy" button.
<img  alt="create-policy-btn" src={useBaseUrl('img/docs/quickstart/create-policy-btn.png')} />  
* Enter a policy name and click next.
* Click the "Mandatory 2FA" switch to enable second factor authentication
<img  alt="2fa-enable" src={useBaseUrl('img/docs/quickstart/2fa-enable.png')} />  
* Click the "Session Recording" menu and enable it.
<img  alt="session-recording-enable" src={useBaseUrl('img/docs/quickstart/session-recording-enable.png')} />  
* Click "Day and Time" to configure weekday and time specific policy.
* Select week days and time range to allow access.
* Click "add" button.
<img  alt="add-day-time-policy" src={useBaseUrl('img/docs/quickstart/add-day-time-policy.png')} />  
* Click next and review the policy to be created. 
* If everything looks good, click the "Submit" button.

:::tip
Go to [policies reference](../policies/basic-policy.md) to know more about static policies
:::

### STEP 4 Create Service
* Open main menu and click Services
* Click "Create new service" button.
<img  alt="create-service-btn" src={useBaseUrl('img/docs/quickstart/create-service-btn.png')} />  
* A drawer will slide from the left. Fill in the details of upstream server you want to connect through TRASA.
<img  alt="create-new-service" src={useBaseUrl('img/docs/quickstart/create-new-service.png')} />  
* Click submit. You will be redirected to the newly created service page.
* Click the "Access Map" tab.
<img  alt="access-map-tab" src={useBaseUrl('img/docs/quickstart/access-map-tab.png')} />  
* Click the "Assign user" button.
<img  alt="assign-user-btn" src={useBaseUrl('img/docs/quickstart/assign-user-btn.png')} />  
* Choose the user and policy you just created.
<img  alt="assign-user-dialog" src={useBaseUrl('img/docs/quickstart/assign-user-dialog.png')} />  

> Right now you are assigning yourself to the service. You can assign other users too when you [create](../users/crud.md) them.

* On privilege, enter the username of the upstream server.
* Click submit.


### STEP 5 Access

* Click the dropdown button on the top right and click on "My Account" menu.
 <img  alt="my-account-btn" src={useBaseUrl('img/docs/quickstart/my-account-btn.png')} />  
* You will see the newly assigned service.
 <img  alt="my-service" src={useBaseUrl('img/docs/quickstart/my-service.png')} />  
* Click connect and then the privilege (username) you configured earlier.
* Allow pop-up if your browser blocks it.
* Enter password.
* Choose TOTP and enter TOTP code from mobile app.
* You will be asked to save new SSH host key, press "y" to do that.

You will be logged in the upstream server through TRASA.

Now go through the docs to configure TRASA according to your needs. 





















