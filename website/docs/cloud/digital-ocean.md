---
id: digital-ocean
title: Setup Digital Ocean Through TRASA
sidebar_label: Digital Ocean
---

import useBaseUrl from '@docusaurus/useBaseUrl';


## Create a firewall rule for TRASA

* Go to Networking-> Firewalls on Main menu 
* Click the "Create Firewall" button
<img  alt="network-firewall-create" src={('/img/docs/cloud/do/network-firewall-create.png')} />


* Fill in the names and description
* Enter TRASA IP as source IP in inbound rules
<img  alt="inbound-rule" src={('/img/docs/cloud/do/inbound-rule.png')} />

* Enter a tag name to apply this rule.
* Click the "Create Firewall" button 
<img  alt="create" src={('/img/docs/cloud/do/create.png')} />

Use this firewall rule to give ssh access to all instances.



