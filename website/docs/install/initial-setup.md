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

- Download TRASA mobile app from [App Store](https://apps.apple.com/us/app/trasa/id1411267389) or [Play Store](https://play.google.com/store/apps/details?id=com.trasa)
- Open [TRASA_HOST](/docs/getting-started/glossary#TRASA_HOST) in your browser.
- Use default credentials to login (root:changeme)

Since this is your first time logging into TRASA, you have not added your 2FA device yet.
QR code will appear on the screen.

- Open TRASA mobile app and press + button on the bottom right and then press the QR icon.

<img alt="enrol device" src={('/img/docs/tutorial/enrol-mobile-device.svg')} />

- Scan the QR code on the browser.
- If everything goes well, you will see the following icon on your app.

<img alt="enrol device" src={('/img/docs/tutorial/device-enroled.svg')} />

- Press the icon to get TOTP codes.

Now 2FA device is added.

- Try logging in again.
- Now, you need to choose TOTP and enter the TOTP code from the mobile app.

<br />

---

<br />

## 2. System Setup (Optional)

To use all features of TRASA, you need to setup

1.  [Secret Store](/docs/providers/vault/tsxvault)
2.  [FCM](../system/fcm-settings.md)
3.  [Email](../system/email-settings.md)

<br />

---

<br />

## 3. Firewall configuration (Recommended)

To use TRASA as an access proxy, you need to enforce remote access to services through TRASA only.
To do that, you need to configure firewall rules accordingly.

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
  <img alt="security-groups-menu" src={('/img/docs/cloud/aws/security-groups-menu.png')} />

- Click the "Create security group" button.
  <img alt="create-security-grp-btn" src={('/img/docs/cloud/aws/create-security-grp-btn.png')} />

- Fill in the name and description.
- On Inbound rules, click the "add rule" button.

- Choose the "SSH" type and "Custom" source.
- Add TRASA IP on the "source IP" field.
  <img alt="inbound-rule-sample" src={('/img/docs/cloud/aws/inbound-rule-sample.png')} />

Now use this security group to allow SSH in all instances.
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

Use these two firewall rules for all instances.

**Configuring SSH keys in Google Cloud**

By default, google cloud uses OS Login, which uses google identity to manage SSH keys.
To use TRASA to manage your SSH keys, you need to disable OS Login.
Then you need to add ssh keys to the instance or project.

- Go to [google cloud compute instances page](https://console.cloud.google.com/compute/instances) and click on the instance you want to configure.
- Click the "Edit" button
  <img alt="edit-instance-btn" src={('/img/docs/cloud/gcp/edit-instance-btn.png')} />

- Generate a new ssh key
  `ssh-keygen -t rsa -b 4096 -f ~/.ssh/[KEY_FILENAME] -C [USERNAME]`
- Scroll down to the "custom metadata" section and add a new key `enable-oslogin`:`FALSE`
- Click the "add item" button under SSH Keys section
- Copy the contents of [KEY_FILENAME].pub into the field
  <img alt="instance-level-metadata" src={('/img/docs/cloud/gcp/instance-level-metadata.png')} />
- Click Save
- [Save the contents of [KEY_FILENAME] in TRASA vault](/docs/providers/vault/tsxvault)

:::tip
If you want to configure this for all instances of a project, go to the [Metadata](https://console.cloud.google.com/compute/metadata) menu on Compute Engine page.
<img alt="project-level-metadata" src={('/img/docs/cloud/gcp/project-level-metadata.png')} />

:::
</TabItem>

<!-- ######################################################################################## -->

<TabItem value="digitalocean">

**Create a firewall rule for TRASA**

- Go to Networking-> Firewalls on the Main menu.
- Click the "Create Firewall" button.
  <img alt="network-firewall-create" src={('/img/docs/cloud/do/network-firewall-create.png')} />

- Fill in the name and description
- Enter TRASA IP as source IP in inbound rules
  <img alt="inbound-rule" src={('/img/docs/cloud/do/inbound-rule.png')} />

- Enter a tag name to apply this rule.
- Click the "Create Firewall" button
  <img alt="create" src={('/img/docs/cloud/do/create.png')} />

Use this firewall rule to give ssh access to all instances.
</TabItem>

<!-- ######################################################################################## -->
</Tabs>
