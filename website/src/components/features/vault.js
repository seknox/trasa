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
  const imgUrl = useBaseUrl('../../../static/img/features/vault.svg');
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={6}>
        <Grid xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={0}>
            <div className={classes.ctaTxt}>
              Password and Secrets Vault
              {/* our unified security features helps you manage security with less effort  */}
            </div>
            <br />
            <Typography variant="body1">
              Storing access credentials in local files, online documents, excel sheets are one of
              the riskiest practice. TRASA vault is one of the safest storage to keep and manage
              secrets for personal or system/service accounts.
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
