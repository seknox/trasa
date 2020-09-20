import { Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import ThemeBase from '../muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 150,
    textAlign: 'center',
  },
}));

export default function SecondaryCTA() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={2} direction="row" justify="center" alignItems="center">
        <Grid item xs={8}>
          <div className={classes.ctaPad}>
            <Typography variant="h2">Secure infrastructure access</Typography>
            <Typography variant="body1" style={{ textAlign: 'center' }}>
              On-premise data center or dynamic cloud infrastructure, dedicated servers or ephemeral
              applications and services, internal team or managed service provider; managing secure
              access to servers or services hosted in these infrastructures is critical to your
              organization's overall security.
            </Typography>
          </div>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
