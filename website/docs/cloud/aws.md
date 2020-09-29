---
id: amazon-web-services
title: Setup AWS Through TRASA
sidebar_label: AWS
---

import useBaseUrl from '@docusaurus/useBaseUrl';



## Create a security group for TRASA

In EC2 management console,
* Go to Security groups
<img  alt="security-groups-menu" src={('/img/docs/cloud/aws/security-groups-menu.png')} />

* Click the "Create security group" button
<img  alt="create-security-grp-btn" src={('/img/docs/cloud/aws/create-security-grp-btn.png')} />

* Fill in the names and description
* On Inbound rules, click the "add rule" button
<img  alt="add-rule" src={('/img/docs/cloud/aws/add-rule.png')} />

* Choose "SSH" type and "Custom" source
* Add TRASA IP on source IP field
<img  alt="add-rule" src={('/img/docs/cloud/aws/add-rule.png')} />


Now use this security group to allow SSH in all instances.





  

