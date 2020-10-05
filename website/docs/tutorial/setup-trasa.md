---
id: setup-trasa
title: Part 1 - Setup TRASA
sidebar_label: Part 1 - Setup TRASA
---

In this first part, we will cover installation and setup of TRASA server.

## Before Installation

1. **Linux server**

   We've created 1 core 2 GB ram 20 GB storage Ubuntu server. We will call this server **Nepsec TRASA server**. Once this server is ready, install and setup Openssh server and [Docker](https://docs.docker.com/engine/install/) in this server.

2. **Domain name**

   We've setup DNS with `A` record `nepsec.trasa.io` which points to our server. Setup a domain in your control.

## Install

:::note
We are using docker install for demonstration. For other installation options, refer to [Install Guides](../install/installation)
:::

SSH to Linux instance (TRASA server) you created in previous step.

```shell script

# Run Postgresql database
sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres

# Run Redis
sudo docker run -d -p 6379:6379 --name redis redis

# Run Guacd Server
sudo docker run -d --rm --name guacd -p 127.0.0.1:4822:4822 -v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac --user root seknox/guacd:v0.0.1

# Run TRASA server
sudo docker run --link db:db \
--link guacd:guacd \
--link redis:redis \
-p 443:443 \
-p 80:80 \
-p 8022:8022 \
-e TRASA.LISTENADDR=<NEPSEC.TRASA.IO> \ # <- Replace it with your preferred trasa domain name.
-v /tmp/trasa/accessproxy/guac:/tmp/trasa/accessproxy/guac \
seknox/trasa:v0.0.1

```

<br />

---

<br />

## Setup Root Account

TRASA server should be ready from previous step.
Enter [TRASA_URL](/docs/getting-started/glossary#trasa_url) in your browser.

In our case, we setup domain `nepsec.trasa.io` so we enter this address in browser.
<img alt="dashboard login" src={('/img/docs/tutorial/dashboard-login.png')} />

<br /><br />

When TRASA is installed, default system account `root` is created for you with default password `changeme`.
Enter username and password (root account) in login box.

:::note
TRASA requires two factor authentication by default and TRASA mobile app is default supported authenticator. Since this is your first login, you need to enrol device first:

Get TRASA authenticator from [Play Store](https://play.google.com/store/apps/details?id=com.trasa&hl=en) or [App Store](https://apps.apple.com/np/app/trasa/id1411267389).

:::

Since this is your first time logging into TRASA, you have not yet added your 2FA device yet.
QR code will appear on screen.
<img alt="qr-code" src={('/img/docs/user-guides/device/qr-code.png')} />

### Enrol Steps:

<img alt="enrol device" src={('/img/docs/tutorial/enrol-mobile-device.svg')} />

1. Press the `+` button (buttom right).
2. Press QR image icon button. This will open in-app camera.
3. Scan the QR image from TRASA dashbaord
4. If everything goes well, you will see the following icon on your app

<img alt="enrol device" src={('/img/docs/tutorial/device-enroled.svg')} />

- Press the icon to get TOTP codes

## Login with root account

Press `login` button in dashboard page (where QR code is shown), you will be redirected to Login page again.

1. Enter root account credentials
2. Once the credentials are validated, you will see **second-step verification page**
   <img alt="enrol device" src={('/img/docs/tutorial/dashboard-totp.svg')} />
3. From your TRASA mobile app, note totp code and enter in dashboard to proceed login.
   <img alt="enrol device" src={('/img/docs/tutorial/enter-totp.svg')} />

4. Server will validate your totp code and will redirecto to dashboard overview page.
   <img alt="enrol device" src={('/img/docs/tutorial/first-dashboard.png')} />

## [Next - Create User accounts](create-users)
