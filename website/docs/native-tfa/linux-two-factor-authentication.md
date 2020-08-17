---
id: linux-two-factor-authentication
title: Linux Two Factor Authentication
sidebar_label: Linux
---

:::note
Linux two factor authentication is supported via TRASA PAM (Pluggable Authentication Modules).
:::


trasaPAM is a PAM module that handle's second factor authentication in *nix systems.

## Installation

:::tip
You may require root privilege while setting up trasaPAM.
:::
:::caution
Keep a separate ssh connection open with root privilege till whole process finishes so that you can always troubleshoot incase configuration goes wrong.
:::


1. Download [TrasaPAM](https://storage.googleapis.com/trasa-public-download-assets/TrasaPAM/TrasaPAM.zip) 

2. Unzip files `# unzip TrasaPAM.zip -d destination_folder`

## Configure
Inside extracted directory, open `# vi trasapam.toml` and configure with following data:

```
[trasaPAM]
trasaServerURL = "<address of trasa server>"
serviceID = "<serviceID(copy from service profile)>"
serviceKey = "<serviceKey(copy from service profile)>"
offlineUsers = "<users to allow in case PAM module cannot contact TRASA server>"
insecureSkipVerify = <boolean value. false by default. set true if TRASA server is using self signed TLS certificate. >
```

4. Copy config file `trasapam.toml` to `/etc/trasa/config/trasapam.toml`

5. Copy `trasapam.so` file to `/lib/security/` for debian or `/lib64/security/`in case of centOS.



## Configure SSH

open `/etc/ssh/sshd_config` file and make sure `UsePam yes` and `ChallengeResponseAuthentication yes` is set.

## Make trasaPAM PAM aware
+ Open `/etc/pam.d/sshd`
+ Add `auth required trasapam.so` for debian or `auth required /lib64/security/trasapam.so` in case centOS at end of the file.


## Finishing
Restart sshd to correctly load pam module: `$ sudo systemctl restart sshd`


## Testing
```
$ ssh root@test-machine
$ password:
$ Enter your trasaID: 
$ Choose TFA method (enter blank for U2F):
```
