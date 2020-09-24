import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../muiTheme';

const useStyles = makeStyles((theme) => ({
  background: {
    // background: '#f5f6ff',
  },
  ctaPad: {
    marginTop: 50,
  },
  ctaTxt: {
    fontSize: '24px',
    color: 'white',
    fontFamily: 'Open Sans, Rajdhani',
    // paddingLeft: '40%',
    padding: theme.spacing(2),
    background: 'linear-gradient(to left, #1a2980, #26d0ce)',

    // l2r:  linear-gradient(to left, #1a2980, #26d0ce)
  },
}));

export default function Features() {
  const imgUrl = useBaseUrl('features/iamm.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.ctaTxt}>
              Privilege Account Management
              {/* our unified security features helps you manage security with less effort  */}
            </div>
            <br />
            <Typography variant="body1">
              Control and monitor privileged accounts, access, and activities in internal servers
              and services. TRASA enables deny-by-default security, allowing you to assign
              privileged access to only required individuals or groups carefully. TRASA integrates
              with existing identity service within the organization(active directory, LDAP
              services, cloud identity services) to automate the privileged account enrolment
              process.
              {/* Manage users and servers identity. TRASA manages complete lifecycle of user management enabling centralized identity throught IT space.
           TRASA can be used as standalone identity server or integrate with existing identity service within organization(active directory, ldap services, cloud identity services) */}
            </Typography>
          </Paper>
        </Grid>
        <Grid xs={12} sm={12} md={6}>
          <img src={imgUrl} alt="identity management" />
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
