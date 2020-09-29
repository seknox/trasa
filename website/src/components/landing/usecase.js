import { Typography, Grid, Paper } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import ThemeBase from '../muiTheme';

const useStyles = makeStyles((theme) => ({
  ctaPad: {
    marginTop: 150,
    marginBottom: 150,
    textAlign: 'center',
  },
  ctaPaper: {
    padding: theme.spacing(2),
    textAlign: 'center',
  },
}));

export default function SecondaryCTA() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <div className={classes.ctaPad}>
        <Grid container spacing={2} direction="row" justify="center" alignItems="center">
          <Grid item xs={4}>
            <Paper className={classes.ctaPaper}>
              <Typography variant="h2">Enterprise Grade</Typography>
              <Typography variant="body1" style={{ textAlign: 'center' }}>
                - Used highly regulated Banks and Fintech - Enterprise Support - Modern Security
                Standards
              </Typography>
            </Paper>
          </Grid>
          <Grid item xs={4}>
            <Paper className={classes.ctaPaper}>
              <Typography variant="h2">Unified Access </Typography>
              <Typography variant="body1" style={{ textAlign: 'center' }}>
                - Used highly regulated Banks and Fintech - Enterprise Support - Modern Security
                Standards
              </Typography>
            </Paper>
          </Grid>
          <Grid item xs={4}>
            <Paper className={classes.ctaPaper}>
              <Typography variant="h2">Flexible Integration</Typography>
              <Typography variant="body1" style={{ textAlign: 'center' }}>
                - Used highly regulated Banks and Fintech - Enterprise Support - Modern Security
                Standards
              </Typography>
            </Paper>
          </Grid>
        </Grid>
      </div>
    </ThemeBase>
  );
}
