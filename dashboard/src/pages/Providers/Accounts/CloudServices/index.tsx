import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
// import SwipeableViews from 'react-swipeable-views';
import { makeStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import axios from 'axios';
// import mixpanel from 'mixpanel-browser';
import React, { useEffect, useState } from 'react';
import AWS from '../../../../assets/cloudiaas/aws.png';
import DO from '../../../../assets/cloudiaas/do.svg';
import GCP from '../../../../assets/cloudiaas/gcp.png';
import Constants from '../../../../Constants';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import { HeaderFontSize, TitleFontSize } from '../../../../utils/Responsive';

function TabPanel(props: any) {
  const { children, value, index, ...other } = props;

  return (
    <Typography
      component="div"
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index}`}
      aria-labelledby={`full-width-tab-${index}`}
      {...other}
    >
      {value === index && <Box p={3}>{children}</Box>}
    </Typography>
  );
}

function a11yProps(index: any) {
  return {
    id: `full-width-tab-${index}`,
    'aria-controls': `full-width-tabpanel-${index}`,
  };
}

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(4),
    textAlign: 'center',
    display: 'flex',
    color: theme.palette.text.secondary,
  },
  paperTrans: {
    backgroundColor: 'transparent',
    padding: theme.spacing(1),
    textAlign: 'center',
  },
  idpPaper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    textAlign: 'center',
    minWidth: 100,
    color: theme.palette.text.secondary,
  },
  marginer: {
    marginBottom: 10,
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {},
  aggHeadersBig: {
    color: '#000066',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: HeaderFontSize(), // window.innerHeight < 750 ? '14px':'18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function StatusPage() {
  const classes = useStyles();

  // function returnIcon(val: string) {
  //   switch (val) {
  //     case 'DO':
  //       return <img src={DO} alt={val} height={70} />;
  //     case 'AWS':
  //       return <img src={GCP} alt={val} height={70} />;
  //     case 'GCP':
  //       return <img src={AWS} alt={val} height={70} />;
  //   }
  // }

  const [value, setValue] = React.useState(1);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setValue(newValue);
  };

  return (
    <div className={classes.root}>
      <Grid container spacing={2} direction="column" justify="center" alignItems="center">
        <Grid item xs={12}>
          <Typography variant="h2">Connect TRASA with your cloud platform</Typography>
        </Grid>

        <br />
        <Grid item xs={12}>
          <Tabs
            value={value}
            onChange={handleChange}
            indicatorColor="primary"
            textColor="primary"
            centered
          >
            <Tab
              label="(Coming Soon)"
              icon={<img src={AWS} alt="AWS" height={70} {...a11yProps(0)} />}
            />
            <Tab
              label="Digital Ocean"
              icon={<img src={DO} alt="DO" height={70} {...a11yProps(1)} />}
            />
            <Tab
              label="(Coming Soon)"
              icon={<img src={GCP} alt="GCP" height={70} {...a11yProps(2)} />}
            />
          </Tabs>
          <Divider light />
        </Grid>

        <Grid item xs={12}>
          <Paper className={classes.paper}>
            <TabPanel value={value} index={0}>
              <AwsSetting />
            </TabPanel>
            <TabPanel value={value} index={1}>
              <DoSetting />
            </TabPanel>
            <TabPanel value={value} index={2}>
              <GcpSetting />
            </TabPanel>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

export default StatusPage;

function DoSetting() {
  return (
    <div>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <IntegrationStats />
        </Grid>

        <Grid item xs={12}>
          <br />
          <Divider light />
        </Grid>
        <Grid item xs={12}>
          <Typography variant="h2">Connect and sync with Digitalocean</Typography>
        </Grid>

        <Grid item xs={12}>
          <DOConnectAndSync />
        </Grid>
      </Grid>
    </div>
  );
}

function AwsSetting() {
  return (
    <div>
      <Grid container spacing={2} direction="column" justify="center" alignItems="center">
        <Grid item xs={12}>
          AWS
        </Grid>
      </Grid>
    </div>
  );
}

function GcpSetting() {
  return (
    <div>
      <Grid container spacing={2} direction="column" justify="center" alignItems="center">
        <Grid item xs={12}>
          GCP
        </Grid>
      </Grid>
    </div>
  );
}

function IntegrationStats() {
  const classes = useStyles();
  const [sshCount, setSshCount] = useState(0);
  const [dbCount, setDbCount] = useState(0);
  const [rdpCount, setRdpCount] = useState(0);
  const [httpCount, setHttpCount] = useState(0);

  useEffect(() => {
    // TODO @sshah make external provider ID dynamic
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/stats/serviceidp/digital-ocean`)

      .then((response) => {
        // console.log(response.data)
        if (response.data.status === 'success') {
          const { data } = response.data;
          data.forEach((a: any) => {
            switch (a.name) {
              case 'db':
                setDbCount(a.value);
                break;
              case 'http':
                setHttpCount(a.value);
                break;
              case 'rdp':
                setRdpCount(a.value);
                break;
              case 'ssh':
                setSshCount(a.value);
                break;
              default:
                break;
            }
          });
        }
      });
  }, []);

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Typography variant="h2">Integration Stats</Typography>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={2}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> Total service </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> {sshCount + rdpCount + httpCount + dbCount} </b>
              </div>
            </Paper>
          </Grid>
          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={2}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> HTTPs </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> {httpCount} </b>
              </div>
            </Paper>
          </Grid>

          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={2}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> SSH </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> {sshCount} </b>
              </div>
            </Paper>
          </Grid>

          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={2}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> DB </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> {dbCount} </b>
              </div>
            </Paper>
          </Grid>

          <Grid item xs={2}>
            <Paper className={classes.paperTrans} elevation={2}>
              <div className={classes.aggHeaders}>
                {' '}
                <b> Last Synced </b>
              </div>
              <div className={classes.aggHeadersBig}>
                {' '}
                <b> Mon. March 2020 </b>
              </div>
            </Paper>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

