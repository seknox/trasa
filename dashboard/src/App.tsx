import axios from 'axios';
// import mixpanel from 'mixpanel-browser';
import React, { Component, Suspense, lazy } from 'react';
import { BrowserRouter as Router, Redirect, Route, Switch } from 'react-router-dom';
import TrasaTheme from './Layout/TrasaTheme';
// import 'semantic-ui-css/semantic.min.css';
import Constants from './Constants';
import Snack from './utils/Components/Snackbar';
// import Manage from './pages/Providers';
// import Monitor from './pages/Monitor';
// import System from './pages/System';
// import Control from './pages/Control';
// import Users from './pages/Users';
// import Services from './pages/Services';

import ProgressBar from './utils/Components/Progressbar';
// var mixpanel = require('mixpanel-browser');
// mixpanel.init('5dcb2f677ad4963e7e6eae978d97b378');
// add this     "semantic-ui-css": "^2.3.1"

type AppProps = {};

type ActionStatus = {
  respStatus: boolean;
  success: boolean;
  loader: boolean;
  statusMsg: string;
};

type Tstate = {
  isLoggedIn: boolean;
  actionStatus: ActionStatus;
};

const LoginPage = lazy(() => import('./pages/Auth/Dashlogin'));
const Overview = lazy(() => import('./pages/Overview'));
const PublicView = lazy(() => import('./pages/Publicpages'));
const MyView = lazy(() => import('./pages/My'));
const Providers = lazy(() => import('./pages/Providers'));
const Control = lazy(() => import('./pages/Control'));
const Users = lazy(() => import('./pages/Users'));
const Services = lazy(() => import('./pages/Services'));
const Monitor = lazy(() => import('./pages/Monitor'));
const System = lazy(() => import('./pages/System'));

export default class APP extends Component<AppProps, Tstate> {
  constructor(props: AppProps) {
    super(props);
    this.state = {
      isLoggedIn: false,
      actionStatus: { respStatus: false, success: false, loader: false, statusMsg: '' },
    };

    axios.interceptors.response.use(
      (resp) => {
        if (resp.config.method !== 'get') {
          this.setState((prevState) => ({
            actionStatus: { ...prevState.actionStatus, loader: false },
          }));
          if (resp.data.status === 'success') {
            this.setState((prevState) => ({
              actionStatus: {
                ...prevState.actionStatus,
                success: true,
                respStatus: true,
                statusMsg: resp.data.reason,
              },
            }));
          } else {
            this.setState((prevState) => ({
              actionStatus: {
                ...prevState.actionStatus,
                respStatus: true,
                statusMsg: resp.data.reason,
                success: false,
              },
            }));
          }
        }
        return resp;
      },
      (error) => {
        if (error && error.response && error.response.status === 403) {
          this.setState((prevState) => ({
            actionStatus: {
              ...prevState.actionStatus,
              respStatus: true,
              statusMsg: 'something went wrong',
              success: false,
            },
          }));
          if (
            !window.location.href.includes('?next=') &&
            !window.location.href.includes('/login') &&
            !window.location.href.includes('/woa')
          ) {
            window.location.href = `/login?next=${window.location.href}`;
          }
        }
        return Promise.reject(error);
      },
    );

    axios.interceptors.request.use((config) => {
      config.headers['X-CSRF'] = localStorage.getItem('X-CSRF');
      this.setState((prevState) => ({
        actionStatus: { ...prevState.actionStatus, respStatus: false, statusMsg: '', loader: true },
      }));

      return config;
    });

      Constants.TRASA_HOSTNAME = `https://${window.location.hostname}`;
      Constants.TRASA_HOSTNAME_WEBSOCKET = `wss://${window.location.hostname}`;

  }

  snackClose = (reason: string) => {
    if (reason === 'clickaway') {
      return;
    }
    this.setState((prevState) => ({
      actionStatus: { ...prevState.actionStatus, respStatus: false },
    }));
  };

  render() {
    const { actionStatus, isLoggedIn } = this.state;

    return (
      <TrasaTheme>
        <Router>
          <Suspense fallback={<ProgressBar />}>
            <Switch>
              <Route path="/woa" component={PublicView} />
              <Route path="/my" component={MyView} />
              <Route path="/overview" component={Overview} />
              <Route path="/login" component={LoginPage} />
              <Route path="/control" component={Control} />
              <Route path="/providers" component={Providers} />
              <Route path="/monitor" component={Monitor} />
              <Route path="/system" component={System} />
              <Route path="/users" component={Users} />
              <Route path="/services" component={Services} />

              <Route
                exact
                path="/"
                render={() => (isLoggedIn ? <Redirect to="/overview" /> : <Redirect to="/login" />)}
              />
            </Switch>
          </Suspense>
        </Router>
        <Snack
          open={actionStatus.respStatus}
          statusMsg={actionStatus.statusMsg}
          success={actionStatus.success}
          snackClose={this.snackClose}
        />
        <br />
      </TrasaTheme>
    );
  }
}
