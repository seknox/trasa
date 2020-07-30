import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
// import Layout from '../../components/Layout/Layout';
import Constants from '../../Constants';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import TrasaGWConsole from '../TrasaGWConsole/classbasedGW';
import FileBrowser from '../TrasaGWConsole/FileBrowser';
import TrasaSSHConsole from '../TrasaGWConsole/TrasaSSHConsole';
import MyAccount from './MyAccount';
// import MyAppSetting from './myappsetting';
import Mydevices from './MyDevices';
// import EventView from './myevents';
import MyEventTable from './MyLoginEvents';
import MyAppsList from './MyServices';
import RegisterDeviceAgent from './RegisterDeviceAgent';

function Mypage() {
  const [userData, setUserData] = useState({});
  const [orgData, setOrgData] = useState({});

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my`)
      .then((response) => {
        setUserData(response.data.data[0].User);
        setOrgData(response.data.data[0].Org);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'My Route', route: '/my' }]}
          tabHeaders={['Services', 'Account', 'Device', 'Access Stats']}
          Components={[
            <MyAppsList />,
            <MyAccount />,
            <Mydevices userData={userData} orgData={orgData} />,
            <MyEventTable />,
          ]}
        />
      </Layout>
    </div>
  );
}
// user this function to pass userDate to TrasaGWConsole
const MyView = () => {
  return (
    <Switch>
      <Route exact path="/my" component={Mypage} />
      <Route exact path="/my/registeragent/:deviceAgentToken" component={RegisterDeviceAgent} />

      <Route exact path="/my/service/connectrdp" component={TrasaGWConsole} />
      <Route exact path="/my/service/connectssh" component={TrasaSSHConsole} />


      <Route
        path="/my/apps/connectdynamic/rdp/:username/:email/:hostname/:rdpProto/"
        component={TrasaGWConsole}
      />
      <Route
        path="/my/apps/connectdynamic/ssh/:username/:email/:hostname"
        component={TrasaSSHConsole}
      />

      <Route path="/my/file/browse" component={FileBrowser} />
    </Switch>
  );
};

export default withRouter(MyView);
