import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
// import { useContainedCardHeaderStyles } from '@mui-treasury/styles/cardHeader/contained';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import CopyIcon from '@material-ui/icons/FileCopy';
import axios from 'axios';
import cx from 'clsx';
import React, { useEffect, useState } from 'react';
import ProgressHOC from '../../../utils/Components/Progressbar';
import { HeaderFontSize, TitleFontSize } from '../../../utils/Responsive';
import Constants from '../../../Constants';
import ServiceSetting from './Settings/ServiceSetting';
import ProxySetting from './Settings/ProxySetting';
import CertSetting from './Settings/CertSetting';

const useStyles = makeStyles((theme) => ({
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

type Anchor = 'top' | 'left' | 'bottom' | 'right';

/// ///////////////////////////////////////////////////////
/// ///////////////////////////////////////////////////////
/// //   ServiceOverview wraps service aggregated data and settings
/// ///////////////////////////////////////////////////////
/// ///////////////////////////////////////////////////////
export default function ServiceOverview(props: any) {
  const classes = useStyles();

  const [serviceDetail, setserviceDetail] = useState({
    ID: '',
    serviceName: '',
    serviceType: '',
    secretKey: '',
    rdpProtocol: '',
    remoteserviceName: '',
    hostname: '',
    adhoc: false,
    passthru: false,
    nativeLog: false,
    proxyConfig: {
      routeRule: '',
      passHostHeader: false,
      upstreamServer: '',
      strictTLSValidation: false,
    },
  });
  const [proxyDetail, setProxyDetail] = useState({});

  useEffect(() => {
    // const serviceName = props.entityID;
    // var url = ({this.props.urlID }).bind(this)

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/services/${props.entityID}`)
      .then((r) => {
        if (r.data.status === 'success') {
          setserviceDetail(r.data.data[0]);
          setProxyDetail(r.data.data[0].proxyConfig);
        }
      })
      .catch((error) => {
        console.log('Error', error);
      });
  }, [props.entityID]);

  // Drawer States
  const [configDrawerState, setConfigDrawerState] = useState({ right: false });
  const [proxyConfigDrawerState, setProxyConfigDrawerState] = useState({ right: false });
  const [certConfigDrawerState, setCertConfigDrawerState] = useState({ right: false });

  const toggleConfigDrawer = (side: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }
    setConfigDrawerState({ ...configDrawerState, [side]: open });
  };

  const toggleProxyConfigDrawer = (side: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }
    setProxyConfigDrawerState({ ...proxyConfigDrawerState, [side]: open });
  };

  const toggleCertConfigDrawer = (side: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }
    setCertConfigDrawerState({ ...certConfigDrawerState, [side]: open });
  };

  const [deleteDlgState, setDeleteDlgState] = useState(false);

  function closeDeleteDlg() {
    setDeleteDlgState(false);
  }

  return (
    <Grid container spacing={2} direction="column" alignItems="center" justify="center">
      <Paper className={cx(classes.card, classes.shadowRise)}>
        <Grid
          container
          direction="row"
          alignItems="center"
          justify="center"
          className={classes.cardHeader}
        >
          <Grid item xs={11}>
            {serviceDetail.serviceName}
          </Grid>

          <Grid item xs={1}>
            <Tooltip title="delete" placement="top">
              <IconButton
                id="deleteBtn"
                style={{ color: 'maroon' }}
                onClick={() => {
                  setDeleteDlgState(true);
                }}
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
          </Grid>
        </Grid>
        <br /> <br />
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <Paper className={classes.backgroundPaper}>
              <ServiceDetail
                serviceDetail={serviceDetail}
                toggleConfigDrawer={toggleConfigDrawer}
              />
              <br />
              <Divider light />
              <br />
              <IntegrationKeys serviceDetail={serviceDetail} />
              {/* <ProxyStatus toggleConfigDrawer={toggleConfigDrawer} />                          */}
            </Paper>
          </Grid>

          <Grid item xs={7}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Paper className={classes.proxyPaper}>
                  <AggStats ID={serviceDetail.ID} />
                </Paper>
              </Grid>

              <Grid item xs={6}>
                <Paper className={classes.proxyPaper}>
                  <ProxyStatus
                    toggleProxyConfigDrawer={toggleProxyConfigDrawer}
                    proxy={proxyDetail}
                    serviceDetail={serviceDetail}
                    hostname={serviceDetail.hostname}
                  />
                </Paper>
              </Grid>
              <Grid item xs={6}>
                <Paper className={classes.proxyPaper}>
                  <CertStatus toggleCertConfigDrawer={toggleCertConfigDrawer} />
                </Paper>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Paper>

      <Drawer
        anchor="right"
        open={configDrawerState.right}
        onClose={toggleConfigDrawer('right', false)}
      >
        <Paper className={classes.paper}>
          <ServiceSetting newApp={false} serviceDetail={serviceDetail} />
        </Paper>
      </Drawer>

      <Drawer
        anchor="right"
        open={proxyConfigDrawerState.right}
        onClose={toggleProxyConfigDrawer('right', false)}
      >
        <Paper className={classes.paper}>
          <ProxySetting proxyConfig={serviceDetail.proxyConfig} serviceDetail={serviceDetail} />
        </Paper>
      </Drawer>

      <Drawer
        anchor="right"
        open={certConfigDrawerState.right}
        onClose={toggleCertConfigDrawer('right', false)}
      >
        <Paper className={classes.paperCertSetting}>
          <CertSetting serviceType={serviceDetail.serviceType} serviceID={serviceDetail.ID} />
        </Paper>
      </Drawer>

      <DeleteServiceDlg open={deleteDlgState} close={closeDeleteDlg} serviceID={serviceDetail.ID} />
    </Grid>
  );
}

/// ////////////////////////////////////////////
/// ////////   service Config Details  /////////////
/// ////////////////////////////////////////////
function ServiceDetail(props: any) {
  const classes = useStyles();

  return (
    <Grid container spacing={1} className={classes.gridPadding}>
      <Grid container direction="row" alignItems="flex-start" justify="flex-start" spacing={2}>
        <Grid item xs={4}>
          <Typography variant="h3">Configurations</Typography>
          {/* <Divider  light /> */}
        </Grid>
        <Grid item xs={3}>
          <Tooltip title="edit" placement="top">
            <IconButton
              id="configEditBtn"
              style={{ color: 'navy' }}
              onClick={props.toggleConfigDrawer('right', true)}
            >
              <EditIcon />
            </IconButton>
          </Tooltip>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>Service Name :</div>
          </Grid>
          <Grid item xs={7}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.serviceName}
            </Typography>
          </Grid>
        </Grid>
        {/* <Divider  light /> */}
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>Adhoc Permission :</div>
          </Grid>
          <Grid item xs={7}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.adhoc ? 'enabled' : 'disabled'}
            </Typography>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>Passthru :</div>
          </Grid>
          <Grid item xs={7}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.passthru ? 'enabled' : 'disabled'}
            </Typography>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>Native Agent Logs :</div>
          </Grid>
          <Grid item xs={7}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.nativeLog ? 'enabled' : 'disabled'}
            </Typography>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>Application Type :</div>
          </Grid>
          <Grid item xs={7}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.serviceType}
            </Typography>
          </Grid>
        </Grid>
      </Grid>

      {props.serviceDetail.serviceType === 'rdp' ? (
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={5}>
              <div className={classes.settingHeader}>Security Protocol : </div>
            </Grid>
            <Grid item xs={7}>
              <Typography component="h4" className={classes.proxyVals}>
                {props.serviceDetail.rdpProtocol}
              </Typography>
            </Grid>

            <Grid item xs={5}>
              <div className={classes.settingHeader}>Remote Service name : </div>
            </Grid>
            <Grid item xs={7}>
              <Typography component="h4" className={classes.proxyVals}>
                {props.serviceDetail.remoteserviceName}
              </Typography>
            </Grid>
          </Grid>
        </Grid>
      ) : null}
      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            {props.serviceDetail.serviceType === 'http' ? (
              <div className={classes.settingHeader}>Domain name : </div>
            ) : (
              <div className={classes.settingHeader}>Hostname : </div>
            )}
          </Grid>
          <Grid item xs={6}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.hostname}
            </Typography>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

/// //////////////////////////////////////////////////////////////
/// ////// Integration data are required for native 2fa trasa agent
/// ///////////////////////////////////////////////////////////////
function IntegrationKeys(props: any) {
  const classes = useStyles();
  return (
    <Grid container spacing={1} className={classes.gridPadding}>
      <Grid container direction="row" alignItems="flex-start" justify="flex-start" spacing={2}>
        <Grid item xs={6}>
          <Typography variant="h3">Integration Values</Typography>
          {/* <Divider  light /> */}
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>ServiceID :</div>
          </Grid>
          <Grid item xs={5}>
            <Typography component="h4" className={classes.proxyVals}>
              {props.serviceDetail.ID
                ? `${props.serviceDetail.ID.substring(0, props.serviceDetail.ID.length - 27)}xxx`
                : ''}
            </Typography>
          </Grid>

          <Grid item xs={2}>
            <Tooltip title="copy" placement="top">
              <IconButton
                size="small"
                color="secondary"
                onClick={() => {
                  navigator &&
                    navigator.clipboard &&
                    navigator.clipboard.writeText(props.serviceDetail.ID);
                }}
                aria-label="copy"
                component="button"
              >
                <CopyIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <div className={classes.settingHeader}>ServiceKey :</div>
          </Grid>
          <Grid item xs={5}>
            <Typography component="h4" className={classes.proxyVals}>
              xxxx-xxxx-xxx
            </Typography>
          </Grid>

          <Grid item xs={2}>
            <Tooltip title="copy" placement="top">
              <IconButton
                size="small"
                color="secondary"
                onClick={() => {
                  navigator &&
                    navigator.clipboard &&
                    navigator.clipboard.writeText(props.serviceDetail.secretKey);
                }}
                aria-label="copy"
                component="button"
              >
                <CopyIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

/// ////////////////////////////////////////////
// AggStats show aggregated data of service details
/// ////////////////////////////////////////////
function AggStats(props: any) {
  const classes = useStyles();
  const [data, setData] = useState({ users: 0, policy: 0, groups: 0, secrets: 0, privileges: 0 });

  useEffect(() => {
    if (props.ID) {
      axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/stats/appperms/${props.ID}`).then((res) => {
        res.data.data && setData(res.data.data[0]);
      });
    }
  }, [props.ID]);

  return (
    <Grid container spacing={2}>
      <Grid item xs={2}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Users </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {data.users} </b>{' '}
          </div>
        </Paper>
      </Grid>

      <Grid item xs={2}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Policies </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {data.policy} </b>{' '}
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
            <b> {data.groups}</b>{' '}
          </div>
        </Paper>
      </Grid>

      <Grid item xs={2}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Secrets </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {data.secrets}</b>{' '}
          </div>
        </Paper>
      </Grid>

      <Grid item xs={2}>
        <Paper className={classes.paperTrans} elevation={0}>
          <div className={classes.aggHeaders}>
            {' '}
            <b> Privileges </b>
          </div>
          <div className={classes.aggHeadersBig}>
            {' '}
            <b> {data.privileges} </b>{' '}
          </div>
        </Paper>
      </Grid>
    </Grid>
  );
}

/// /////////////////////////////////////////
/// //// Proxy Status shows htt proxy details
/// /////////////////////////////////////////
function ProxyStatus(props: any) {
  const classes = useStyles();
  return (
    <div>
      {props.serviceDetail.serviceType === 'http' ? (
        <Grid container spacing={2} className={classes.gridPadding}>
          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography variant="h3">Proxy Setting</Typography>
            </Grid>
            <Grid item xs={4}>
              <Tooltip title="edit" placement="top">
                <IconButton
                  style={{ color: 'navy' }}
                  onClick={props.toggleProxyConfigDrawer('right', true)}
                >
                  <EditIcon />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <div className={classes.settingHeader}> Domain Name :</div>
              </Grid>
              <Grid item xs={6}>
                <div className={classes.proxyVals}> {props.hostname} </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <div className={classes.settingHeader}> Route Rule :</div>
              </Grid>
              <Grid item xs={6}>
                <div className={classes.proxyVals}> matching host header </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <div className={classes.settingHeader}> Pass host header : </div>
              </Grid>
              <Grid item xs={6}>
                <div className={classes.proxyVals}>
                  {' '}
                  {props.proxy.passHostHeader ? 'forwarding' : 'not forwarding'}{' '}
                </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <div className={classes.settingHeader}> Strict TLS validation :</div>
              </Grid>
              <Grid item xs={6}>
                <div className={classes.proxyVals}>
                  {' '}
                  {props.proxy.strictTLSValidation ? 'enabled' : 'disabled'}{' '}
                </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <div className={classes.settingHeader}> Upstream Server :</div>
              </Grid>
              <Grid item xs={6}>
                <div className={classes.proxyVals}> {props.proxy.upstreamServer} </div>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      ) : (
        <Grid container spacing={2} className={classes.gridPadding}>
          <Grid container spacing={2}>
            <Grid item xs={6}>
              <Typography variant="h3">Proxy Setting</Typography>
            </Grid>
            <Grid item xs={4}>
              <Tooltip title="edit" placement="top">
                <IconButton style={{ color: 'grey' }} disabled>
                  <EditIcon />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <div className={classes.settingHeader}>
                  {' '}
                  Proxy setting is not required for this service{' '}
                </div>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      )}
    </div>
  );
}

/// /////////////////////////////////////////
/// //// Cert Status shows CA or TLS certs
/// /////////////////////////////////////////
//  TODO @sshahcodes certi setting is not required for https for now.
function CertStatus(props: any) {
  const classes = useStyles();
  return (
    <Grid container spacing={2} className={classes.gridPadding}>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <Typography variant="h3">Certificates</Typography>
        </Grid>
        <Grid item xs={4}>
          <Tooltip title="edit" placement="top">
            <IconButton
              style={{ color: 'navy' }}
              onClick={props.toggleCertConfigDrawer('right', true)}
            >
              <EditIcon />
            </IconButton>
          </Tooltip>
        </Grid>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <div className={classes.settingHeader}>
              {' '}
              Configure certificate for this application{' '}
            </div>
          </Grid>
          {/* <Grid item xs={7}>
          <div className={classes.proxyVals}> app.trasa.io </div>
          </Grid> */}
        </Grid>
      </Grid>
    </Grid>
  );
}

/// ////////////////////////////////////////////////////
const DeleteServiceDlg = (props: any) => {
  const classes = useStyles();

  const [loader, setLoader] = useState(false);

  const DeleteApp = () => {
    setLoader(true);

    const deleteAppData = { ID: props.serviceID };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/delete`, deleteAppData)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success') {
          window.location.href = '/services';
        }
        props.closeDlg();
      })
      .catch((error) => {
        console.log('catched err: ', error);
      });
  };

  return (
    <div>
      <Dialog
        open={props.open}
        onClose={props.close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {' '}
          <div className={classes.Warning}> !!! WARNING !!! </div>
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Deleting service will revoke access to system. If this Service is deleted before
            uninstalling TRASA agent, you may be locked out from that system. Make sure you know
            what you are doing.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            id="deleteConfirmBtn"
            onClick={DeleteApp}
            className={classes.WarningButton}
            variant="contained"
          >
            Yes, Delete this app.
          </Button>
          <Button onClick={props.close} color="primary" variant="contained" autoFocus>
            No
          </Button>
          <br />
        </DialogActions>
        {loader ? <ProgressHOC /> : null}
      </Dialog>
    </div>
  );
};
