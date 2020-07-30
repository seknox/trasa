import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import { ResourceStatsFilterProps } from '../../../../types/analytics';
import { HeaderFontSize, TitleFontSize } from '../../../../utils/Responsive';
import BrowserByType from './BrowserTypes';
import MobileType from './MobileTypes';
import WorkstationByType from './WorkstationType';

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
  paperTrans1: {
    backgroundColor: 'transparent',
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
  headerMain: {
    color: '#1b1b32',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeadersBig: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: HeaderFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function DeviceAggregatedData(props: ResourceStatsFilterProps) {
  const [alldevices, setAlldevices] = useState({
    totalUserdevices: 0,
    browserByType: [],
    mobileByType: [],
    workstationsByType: [],
    totalBrowsers: 0,
    totalWorkstations: 0,
    totalMobiles: 0,
  });

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/devices/${props.entityType}/${props.entityID}`;

    axios
      .get(reqPath)
      .then((r) => {
        setAlldevices(r.data.data[0]);
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
            <b> User Devices (endpoints) </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Total Devices </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {alldevices.totalUserdevices} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Browsers </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {alldevices.totalBrowsers} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Workstations </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {alldevices.totalWorkstations} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Mobile Devices </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {alldevices.totalMobiles} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Divider light />
        </Grid>

        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={4}>
              {' '}
              <MobileType mobilesByType={alldevices.mobileByType} />{' '}
            </Grid>
            <Grid item xs={4}>
              {' '}
              <BrowserByType browsersByType={alldevices.browserByType} />{' '}
            </Grid>
            <Grid item xs={4}>
              {' '}
              <WorkstationByType workstationsByType={alldevices.workstationsByType} />{' '}
            </Grid>
          </Grid>
        </Grid>

        {/* <Grid item xs={12}>
            <UserPie />
          </Grid> */}
      </Grid>
    </div>
  );
}

export default DeviceAggregatedData;
// export default withRouter(Overvieww);
