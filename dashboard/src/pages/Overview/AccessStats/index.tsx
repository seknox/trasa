import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import AccessLocation from './AccessLocation';
import AggregatedStats from './AggStats';
import LoginFailedReasonAggs from './LoginFailedReasonAggs';
import LoginFailedAggrBar from './LoginHours';
import TodaysEvent from './TodaysEvent';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
    // background: '#030417',
  },

  paperSmaller: {
    maxHeight: 350,
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paper: {
    // [theme.breakpoints.down('lg')]: {
    //   maxHeight: 300,
    // },
    boxShadow: '0 0 5px 0 rgba(0,0,0,0.1)',
    //  minHeight: 370,
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paperLong: {
    [theme.breakpoints.down('lg')]: {
      maxHeight: 580,
    },
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paperNopadd: {
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
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
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    // height: '100%',
    // maxWidth: '100%',
    // maxHeight: '100%',
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
}));

export type AccessStatsFilterProps = {
  timeFilter: string;
  statusFilter: string;
  entityType: string;
  entityID: string;
};

export default function AccessStats(props: AccessStatsFilterProps) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={3}>
        <Grid item xs={12} sm={6}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={12} md={12} lg={7}>
              <Paper className={classes.paper} elevation={2}>
                {' '}
                <TodaysEvent
                  entityType={props.entityType || 'org'}
                  entityID={props.entityID}
                  timeFilter="Today"
                  statusFilter="All"
                />{' '}
              </Paper>
            </Grid>

            <Grid item xs={12} sm={12} md={12} lg={5}>
              <Paper className={classes.paper} elevation={2}>
                {' '}
                <AggregatedStats
                  entityType={props.entityType || 'org'}
                  entityID={props.entityID}
                  timeFilter="Today"
                  statusFilter="All"
                />{' '}
              </Paper>
            </Grid>

            <Grid item xs={12} sm={12} md={12} lg={7}>
              <Paper className={classes.paper} elevation={2}>
                {' '}
                <LoginFailedReasonAggs
                  entityType={props.entityType || 'org'}
                  entityID={props.entityID}
                  timeFilter="Today"
                  statusFilter="All"
                />{' '}
              </Paper>
            </Grid>

            <Grid item xs={12} sm={12} md={12} lg={5}>
              <Paper className={classes.paper} elevation={2}>
                {' '}
                <LoginFailedAggrBar
                  entityType={props.entityType || 'org'}
                  entityID={props.entityID}
                  timeFilter="Today"
                  statusFilter="All"
                />{' '}
              </Paper>
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={12} sm={6}>
          <Paper className={classes.paper} elevation={2}>
            {' '}
            <AccessLocation
              entityType={props.entityType || 'org'}
              entityID={props.entityID}
              timeFilter="Today"
              statusFilter="All"
            />{' '}
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}
