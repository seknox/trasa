---
id: dynamic-access
title: Dynamic Access
sidebar_label: Dynamic Access
---
import useBaseUrl from '@docusaurus/useBaseUrl';


Dynamic access feature allows you to create services and accessmaps automatically when user tries to access them.

To enable dynamic access,
* Go to Services page.
* Click on "Settings" tab.
<img  alt="service-settings" src={useBaseUrl('img/docs/access-map/service-settings.png')} />

* Click on "Dynamic Access" section
* Enable the status checkbox.
* Select user group and policy to authorise dynamic access.
<img  alt="dynamic-access-settings" src={useBaseUrl('img/docs/access-map/dynamic-access-settings.png')} />

Now when users enters an IP/hostname which is not yet used in any service, a new service will be automatically created.
And also, user will be automatically assigned to that service.
