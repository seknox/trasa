import useBaseUrl from '@docusaurus/useBaseUrl';
import { Paper, Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import ThemeBase from '../../../muiTheme';

const useStyles = makeStyles((theme) => ({
  ctaPad: {
    marginTop: 140,
    textAlign: 'center',
  },
  contained: {
    color: 'white',
    backgroundColor: '#000080',
    fontWeight: 400,
    //  fontSize: '14px',
    boxShadow: 'none',
    '&:hover, &:focus': {
      color: 'white',
    },
  },
  paper: {
    backgroundColor: 'transparent',
    padding: theme.spacing(2),

    borderColor: '#FFFFFF',
    //   boxShadow: '0 3px 5px 2px white',
    color: 'white',
    textAlign: 'center',
  },
  featuresList: {
    color: 'black',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function Enterprise() {
  const classes = useStyles();
  const IDPs = useBaseUrl('brands/identity-providers.svg');
  const SDPs = useBaseUrl('brands/service-providers.svg');

  return (
    <ThemeBase>
      <Grid container spacing={2} direction="row" justify="center" alignItems="center">
        <Grid item xs={12}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h4">User Identity Providers</Typography>
            <div className={classes.featuresList}>
              <img src={IDPs} alt="Identity Providers" />
            </div>
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Paper className={classes.paper} elevation={0}>
            <Typography variant="h4">Service Identity Providers</Typography>
            <div className={classes.featuresList}>
              <img src={SDPs} alt="Service Identity Providers" />
            </div>
          </Paper>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
