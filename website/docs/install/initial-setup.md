---
id: initial-setup
title: Initial Setup
sidebar_label: Initial Setup
---

<!-- > If you signed up for TRASA Cloud service, you can skip this guide. -->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

<br />



## 1. Root account setup



Once TRASA server is up and running, you can open web browser to access TRASA dashboard (listening at [TRASA_HOST](/docs/getting-started/glossary#TRASA_HOST) in browser which is an IP address or domain name.)

By default, a `root` user account will be created for you with default password `changeme`. As an administraotr, you will need to setup this root account before you can access TRASA.




Follow account setup steps detailed at [account setup guide](/docs/users/account-setup)


<br />



## 2. Initialize encrypted storage - [TsxVault](/docs/getting-started/glossary#tsxvault)

TRASA has secure encrypted storage to which is used to store service credentials like password, private keys and integration keys.
Follow steps at [initializing TsxVault](/docs/providers/vault/tsxvault) to enable secret storage.

<br />

## 3. Setup FCM with TRASA FCM proxy (Optional, Recommended)

TRASA push U2F is very convinient way of authorizing 2FA process. With push U2F, user's do not need to enter 6 digit totp code every time they need to verify second step verification process and is also immune to phishing attacks on totp codes. 

This feature requires sending push notification to user's mobile device. To enable it, register with [TRASA FCM Proxy](/docs/system/fcm-setup)

To use all features of TRASA, you need to setup

## 4.  Email setup (Optional, Recommended)

To receive emails and security alerts from TRASA you will need to integrate TRASA with your existing email provider. Follow [Email setup](/docs/system/email-setup) guide to setup email.

<br />



## 5. Configuring Network Firewall  (Optional, Recommended)

TRASA access proxy can only control access if traffic passes through it. To ensure security policy is enforced on access proxy, you should configure a network firewall so that every remote access to your server and services is only routed and allowed from the IP address of the TRASA server. 

:::note
TRASA also supports native two-factor authentication integration (with installable agents that protect windows server, Linux server, and hardware appliance). If you are using TRASA just for native two-factor authentication, you can skip configuring your network because agents will communicate with TRASA server for authorization.
:::

<!-- ######################################################################################## -->

<Tabs
defaultValue="aws"
values={[
{label: 'AWS', value: 'aws'},
{label: 'GCP', value: 'gcp'},
{label: 'DigitalOcean', value: 'digitalocean'},
]}>

<!-- ######################################################################################## -->
<TabItem value="aws">

**Create a security group for TRASA**

In EC2 management console,

- Go to Security groups.
- Click the "Create security group" button.
  <img alt="create-security-grp-btn" src={('/img/docs/cloud/aws/create-security-grp-btn.png')} />

- Fill in the name and description.
- In inbound rules section, click the "Add rule" button.

  <img alt="create-security-grp-btn" src={('/img/docs/cloud/aws/create-rule-btn.png')} />

- Choose the "SSH" type and "Custom" source.
- Add TRASA IP on the "source IP" field.
  <img alt="inbound-rule-sample" src={('/img/docs/cloud/aws/inbound-rule-sample.png')} />

Now use this security group to allow SSH in all instances you want to protect with TRASA.
</TabItem>

<!-- ######################################################################################## -->

<TabItem value="gcp">

**Create a firewall rule for TRASA**

We will make two rules, one to block all remote access requests (`block-all-remote`) and one to allow requests from TRASA (`allow-from-trasa`).

If two rules have the same priority value, the blocking rule will override, so we will give `allow-from-trasa` more priority than `block-all-remote`

- Go to VPC Network-> Firewall on Main menu
  <img alt="firewall-menu" src={('/img/docs/cloud/gcp/firewall-menu.png')} />

- Click the "Create Firewall Rule" button
  <img alt="create-firewall-rule-btn" src={('/img/docs/cloud/gcp/create-firewall-rule-btn.png')} />

- Fill in the following parameters
  - "Target tag name" : `allow-from-trasa`
  - "Action on match" : "Allow"
  - "Source IP range" : [TRASA_IP]
  - "Ports and protocols" : "tcp:22, tcp:3389"
  - "Priority" : 999

<img alt="new-firewall-rule" src={('/img/docs/cloud/gcp/new-firewall-rule.png')} />

- Create another rule with the following parameters:
  - "Target tag name" : `allow-from-trasa`
  - "Action on match" : "Deny"
  - "IP source" : "0.0.0.0/0"
  - "Ports and protocols" : "tcp:22, tcp:3389"
  - "Priority" : 1000

<img alt="blocking-firewall-rule" src={('/img/docs/cloud/gcp/blocking-firewall-rule.png')} />

:::tip
\*You can give different priority and target tag name.
Just make sure `allow-from-trasa` has more priority than `block-all-remote`.
:::

:::tip
Learn more about priority [here](https://cloud.google.com/vpc/docs/firewalls#priority_order_for_firewall_rules)
:::

Use these two firewall rules for all instances you want to protect with TRASA.




</TabItem>

<!-- ######################################################################################## -->

<TabItem value="digitalocean">

**Create a firewall rule for TRASA**

- Go to "Networking" section from the main menu.
- Go to "Firewalls" tab. 
- Click the "Create Firewall" button.
  <img alt="network-firewall-create" src={('/img/docs/cloud/do/network-firewall-create.png')} />

- Fill in the name and description
- Enter TRASA IP as source IP in inbound rules
  <img alt="inbound-rule" src={('/img/docs/cloud/do/inbound-rule.png')} />

- Enter a tag name to apply this rule.
- Click the "Create Firewall" button
  <img alt="create" src={('/img/docs/cloud/do/create.png')} />

Use this firewall rule to give ssh access to all droplets you want to protect with TRASA.
</TabItem>

<!-- ######################################################################################## -->
</Tabs>
