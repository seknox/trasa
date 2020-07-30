import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import AccessStats from './AccessStats';
import ResourceStats from './ResourceStats';
import SystemOverview from './SystemStats';

function MainPage() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Overview', route: '/overview' }]}
          tabHeaders={['Access Stats', 'Resource Stats', 'System Stats']}
          Components={[
            <AccessStats entityType="org" entityID="org" timeFilter="Today" statusFilter="All" />,
            <ResourceStats />,
            <SystemOverview />,
          ]}
        />
      </Layout>
    </div>
  );
}

const AppView = () => {
  return (
    <Switch>
      <Route exact path="/overview" component={MainPage} />
    </Switch>
  );
};

export default withRouter(AppView);
