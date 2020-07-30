// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import Constants from '../../../Constants';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
// import servicesettingPage from './servicesetting.js.txt';
import AccessMaps from '../AccessMap';
import AccessStats from '../AccessStats';
import ManageCreds from './ManageCredentials';
import ServiceOverview from './ServiceOverview';

function AppPageIndex(props: any) {
  const [serviceName, setserviceName] = useState('');
  const [appData, setAppData] = useState({});
  const [allUsers, setallUsers] = useState([]);
  // const [managedUsers, setManagedUsers] = useState([]);
  // const [proxyConfig, setproxyConfig] = useState({});

  const getserviceDetail = (serviceID: string) => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/services/${serviceID}`)
      .then((response) => {
        if (response.data.status === 'success') {
          // const managedUsers = response.data.data[0].managedAccounts.split(',');
          // response.data.App.managedAccounts = managedUsers

          // let accounts = []
          // const accounts = managedUsers.map(function (v: any) {
          //   return [v, '**********'];
          // });
          // accounts.shift();

          setAppData(response.data.data[0]);
          // setManagedUsers(accounts);
          setserviceName(response.data.data[0].serviceName);
          // setproxyConfig(response.data.data[0].proxyConfig);
        }
      })
      .catch((error) => {
        console.log(error);
      });

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/all`)
      .then((response) => {
        if (response.status === 403) {
          window.location.href = '/login';
        }
        setallUsers(response.data.data[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    getserviceDetail(props.ID);
  }, [props.ID]);

  const UpdateserviceName = (val: string) => {
    setserviceName(val);
  };

  const getRoute = (staticVal: any, dynamicVal: any) => {
    return staticVal + dynamicVal;
  };

  return (
    <Layout>
      <Headers
        pageName={[
          { name: 'Services', route: '/services/groups#All%20services' },
          {
            name: serviceName,
            route: getRoute('/services/service#serviceID=', props.ID),
          },
        ]}
        tabHeaders={['Overview', 'Access Stats', 'Access Maps', 'Manage Credentials']}
        Components={[
          <ServiceOverview
            entityType="service"
            entityID={props.ID}
            appData={appData}
            allUsers={allUsers}
            UpdateserviceName={UpdateserviceName}
            serviceName={serviceName}
          />,

          <AccessStats entityType="service" entityID={props.ID} />,

          <AccessMaps
            entityType="service"
            entityID={props.ID}
            appData={appData}
            allUsers={allUsers}
            UpdateserviceName={UpdateserviceName}
          />,

          <ManageCreds
            urlID={props.ID}
            appData={appData}
            allUsers={allUsers}
            UpdateserviceName={UpdateserviceName}
          />,
        ]}
      />
    </Layout>
  );
}

type TParams = { ID: string };

const AppPage = ({ match }: RouteComponentProps<TParams>) => {
  return (
    <div>
      <AppPageIndex ID={match.params.ID} />
    </div>
  );
};

export default AppPage;
