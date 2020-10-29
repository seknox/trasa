---
id: rdp-connection-via-proxy
title: Windows/RDP
---

import useBaseUrl from '@docusaurus/useBaseUrl';


 When your administrator configures the access proxy, you can access RDP service only via the TRASA dashboard.

 If the access proxy is not configured and only the TRASAWIN 2fa agent is configured, you can access remote computers with any regular RDP client you have been using.

## Direct Access
When the windows server or workstation is protected with TRASA 2FA agent, just after the username(privilege) and password is validated in the windows logon screen, a prompt will appear on your screen. 

<img alt="trasa-credprov" src={('/img/docs/user-guides/access/trasa-credprov.png')} />  
You will have to enter your trasaID(email or username) and select the TFA method to perform second step verification.
<img alt="trasa-tfa-prompt" src={('/img/docs/user-guides/access/trasa-tfa-prompt.png')} />  



## Via TRASA access proxy
* Login in to your TRASA account.
* Search for the service you want to connect.
* Click the "Connect" button and choose service username.
* Enter password and TOTP. 


### RDP Console Menu
When accessing RDP using a browser, you can open the console menu by clicking the gear icon on top.
<img  alt="console-toggle" src={('/img/docs/user-guides/access/console-toggle.png')} />


You can see "Clipboard","File Transfer" and "Keyboard Events" section
<img  alt="console-menu" src={('/img/docs/user-guides/access/console-menu.png')} />

### Clipboard
> During an RDP session, your keyboard input is captured by the RDP console. So Ctrl+C/Ctrl+V only works on the remote RDP server.

To copy from the RDP server to your local PC,
* Copy something in the RDP server (You can use Ctrl+C).
* Open up the console menu. The copied text will appear in the clipboard.
* Select the text in the clipboard and right-click on it and select copy (You can't use Ctrl+C). 
* Now you can paste the text in the local PC.

To Copy from your local PC to remote RDP server,
* Copy something on the local PC.
* Open up the console menu.
* Right-click on the clipboard and paste it.
* Now you can paste the text on the remote computer.

### File Transfer 
If you have file transfer access, TRASA shared drive will be mounted on the remote RDP server when you access them.
* In remote RDP server, go to "This PC".
* You will see the TRASA shared drive.
<img  alt="trasa-shared-drive" src={('/img/docs/user-guides/access/trasa-shared-drive.png')} />



To upload a file from a local PC to a remote RDP server,
* Open up the console menu.
* Choose a file to upload, and click upload
* Open the TRASA shared drive and copy the file from there.


To download a file from a remote RDP server to a local PC,
* Copy the file into TRASA shared drive.
* Open the console menu.
* Click download.
* You will see the file list in your TRASA shared drive.
* Click the download icon to download the file.
<img  alt="file-download-btn" src={('/img/docs/user-guides/access/file-download-btn.png')} />



