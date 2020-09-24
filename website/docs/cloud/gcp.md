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


## Getting SSH private key
If you are using google cloud, chances are you are using gcloud cli tool to access ssh.
gcloud stores your private key as `~/.ssh/google_compute_engine`. 

