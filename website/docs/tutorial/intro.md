---
id: intro
title: Tutorial Introduction
sidebar_label: Introduction
---

:::note
This six-part tutorial guide will provide you with a quick recipe on how TRASA can practically protect remote access in real-world scenarios.
:::

For this tutorial, we have picked a hypothetical organization. We will call it Nepsec, a not-for-profit security organization from Nepal which offers security training and awareness services.

## Nepsec internal infrastructure

Nepsec has the following server and service, which makes up all its internal services.

- Server 1 - Windows server 2016.
  - Hosted in AWS
  - RDP server listening in port 3389
- Server 2 - Centos 7 server.

  - Hosted in Digital ocean
  - SSH listening in port 22

- Gitlab Enterprise (Web service) hosted in GCP Kubernetes

## Nepsec employees

Nepsec has one system administrator and has hired one 3rd party contractor to manage their AWS and GCP account. Together, they are the only ones who access their internal infrastructure frequently for maintenance and upgrade.

Since this is your first time logging into TRASA, you have not added your 2FA device yet.

## [Next, we will setup TRASA](setup-trasa)
