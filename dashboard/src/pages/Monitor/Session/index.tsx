import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
import TrasaSSHConsole from '../../TrasaGWConsole/TrasaSSHConsole';
import LiveSessions from './LiveSessions';
import LogTable from './LogTableV2';
import HTTPSession from './RecordedSession';
import classbasedGW from '../../TrasaGWConsole/classbasedGW';

// import Layout from '../../components/Layout/Layout';

function SessionPage() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[
            { name: 'Monitor', route: '/monitor' },
            { name: 'Sessions', route: '/monitor/sessions' },
          ]}
          tabHeaders={['Live Sessions', 'History']}
          Components={[<LiveSessions />, <LogTable entityType="org" entityID="org" />]}
        />
      </Layout>
    </div>
  );
}

const SessionPageHOC = () => {
  return (
    <Switch>
      <Route exact path="/monitor/sessions" component={SessionPage} />
      <Route exact path="/monitor/sessions/history" component={LogTable} />
      <Route
        path="/monitor/sessions/history/:type/:year/:month/:day/:sessionID"
        component={HTTPSession}
      />
      <Route
        path="/monitor/sessions/history/:type/:year/:month/:day/:sessionID"
        component={HTTPSession}
      />
      <Route
        path="/monitor/sessions/history/:type/:year/:month/:day/:sessionID"
        component={HTTPSession}
      />
       <Route
        path="/monitor/sessions/join/:serviceID/:username/:email/:connID/:hostname"
        component={classbasedGW}
      />
       <Route
        path="/monitor/sessions/joinxterm/:serviceID/:username/:email/:connID/:hostname"
        component={TrasaSSHConsole}
      />
      {/* <Route path='/rc/users/:username' component={UserPage} />            */}
    </Switch>
  );
};

export default withRouter(SessionPageHOC);
