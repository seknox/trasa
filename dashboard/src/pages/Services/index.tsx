import React from 'react';
// import Layout from '../../components/Layout/Layout';
import { Route, Switch, withRouter } from 'react-router-dom';
// import Layout from '../../components/Layout/DashboardBase'
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import GlobalServieSetting from './GlobalServicesSetting';
import ServiceGroups from './Groups';
import GroupPage from './Groups/GroupPage';
import ServicePage from './Service';
import ServicesList from './ServiceList';

function AllServices() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Services', route: '/services' }]}
          tabHeaders={['All services', 'Service Groups (clusters)', 'Settings']}
          // 'TODO revamp service global setting'
          Components={[<ServicesList />, <ServiceGroups />, <GlobalServieSetting />]}
        />
      </Layout>
    </div>
  );
}

const ServiceView = () => {
  return (
    <Switch>
      <Route exact path="/services" component={AllServices} />
      <Route path="/services/groups/group/:ID" component={GroupPage} />
      <Route path="/services/service/:ID" component={ServicePage} />
    </Switch>
  );
};

export default withRouter(ServiceView);
