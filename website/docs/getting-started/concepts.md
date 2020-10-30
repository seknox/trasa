---
id: concepts
title: Concepts
sidebar_label: Concepts
description: Basic concepts for using TRASA
---






TRASA Access Proxy is a drop in upgrade for your homegrown bastian/jump server.
Whether you are using a linux server (configured as Jump server) or Microsoft Remote Desktop Gateway, TRASA offers all those features along with best practices enabled and configurable by default.

### TRASA server as an access point for your internal infrastructure
Access Proxy is just a reverse proxy that also understans RDP, SSH, HTTP, and Database protocols and makes a forwarding decision based on the administrator's policy. For TRASA access proxy to work, the network must be configured in a way such that every remote access to your server and services are only allowed from the server IP address in which TRASA is installed.
<!-- <img alt="enrol device" src={('/img/docs/tutorial/all-users.png')} /> -->

### Configuring your network
+ Make changes in your network firewall such that ingress traffic to your server and services listening for SSH, RDP, HTTPs, and DB traffic are only allowed from TRASA server IP address.
 
