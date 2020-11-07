import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import SecurityIcon from '@material-ui/icons/Security';
import CheckIcon from '@material-ui/icons/CheckCircle';
import TrustedIcon from '@material-ui/icons/VerifiedUser';
import axios from 'axios';
import React, { useState } from 'react';
import Switch from '@material-ui/core/Switch';
import AndroidIcon from '../../../assets/devices/android.png';
import FirefoxIcon from '../../../assets/devices/firefox.png';
import IphoneIcon from '../../../assets/devices/iphone.svg';
import LinuxLaptopIcon from '../../../assets/devices/linuxlaptop.png';
import MacbookIcon from '../../../assets/devices/mac.png';
import WinLaptopIcon from '../../../assets/devices/winlaptop.png';
import Constants from '../../../Constants';
import { DeleteConfirmDialogue } from '../../../utils/Components/Confirms';
import DialogWrapper from '../../../utils/Components/DialogueWrapComponent';
import { HeaderFontSize, TitleFontSize } from '../../../utils/Responsive';

const useStyles = makeStyles((theme) => ({
  idpPaper: {
    // backgroundColor:  '#fdfdfd',
    // padding: theme.spacing(2),
    textAlign: 'center',
    minWidth: 100,
    maxWidth: 250,
    minHeight: 190,
    color: theme.palette.text.secondary,
    boxShadow: '0 0 5px 0 rgba(0,0,0,0.1)',
  },
  wrappaper: {
    borderColor: 'navy',
  },
  backgroundPaper: {
    padding: theme.spacing(2),
    height: '100%',
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
    transition: '0.3s',
    marginLeft: 50,
    // minWidth: 500,
    // '&:hover': {
    //   transform: 'translateY(-3px)',
    //   boxShadow: '0 4px 20px 0 rgba(0,0,0,0.12)',
    // },
  },
  gridPadding: {
    paddingLeft: theme.spacing(2),
  },
  card: {
    marginTop: 40,
    borderRadius: 0.5, // theme.spacing(0.5),
    transition: '0.3s',
    // width: '90%',
    // maxWidth: 600,
    overflow: 'initial',
    minHeight: 600,
    height: '100%',
    background: '#ffffff',
    // overflowX: 'auto',
    textAlign: 'center',
    paddingRight: theme.spacing(2),
  },
  content: {
    textAlign: 'center',
    padding: theme.spacing(2),
  },
  shadowRise: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
    transition: '0.3s',
    //   '&:hover': {
    //     transform: 'translateY(-3px)',
    //     boxShadow: '0 4px 20px 0 rgba(0,0,0,0.12)',
    //   },
  },
  shadowFaded: {
    boxShadow: '0 2px 4px -2px rgba(0,0,0,0.24), 0 4px 24px -2px rgba(0, 0, 0, 0.2)',
  },
  cardHeader: {
    boxShadow: '0 2px 4px -2px rgba(0,0,0,0.24), 0 4px 24px -2px rgba(0, 0, 0, 0.2)',
    background: 'navy',
    borderRadius: 8,
    margin: '-20px auto 0',
    width: '38%',
    color: 'white',
    fontSize: '18px',
    fontWeight: 'bold',
    minHeight: 50,
  },
  title: {
    color: 'white',
    fontWeight: 'bold',
  },
  subheader: {
    color: 'rgba(255, 255, 255, 0.76)',
  },
  settingHeader: {
    textAlign: 'left',
    color: '#1b1b32',
    fontSize: HeaderFontSize(), // window.innerHeight < 750 ? '14px':'18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingValue: {
    textAlign: 'center',
    color: 'black',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeadersBig: {
    color: '#000066',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    // textAlign: 'left',
    color: '#1b1b32',
    fontSize: HeaderFontSize(), // window.innerHeight < 750 ? '14px':'18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  paperTrans: {
    backgroundColor: 'transparent',
    // padding: 4,
    textAlign: 'center',
  },
  drawer: {
    // padding: 20,
    background: 'white',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperCertSetting: {
    minWidth: 600,
    maxWidth: 600,
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
    minHeight: 700,
  },
  proxyPaper: {
    boxShadow: '0 0 20px 0 rgba(0,0,0,0.12)',
    padding: theme.spacing(2),
    textAlign: 'left',
  },
  proxyVals: {
    color: '#000066',
    textAlign: 'left',
    fontSize: 14,
    fontFamily: 'Open Sans, Rajdhani',
  },
  Warning: {
    color: 'maroon',
    fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  WarningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
}));

const getIcon = (device: any) => {
  if (device.deviceType === 'browser') {
    return <img src={FirefoxIcon} alt="FirefoxIcon" height={70} />;
  }
  if (device.deviceType === 'workstation') {
    const devOs = device?.deviceHygiene?.deviceOS;
    if (devOs.kernelType === 'windows') {
      return <img src={WinLaptopIcon} alt="FirefoxIcon" height={80} />;
    }
    if (devOs.kernelType === 'linux') {
      return <img src={LinuxLaptopIcon} alt="FirefoxIcon" height={80} />;
    }
    return <img src={MacbookIcon} alt="FirefoxIcon" height={80} />;
  }

  const devHyg = device.deviceHygiene;
  if (devHyg && devHyg.deviceOS.OSName === 'Android') {
    return <img src={AndroidIcon} alt="androidIcon" height={70} />;
  }
  if (devHyg && devHyg.deviceOS.OSName === 'ios') {
    return <img src={IphoneIcon} alt="IphoneIcon" height={70} />;
  }
  return <img src={AndroidIcon} alt="androidIcon" height={70} />;
};

const getDeviceString = (device: any) => {
  if (device && device.deviceType === 'mobile') {
    if (device.deviceHygiene && device?.deviceHygiene?.deviceInfo) {
      return `${device.deviceHygiene.deviceInfo.manufacturer} ${device.deviceHygiene.deviceInfo.deviceModel}`;
    }
    return 'mobile';
  }
  // TODO
  // if (device && device.deviceType === 'browser') {
  //   return device.deviceHygiene.deviceBrowser.browserName || 'Browser';
  // }
  if (device && device.deviceType === 'workstation') {
    return device.deviceHygiene && device.deviceHygiene?.deviceOS?.OSName;
  }
  return 'Unknown Device';
};

// ////////////////////////////////////
// DeviceDetail is a card for device.
// ///////////////////////////////////
export default function DeviceDetail(props: any) {
  const [open, setOpen] = React.useState(false);

  const openDeleteDlg = () => {
    setOpen(!open);
  };

  function callDeleteFunc(deviceDetail: any, index: any) {
    props.deleteDevice({
      deviceID: deviceDetail.deviceID,
      deviceType: deviceDetail.deviceType,
      deviceIndex: index,
    });
  }

  const [singleDeviceDetail, setSingleDeviceDetail] = useState({
    deviceID: '',
    trusted: false,
    deviceHygiene: {},
  });
  // const [singleDeviceHygiene, setSingleDeviceHygiene] = useState({});
  const [viewHygieneDlgState, setviewHygieneDlgState] = useState(false);
  const [singleDeviceType, setsingleDeviceType] = useState('');
  const [selectedDeviceIndex, setselectedDeviceIndex] = useState(0);

  function viewDeviceHygiene(dd: any, index: number) {
    // console.log('device:: ', dd);
    setSingleDeviceDetail(dd);
    // setSingleDeviceHygiene(dd.deviceHygiene);
    setsingleDeviceType(props.deviceType);
    setselectedDeviceIndex(index);
    setviewHygieneDlgState(true);
  }

  function handleTrust(deviceID: string, trusted: boolean) {
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/user/devices/trust`, { deviceID, trusted })
      .then((r) => {
        if (r.data.status === 'success') {
          // TODO
        }
      });
  }

  const classes = useStyles();
  return (
    <Grid container spacing={2} direction="row">
      {props.userDevices.map((d: any, i: any) => (
        <Grid item xs={3}>
          <Paper className={classes.idpPaper} elevation={1}>
            <Grid container direction="row" alignItems="center" justify="center" spacing={1}>
              <Grid item xs={1}>
                {d.trusted ? (
                  <Tooltip title="This device is trusted." placement="top-start">
                    <div>
                      <TrustedIcon style={{ fontSize: 20, color: 'green' }} />
                    </div>
                  </Tooltip>
                ) : (
                  <Tooltip title="This device is not trusted." placement="top-start">
                    <div>
                      <TrustedIcon style={{ fontSize: 20, color: 'grey' }} />
                    </div>
                  </Tooltip>
                )}
              </Grid>
              <Grid item xs={9}>
                <Typography variant="h3">{getDeviceString(d)}</Typography>
              </Grid>
              <Divider light />
              <Grid item xs={12}>
                {getIcon(d)}
              </Grid>
              <Grid item xs={6}>
                <Button
                    name={"viewDetailBtn"}
                  size="small"
                  color="secondary"
                  variant="outlined"
                  onClick={() => viewDeviceHygiene(d, i)}
                >
                  View Detail
                </Button>
              </Grid>
            </Grid>
          </Paper>
        </Grid>
      ))}
      <DeleteConfirmDialogue
        open={open}
        close={openDeleteDlg}
        deleteFunc={props.deleteFunc}
        confirmMessage="Are You Sure You Want To Remove this device?"
      />

      <DialogWrapper
        open={viewHygieneDlgState}
        handleClose={() => setviewHygieneDlgState(false)}
        title="Device Hygiene"
        maxWidth="lg"
        fullScreen={false}
      >
        <DeviceSettings
          singleDeviceDetail={singleDeviceDetail}
          selectedDeviceIndex={selectedDeviceIndex}
          openDeleteDlg={openDeleteDlg}
          callDeleteFunc={callDeleteFunc}
          trusted={singleDeviceDetail.trusted}
          handleTrustDevice={handleTrust}
          renderFor={props.renderFor}
        />

        <DeviceHygiene dh={singleDeviceDetail.deviceHygiene} dt={singleDeviceType} />
      </DialogWrapper>
    </Grid>
  );
}

function DeviceHygiene(props: any) {
  const classes = useStyles();
  const { dh } = props;
  return (
    <Grid container spacing={2}>
      <Grid item xs={4}>
        <Paper className={classes.proxyPaper}>
          <DeviceInfo di={dh.deviceInfo} db={dh.deviceBrowser} dt={props.dt} />
        </Paper>
      </Grid>
      <Grid item xs={4}>
        <Paper className={classes.proxyPaper}>
          <OSInfo oi={dh.deviceOS} dt={props.dt} />
        </Paper>
      </Grid>

      <Grid item xs={4}>
        <Paper className={classes.proxyPaper}>
          <NetworkInfo ni={dh.networkInfo} dt={props.dt} />
        </Paper>
      </Grid>

      <Grid item xs={4}>
        <Paper className={classes.proxyPaper}>
          <LoginSecurity ls={dh.loginSecurity} dt={props.dt} />
        </Paper>
      </Grid>

      <Grid item xs={4}>
        <Paper className={classes.proxyPaper}>
          <EndpointSecurity es={dh.endpointSecurity} dt={props.dt} />
        </Paper>
      </Grid>
    </Grid>
  );
}

function DeviceSettings(props: any) {
  const [checked, setChecked] = useState(false);

  React.useEffect(() => {
    setChecked(props.singleDeviceDetail.trusted);
  }, [props.singleDeviceDetail]);

  function handleCheck(e: React.ChangeEvent<HTMLInputElement>) {
    props.handleTrustDevice(props.singleDeviceDetail.deviceID, e.target.checked);
    setChecked(e.target.checked);
  }

  return (
    <Grid container spacing={2} direction="row">
      {props.renderFor === 'myRoute' ? (
        ''
      ) : (
        <Grid item xs={3}>
          <Grid container spacing={2}>
            <Grid item xs={8}>
              <Typography variant="h3">Trust this device : </Typography>
            </Grid>
            <Grid item xs={4}>
              <Switch name={"deviceTrustSwitch"} color="secondary" onChange={handleCheck} checked={checked} />
            </Grid>
          </Grid>
        </Grid>
      )}
      <Grid item xs={2}>
        <Typography variant="h3">Delete this device : </Typography>
      </Grid>
      <Grid item xs={1}>
        <IconButton
          aria-label="delete"
          onClick={() => {
            props.openDeleteDlg();
            props.callDeleteFunc(props.singleDeviceDetail, props.selectedDeviceIndex);
          }}
        >
          <DeleteIcon color="error" />
        </IconButton>
      </Grid>
    </Grid>
  );
}

function DeviceInfo(props: any) {
  const classes = useStyles();
  const { di } = props;

  // note all variables names are shorthanded. device info as di, login security as ls,
  // device type as dt, hygiene data as hd etc..
  // function getDeviceValue(hd){
  //     switch(true){
  //         case dt==='mobile' && hd==='dv'{
  //             return
  //         }
  //     }
  // }

  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <Typography variant="h3">Device Info</Typography>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Brand </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{di.brand}</div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Device Model </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{di.deviceName}</div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Device Version </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{di.deviceVersion}</div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Manufacturer </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{di.manufacturer}</div>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

function OSInfo(props: any) {
  const classes = useStyles();

  const { oi } = props;
  const { dt } = props;
  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <Typography variant="h3">OS Info</Typography>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> OS Name </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}> {oi.osName} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> OS Version </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}> {oi.osVersion} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> OS Autoupdate </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}>{oi.autoUpdate ? 'enabled' : 'disabled'}</div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Debug Mode Enabled </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}>
              {dt.deviceType === 'mobile' ? (oi.debugModeEnabled ? 'enabled' : 'disabled') : 'n.a'}
            </div>
          </Grid>
        </Grid>
      </Grid>

      {/* <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Pending Updates </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>

              {oi.debugModeEnabled ? 'enabled' : 'disabled'}
            </div>
          </Grid>
        </Grid>
      </Grid> */}
    </Grid>
  );
}

function EndpointSecurity(props: any) {
  const classes = useStyles();
  const { es } = props;
  const { dt } = props;
  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant="h3">Endpoint Security</Typography>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Device Encryption </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>
              {dt === 'mobile' ? 'n.a' : es.deviceEncryptionEnabled ? 'enabled' : 'disabled'}
            </div>
          </Grid>
        </Grid>
      </Grid>

      {/* <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Device Encryption Metadata </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}> n.a </div>
          </Grid>
        </Grid>
      </Grid> */}

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Endpoint Security </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}> {es.epsVendorName} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Firewall Enabled </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}> {es.firewallEnabled ? 'enabled' : 'disabled'} </div>
          </Grid>
        </Grid>
      </Grid>

      {/* <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Firewall Policy </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}> n.a </div>
          </Grid>
        </Grid>
      </Grid> */}
    </Grid>
  );
}

function LoginSecurity(props: any) {
  const classes = useStyles();

  const { ls } = props;
  const { dt } = props;
  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant="h3">Login Security</Typography>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Auto Login </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}>
              {/* TODO imp @bhrg3se this data is incorrect from device hygiene */}
              {ls.autologinEnabled ? 'enabled' : 'disabled'}
            </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Login Method </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}> n.a </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Screen Lock </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}>
              {/* TODO imp @bhrg3se this data is incorrect from device hygiene */}
              {ls.idleDeviceScreenLock ? 'enabled' : 'disabled'}
            </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Idle-screen lockout </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}> {`${ls.idleDeviceScreenLockTime}(seconds)`} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={6}>
            <div className={classes.settingHeader}> Password Updated </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.proxyVals}>
              {dt === 'mobile' ? 'n.a' : ls.passwordLastUpdated}
            </div>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

function NetworkInfo(props: any) {
  const classes = useStyles();

  function pIfc(d: any) {}

  const { ni } = props;
  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant="h3">Network Security</Typography>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Hostname </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}> {ni.hostname} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Interface Name </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{ni.interfaceName}</div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> IP Address </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}> {ni.ipAddress} </div>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}> Wireless Network </div>
          </Grid>
          <Grid item xs={7}>
            <div className={classes.proxyVals}>{ni.wirelessNetwork ? 'yes' : 'no'}</div>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}
