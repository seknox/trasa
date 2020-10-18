import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../../muiTheme';

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
  paper: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(2),

    borderColor: '#FFFFFF',
    //   boxShadow: '0 3px 5px 2px white',
    color: 'white',
    textAlign: 'center',
  },
  featuresList: {
    paddingTop: 50,
    color: 'black',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function Features() {
  const imgUrl = useBaseUrl('arch/geo-spanned.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h2"> Geo spanned access point.</Typography>
            <br />
            <Typography variant="body1">
              Distributed teams and distributed edge computing is the new norm for modern
              tech-driven teams. TRASA can be deployed as Geo spanned access points, protecting
              internal infrastructure access with the highest availability. Configure from a single
              dashboard, and policy will be synced automatically to every access point.
            </Typography>
          </Paper>
        </Grid>

        <Grid item xs={12} sm={12} md={6} className={classes.image}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.featuresList}>
              <img src={imgUrl} alt="Geo spanned TRASA" />
            </div>
          </Paper>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
