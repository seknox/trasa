import { Typography, Paper } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import ThemeBase from '../../../muiTheme';

// import MFA from '../../../static/features/2faicon.svg';
// import Audit from '../../../static/features/auditicon.svg';
// import CredVault from '../../../static/features/vaulticon.svg';
// import PAM from '../../../static/features/pamicon.svg';
// import PDacl from '../../../static/features/pdaclicon.svg';
// import SSO from '../../../static/features/ssoicon.svg';
// import IM from '../../../static/features/usericon.svg';

const useStyles = makeStyles((theme) => ({
  ctaPad: {
    marginTop: 140,
    textAlign: 'center',
  },
  contained: {
    color: 'white',
    backgroundColor: '#000080',
    fontWeight: 400,
    //  fontSize: '14px',
    boxShadow: 'none',
    '&:hover, &:focus': {
      color: 'white',
    },
  },
  paper: {
    padding: theme.spacing(2),

    // borderColor: '#FFFFFF',

    // backgroundColor: 'transparent',
    // color: 'white',
    textAlign: 'center',
  },
  featuresList: {
    color: 'black',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  icons: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
  },
}));

export default function Enterprise() {
  const classes = useStyles();
  const MFA = useBaseUrl('features/2faicon.svg');
  const Audit = useBaseUrl('features/auditicon.svg');
  const CredVault = useBaseUrl('features/vaulticon.svg');
  const PAM = useBaseUrl('features/pamicon.svg');
  const PDacl = useBaseUrl('features/pdaclicon.svg');
  const SSO = useBaseUrl('features/ssoicon.svg');
  const IM = useBaseUrl('features/usericon.svg');
  return (
    <ThemeBase>
      <Grid container spacing={2} direction="row" justify="center" alignItems="center">
        {/* <Grid item xs={12}>
          <div className={classes.ctaPad}>
            <Typography variant="h1">Features</Typography>
          </div>
        </Grid> */}

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={PAM} alt="Privileged Access Management" className={classes.icons} />
              <br />
              <br />
              Privileged Account Management
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={SSO} alt="Single Sign On" className={classes.icons} />
              <br />
              <br />
              Single Sign On for HTTPs, RDP, SSH
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={PDacl} alt="Policy Defind Access Control" className={classes.icons} />
              <br />
              <br />
              Policy Defined Access Control
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={CredVault} alt="Credential Vault" className={classes.icons} />
              <br />
              <br />
              Password and Keys Vault
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={Audit} alt="Session Audit" className={classes.icons} />
              <br />
              <br />
              Session Audit
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6} sm={4}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={MFA} alt="Dual Factor Authentication" className={classes.icons} />
              <br />
              <br />
              Two Factor Authentication
            </div>
          </Paper>
        </Grid>

        {/* <Grid item xs={12} style={{ textAlign: 'center' }}>
          <Link
            className={clsx('button  button--lg', classes.contained)}
            to={useBaseUrl('features/')}
          >
            Learn more about features
          </Link>
        </Grid> */}
      </Grid>
    </ThemeBase>
  );
}
