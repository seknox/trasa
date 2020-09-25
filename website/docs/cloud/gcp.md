---
id: google-cloud
title: Setup GCP Through TRASA
sidebar_label: GCP
---
import useBaseUrl from '@docusaurus/useBaseUrl';


## Create a firewall rule for TRASA

* Go to VPC Network-> Firewall on Main menu 
<img  alt="firewall-menu" src={useBaseUrl('img/docs/cloud/gcp/firewall-menu.png')} />

* Click the "Create Firewall Rule" button
<img  alt="create-firewall-rule-btn" src={useBaseUrl('img/docs/cloud/gcp/create-firewall-rule-btn.png')} />

* Fill in the names and description
* Give a target tag name
* Enter TRASA IP as source IP
* Specify tcp 22 port
<img  alt="new-firewall-rule" src={useBaseUrl('img/docs/cloud/gcp/new-firewall-rule.png')} />

Use this firewall rule to give ssh access to all instances.


## Configuring SSH keys in Google Cloud
By default google cloud uses OS Login which uses google identity to manage SSH keys.
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


