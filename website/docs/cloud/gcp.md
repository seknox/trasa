---
id: google-cloud
title: Setup GCP Through TRASA
sidebar_label: GCP
---
import useBaseUrl from '@docusaurus/useBaseUrl';


## Create a firewall rule for TRASA

To use TRASA properly, we need to allow ssh access from the TRASA server only. To do that, we are going to create firewall rules in GCP.

We will make two rules, one to block all remote access requests (`block-all-remote`) and one to allow requests from TRASA (`allow-from-trasa`).

If two rules have the same priority value, the blocking rule will override, so we will give `allow-from-trasa` more priority than `block-all-remote`

* Go to VPC Network-> Firewall on Main menu 
<img  alt="firewall-menu" src={useBaseUrl('img/docs/cloud/gcp/firewall-menu.png')} />

* Click the "Create Firewall Rule" button
<img  alt="create-firewall-rule-btn" src={useBaseUrl('img/docs/cloud/gcp/create-firewall-rule-btn.png')} />

* Fill in the following parameters
    + "Target tag name" : `allow-from-trasa`
    + "Action on match" : "Allow"
    + "Source IP range" : [TRASA_IP]
    + "Ports and protocols" : "tcp:22, tcp:3389"
    + "Priority" : 999

<img  alt="new-firewall-rule" src={useBaseUrl('img/docs/cloud/gcp/new-firewall-rule.png')} />

* Create another rule with the following parameters: 
    + "Target tag name" : `allow-from-trasa`
    + "Action on match" : "Deny"
    + "IP source" : "0.0.0.0/0"
    + "Ports and protocols" : "tcp:22, tcp:3389"
    + "Priority" : 1000
    
<img  alt="blocking-firewall-rule" src={useBaseUrl('img/docs/cloud/gcp/blocking-firewall-rule.png')} />

:::tip
*You can give different priority and target tag name. 
Just make sure `allow-from-trasa` has more priority than `block-all-remote`.
:::

:::tip
Learn more about priority [here](https://cloud.google.com/vpc/docs/firewalls#priority_order_for_firewall_rules) 
:::

Use these two firewall rules for all instances.


## Configuring SSH keys in Google Cloud
By default, google cloud uses OS Login, which uses google identity to manage SSH keys.
To use TRASA to manage your SSH keys, you need to disable OS Login.
Then you need to add ssh keys to the instance or project.


* Go to [google cloud compute instances page](https://console.cloud.google.com/compute/instances) and click on the instance you want to configure. 
* Click the "Edit" button
<img  alt="edit-instance-btn" src={useBaseUrl('img/docs/cloud/gcp/edit-instance-btn.png')} />

* Generate a new ssh key
`ssh-keygen -t rsa -b 4096  -f ~/.ssh/[KEY_FILENAME] -C [USERNAME]`
* Scroll down to the "custom metadata" section and add a new key `enable-oslogin`:`FALSE`
* Click the "add item" button under SSH Keys section
* Copy the contents of [KEY_FILENAME].pub into the field
<img  alt="instance-level-metadata" src={useBaseUrl('img/docs/cloud/gcp/instance-level-metadata.png')} />
* Click Save
* [Save the contents of [KEY_FILENAME] in TRASA vault](../providers/secret-vault/index.md#storing-service-credentials)

:::tip
If you want to configure this for all instances of a project, go to [Metadata](https://console.cloud.google.com/compute/metadata) menu in Compute Engiene page
 <img  alt="project-level-metadata" src={useBaseUrl('img/docs/cloud/gcp/project-level-metadata.png')} />

:::


