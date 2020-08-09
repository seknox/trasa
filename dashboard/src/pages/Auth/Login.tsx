import { TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import Typography from '@material-ui/core/Typography';
import CardContent from '@material-ui/core/CardContent';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React, { ReactElement, useEffect, useState } from 'react';
import AzureAD from '../../assets/idp/azuread.png';
import GSuite from '../../assets/idp/gsuite.png';
import JumpCloud from '../../assets/idp/jumpcloud.png';
import Okta from '../../assets/idp/okta.png';
import OpenLdapIcon from '../../assets/idp/openldap.png';
import TrasaLogo from '../../assets/trasa-ni.svg';
import LinearProgress from '../../utils/Components/Progressbar';
import { FetchExternalIdps } from './api/auth';

// STYLES
const useStyles = makeStyles((theme) => ({
  card: {
    textAlign: 'center',
    minWidth: 450,
    minHeight: 450,
    // minWidth: 750,
    // minHeight: 200,
    backgroundColor: 'white', // '#d0d3d4', // //rgba(1,1,35,1)
  },
  padMiddle: {
    textAlign: 'center',
  },
  heading: {
    // color: '#0b1728ff',
    // fontSize: '20px',
    // fontFamily: 'Open Sans, Rajdhani',
  },

  fpBtn: {
    // backgroundColor: ' #e6eaea',
    // color: 'black',
    // minWidth: 17,
    // marginRight: '35%',
  },
  cprightText: {
    // color: 'black',
    // fontSize: '12px',
    // fontFamily: 'Open Sans, Rajdhani',
    // marginLeft: '10%'
  },
}));

type LoginProps = {
  autofillEmail: boolean;
  userData: UserData;
  intent: string;
  title: string;
  loader: boolean;
  showForgetPass: boolean;
  loginData: { email: string; password: string; orgID: string; intent: string; idpName: string };
  sendLoginRequest: (
    e: React.FormEvent<Element>,
    intent: string,
    idpName: string,
    orgID: string,
  ) => void;
  handleLoginDataChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
};

type UserData = {
  idpName: string;
  userName: string;
  email: string;
  orgID: string;
};

export default function LoginPage(props: LoginProps): ReactElement {
  const classes = useStyles();

  const {
    handleLoginDataChange,
    loginData,
    autofillEmail,
    title,
    sendLoginRequest,

    showForgetPass,
  } = props;

  return (
    <Card className={classes.card}>
      <CardContent>
        <div className={classes.padMiddle}>
          <img src={TrasaLogo} height={100} width={200} alt="trasa logo" />
          <Typography variant="h3"> {title} </Typography>
        </div>
        <br />
        <form
          onSubmit={(e) => sendLoginRequest(e, '', 'trasa', '')}
          noValidate
          autoComplete="off"
          name="loginform"
        >
          {autofillEmail ? null : (
            <TextField
              fullWidth
              label="  Email or Username"
              // defaultValue={loginData.email}
              onChange={handleLoginDataChange}
              autoFocus
              id="email"
              name="email"
              value={loginData.email}
              variant="outlined"
              size="small"
            />
          )}

          <br />
          <br />

          <TextField
            fullWidth
            label="  Password"
            // defaultValue={loginData.password}
            onChange={handleLoginDataChange}
            id="password"
            name="password"
            type="password"
            value={loginData.password}
            variant="outlined"
            size="small"
          />

          <br />
          <br />

          <Grid container spacing={2} direction="row" justify="center" alignItems="center">
            <Grid item xs={4}>
              <Button
                variant="contained"
                name="submit"
                id="submit"
                type="submit"
                onClick={(e) => sendLoginRequest(e, '', 'trasa', '')}
              >
                {' '}
                Sign In{' '}
              </Button>
            </Grid>
            {showForgetPass === false ? null : (
              <Grid item xs={8}>
                <Button
                  variant="outlined"
                  onClick={(e) => sendLoginRequest(e, 'AUTH_REQ_FORGOT_PASS', 'trasa', '')}
                >
                  Forgot Password
                </Button>
              </Grid>
            )}
            {props.loader && <LinearProgress />}
          </Grid>
        </form>

        <br />

        <Divider light />
      </CardContent>

      <CardContent>
        <div className={classes.heading}> External Identity Providers</div>
        <ExternalIdps />
        <br />
        <Divider light />
      </CardContent>

      <div className={classes.padMiddle}>
        <div className={classes.cprightText}> Trasa Dashboard v20.6.1 by Seknox </div>{' '}
      </div>

      {/* <OrgSelect orgs={orgs} submitLoginRequest={sendLoginRequest} /> */}
    </Card>
  );
}

type Idp = {
  idpName: string;
  endpoint: string;
};

function ExternalIdps() {
  const [idps, setIdps] = useState([]);
  useEffect(() => {
    FetchExternalIdps(setIdps);
  }, []);

  function onClick(val: Idp) {
    //  //console.log(val)
    switch (val.idpName) {
      case 'okta':
        window.location.replace(val.endpoint);
        break;
      default:
        // props.ldapAuth('freeipa', '')
        break;
    }
  }

  function returnIdp(v: Idp, k: number): unknown {
    //  //console.log(val)
    switch (v.idpName) {
      case 'okta':
        return (
          <Button onClick={() => onClick(v)} key={k}>
            <img src={Okta} alt={v.idpName} height={50} />
          </Button>
        );
      case 'freeipa':
        return '';
      case 'gsuite':
        return (
          <Button onClick={() => onClick(v)} key={k}>
            <img src={GSuite} alt={v.idpName} height={50} />
          </Button>
        );

      case 'jumpcloud':
        return (
          <Button onClick={() => onClick(v)} key={k}>
            <img src={JumpCloud} alt={v.idpName} height={50} />
          </Button>
        );

      case 'azuread':
        return (
          <Button onClick={() => onClick(v)} key={k}>
            <img src={AzureAD} alt={v.idpName} height={50} />
          </Button>
        );
      default:
        return (
          <Button onClick={() => onClick(v)} key={k}>
            <img src={OpenLdapIcon} alt={v.idpName} height={50} />
          </Button>
        );
    }
  }

  return <div>{idps.map((v, k) => returnIdp(v, k))}</div>;
}
