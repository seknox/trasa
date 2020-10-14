---
id: ssh-service
title: Secure Shell (SSH)
sidebar_label: SSH
---

import useBaseUrl from '@docusaurus/useBaseUrl';

## Setup Methods

There are two ways to protect SSH servers.

1. Native 2FA agents
2. Via TRASA access proxy

Either way, you need to [create a service](../index.md#creating-a-new-service) first.

### 1. Native 2FA agents

You need to install and configure 2fa agents in all SSH servers you want to protect.
[This guide](../../native-tfa/linux-two-factor-authentication.md) will help you configure native agents.

### 2. SSH Access Proxy

To use TRASA as an SSH proxy, you need to [configure firewall rules](../../install/initial-setup.md#3-firewall-configuration-recommended) to enforce ssh access from TRASA only.

Now users need to SSH into TRASA proxy instead of the upstream server.

```shell script
ssh user@TRASA_HOST -p 8022
```

Here the port 8022 is the default port of TRASA proxy.
:::tip
You can change the default port in [config](../../system/config-reference.md#sshlistenaddr) if you want.
:::

Learn more about accessing SSH proxy [here](../../guides/user/access/ssh-connection-via-proxy.md)

## Store Password/Keys in vault

If you save password or ssh keys in the TRASA vault, users don't need to enter the upstream server password while accessing through the TRASA proxy.  
Follow [this guide](/docs/providers/vault/tsxvault#storing-service-credentials) to configure and store credentials in the vault.

## SSH Certificates

You can use TRASA as an SSH CA.

CA private keys are stored in the vault, so the vault must be [initialized](/docs/providers/vault/tsxvault#initialize-vault-one-time-only) and [decrypted](/docs/providers/vault/tsxvault#decrypt-the-vault) to use TRASA CA.

### Initialize CA

- Go to Providers page.

<img alt="download-user-ca" src={('/img/docs/providers/providers-menu.svg')} />

- Click "Certificate Authority" tab.
- Click the "Generate certs" button.

<img alt="ca-tab" src={('/img/docs/providers/ca/ca-tab.svg')} />

- Generate both "SSH User CA" and "SSH Host CA"

<img alt="generate-ca-dialog" src={('/img/docs/providers/ca/generate-ca-dialog.png')} />

### User Certificates

User certificates are used to authenticate ssh users. It can be used instead of a password or a private key.

If you configure user certificates, you don't need to store password or private keys in the vault.
During SSH access through TRASA access proxy, a temporary certificate is used to make an upstream connection. This makes remote access very easy and secure since the user doesn't need to know password or store keys.

Now we are going to tell each upstream server to trust any certificate signed by TRASA CA.  
To do that,

- Go to Providers page and click the "Certificate Authority" tab.
- Download client CA public key.

<img alt="download-user-ca" src={('/img/docs/providers/ca/download-user-ca.png')} />

- Copy the downloaded public key into upstream servers.
- Edit /etc/ssh/sshd_config of upstream server and add the following.
  `TrustedUserCAKeys <path to ca public key>`
- Restart ssh daemon.
  `sudo systemctl restart sshd`

### Host Certificates

Host certificates are used to authenticate ssh servers (hosts).

TRASA proxy will automatically validate host keys and certificates when accessing through the TRASA proxy.

#### Configure Client Device

Configuring client device is applicable when accessing SSH servers directly instead through the TRASA proxy.

- Go to Providers page and click the "Certificate Authority" tab.

- Download host CA public key.
  <img alt="download-host-ca" src={('/img/docs/providers/ca/download-host-ca.png')} />

- Copy its contents to /etc/ssh/ssh_known_hosts in following format.  
  `@cert-authority * <public key content>`

#### Configure Upstream Server

- Go to the service page in TRASA dashboard.
- Click the Edit icon in "Certificates" section.
  <img alt="services-page" src={('/img/docs/providers/ca/services-page.png')} />

- A drawer will slide from right, click the "Generate and Download" button.
  <img alt="service-certificate-slider" src={('/img/docs/providers/ca/service-certificate-slider.png')} />

- Copy the downloaded zip file to upstream server.
- Extract the files into /etc/ssh.
- Edit /etc/ssh/sshd_config and add the following.  
  `HostKey /etc/ssh/id_rsa`
  `HostCertificate /etc/ssh/id_rsa-cert.pub`
- Restart sshd daemon.  
  `sudo systemctl restart sshd`
