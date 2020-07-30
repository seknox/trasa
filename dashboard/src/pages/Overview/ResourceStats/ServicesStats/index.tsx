import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import { ResourceStatsFilterProps } from '../../../../types/analytics';
import { HeaderFontSize, TitleFontSize, TitleTextFontSize } from '../../../../utils/Responsive';
import ServiceType from './ServiceTypes';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paperTrans: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(1),
    textAlign: 'center',
  },

  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 400,
    minHeight: 500,
    // padding: theme.spacing(2),
    // textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  aggHeadersBig: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  headerMain: {
    color: '#1b1b32',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  titleText: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleTextFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: HeaderFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function AppAggregatedData(props: ResourceStatsFilterProps) {
  const [allservices, setAllservices] = useState({
    totalServices: 0,
    servicesByType: [],
    totalGroups: 0,
    dynamicService: false,
    sessionRecordingDisabledCount: 0,
  });

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/services/${props.entityType}/${props.entityID}`;

    axios
      .get(reqPath)
      .then((r) => {
        setAllservices(r.data.data[0]);
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('Error', error.message);
        }
      });
  }, [props.entityType, props.entityID]);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Services (services) </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Total </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allservices.totalServices} </b>
            </div>
          </Paper>
        </Grid>
        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Groups </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allservices.totalGroups} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Dyn. Enroll </b>
            </div>
            <div className={classes.titleText}>
              {' '}
              <b> {allservices.dynamicService ? 'Active' : 'Disabled'} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Recording Disabled </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allservices.sessionRecordingDisabledCount} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Imported </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              {/* TODO @bhrg3se change value here for imported Services count? */}
              <b> {allservices.sessionRecordingDisabledCount} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Divider light />
        </Grid>

        <Grid item xs={12}>
          {' '}
          <ServiceType ServiceByType={allservices.servicesByType} />{' '}
        </Grid>
      </Grid>
    </div>
  );
}

export default AppAggregatedData;
