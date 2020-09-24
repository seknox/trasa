---
id: quickstart-cloud-ssh
title: Setup SSH Through TRASA
sidebar_label: Setup SSH
---
import useBaseUrl from '@docusaurus/useBaseUrl';



## Prerequisite

+ User profile in TRASA
+ Service profile in TRASA
+ Vault initialised and decrypted (Optional)
+ SSH CA initialised (Optional)


:::tip
Most cloud providers provide you ssh private key instead of password by default.
So it is recommended to initialise vault to store private-key
:::


## Set security rules

To use TRASA as access proxy, you must setup firewall/ network security rules to allow ssh from TRASA server.



import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';






<Tabs
    defaultValue="openssh"
    values={[
  
            {label: 'AWS', value: 'aws'},
            {label: 'GCP', value: 'gcp'},
            {label: 'Digital Ocean', value: 'do'},
            {label: 'Azure', value: 'azure'},
        ]}
>


<TabItem value="do">


DO 

</TabItem>

<TabItem value="gcp">

GCP

</TabItem>


<TabItem value="aws">

In EC2 management console,
* Go to Security groups
<img  alt="main-menu" src={useBaseUrl('img/docs/main-menu.png')} />  


</TabItem>

<TabItem value="aws">
TODO
</TabItem>




</Tabs>

