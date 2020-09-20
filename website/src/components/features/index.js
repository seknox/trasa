import { Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../muiTheme';
import IAM from './iam';
import Sso from './sso';
import Mfa from './mfa';
import Vault from './vault';
import Audit from './audit';

const useStyles = makeStyles(() => ({
  root: {
    marginTop: 100,
  },
}));

export default function Features() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={8} direction="column" className={classes.root}>
        <Grid item xs={12} sm={12} md={12}>
          <IAM />
        </Grid>
        <br /> <br />
        <Grid item xs={12} sm={12} md={12}>
          <Mfa />
        </Grid>
        <br /> <br />
        <Grid item xs={12} sm={12} md={12}>
          <Vault />
        </Grid>
        <br /> <br />
        <Grid item xs={12} sm={12} md={12}>
          <Sso />
        </Grid>
        <br /> <br />
        <Grid item xs={12} sm={12} md={12}>
          <Audit />
        </Grid>
        <br /> <br />
      </Grid>
    </ThemeBase>
  );
}
