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
  paper: {
    background: '#fafafa',
  },

  image: {
    // boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
  },
}));

export default function Features() {
  const imgUrl = useBaseUrl('dash/tfaapp.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6} className={classes.image}>
          <img src={imgUrl} alt="identity management" />
        </Grid>

        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h2"> Two Factor Authentication</Typography>
            <br />
            <Typography variant="body1">
              TRASA enables native agent-based (most-secure) and agentless (controlled by TRASA
              access-proxy) two-factor authentication in web, RDP, SSH, and Database services. TRASA
              supports TOTP, Push U2F and Yubikey as authentorization method.
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
