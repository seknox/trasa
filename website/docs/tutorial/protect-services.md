---
id: protect-services
title: Part 4 - Protect Services
sidebar_label: Part 4 - Protect Services
---

:::note

This tutorial will show you how to protect ssh service with TRASA. Especially, you will learn to

- [Create SSH service profile in TRASA and map user access](#create-ssh-service-profile-in-trasa-and-map-user-access)
- [Create RDP service profile in TRASA and map user access](#create-rdp-service-profile-in-trasa-and-map-user-access)
- [Create Web service profile in TRASA and map user access](#create-web-service-profile-in-trasa-and-map-user-access)
  - [Service profile for Gitlab](#service-profile-for-gitlab)
  - [Service profile for Discourse admin](#service-profile-for-discourse-admin)

:::

<br /><br />

---

<br />

## Create SSH service profile in TRASA and map user access

In the video below,

1. We will create a service profile for Centos7 which is hosted in digital ocean.
2. Assign access to administrator with `full access` policy and support with `trusted device` policy.


<iframe width="100%" height='600' src="https://www.youtube.com/embed/M7pmL8h-OXI?list=PLZOFebo-o2K44zdkUPWnGO_cTz6KjNAnN" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<br /><br />

---

<br />

## Create RDP service profile in TRASA and map user access

In the video below,

1. We will create a service profile for windows 2016 server which is hosted in AWS.
2. Assign access to administrator with `full access` policy and support with `friday` policy.


<iframe width="100%" height='600' src="https://www.youtube.com/embed/EOUzO9vymug?list=PLZOFebo-o2K44zdkUPWnGO_cTz6KjNAnN" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<br /><br />

---

<br />

## Create Web service profile in TRASA and map user access

Configuring HTTP (web) service is a bit different from SSH or RDP.

1. First, you will need to configure DNS to point web traffic to the TRASA server.
2. Then, configure proxy detail in the TRASA dashboard to forward incoming web traffic to the upstream web application.

:::important
In the video guide below, we only show how to configure HTTP service in TRASA. For a fully working setup, you will need to configure a DNS record that points to the TRASA server.
:::

### Service profile for Gitlab

In the video below,

1. We will create a service profile for Gitlab ce which is hosted in GCP.
2. Assign access to administrator and security professional with `trusted device` policy.


<iframe width="100%" height='600' src="https://www.youtube.com/embed/7cxquNp7Lm8?list=PLZOFebo-o2K44zdkUPWnGO_cTz6KjNAnN" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<br /><br />

---

<br />

### Service profile for Discourse admin

In the video below,

1. We will create a service profile for Gitlab ce which is hosted in GCP.
2. Assign access to administrator and security professional with `trusted device` policy.

   <iframe width="100%" height='600' src="https://www.youtube.com/embed/VzwARHXnLWI?list=PLZOFebo-o2K44zdkUPWnGO_cTz6KjNAnN" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
