import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../muiTheme';

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
  const imgUrl = useBaseUrl('../../../static/img/features/tfa.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6}>
          <img src={imgUrl} alt="identity management" />
        </Grid>

        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.ctaTxt}>Two Factor Authentication</div>
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
