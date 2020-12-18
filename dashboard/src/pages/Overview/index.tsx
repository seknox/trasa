import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import AccessStats from './AccessStats';
import ResourceStats from './ResourceStats';
import SystemOverview from './SystemStats';
import axios from 'axios';
import Constants from '../../Constants';
import { ShowTokens } from '../Providers/CryptoOps/vault'

function Overview() {
  const [open, setopen] = React.useState(false);
  const [decryptKeys, setDecryptKeys] = React.useState(['', ''])



  React.useEffect(() => {
    // if (window.Location === "overview")
    let path = window.location.pathname;
    if (path === '/overview'){
      const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/system/welcome-note`;

      axios
        .get(reqPath)
        .then((response) => {
          const resp = response.data.data[0];
        
         
          for (let i = 0; i < resp.length; i++){
            switch (resp[i].intent){
              case 'SHOW_VAULT_KEYS':
                if (resp[i].show){
                  setDecryptKeys(resp[i].data);
                  setopen(true);
                }
              
              default:
                return
            }
          }
         
  
        })
        .catch((error) => {
          console.error('catched ', error);
        });
    }
  }, [])

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
         <ShowTokens
                  open={open}
                  decryptKeys={decryptKeys}
                  handleClose={()=> {setopen(false)}}
                />
      </Layout>
    </div>
  );
}

const OverviewWithRouter = () => {
  return (
    <Switch>
      <Route exact path="/overview" component={Overview} />
    </Switch>
  );
};

export default withRouter(OverviewWithRouter);