function DOConnectAndSync() {
  const classes = useStyles();
  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
  });

  const [storedKey, setStoredKey] = useState({ keyTag: '', keyName: 'KEY_DOAPI', keyVal: '' });
  // const [lastSynced, setLastSynced] = useState({ time: '' });

  const handleChange = (name: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
    const val = event.target.value;
    setStoredKey({ ...storedKey, [name]: val });
  };

  const storeApiKey = () => {
    updateActionStatus({ ...actionStatus, respStatus: false, statusMsg: '', loader: true });
    // mixpanel.track('manage-accounts-dointegration');
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/store/key`, storedKey)
      .then((r) => {
        updateActionStatus({ ...actionStatus, loader: false });
        if (r.data.status === 'success') {
          updateActionStatus({
            ...actionStatus,
            success: true,
            respStatus: true,
            statusMsg: r.data.reason,
          });
        } else {
          updateActionStatus({
            ...actionStatus,
            respStatus: true,
            statusMsg: r.data.reason,
            success: false,
          });
        }
      })
      .catch((error) => {
        updateActionStatus({
          ...actionStatus,
          respStatus: true,
          statusMsg: 'something went wrong',
          success: false,
        });
      });
  };

  const syncNow = () => {
    updateActionStatus({ ...actionStatus, respStatus: false, statusMsg: '', loader: true });
    // mixpanel.track('manage-accounts-dointegration');
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .post(
        `${Constants.TRASA_HOSTNAME}/api/v1/providers/sidp/syncnow/do`,
        storedKey,
        config,
      )
      .then((r: any) => {
        updateActionStatus({ ...actionStatus, loader: false });
        if (r.data.status === 'success') {
          updateActionStatus({
            ...actionStatus,
            success: true,
            respStatus: true,
            statusMsg: r.data.reason,
          });
        } else {
          updateActionStatus({
            ...actionStatus,
            respStatus: true,
            statusMsg: r.data.reason,
            success: false,
          });
        }
      })
      .catch((error) => {
        updateActionStatus({
          ...actionStatus,
          respStatus: true,
          statusMsg: 'something went wrong',
          success: false,
        });
      });
  };

  useEffect(() => {
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/key/KEY_DOAPI`, config)

      .then((response) => {
        // console.log(response.data)
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          setStoredKey(data);
        }
      });
  }, []);

  return (
    <div>
      <div className={classes.root}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Grid container spacing={3}>
              <Grid item xs={3}>
                <Typography variant="h3">API Key : </Typography>
              </Grid>
              <Grid item xs={7}>
                <TextField
                  fullWidth
                  // label="Service name"
                  onChange={handleChange('keyVal')}
                  name="keyVal"
                  // variant = 'outlined'
                  value={storedKey.keyVal}
                  defaultValue={storedKey.keyVal}
                  InputProps={{
                    disableUnderline: true,
                    classes: {
                      root: classes.textFieldRoot,
                      input: classes.textFieldInputBig,
                    },
                  }}
                  InputLabelProps={{
                    shrink: true,
                    className: classes.textFieldFormLabel,
                  }}
                />
              </Grid>

              <Grid item xs={1}>
                <Button variant="contained" onClick={storeApiKey}>
                  Save
                </Button>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={3}>
              <Grid item xs={3}>
                <Typography variant="h3">Sync Now :</Typography>
              </Grid>

              <Grid item xs={8}>
                <Button fullWidth variant="contained" onClick={syncNow}>
                  Sync Now
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
        <br />
        {actionStatus.loader ? <ProgressHOC /> : ''}
      </div>
      <Divider light />
    </div>
  );
}
