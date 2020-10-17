import { Grid, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import HubspotForm from 'react-hubspot-form';
import ThemeBase from '../../muiTheme';

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
        spacing={2}
        direction="column"
        justify="center"
        alignItems="center"
        className={classes.root}
      >
        <Grid item xs={12}>
          <Typography variant="h2" component="h2">
            {' '}
            Open source and Enterprise
          </Typography>
        </Grid>
        <Grid item xs={7}>
          <Typography variant="body1" component="p">
            TRASA is fully free and open source project distributed under Mozilla Public License
            (MPLv2). For enterprise deployments and support, send us your interest using form ahead.
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <HubspotForm
            portalId="5642830"
            formId="db368da0-2c0a-40c5-ac75-dc985479d415"
            onSubmit={() => console.log('Submit!')}
            onReady={(form) => console.log('Form ready!')}
            loading={<div>Loading...</div>}
          />
        </Grid>
        <br /> <br />
      </Grid>
    </ThemeBase>
  );
}
