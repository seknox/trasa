import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import QueryString from 'query-string';
import LoginBox from './index';
// import Particles from '../../utils/particle';

// STYLES
const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    minHeight: '100vh',
    background: 'rgba(1,1,35,1)',
    // width: '100%',
    // Height: '100%',
    // minHeight: 800,
    // margin: '0',
  },
  login: {
    boxShadow: '0 0 15px 0 #12c2e9', // #6DD5FA  , rgba(0,0,0,0.12) , #302b63
    transition: '0.3s',
    position: 'absolute',
    top: '45%',
    right: '50%',
    transform: 'translate(50%,-50%)',
  },
  tfa: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
    transition: '0.3s',
    position: 'absolute',
    top: '45%',
    right: '50%',
    transform: 'translate(50%,-50%)',
  },
  particle: {
    height: '100%',
    minHeight: 900,
    background: 'rgba(1,1,35,1)',
  },
}));

export default function MainLogin(props: any) {
  const [intent, setIntent] = React.useState<
    | 'AUTH_REQ_DASH_LOGIN'
    | 'AUTH_REQ_CHANGE_PASS'
    | 'AUTH_REQ_ENROL_DEVICE'
    | 'AUTH_HTTP_ACCESS_PROXY'
  >('AUTH_REQ_DASH_LOGIN');

  const [proxyDomain, setProxyDomain] = React.useState<string | string[] | null | undefined>('');

  React.useEffect(() => {
    const hashed = QueryString.parse(props.location.hash);
    if (hashed.httphost) {
      setIntent('AUTH_HTTP_ACCESS_PROXY');
      setProxyDomain(hashed.httphost);
      console.log('ProxyDomain: ', hashed.httphost);
    }
  }, []);
  const userData = {
    email: '',
    password: '',
    orgID: '',
    intent: '',
    idpName: 'freeipa',
    userName: '',
  };
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {/* <Particles /> */}
      <div className={classes.login}>
        <LoginBox
          autofillEmail={false}
          intent={intent}
          title="Dashboard Login"
          showForgetPass
          userData={userData}
          setData={() => 0}
          proxyDomain={proxyDomain}
          changeHasAuthenticated={() => void 0}
        />
      </div>
    </div>
  );
}
