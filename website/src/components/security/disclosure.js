import { Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import ThemeBase from '../muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
  },
}));

export default function MainCta() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={2} className={classes.background}>
        <Grid item xs={12} sm={12}>
          <div className={classes.ctaPad}>
            <Typography variant="h1">Responsible disclosure</Typography>
          </div>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
