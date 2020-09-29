import { Grid, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
// import DashImage from '../../static/img/trasa-bluebg.svg';
import ThemeBase from '../../muiTheme';
import HubspotForm from 'react-hubspot-form';

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
        <Grid item xs={12}>
          <Typography variant="body1" component="p">
            TRASA is opensource project with Mozilla Public License (MPLv2). <br />
            For Enterprise deployments and support, ping us using form below.
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
