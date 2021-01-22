import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';

import AccessPolicies from './AccessPolicies';
import SecurityRules from './SecurityRules';
import AdhocAccess from './AdhocAccess';

function PolicyView() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Policies', route: '/policies' }]}
          tabHeaders={['Adhoc Access', 'Access Policies', 'Security Rules']}
          Components={[<AdhocAccess />, <AccessPolicies />, <SecurityRules />]}
        />
      </Layout>
    </div>
  );
}

const Policies = () => {
  return (
    <Switch>
      <Route exact path="/policies" component={PolicyView} />
    </Switch>
  );
};

export default withRouter(Policies);
