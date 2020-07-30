import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { EchartElementHeight, AGHeaderFontSize, AGTitleFontSize, Spacing } from '../../../utils/Responsive';
import Constants from '../../../Constants';

import { AccessStatsFilterProps } from '../../../types/analytics';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    height: EchartElementHeight(),
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
    color: '#000066',
    fontSize: AGTitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: AGHeaderFontSize(), // window.innerHeight < 750 ? '14px':'18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  selectLabel: {
    fontSize: 12,
    //  marginBotton: 5,
    padding: '1px',
    color: 'grey',
  },
  selectCustom: {
    fontSize: 12,
    fontFamily: 'Open Sans, Rajdhani',
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 25,
    marginTop: 5,
    // padding: '10px 100px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

function AppAggregatedData(props: AccessStatsFilterProps) {
  const [filter, setFilter] = useState('Today');
  const [successAndFailed, setSuccessAndFailed] = useState({
    successfulLogins: 0,
    failedLogins: 0,
    totalLogins: 0,
  });
  const [dashboardCount, setDashboardCount] = useState(0);
  const [sshCount, setSshCount] = useState(0);
  const [rdpCount, setRdpCount] = useState(0);
  const [dbCount, setDbCount] = useState(0);
  const [remoteAppCount, setRemoteAppCount] = useState(0);
  const [httpCount, setHttpCount] = useState(0);
  const [statusFilter, setStatusFilter] = useState('All');

  const time = ['Today', 'Yesterday', '7 Days', '30 Days', '90 Days', 'All'];
  const status = ['All', 'Success', 'Failed'];

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/loginsbytype/${props.entityType}/${props.entityID}/${filter}/${statusFilter}`;

    axios
      .get(reqPath)
      .then((r) => {
        if (r.data.status === 'success' && r.data.data) {
          r.data.data[0].forEach((a: any) => {
            switch (a.name) {
              case 'dashboard':
                setDashboardCount(a.value);
                return;

              case 'rdp':
                setRdpCount(a.value);
                return;
              case 'db':
                setDbCount(a.value);
                return;
              case 'http':
                setHttpCount(a.value);
                return;
              default:
                setSshCount(a.value);
                break;
            }
          });
          setRemoteAppCount(r.data.data[1]);
        }
        // let totalAuths = response.data.data[0].totalAuthEvents
        // this.setState({ totalAuthEvents: totalAuths, totalAppUsers: response.data.data[0].totalAppUsers })
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('Error', error.message);
        }
      });

    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/stats/total/${props.entityType}/${props.entityID}/${filter}`,
      )
      .then((r) => {
        if (r.data.status === 'success') {
          setSuccessAndFailed(r.data.data[0]);
        }
      });
  }, [filter, statusFilter, props.entityType, props.entityID]);

  function handleChange(event: React.ChangeEvent<HTMLSelectElement>) {
    setFilter(event.target.value);
  }

  // let kval: any = 'k'
  function kFormatter(num: number) {
    return Math.abs(num) > 999
      ? `${Math.sign(num) * +(Math.abs(num) / 1000).toFixed(1)} k`
      : Math.sign(num) * Math.abs(num);
  }

  function changeStatusFilter(event: React.ChangeEvent<HTMLSelectElement>) {
    setStatusFilter(event.target.value);
  }

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={Spacing()}>
        <Grid item xs={12}>
          {/* <Divider light />
            <br /> */}
          <div className={classes.headerMain}>
            {' '}
            <b> Active Sessions </b>{' '}
          </div>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> SSH </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(sshCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> RDP </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(rdpCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> R. App </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(remoteAppCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> HTTPs</b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(httpCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>
        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> DB </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(dbCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Dash. </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(dashboardCount)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Total </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b>
                {' '}
                {kFormatter(
                  successAndFailed.successfulLogins + successAndFailed.failedLogins,
                )}{' '}
              </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Success </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(successAndFailed.successfulLogins)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Failed </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {kFormatter(successAndFailed.failedLogins)} </b>{' '}
            </div>
          </Paper>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={6}>
            <select className={classes.selectCustom} name="time" onChange={(e) => handleChange(e)}>
              {time.map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </select>
          </Grid>
          <Grid item xs={6}>
            <select className={classes.selectCustom} name="status" onChange={changeStatusFilter}>
              {status.map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </select>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

export default AppAggregatedData;
// export default withRouter(Overvieww);
