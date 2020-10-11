---
id: test-access-to-services
title: Part 5 - Test access to services
sidebar_label: Part 5 - Test access to services
---

We've configured service profile and access routes in [Part 4](protect-services) or this tutorial.

Now, we will test access to services.

:::important

In TRASA,

- SSH service can be accessed from both dashboard (no external client required) or Linux, Mac, Windows terminal clients.
- RDP services are only accessible from the dashboard. I.e. does not support RDP from other clients.
- Https service can only be accessed from the browser and only if the TRASA browser extension is installed in the browser.

:::

## Viewing available service access

You can view available services that can be accessed from [My Page](../getting-started/glossary#my-route). Every service that is assigned to you will be listed on that page.

Below is a screenshot of available access applications and services for the Nepsec administrator.

<img alt="my services" src={('/img/docs/tutorial/my-services.png')} />

## Testing SSH access

<iframe src="/img/docs/tutorial/ssh-test.mp4" frameborder="0" allowfullscreen width="100%" height='600'></iframe>

## Testing Windows RDP access

<iframe src="/img/docs/tutorial/rdp-test.mp4" frameborder="0" allowfullscreen width="100%" height='600'></iframe>

## Testing Web(Gitlab) access

<iframe src="/img/docs/tutorial/http-test.mp4" frameborder="0" allowfullscreen width="100%" height='600'></iframe>
