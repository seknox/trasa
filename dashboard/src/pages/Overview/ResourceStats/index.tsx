import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles, Theme } from '@material-ui/core/styles';
import React from 'react';
import Servicestats from './ServicesStats';
import BackupStats from './BackupStats';
import CAStats from './CAStats';
import DeviceStats from './DeviceStats';
import PolicyStats from './PolicyStats';
import UserStats from './UserStats';
import VaultStats from './VaultStats';

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
    // background: '#030417',
  },

  paperSmaller: {
    maxHeight: 200,
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paper: {
    // minHeight: 400,
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
    boxShadow: '0 0 5px 0 rgba(0,0,0,0.1)',
  },
  paperNopadd: {
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    minWidth: 400,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
}));

export default function ResourseStatsIndex() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={3}>
        {/* //TODO fix grid size @sshah */}
        <Grid item xs={12} sm={6} md={2}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <PolicyStats />{' '}
          </Paper>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <CAStats />{' '}
          </Paper>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <VaultStats />{' '}
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <BackupStats />{' '}
          </Paper>
        </Grid>

  

        <Grid item xs={12} sm={12} md={6}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <UserStats entityType="org" entityID="org" />{' '}
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={6}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <Servicestats entityType="org" entityID="org" />{' '}
          </Paper>
        </Grid>

        <Grid item xs={12} sm={6} md={12}>
          <Paper className={classes.paper} elevation={3}>
            {' '}
            <DeviceStats entityType="org" entityID="org" />{' '}
          </Paper>
        </Grid>

        {/* <Grid item xs={12} sm={6} md={3}>
          <Paper className={classes.paper}>  HA setup</Paper>
        </Grid> */}
      </Grid>
    </div>
  );
}
