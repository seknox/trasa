import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
import CloudServices from './CloudServices';

function AccountIndex() {
  return (
    <div>
      <Layout>
        {/* 'Identity Provider', <IdP />, */}
        <Headers
          pageName={[
            { name: 'Manage', route: '/providers' },
            { name: 'Accounts', route: '/providers/accounts' },
          ]}
          tabHeaders={['Cloud Connect']}
          Components={[<CloudServices />]}
        />
      </Layout>
    </div>
  );
}

const AccountsMain = () => {
  return (
    <Switch>
      <Route exact path="/providers/accounts" component={AccountIndex} />
    </Switch>
  );
};

export default withRouter(AccountsMain);
