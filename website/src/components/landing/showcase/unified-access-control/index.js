import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../../muiTheme';
import AllFeatures from './all-features';

const useStyles = makeStyles((theme) => ({
  background: {
    // background: '#f5f6ff',
  },
  paper: {
    background: '#fafafa',
    textAlign: 'center',
  },
  ctaPad: {
    marginTop: 50,
  },
  ctaTxt: {
    fontSize: '24px',
    // color: 'white',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function Unified() {
  const imgUrl = useBaseUrl('dash/access-stats.png');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h2"> Unified Access Control Platform</Typography>
            <br />
            <Typography variant="body1" style={{ textAlign: 'center' }}>
              A unified security platform that can address every access control security
              requirements for your team. TRASA enbles best practice security by default.
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12}>
          <AllFeatures />
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
