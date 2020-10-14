---
id: intro
title: Tutorial Introduction
sidebar_label: Introduction
---

:::note
This six-part tutorial guide will provide you with a quick recipe on how TRASA can practically protect remote access in real-world scenarios.
:::

For this tutorial, we have picked a hypothetical organization. We will call it Nepsec, a not-for-profit security organization from Nepal which offers security training and awareness services.

## Nepsec employees

Nepsec has one system administrator, one security professional, and has hired one 3rd party contractor (we will call them support) to manage their AWS and GCP account. Together, they are the only ones who access their internal infrastructure frequently for maintenance and upgrade.

## Nepsec internal infrastructure

Nepsec has one windows server, one Linux server, and two web services (GitLab and discourse), making up all its internal services.

### Access requirements

1. **The administrator (looks after all Nepsec resources)** should be able to access all servers and services.
2. **Security professional (provides training to Nepsec followers)** should only access Discourse and GitLab website.
3. **3rd party vendor (does server maintaining and supports)** should only access the centos server via SSH and windows server via RDP.

## Profiling Server, services and access

:::important
It is essential to maintain service profiles and access requirements, which can then be used as a reference to assign access control policies. This is recommended regardless of you use systems like TRASA.
:::

- Server 1 : Windows server 2016.
  - Hosted in AWS
  - Server is accessed using RDP.
  - Currently directly accessed using RDP client and server IP address.
  - Existing user account in this server (privilege): `administrator` , `support`
  - Should be accessible by Nepsec `administrator` and `support` vendor
- Server 2 : Centos 7 server.

  - Hosted in Digital ocean
  - Server is accessed using ssh
  - Currently directly accessed using SSH client and server IP address.
  - Existing user account in this server (privilege): `root` , `support`
  - Should be accessible by Nepsec `administrator` and `support` vendor

- Web services:
  - Gitlab Enterprise (Web service) hosted in GCP. Currently directly accessed using IP `https://34.122.5.10`.
  - Discourse hosted in GCP. Currently directly accessed using domain `http://discuss.nepsec.io`.
  - Should be accessible by Nepsec `administrator` and `security engineer`

## [Next, we will setup TRASA](setup-trasa)
