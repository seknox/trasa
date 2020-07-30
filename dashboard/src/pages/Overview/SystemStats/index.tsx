// import Storage from '@material-ui/icons/Storage';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
// import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Album from '@material-ui/icons/Album';
import StorageIcon from '@material-ui/icons/AllInbox';
import MemIcon from '@material-ui/icons/DeveloperBoard';
import DNS from '@material-ui/icons/Dns';
import Memory from '@material-ui/icons/Memory';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    //  marginLeft: '27%',
    //  marginRight: '30%',
    justify: 'center',
  },

  paperSmaller: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    marginTop: '20%',

    // marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)' // #011019
    padding: theme.spacing(2),
    textAlign: 'center',
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
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  StatusTextHeader: {
    color: '#011019',
    fontSize: '24px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  StatusTexSubtHeader: {
    color: '#011019',
    fontSize: '17px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  StatusText: {
    color: '#011019',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  Center: {},
}));

export default function SystemStats() {
  const [hostStatus, setHoststatus] = useState({});
  const [diskStatus, setDiskStatus] = useState({});

  const [memStatus, setMemStatus] = useState({});

  const [cpuStatus, setCpuStatus] = useState({});

  const dateConverter = (data: any) => {
    const date = new Date(data * 1000);
    return date.toDateString();
  };

  const secondsToHour = (data: any) => {
    const d = Number(data);
    const h = Math.floor(d / 3600);
    return h;
  };

  const byteToGB = (data: any) => {
    const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    let l = 0;
    let n = parseInt(data, 10) || 0;
    while (n >= 1024 && ++l) n /= 1024;

    return `${n.toFixed(n < 10 && l > 0 ? 1 : 0)} ${units[l]}`;
  };

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/status`)
      .then((response) => {
        setHoststatus(response.data.data[0].hostStatus);
        setDiskStatus(response.data.data[0].diskStatus);
        setMemStatus(response.data.data[0].memStatus);
        setCpuStatus(response.data.data[0].cpuStat);
      })
      .catch((error) => {
        if (error.response.status === 403) {
          window.location.href = '/login';
        }
        if (error.response) {
          console.log(error.response.data);
        } else {
          console.log('Error', error.message);
        }
      });
  }, []);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2} alignItems="center" justify="center">
        <Paper className={classes.paper}>
          <Grid item xs={12} sm={12} md={12}>
            <UpStats
              hostStatus={hostStatus}
              secondsToHour={secondsToHour}
              dateConverter={dateConverter}
            />
          </Grid>

          <Divider light />
          <br />
          <Grid item xs={12} sm={12} md={12}>
            <HostStats
              hostStatus={hostStatus}
              secondsToHour={secondsToHour}
              dateConverter={dateConverter}
            />
          </Grid>

          <Divider light />
          <br />
          <Grid item xs={12} sm={12} md={12}>
            <CpuStats cpuStatus={cpuStatus} byteToGB={byteToGB} />
          </Grid>

          <Divider light />
          <br />
          <Grid item xs={12} sm={12} md={12}>
            <MemStats memStatus={memStatus} byteToGB={byteToGB} />
          </Grid>

          <Divider light />
          <br />
          <Grid item xs={12} sm={12} md={12}>
            <DiskStats diskStatus={diskStatus} byteToGB={byteToGB} />
          </Grid>
        </Paper>
      </Grid>
    </div>
  );
}

function UpStats(props: any) {
  const classes = useStyles();
  return (
    <div>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Album style={{ fontSize: 70, color: '#000080' }} />
        </Grid>
        <Grid item xs={12}>
          <div className={classes.StatusTextHeader}>
            {' '}
            Uptime since <b>{props.secondsToHour(props.hostStatus.uptime)}</b> hours.{' '}
          </div>
          <div className={classes.StatusText}>
            {' '}
            Server Started on {props.dateConverter(props.hostStatus.bootTime)}{' '}
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function HostStats(props: any) {
  const classes = useStyles();
  return (
    <div>
      <Grid container spacing={2} alignItems="center" justify="center">
        <Grid item xs={3}>
          <ListItem>
            <ListItemIcon>
              <DNS style={{ fontSize: 30, color: '#000080' }} />
            </ListItemIcon>
            <ListItemText>
              <div className={classes.StatusTexSubtHeader}> Platform </div>{' '}
            </ListItemText>
          </ListItem>
        </Grid>
        <Grid item xs={12}>
          <div className={classes.StatusText}>
            <b style={{ color: 'grey' }}>OS </b>{' '}
            {`${props.hostStatus.os}(${props.hostStatus.kernelVersion})`} {', '}
            <b style={{ color: 'grey' }}>Platform </b> {props.hostStatus.platform}{' '}
            {props.hostStatus.platformVersion}
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function CpuStats(props: any) {
  const classes = useStyles();
  return (
    <div>
      <Grid container spacing={2} alignItems="center" justify="center">
        <Grid item xs={3}>
          <ListItem>
            <ListItemIcon>
              <Memory style={{ fontSize: 30, color: '#000080' }} />
            </ListItemIcon>
            <ListItemText>
              <div className={classes.StatusTexSubtHeader}> CPU </div>{' '}
            </ListItemText>
          </ListItem>
        </Grid>

        <Grid item xs={12}>
          <div className={classes.StatusText}>
            <b style={{ color: 'grey' }}> CPU(s) </b> {props.cpuStatus.cpuCount} {', '}
            <b style={{ color: 'grey' }}> Used </b>{' '}
            {Math.round(props.cpuStatus.cpuStat * 100) / 100} %
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function MemStats(props: any) {
  const classes = useStyles();
  return (
    <div>
      <Grid container spacing={2} alignItems="center" justify="center">
        <Grid item xs={3}>
          <ListItem>
            <ListItemIcon>
              <MemIcon style={{ fontSize: 27, color: '#000080' }} />
            </ListItemIcon>
            <ListItemText>
              <div className={classes.StatusTexSubtHeader}> Memory </div>{' '}
            </ListItemText>
          </ListItem>
        </Grid>
        <Grid item xs={12}>
          <div className={classes.StatusText}>
            <b style={{ color: 'grey' }}> Total </b> {props.byteToGB(props.memStatus.total)} {', '}
            <b style={{ color: 'grey' }}> Used </b> {props.byteToGB(props.memStatus.used)} {', '}
            <b style={{ color: 'grey' }}> Available </b> {props.byteToGB(props.memStatus.available)}
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function DiskStats(props: any) {
  const classes = useStyles();
  return (
    <div>
      <Grid container spacing={2} alignItems="center" justify="center">
        <Grid item xs={3}>
          <ListItem>
            <ListItemIcon>
              <StorageIcon style={{ fontSize: 30, color: '#000080' }} />
            </ListItemIcon>
            <ListItemText>
              <div className={classes.StatusTexSubtHeader}> Disk </div>{' '}
            </ListItemText>
          </ListItem>
        </Grid>

        <Grid item xs={12}>
          <div className={classes.StatusText}>
            <b style={{ color: 'grey' }}> Total </b> {props.byteToGB(props.diskStatus.total)} {', '}
            <b style={{ color: 'grey' }}> Used </b> {props.byteToGB(props.diskStatus.used)} {', '}
            <b style={{ color: 'grey' }}> Available </b> {props.byteToGB(props.diskStatus.free)}
          </div>
        </Grid>
      </Grid>
    </div>
  );
}
