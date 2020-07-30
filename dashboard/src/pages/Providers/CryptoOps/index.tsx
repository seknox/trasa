// import Layout from '../../components/Layout/Layout';
import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
import CA from './CertificateAuthority';
import VaultConfig from './vault';

function CryptoOpsI() {
  return (
    <div>
      <Layout>
        <Headers
          pageName={[
            { name: 'Manage', route: '/providers' },
            { name: 'CryptoOps', route: '/providers/crypto-ops' },
          ]}
          tabHeaders={['Certificate Authority', 'Secret Storage']}
          Components={[<CA />, <VaultConfig />]}
        />
      </Layout>
    </div>
  );
}

const CryptoOps = (props: any) => {
  return (
    <Switch>
      <Route exact path="/providers/crypto-ops" component={CryptoOpsI} />
      {/* <Route path='/crypto/vault' component={Myservicesetting} /> */}
    </Switch>
  );
};

export default withRouter(CryptoOps);
