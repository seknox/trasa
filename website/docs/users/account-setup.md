---
id: account-setup
title: Account Setup
sidebar_label: Account Setup
---



:::note
### To start the setup process
#### For root account: 
Open [TRASA_HOST](/docs/getting-started/glossary#TRASA_HOST) in your browser and authenticate with default credential (`root:changeme`). It is either the IP or domain address you used during installation. 

#### For all other accounts:  
Open verification link (either received via email or directly from your administrator) in your browser and setup new password.
:::

:::note
### Account setup flow:
**1. Setup password (Required).**

**2. Enrol mobile device (Required):** Once you authenticate with changed password, you will be redirected to enrol mobile 2FA device.

**3. Enrol workstation device (Optional):** Once you enrol mobile device, you can now access dashboard pages.
:::

<br />



<br />


### Accessing TRASA dashboard


- Open [TRASA_HOST](/docs/getting-started/glossary#TRASA_HOST) in your browser.

Following image shows openening [TRASA_HOST](/docs/getting-started/glossary#TRASA_HOST) in browser which is an IP address.
<img alt="dashboard login" src={('/img/docs/users/setup-account/trasa_host-ip.png')} />

If you are root user, use default credentials `(root:changeme)` to authenticate. Otherwise use verification link (which will be in url format: `https://<trasa-server-host>/woa/verify#token=tokenval` )

 <br />

### Set your password

If the previous step is successfull, you will be redirected and forced to set new password.

<img alt="setup password" src={('/img/docs/users/setup-account/setup-password.png')} />



After you set your password, you will be redirected to login page again.

<br /><br />




## Enrol Mobile Device


### Prepare mobile device
Download TRASA mobile app from [App Store](https://apps.apple.com/us/app/trasa/id1411267389) or [Play Store](https://play.google.com/store/apps/details?id=com.trasa).


:::important
TRASA requires two factor authentication by default and TRASA mobile app is default supported authenticator. Since this is your first login, you need to enrol device first:

:::



### Enrol Steps:

Authentcate with your username (email address or `root` username) and freshly set password.
You will be redirected to a page to enrol your mobile device.

<img alt="qr-code" src={('/img/docs/user-guides/device/qr-code.png')} />

<img alt="enrol device" src={('/img/docs/tutorial/enrol-mobile-device.svg')} />

1. Press the `+` button (buttom right).
2. Press QR image icon button. This will open in-app camera.
3. Scan the QR image from TRASA dashbaord
4. If everything goes well, you will see the following icon on your app

<img alt="enrol device" src={('/img/docs/tutorial/device-enroled.svg')} />

You can click the icon to view your totp codes.

<br />

---

<br />

## Test dashboard access

Press `login` button in dashboard page (where QR code is shown), you will be redirected to Login page again.

1. Authenticate with your credentialss
2. Once the credentials are validated, you will see **second-step verification page**
   <img alt="enrol device" src={('/img/docs/tutorial/dashboard-totp.svg')} />
3. From your TRASA mobile app, note totp code and enter in dashboard to proceed login.
   <img alt="enrol device" src={('/img/docs/tutorial/enter-totp.svg')} />

4. Server will validate your totp code and will redirecto to dashboard overview page.
   <img alt="dashbaord overview" src={('/img/docs/tutorial/first-dashboard.png')} />


<br />

---

<br />



