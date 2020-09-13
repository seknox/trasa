---
id: initial-setup
title: Initial Setup
sidebar_label: Initial Setup
---

> If you signed up for TRASA Cloud service, you can skip this guide.

## Setting up TRASA vault ([TsxVault](docs-link-tsxvault "TsxVault"))
TRASA stores all keys and secrets in secure vault known is [TsxVault](docs-link-tsxvault "TsxVault").

  _Passwords_, _Secret keys_, _API tokens_ etc. are needed by TRASA to integrate with 3rd party services. For example, FCM tokens, Email config settings, IDP integration keys.

So before we begin, the first after installtion is to initialize [TsxVault](docs-link-tsxvault "TsxVault") 

## Setting up FCM configs
FCM is mobile notification service offered by Google. In order to perform TRASA U2F based 2FA. 
By default, TRASA proxies U2F request via Seknox's Cloud server. It means that by default, when users chooses to perform U2F via TRASA mobile app, the U2F request is forwarded to Seknox's Cloud server and proxied back when user confirms response (authorize or deny).

If you do not want to route U2F requests via Seknox server, you can configure your own FCM service. 

## Setting up Email
TRASA supports mailgun and authenticated SMTP protocol to forwarded emails generated within TRASA.
