import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import DeviceList from './DeviceList';

const useStyles = makeStyles((theme) => ({
  idpPaper: {
    backgroundColor: '#fdfdfd',
    // minWidth: 200,
    // padding: theme.spacing(2),
    textAlign: 'center',
    borderColor: 'navy',
    //   minWidth: 100,
    //   maxWidth: 200,
    minHeight: 190,
    color: theme.palette.text.secondary,
  },
  wrappaper: {
    borderColor: 'navy',
  },
  paperTrans: {
    backgroundColor: 'transparent',
  },
  aggHeadersBig: {
    color: '#000066',
    fontSize: '21px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: '19px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

type DeviceProps = {
  userID: string,
  renderFor: 'myRoute' | 'userRoute',

}

export default function Device(props: DeviceProps) {
  const [deviceDeleteState, setDeviceDeleteState] = useState(false);
  // const [deviceInfoDialogueOpen, setdeviceInfoDialogueOpen] = useState(false);
  const [deleteDevice, setDeleteDevice] = useState({
    deviceID: '',
    deviceType: '',
    deviceIndex: 0,
  });
  // const [selDeviceFinger, setselDeviceFinger] = useState({});
  // const [selDevice, setselDevice] = useState({});

  const [userDevices, setUserDevices] = useState({ mobile: [], workstation: [], browser: [] });

  useEffect(() => {
    let reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/my/devices`;
    if (props.renderFor !== 'myRoute') {
      reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/user/devices/all/${props.userID}`;
    }

    axios
      .get(reqPath)
      .then((r) => {
        setUserDevices(r?.data?.data?.[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [props.renderFor, props.userID]);

  // removeUserDevice sends delete device request to trasacore and removes deleted device from local deivce state after successful response
  const removeUserDevice = () => {

    let apiURL=`${Constants.TRASA_HOSTNAME}/api/v1/user/devices/delete/${deleteDevice.deviceID}`;
    if(props.renderFor=="myRoute"){
      apiURL=`${Constants.TRASA_HOSTNAME}/api/v1/my/devices/delete/${deleteDevice.deviceID}`;
    }

    axios
      .post(apiURL)
      .then((response) => {
        if (response.data.status === 'success') {
          let device = userDevices.browser;
          switch (true) {
            case deleteDevice.deviceType === 'browser':
              device = userDevices.browser;
              device.splice(deleteDevice.deviceIndex, 1);
              setUserDevices({ ...userDevices, browser: device });
              break;
            case deleteDevice.deviceType === 'workstation':
              device = userDevices.workstation;
              device.splice(deleteDevice.deviceIndex, 1);
              setUserDevices({ ...userDevices, workstation: device });
              break;
            case deleteDevice.deviceType === 'mobile':
              device = userDevices.mobile;
              device.splice(deleteDevice.deviceIndex, 1);
              setUserDevices({ ...userDevices, mobile: device });
              break;
            default:
              break;
          }
        }
      })
      .catch((error) => {
        console.log('ERRRRR: ', error);
      });
  };

  function openDeviceDeleteDlg() {
    setDeviceDeleteState(true);
    // setToDeleteVal(device)
  }

  function closeDeviceDeleteDlg() {
    setDeviceDeleteState(false);
  }

  function showDeviceDetail(dev: any, devFinger: any) {
    // setselDevice(dev);
    // setselDeviceFinger(devFinger);
    //  setdeviceInfoDialogueOpen(true);
  }

  function setDeleteDeviceFunc(v: any) {
    setDeleteDevice(v);
  }

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <DeviceStats
          total={
            userDevices.mobile.length + userDevices.browser.length + userDevices.workstation.length
          }
          mobile={userDevices.mobile.length}
          workstation={userDevices.workstation.length}
          browser={userDevices.browser.length}
        />
        {/* <Divider light /> */}
        <br />
      </Grid>
      <Grid item xs={12}>
        <Typography variant="h3">Mobile Devices</Typography>
        {/* <Divider light /> */}
        <br />
        <DeviceList
          userDevices={userDevices.mobile}
          deviceType="mobile"
          deleteUserDlgOpen={openDeviceDeleteDlg}
          open={deviceDeleteState}
          close={closeDeviceDeleteDlg}
          deleteFunc={removeUserDevice}
          deleteDevice={setDeleteDeviceFunc}
          showDeviceDetail={showDeviceDetail}
          renderFor={props.renderFor}
        />
      </Grid>

      <Grid item xs={12}>
        <Typography variant="h3">Workstations</Typography>
        {/* <Divider light /> */}
        <br />
        <DeviceList
          userDevices={userDevices.workstation}
          deviceType="workstation"
          deleteUserDlgOpen={openDeviceDeleteDlg}
          open={deviceDeleteState}
          close={closeDeviceDeleteDlg}
          deleteFunc={removeUserDevice}
          deleteDevice={setDeleteDeviceFunc}
          showDeviceDetail={showDeviceDetail}
          renderFor={props.renderFor}
        />
      </Grid>
    </Grid>
  );
}

function DeviceStats(props: any) {
  const classes = useStyles();

  return (
    <Grid container spacing={2}>
      <Grid item xs={3}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Total </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {props.total} </b>{' '}
          </div>
        </Paper>
      </Grid>

      <Grid item xs={3}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b>Mobile Phone </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {props.mobile} </b>{' '}
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
            <b> {props.workstation} </b>{' '}
          </div>
        </Paper>
      </Grid>

      {/* <Grid item xs={3}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Browsers </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {props.browser} </b>{' '}
          </div>
        </Paper>
      </Grid> */}
    </Grid>
  );
}
