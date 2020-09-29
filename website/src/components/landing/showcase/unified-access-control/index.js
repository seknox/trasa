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
  },
  ctaPad: {
    marginTop: 50,
  },
  ctaTxt: {
    fontSize: '24px',
    // color: 'white',
    fontFamily: 'Open Sans, Rajdhani',
    // paddingLeft: '40%',
    // padding: theme.spacing(2),
    // background: 'linear-gradient(to left, #1a2980, #26d0ce)',
  },
  image: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
  },
}));

export default function Unified() {
  const imgUrl = useBaseUrl('dash/access-stats.png');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={12} md={5}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h2"> Unified Access Control Platform</Typography>
            <br />
            <Typography variant="body1">
              A unified security platform that can address every access control security
              requirements for your team. TRASA enbles best practice security by default.
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} sm={12} md={7} className={classes.image}>
          <AllFeatures />
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
