import React from 'react';
// import HttpCertSettings from './HttpCertSettings';
import SSLCerts from './SSLCerts';
import HostCerts from './HostCert';

export default function (props: any) {
  console.info(props.serviceID)
  switch (props.serviceType) {
    // case 'http':
    //   return <HttpCertSettings serviceID={props.serviceID} />;
    case 'db':
      return <SSLCerts serviceID={props.serviceID} />;
    case 'ssh':
      return <HostCerts serviceID={props.serviceID} />;
    default:
      return <div>Not Applicable for this type of Service</div>;
  }
}
