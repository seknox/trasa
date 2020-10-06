---
id: protect-ssh
title: Part 4 - Protect SSH Access
sidebar_label: Part 4 - Protect SSH Access
---

:::note

This tutorial will show you how to protect ssh service with TRASA. Especially, you will learn to

- [STEP 1: Install and setup TRASA](#step-1-install-and-setup-trasa)
- [STEP 2: Create Policy](#step-2-create-policy)
- [STEP 3: Create Service](#step-3-create-service)
- [STEP 4: Assess Mapping](#step-4-assess-mapping)
- [STEP 4: Access](#step-4-access)

:::

### STEP 1: Install and setup TRASA

Follow [installation guide](../install/installation.md) and [Initial Setup](../install/initial-setup.md).

### STEP 2: Create Policy

- Click the menu button on top left to open main menu drawer.
  <img alt="main-menu" src={('/img/docs/main-menu.png')} />
- Click the Control menu to open policy page.
- Click the "Create new policy" button.
  <img alt="create-policy-btn" src={('/img/docs/quickstart/create-policy-btn.png')} />
- Enter a policy name and click next.
- Click the "Mandatory 2FA" switch to enable second factor authentication
  <img alt="2fa-enable" src={('/img/docs/quickstart/2fa-enable.png')} />
- Click the "Session Recording" menu and enable it.
  <img alt="session-recording-enable" src={('/img/docs/quickstart/session-recording-enable.png')} />
- Click "Day and Time" to configure weekday and time specific policy.
- Select week days and time range to allow access.
- Click "add" button.
  <img alt="add-day-time-policy" src={('/img/docs/quickstart/add-day-time-policy.png')} />
- Click next and review the policy to be created.
- If everything looks good, click the "Submit" button.

:::tip
Go to [policies reference](../policies/basic-policy.md) to know more about static policies
:::

<br />

---

<br />

### STEP 3: Create Service

- Open main menu and click Services
- Click "Create new service" button.
  <img alt="create-service-btn" src={('/img/docs/quickstart/create-service-btn.png')} />
- A drawer will slide from the left. Fill in the details of upstream server you want to connect through TRASA.
  <img alt="create-new-service" src={('/img/docs/quickstart/create-new-service.png')} />
- Click submit. You will be redirected to the newly created service page.

<br />

---

<br />

### STEP 4: Assess Mapping

> Right now you are assigning yourself to the service. You can assign other users too when you [create](../users/crud.md) them.

- Click the "Access Map" tab.
  <img alt="access-map-tab" src={('/img/docs/quickstart/access-map-tab.png')} />
- Click the "Assign user" button.
  <img alt="assign-user-btn" src={('/img/docs/quickstart/assign-user-btn.png')} />
- Choose the user and policy you just created.
  <img alt="assign-user-dialog" src={('/img/docs/quickstart/assign-user-dialog.png')} />

- On privilege, enter the username of the upstream server.
- Click submit.

<br />

---

<br />

### STEP 4: Access

- Click the dropdown button on the top right and click on "My Account" menu.
  <img alt="my-account-btn" src={('/img/docs/quickstart/my-account-btn.png')} />
- You will see the newly assigned service.
  <img alt="my-service" src={('/img/docs/quickstart/my-service.png')} />
- Click connect and then the privilege (username) you configured earlier.
- Allow pop-up if your browser blocks it.
- Enter password.
- Choose TOTP and enter TOTP code from mobile app.
- You will be asked to save new SSH host key, press "y" to do that.

You will be logged in the upstream server through TRASA.

Now go through the docs to configure TRASA according to your needs.
