import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../../muiTheme';
import AllIntegrations from './all-integrations';

const useStyles = makeStyles((theme) => ({
  background: {
    // background: '#f5f6ff',
  },
  paper: {
    background: '#fafafa',
  },
  ctaPad: {
    marginTop: 50,
  },
  ctaTxt: {
    fontSize: '24px',
    // color: 'white',
    fontFamily: 'Open Sans, Rajdhani',
    // paddingLeft: '40%',
    padding: theme.spacing(2),
    // background: 'linear-gradient(to left, #1a2980, #26d0ce)',
  },
  image: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
  },
}));

export default function Features() {
  const imgUrl = useBaseUrl('features/tfa.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6} className={classes.image}>
          <AllIntegrations />
        </Grid>

        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h2"> Integrations with services you already use </Typography>
            <br />
            <Typography variant="body1">
              TRASA securely controls user's access to protected services. In most cases, user's
              identities are managed by identity providers like active directory, LDAP servers, and
              service identities are managed in cloud service providers or data center
              virtualization platforms. TRASA integrates with these service platforms (providers) to
              seamlessly integrate with existing workflows.
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
