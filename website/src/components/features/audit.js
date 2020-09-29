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
  const imgUrl = useBaseUrl('features/audit2.png');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.ctaTxt}>Session Audit and Analytics</div>
            <br />
            <Typography variant="body1">
              TRASA offers session monitoring features (live or recorded), allowing to inspect
              malicious sessions for authorized user accounts. It's a powerful feature to track
              session history, dig for audit trails, and identify hostile sessions. TRASA can record
              secure shells(SSH), remote desktop(RDP), and https(web services).
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} sm={12} md={6}>
          <img src={imgUrl} alt="identity management" />
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
