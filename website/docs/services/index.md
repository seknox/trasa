---
id: introduction
title: Services
sidebar_label: Overview
---


Services are profiles of web, SSH, RDP, database and are the fundamental component of TRASA. 



In fact, The sole purpose of using TRASA is to protect these services from unauthorized access. Configuration profiles of these services are used by TRASA to authenticate and authorize users' access. TRASA access proxy uses these profiles to determine how to forward incoming traffic to back-end services.



## Creating a new Service
You need to create a profile of SSH,RDP of HTTP service to start protecting them with TRASA.

* Go to Services page and click "Create new service" button.
<img alt="create-a-service" src={('/img/docs/quickstart/create-new-service.png')} />  
* Give it a friendly name.
* Choose service type.
* Enter the service hostname.
* Click Submit.

You will be redirected to the newly created service page.

![service profile](./service-profile.png 'Example of Web Service Profile')
