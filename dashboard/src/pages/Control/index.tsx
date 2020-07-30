import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';

import Policies from './policies';
import SecurityRules from './SecurityRules';

function ControlView() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Control', route: '/control' }]}
          tabHeaders={[ 'Access Policies', 'Security Rules']}
          Components={[ <Policies />, <SecurityRules />]}
        />
      </Layout>
    </div>
  );
}

const Control = () => {
  return (
    <Switch>
      <Route exact path="/control" component={ControlView} />
      <Route path="/control/access-policies" component={Policies} />
    </Switch>
  );
};

export default withRouter(Control);
