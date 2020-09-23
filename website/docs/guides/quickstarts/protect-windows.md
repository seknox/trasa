---
id: protect-windows
title: Adding 2FA to Windows Logon
sidebar_label: Protect Windows Logon
---



In this quick start, we start from account signup to protecting your first windows machine with TRASA.  
  

Precisely,

1. We demonstrate the signup process, 
2. Login to TRASA dashboard, 
3. Create and setup the TRASA authentication agent, assign a user to this app.
4. Register Trasa Authenticator (mobile app)
5. Install and configure TrasaWin (windows authentication package by TRASA)
6. Feel secure forever after...  
### 1. Let's first signup.

Head to [signup page](https://seknox.com/trasa/signup) and provide your details.

<!-- ![](https://storage.googleapis.com/trasa-website-static/signup.png)  -->

Once you submit the form, you should receive a welcome invitation by the TRASA team. 
Head to your inbox and check the email.

The email you receive should look like the following.

<!-- ![](https://storage.googleapis.com/trasa-website-static/signup-email.png) -->

Click the "verify my account" button, and you will be redirected to the TRASA dashboard to setup your credentials. 

### 2. Login to trasa dashboard
Head over to [login page](https://app.trasa.io/login) and enter your credentials which you receive earlier in email.
<!-- ![](https://storage.cloud.google.com/trasa-website-static/login.png) -->



### 3. Create TRASA service  
Once you log in, proceed to creating [auth agents](https://app.trasa.io/services) . Auth agents are lightwieght authenticators which protects your system by adding two factor authentication. Auth agents are also responsible to collect event metrics which are populated in analytical dashboard.
<!-- ![](https://storage.googleapis.com/trasa-website-static/create-auth-app.png) -->
1) Click on auth Agents, which can be found in the side menu of the dashboard
2) Click the create Auth Agent button
3) Enter the name (or any identifier specific to your organization) to this new agent.

Once we are done creating the service, we should now assign a user to this app. when we install and configure TRASA authentication agents in the system we want to protect. Authentication agents only allow users to login who are assigned to that particular application. Even if the user has a local account but is not assigned within service,
TRASA will block the user from logging in. This adds flexibility to user administration, and access can be granted or revoked with the simple process of user permission assignment and revocation.

 We assign user to app by clicking on assign user button. Assinging user to service is a three step process.
 1) Select user which you want to assign to this app.
 2) Assign permission. Permission is based on time, day and expiry date.
 3) Finally, review user and permission which you are about to assign. Hit submit to submit the request. If in case you need to change any permissions, you can always go back and edit it again.


### 4. Register Trasa Authenticator (mobile app)
Download trasa authenticator app from [Android play store](https://play.google.com/store/apps/details?id=com.trasa), [IOS app store](https://itunes.apple.com/us/app/trasa/id1411267389?mt=8)  

##### | For TRASA SaaS Clients



Log into trasa authenticator. you will need to enter OTP code in the process. check your inbox once you enter email and password in login page.

##### | For TRASA Onpremise Clients



Visit [Enrol Device Page](https://app.trasa.io/woa/enrol/device). Once you enter your email and password, you will be presented with QR code. Scan the QR code with *Login with QR code option* and your device will be automatically synced with your account.




### 5. Install and configure TrasaWin (windows authentication package by TRASA)

- Dependency required: 
[visual c++ redestributable](https://aka.ms/vs/15/release/vc_redist.x64.exe)

 *WARNING:    Do not restart your computer untill you configure trasaWIN. Broken configuration may forever lock out your access to operating system. you should check the configurations using trasa config application which is avialable once you install trasaWIN(refer to image below)*

Download [TrasaWIN](https://storage.googleapis.com/trasa-public-download-assets/trasa-installers/TrasaWINv2.5.msi). Install it into windows workstation/server you want to protect with Trasa.
Once you finish the installation, Open the trasaWIN application as administrator.
Head back to your dashboard and note the serviceID and the serviceSecret of the authentication app from the service page (which we created earlier in step 3).
You can see serviceID ,  serviceSecret, and API hostname entered as per service created earlier in the following image. Note API hostname "app.trasa.io" always remains the same for TRASA SaaS users and can be a custom URL for self-hosted (On-Premise) TRASA users.

<!-- https://storage.googleapis.com/trasa-website-static/quickstart/configure-trasawin.png -->

Once you fill in required value, click check button. This verifies our values, appID, appSecret and api hostname to TRASA server to check if everything is correct. CLick save if reponse from server does not alert any error. If you validation fails, re-check the credentials you enter are correct.

Restart your operatin system after you save the configurations. 

### 6. Authorize Login with Trasa Authenticator
Now as soon as your operating system boots, in this case windows boots up, you will need to enter your username and password as you do normally. The only difference, you should get notification in your mobile app to authorize login.   To complete you login process, press authorize and to deny it, press cancle. we will press authorize to login to windows.

<!-- https://storage.googleapis.com/trasa-website-static/quickstart/mobileapp.png -->

Congratulations, you just added second-factor authentication to your normal login process. Every time you try to login to your account, you will need to authorize the login attempt with your mobile app. 

Refer to other docs to better understand the TRASA Platform. 
