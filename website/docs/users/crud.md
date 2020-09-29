---
id: creating-updating-users
title: Creating and updating Users
sidebar_label: Creating and updating users
---

export const Hlt = ({children, color}) => ( <span style={{
      backgroundColor: color,
      borderRadius: '2px',
      color: '#fff',
      padding: '0.2rem',
    }}>{children}</span> );


## Creating User

TRASA IDP is builtin default user identity provider.

:::important
TRASA expects unique email and username across the organization.
:::

:::tip
Assign short and easy to use username. Emails are usually lengthy, so short usernames come in handy while signing in.
:::

+ Head over to users page and click <Hlt  color="#1877F2">Create User</Hlt> button.

![Create User](/img/docs/users/trasa/create-user.png 'Create User')

+ If your request is successful, the user will receive an account activation link via email. Though immediately, you will also receive a link that you can share with the user.

![Activation Link](/img/docs/users/trasa/verification-link.png 'Activation Link')

Sample of a user profile
![User Profile](/img/docs/users/user-profile.png 'User Profile')


##  Updating or deleting user

+ You can update the user profile by clicking on `:` menu icon in user profile card.

![Update User](/img/docs/users/trasa/update-user.png 'Update User')
