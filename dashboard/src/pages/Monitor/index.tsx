import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import LiveSessions from './Session/LiveSessions';
import LogTable from './Session/LogTableV2';
import RecordedSession from './Session/RecordedSession';
import TrasaSSHConsole from "../TrasaGWConsole/TrasaSSHConsole";
 import ClassbasedGW from '../TrasaGWConsole/classbasedGW';

function MonitorView() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Monitor Overview', route: '/monitor' }]}
          tabHeaders={['Live Sessions', 'History']}
          Components={[<LiveSessions />, <LogTable entityType="org" />]}
        />
      </Layout>
    </div>
  );
}

function MonitorMain() {
  return (
    <Switch>
      <Route exact path="/monitor/sessions" component={MonitorView} />
      {/* <Route exact path='/monitor/compliance' component={ComplianceView} />   */}
      {/* <Route path="/monitor/sessions" component={SessionView} /> */}
      {/* <Route exact path="/monitor/sessions/history" component={LogTable} /> */}
      <Route
        exact
        path="/monitor/sessions/view"
        component={RecordedSession}
      />

      <Route
      exact
        path="/monitor/sessions/joinxterm"
        component={TrasaSSHConsole}
      />

    <Route
    exact
        path="/monitor/sessions/join"
        component={ClassbasedGW}
      />

  
    </Switch>
  );
}

export default withRouter(MonitorMain);
