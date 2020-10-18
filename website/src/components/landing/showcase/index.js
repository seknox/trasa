import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../muiTheme';
import RemoteAccess from './secure-remote-access';
import Integrations from './integrations';
import UACL from './unified-access-control';
import GeoSpanned from './geo-spanned';

const useStyles = makeStyles((theme) => ({
  root: {
    marginTop: 100,
    padding: theme.spacing(2),
    // backgroundColor: 'white',
  },
}));

export default function Features() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid
        container
        spacing={8}
        direction="row"
        aligntItems="center"
        justify="center"
        className={classes.root}
      >
        <Grid item xs={12} sm={12} md={12}>
          <RemoteAccess />
          <br /> <br />
          <br /> <br />
        </Grid>
        <Grid item xs={12} sm={12} md={12}>
          <GeoSpanned />
          <br /> <br />
          <br /> <br />
        </Grid>

        <Grid item xs={12} sm={12} md={12}>
          <Integrations />
          <br /> <br />
          <br /> <br />
        </Grid>

        <Grid item xs={12} sm={12} md={12}>
          <UACL />
          <br /> <br />
          <br /> <br />
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
