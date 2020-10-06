---
id: glossary
title: Glossary
sidebar_label: Glossary
description: Words and their meaning used in TRASA
---

## User

A user is a profile of anyone who accesses remote servers and services. Typically, they are your employees' profile, 3rd party contractors, DevOps teams, or service accounts.

## User Identity Provider (uIDP)

A User IDP is a source that TRASA uses to import, manage, and authenticate users. TRASA has a built-in identity provider (TrasaIDP) and can be connected with other Identity providers, including G Suite, Okta, LDAP, and FreeIPA.

## Service

A Service is a profile of any network applications which sit behind TRASA proxy. TRASA supports Web application, SSH, Remote Desktop (RDP), and Database services.

## Upstream Service

An Upstream Service is an actual service that runs behind TRASA. E.g., include SSH or RDP service for your internal servers, Web application of your firewall console.

## Service Identity Provider (sIDP)

A service IDP is the source which TRASA uses to import and manage service profile. A service profile can be directly created in TRASA or imported from cloud service providers such as AWS, GCP, Digital Ocean.

## Privilege

A privilege is a username or role which your users use to access service.
Privilege typically has roles, permissions, and policies attached to them managed by the server and service that privilege belongs.
e.g.:

- `root` is a privilege usually found in Linux OS. A user with email _james@nepsec.com_ logs into the centOS server as `root` privilege.
- `Administrator` is a privilege usually found in Windows systems. A user with email _james@nepsec.com_ access the windows server as `Administrator` privilege.

## Access Proxy

Access Proxy is a reverse proxy server which has two purposes:

- Manage access to your server and services,
- Control and block unauthorized access
  TRASA access proxy currently supports HTTPs, SSH, RDP, and Database protocols.

## Device Hygiene

Device hygiene is a security state of user's devices (workstations and mobile devices) that users use to access servers and services.

## Policy

A policy defines constraints and restriction which directs TRASA to either block or allow access to protected service.

### Static Policy

Based on time, location, and health of device hygiene.

### Dynamic Policy

Based on risk scoring scored by TRASA AI (only available in enterprise edition)

## Adhoc permission

Adhoc permission enables explicit permission management to access service.

## Access Map

Access Map defines how a user can access a service with specific privileges. The how is directed via policy. Before a user can access a service, the TRASA administrator must assign a user to that service along with privilege and policy.

## Vault

A vault is where TRASA stores secrets. There are two types of secrets that TRASA manages.

- **Upstream Service Secrets -** Passwords and Keys of upstream services.
- **Integration Keys -** Secrets and Keys of external services that TRASA connects during the process of integration.

TRASA has a built-in vault named `TsxVault`, which can be used to store both of the secrets mentioned above. While Integration Keys are always stored in `TsxVault`, Upstream Service Secrets can also be stored in external secret storage providers such as HashiCorp Vault, AWS KMS, and GCP KMS.

## TRASA_URL

URL address (usually a domain name) which points to TRASA server.
