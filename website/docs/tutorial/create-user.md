---
id: create-users
title: Part 2 - Create User
sidebar_label: Part 2 - Create User
---

export const Hlt = ({children, color}) => ( <span style={{
      backgroundColor: color,
      borderRadius: '2px',
      color: '#fff',
      padding: '0.2rem',
    }}>{children}</span> );

:::tip
TRASA IDP is builtin default user identity provider. Creating user profile directly in TRASA is one way to create user profile. You can also use available user identity providers to sync user profiles in TRASA.
 Checkout [identity providers](../providers/providers.md)
:::

By default, root account (which you used to login to TRASA) is already created during installtion.
To add other user's

<img alt="enrol device" width='1000' src={('/img/docs/tutorial/default-user.png')} />

## Creating profile for Nepsec security professional

- Head over to users page and click <Hlt  color="#1877F2">Create User</Hlt> button.

![Create User](/img/docs/users/trasa/create-user.png 'Create User')

- If your request is successful, the user will receive an account activation link via email. Though immediately, you will also receive a link that you can share with the user.

![Activation Link](/img/docs/users/trasa/verification-link.png 'Activation Link')

## Creating profile for 3rd party managed service provider

- Head over to users page and click <Hlt  color="#1877F2">Create User</Hlt> button.

![Create User](/img/docs/users/trasa/create-user.png 'Create User')

- If your request is successful, the user will receive an account activation link via email. Though immediately, you will also receive a link that you can share with the user.

![Activation Link](/img/docs/users/trasa/verification-link.png 'Activation Link')

after creating all user's, our user dashbaord looks like this:
<img alt="enrol device" src={('/img/docs/tutorial/all-users.png')} />
