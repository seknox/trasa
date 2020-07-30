import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import Accounts from './Accounts';
import CloudConnect from './Accounts/CloudServices';
import IdP from './Accounts/IdP';
import CryptoOps from './CryptoOps';
import CA from './CryptoOps/CertificateAuthority';
import Vault from './CryptoOps/vault';
// import Services from './Services';
// import Users from './Users';

function ManageView() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Providers', route: '/providers' }]}
          tabHeaders={[
            'User Identity Provider',
            'Service Identity Provider',
            'Certificate Authority',
            'Secret Storage',
          ]}
          Components={[<IdP />, <CloudConnect />, <CA />, <Vault />]}
        />
      </Layout>
    </div>
  );
}

function ManageIndex() {
  return (
    <Switch>
      <Route exact path="/providers" component={ManageView} />
      <Route path="/providers/accounts" component={Accounts} />
      <Route path="/providers/crypto-ops" component={CryptoOps} />
      {/* <Route path="/services" component={Services} />
      <Route path="/users" component={Users} /> */}
    </Switch>
  );
}

export default withRouter(ManageIndex);
