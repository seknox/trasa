import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';

import Policies from './policies';
import SecurityRules from './SecurityRules';
import AllRequestsHOC from './AdhocAccess';


function ControlView() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Control', route: '/control' }]}
          tabHeaders={['AdHoc Requests', 'Access Policies', 'Security Rules']}
          Components={[<AllRequestsHOC />,  <Policies />, <SecurityRules />]}
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
      <Route path="/control/access-requests" component={AllRequestsHOC} />
    </Switch>
  );
};

export default withRouter(Control);
