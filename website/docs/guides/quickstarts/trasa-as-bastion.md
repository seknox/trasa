---
id: trasa-as-bastian
title: TRASA Access Proxy
sidebar_label: TRASA Access Proxy
---

TRASA Access Proxy is a drop in updgrade for your homegrown bastian/jump server.
Whether you are using linux server (configured as Jump server) or Microsoft Remote Desktop Gateway, TRASA offers all those features along with best practices enabled and configurable by default.

## TRASA server as an access point for your internal infrastructure
Access Proxy is just a reverse proxy that happens to know RDP, SSH, HTTP and Database protocols and makes forwarding decision based on policy defined by administrator. For TRASA access proxy to work, network must be configured in a way such that every remote access to your server and services are only allowed from server IP address in which TRASA is installed.

## Configuring your network
+ Make changes in your network firewall such that ingress traffic to your server and services listening for SSH, RDP, HTTPs and DB traffic are only allowed from TRASA server IP address.
 