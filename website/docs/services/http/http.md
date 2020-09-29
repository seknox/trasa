---
id: http-service
title: Web Service (HTTPs)
sidebar_label: Web Service (HTTPs)
---

TRASA web proxy controls access to web applications (HTML5) with policies, two-factor authentication without changing a single line of code in your application.

## How does it work?

TRASA web access-proxy works a bit differently than SSH and RDP services.
Your browser should send all HTTPs traffic of your upstream service to the TRASA server. Once policy validations are completed, TRASA forwards those HTTPs requests to the upstream webserver. For this, you will need to obtain a domain name for upstream service and point DNS to the TRASA server IP address.

### Prerequisites

- A domain name for web service which points to TRASA (should be a subdomain of TRASA server A record)
- IP or domain name of primary web service (where TRASA access proxy will forward incoming requests)

:::note
Users will need to install the TRASA browser extension in their browsers to access the HTTPs service.
:::

### In this guide

- Gitlab is hosted in IP `10.10.0.10`.
- We have set the domain name `gitlab.trasa.io`, which points to the IP of the TRASA server.
- When users visit `gitlab.trasa.io`, their browser sends a request to the TRASA server. When the TRASA server receives a request for `gitlab.trasa.io` It checks policies and forwards traffic to GitLab server's origin IP `10.10.0.10`.

## Creating Web Service Profile

Create a new service selecting `http` as the service type.
![create web service](./create-http.png 'Integrate new web service')

## Configuring HTTP service.

Configure Proxy
