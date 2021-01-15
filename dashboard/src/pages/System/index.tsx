import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';

import InAppTrails from './InappTrail';
// import PasswordSettings from './Security/PasswordPolicy'
import Setting from './Settings';
import Bnr from './Bnr';

function SystemView() {
  return (
    <Layout>
      <Headers
        pageName={[{ name: 'System', route: '/system' }]}
        tabHeaders={['Audit Trails',  'Backup', 'Settings']}
        Components={[<InAppTrails />,<Bnr />,<Setting />]}
      />
    </Layout>
  );
}

const SystemViewMain = () => {
  return (
    <Switch>
      <Route exact path="/system" component={SystemView} />
    </Switch>
  );
};

export default withRouter(SystemViewMain);
