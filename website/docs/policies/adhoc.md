---
id: adhoc-policy
title: Adhoc Policy
sidebar_label: Adhoc Policy
---
import useBaseUrl from '@docusaurus/useBaseUrl';

Ad Hoc policy is a temporary policy given by administrator when requested. 

Users need to request ad hoc policy when
* A service is configured to enforce ad hoc policy.
* Or when normal policy denies them access.

### Enforce ad hoc policy
If you enforce ad hoc policy on a certain service, then users need to request administrator every time they access the service.
To do that,
* Go to service settings
* Click edit icon in "Configuration" card
<img  alt="edit-service" src={useBaseUrl('img/docs/policies/edit-service.png')} />
* Enable the ad hoc switch
<img  alt="adhoc-switch" src={useBaseUrl('img/docs/policies/adhoc-switch.png')} />

* Click submit.


### Grant users ad hoc access
To grant users temporary ad hoc access,

* Click  "Control" on main menu
* Click on "Ad Hoc" tab
* You will see pending adhoc requests and history of older requests.
* Click "View" on a certain pending request
* Choose expiry time
* Click on "Grant" or "Reject" 




