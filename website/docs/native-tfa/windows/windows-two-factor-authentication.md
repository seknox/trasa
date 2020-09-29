---
id: windows-two-factor-authentication
title: Windows Two Factor Authentication
sidebar_label: Windows
---

export const Highlight = ({children, color}) => ( <span style={{
      backgroundColor: color,
      borderRadius: '2px',
      color: '#fff',
      padding: '0.2rem',
    }}>{children}</span> );

:::note
Window two factor authentication is supported via TRASA Windows Credential Provider.
:::

## Prerequisite

- User profile in TRASA
- Service profile in TRASA
- TRASA Windows tfa agent installer.
- Windows OS (windows 7 and above)
- [visual c++ redestributable](https://aka.ms/vs/15/release/vc_redist.x64.exe) (Optional)

## Installation

Download [TrasaWIN](https://storage.googleapis.com/trasa-public-download-assets/trasaWIN/TRASA-TFA-20-04.msi) and proceed installation.

:::caution
Do not reboot or signout from your computer until you configure agent.
Broken configuration may lock your access to operating system.
:::

After installation and before you close the installer, it is very important to configure agent.

Check on <Highlight color="#1877F2">Launch TrasaWIn to configure now</Highlight> checkbox which will launch configuration panel.

![Install TRASA TFA agent](/img/docs/tfa/windows/install-trasa-win-tfa.png 'Install TRASA TFA agent')

## Configuration

If you checked on "Launch TrasaWIn to configure now" checkbox, configuration application will open. You will need to input configuration values in required field.

### What values goes in input fields?

- ServiceID: Copy from service profile page
- ServiceKey: Copy from service profile page
- TRASA server address: IP or domain of where TRASA server is hosted.
- Offline Users: Usernames which are allowed to login if the agent could not contact TRASA server (eg. network failure)
- Skip TLS verification: Allows to connect to TRASA server if self signed certificate is used at port 443.
  ![Configuring value from service profile](/img/docs/tfa/windows/config-from-trasa.png 'Configuring TRASA TFA agent')

In following image, you can see appID appsecret and api hostname entered as per auth app created earlier. Note api host name "app.trasa.io" always remains same for TRASA SaaS users and can be custom url for self hosted (On-Premise) TRASA users.

Below is example on how configuration would look like.
![Configuring TRASA TFA agent](/img/docs/tfa/windows/configure-trasawin-tfa.png 'Configuring TRASA TFA agent')

Once you are ready with required configuration values, click <Highlight color="#1877F2">Save Configuration</Highlight> button.
It will

1. verify configuration values
2. save it in file if verification is successful.
   ![Confirming verification](/img/docs/tfa/windows/check-config.png 'Confirming verification')

## Finishing

If your verification was success, you will be prompted for TRASA tfa process in next login.

To check, you can try swithing user (from alt+F4 key).
![TRASA Credential Provider](/img/docs/tfa/windows/trasa-credprov.png 'TRASA Credential Provider')

If your username and password validation was success, you will be prompted with TRASA tfa prompt.

- You will need to enter your TRASA username or email address.
- On Choose 2FA method, you can leave it empty for push U2F or select TOTP option for TOTP.
  ![TRASA TFA Prompt](/img/docs/tfa/windows/trasa-tfa-prompt.png 'TRASA TFA Prompt')

## FAQ

### What happens if agent could not resolve TRASA server?

When user tries to login to windows protected with TRASA TFA agent, agent will contact TRASA server for 2FA verification.

_What happens when agent cannot contact TRASA server?_

In case that you have a network problem, and the agent cannot resolve TRASA server address, your access will be blocked. To overcome this situation, TRASA allows you to set an emergency access account or offline user.

Offline user account can be any currently being used user account in your windows logon (domain or local account). Do note that the username must be exactly matched to the existing user account.

:::tip
If the user is local account, but windows is domain joined, you will need to assign full user path in format
`local-workgroup-name\username`
:::
