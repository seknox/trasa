import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import Stepper from '@material-ui/core/Stepper';
import { makeStyles, Theme, withStyles } from '@material-ui/core/styles';
import Tab from '@material-ui/core/Tab';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import Send from '@material-ui/icons/Send';
// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
// import mixpanel from 'mixpanel-browser';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import ProgressBar from '../../../utils/Components/Progressbar';
import DevicePolicy from './DevicePolicy/DevicePolicy';
import TrasaUAC from './TrasUAC';
import { DevicePolicyProps } from './index';

const StyledTableCell = withStyles((theme) => ({
  head: {
    backgroundColor: '#ffffff', // '#e0e0e0',
    color: 'black',
    height: 5,
    fontSize: '14px',
    fontWeight: 700,
    fontFamily: 'Open Sans, Rajdhani',
  },
  body: {
    fontSize: '13px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}))(TableCell);

const useStyles = makeStyles((theme: Theme) => ({
  expiry: {
    marginLeft: 10,
    width: 100,
  },
  useradd: {
    flexGrow: 1,
  },
  dialogue: {
    // minHeight: 780,
    // width: '500%',
    // maxWidth: '500%',
  },
  formControl: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 31,
    marginTop: 5,
    // padding: '10px 100px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  stepperContent: {
    flexGrow: 1,
    marginTop: 20,
    marginBottom: 20,
  },
  chips: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  chip: {
    margin: theme.spacing(1) / 4,
  },

  backButton: {
    marginRight: theme.spacing(1),
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },

  timePickerTextField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 100,
  },

  timePickerContainer: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  button: {
    margin: theme.spacing(1),
    marginTop: '5%',
  },
  rightIcon: {
    marginRight: theme.spacing(1),
  },
  snackbar: {
    position: 'absolute',
    background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 48,
    padding: '0 30px',
    boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .30)',
  },
  checkInbox: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  stepperButtonLeft: {
    // marginTop:'15%',
    marginLeft: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },

  stepperButtonRight: {
    marginTop: '15%',
    marginRight: '10%',
  },
  newPaper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),

    color: theme.palette.text.secondary,
  },
  table: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  manageUser: {
    padding: theme.spacing(2),
  },
  menuHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    padding: '10px 12px',
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
    fontSize: 16,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
  },
  policyName: {
    marginLeft: '20%',
  },
  toolTip: {
    padding: 5,
    // maxWidth: 220,
    fontSize: 16,
    // fontFamily: 'Open Sans, Rajdhani',
    border: '1px solid #dadde9',
  },
  selectCustom: {
    fontSize: 15,
    fontFamily: 'Open Sans, Rajdhani',
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 31,
    marginTop: 5,
    // padding: '10px 100px',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  tabRoot: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
    display: 'flex',
    // height: 224,
  },
  tabs: {
    borderRight: `1px solid ${theme.palette.divider}`,
  },
  stepperRoot: {},
}));

const CTooltip = withStyles((theme) => ({
  tooltip: {
    // backgroundColor: '#f5f5f9',
    color: 'rgba(0, 0, 0, 0.87)',
    maxWidth: 420,
    fontSize: 16,
    border: '1px solid #dadde9',
  },
}))(Tooltip);

function getSteps() {
  return ['Policy Name', 'Add Permissions', 'Review and Create'];
}

type createPolicyProps = {
  open: boolean;
  updateData: singlePolicytype;
  update: boolean;
  handleClose: () => void;
};

type singlePolicytype = {
  policyID: string;
  policyName: string;
  dayAndTime: any;
  tfaRequired: boolean;
  recordSession: boolean;
  fileTransfer: boolean;
  ipSource: string;
  expiry: string;
  devicePolicy: any;
};

type devicePolicy = {};

export default function CreatePolicy(props: createPolicyProps) {
  const [policyName, setPolicyName] = useState(''); // useState(props.update? props.updateData.policyName: '')
  const [dayAndTime, setDayAndTime] = useState<any>([]); // useState(props.update? props.updateData.dayAndTime: [])
  const [tfaRequired, settfaRequired] = useState(false); // useState(props.update? props.updateData.tfaRequired: false)
  const [recordSession, setRecordSession] = useState(false); // useState(props.update? props.updateData.recordSession: false)
  const [fileTransfer, setFileTransfer] = useState(false); // useState(props.update? props.updateData.fileTransfer: false)
  const [ipSource, setIPSource] = useState('0.0.0.0/0'); // useState(props.update? props.updateData.ipSource : '0.0.0.0/0')
  const [expiry, setExpiry] = useState('2021-03-01'); // useState(props.update? props.updateData.expiry : '2021-03-01')
  const [activeStep, setActiveStep] = useState(0);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (props.update === true) {
      setPolicyName(props.updateData.policyName);
      setDayAndTime(props.updateData.dayAndTime);
      settfaRequired(props.updateData.tfaRequired);
      setRecordSession(props.updateData.recordSession);
      setFileTransfer(props.updateData.fileTransfer);
      setIPSource(props.updateData.ipSource);
      setExpiry(props.updateData.expiry);
    }
    return () => {
      setPolicyName('');
      setDayAndTime([]);
      settfaRequired(false);
      setRecordSession(false);
      setFileTransfer(false);
      setIPSource('0.0.0.0/0');
      setExpiry('2021-03-01');
    };
  }, [props.update]);

  function closer() {
    setDayAndTime([{ days: [], fromTime: '', toTime: '' }]);
    props.handleClose();
  }

  const handlePolicyNameChange = (name: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setPolicyName(event.target.value);
  };

  function handle2FAChange() {
    settfaRequired(!tfaRequired);
  }

  function handleSessionRecordChange() {
    setRecordSession(!recordSession);
  }

  function handleFileTransferChange() {
    setFileTransfer(!fileTransfer);
  }

  function handleIPSourceChange(e: React.ChangeEvent<HTMLInputElement>) {
    setIPSource(e.target.value);
  }

  function handleExpiryChange(e: React.ChangeEvent<HTMLInputElement>) {
    setExpiry(e.target.value);
  }

  function updateDayAndTime(days: any, fromTime: string, toTime: string, index: number) {
    setDayAndTime(dayAndTime.concat({ days, fromTime, toTime }));
  }

  const deleteDayAndTime = (index: number) => {
    dayAndTime.splice(index, 1);
    setDayAndTime([...dayAndTime]);
  };

  /// /////////////////////////////
  /// /////   DevicePolicy
  /// /////////////////////////////
  const [dhBlocking, setDHBlocking] = useState({
    blockAutologinEnabled: false,
    blockTfaNotConfigured: false,
    blockIdleScreenLockDisabled: false,
    blockRemoteLoginEnabled: false,
    blockJailBroken: false,
    blockDebuggingEnabled: false,
    blockEmulated: false,
    blockEncryptionNotSet: false,
    blockOpenWifiConn: false,
    blockUntrustedDevices: false,
    blockAntivirusDisabled: false,
    blockFirewallDisabled: false,
  });

  const handleDHBlockingChange = (name: any) => (event: any) => {
    setDHBlocking({ ...dhBlocking, [name]: event.target.checked });
  };

  useEffect(() => {
    setDHBlocking(props.updateData.devicePolicy);
  }, [props.update]);

  /// /////////////////////////////////////////
  /// ////// Stepper component functions
  /// /////////////////////////////////////////

  const handleNext = () => {
    console.log('next: ', policyName);
    if (activeStep === 0 && policyName.length === 0 && props.update) {
      setPolicyName(props.updateData.policyName);
    }
    const counter = activeStep + 1;
    if (counter === 1 && dayAndTime.length === 0 && props.update) {
      setDayAndTime(props.updateData.dayAndTime);
    }

    setActiveStep(counter);
  };

  const handleBack = () => {
    setActiveStep(activeStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  const SubmitPolicy = () => {
    // mixpanel.track('control-policies-createpolicy');
    setLoading(true);
    const devicePolicy = { ...dhBlocking };
    const basicPolicy = {
      policyID: props.update ? props.updateData.policyID : '',
      policyName,
      dayAndTime,
      tfaRequired,
      recordSession,
      fileTransfer,
      ipSource,
      devicePolicy,
      expiry,
    };

    const url = props.update
      ? `${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/update`
      : `${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/create`;

    axios.post(url, basicPolicy).then((r) => {
      setLoading(false);
      if (r.data.status === 'success') {
        window.location.reload();
        closer();
      } else {
        console.log('error');
      }
      console.log(r.data);
    });
  };

  const getStepContent = (stepIndex: number) => {
    switch (stepIndex) {
      case 0:
        return (
          <PolicyName
            handlePolicyNameChange={handlePolicyNameChange}
            update={props.update}
            policyName={props.updateData.policyName}
          />
        );
      case 1:
        return (
          <PolicyTab
            dayAndTime={dayAndTime}
            deleteDayAndTime={deleteDayAndTime}
            updateDayAndTime={updateDayAndTime}
            tfaRequired={tfaRequired}
            handle2FAChange={handle2FAChange}
            recordSession={recordSession}
            handleSessionRecordChange={handleSessionRecordChange}
            fileTransfer={fileTransfer}
            handleFileTransferChange={handleFileTransferChange}
            ipSource={ipSource}
            handleIPSourceChange={handleIPSourceChange}
            expiry={expiry}
            handleExpiryChange={handleExpiryChange}
            /// ///////////////////////////////
            dhBlocking={dhBlocking}
            handleDHBlockingChange={handleDHBlockingChange}
          />
        );
      default:
        return (
          <Grid container spacing={2} alignItems="flex-end" direction="row" justify="center">
            <ReviewAccess
              policyName={policyName}
              dayAndTime={dayAndTime}
              tfaRequired={tfaRequired}
              recordSession={recordSession}
              fileTransfer={fileTransfer}
              ipSource={ipSource}
              expiry={expiry}
              // Device hygiene below
              devicePolicy={dhBlocking}
            />
            <Button
              id="submitBtn"
              className={classes.button}
              variant="contained"
              color="secondary"
              onClick={SubmitPolicy}
            >
              Submit
              <Send className={classes.rightIcon} />
            </Button>
          </Grid>
        );
    }
  };

  const classes = useStyles();
  const steps = getSteps();
  return (
    <div>
      <Dialog
        classes={{
          root: classes.dialogue,
        }}
        fullWidth
        maxWidth="lg"
        open={props.open}
        onClose={closer}
        aria-labelledby="form-dialog-title"
        disableBackdropClick
        disableEscapeKeyDown
      >
        <DialogContent>
          <Typography variant="h2"> Create New Policy </Typography>
          <DialogContentText>
            Policies created here will be avaialable when assigning user or group to Service or
            Servicegroup.
          </DialogContentText>

          <Divider light />
          <Grid container>
            <div className={classes.stepperContent}>{getStepContent(activeStep)}</div>
          </Grid>

          <Divider light />
        </DialogContent>
        <div className={classes.stepperRoot}>
          <Grid container>
            <Grid item xs={1} sm={1} md={1}>
              <Button
                className={classes.stepperButtonLeft}
                disabled={activeStep === 0}
                onClick={handleBack}
                variant="contained"
                color="secondary"
              >
                Back
              </Button>
            </Grid>

            <Grid item xs={9} sm={9} md={9}>
              {activeStep === steps.length ? (
                <div>
                  <Typography className={classes.instructions}>
                    All steps completed - you're finished
                  </Typography>
                  <Button onClick={handleReset}>Reset</Button>
                </div>
              ) : (
                <Stepper activeStep={activeStep} alternativeLabel>
                  {steps.map((label) => {
                    return (
                      <Step key={label}>
                        <StepLabel>{label}</StepLabel>
                      </Step>
                    );
                  })}
                </Stepper>
              )}
              {loading ? <ProgressBar /> : ''}
            </Grid>

            <Grid item xs={1} sm={1} md={1}>
              <Button
                id="nextBtn"
                className={classes.stepperButtonRight}
                variant="contained"
                color="secondary"
                onClick={handleNext}
              >
                {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
              </Button>
            </Grid>

            <Grid item xs={1} sm={1} md={1}>
              <Button
                className={classes.stepperButtonRight}
                variant="contained"
                color="secondary"
                onClick={() => {
                  closer();
                  setActiveStep(0);
                }}
              >
                Close
              </Button>
            </Grid>
          </Grid>
        </div>
      </Dialog>
    </div>
  );
}

/*  ////////////////////////////////////////////////////////////////////////////   */

interface styledTabsProps {
  value: number;
  centered: boolean;
  onChange: (event: React.ChangeEvent<{}>, newValue: number) => void;
}

interface styledTabProps {
  label: string;
}

const StyledTabs = withStyles({
  root: {
    maxHeight: 10,
    borderBottom: '1px solid #e8e8e8',
  },
})((props: styledTabsProps) => <Tabs {...props} TabIndicatorProps={{ children: <div /> }} />);

const StyledTab = withStyles((theme) => ({
  root: {
    textTransform: 'none',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
    color: 'black',
    marginRight: theme.spacing(1),
    '&:focus': {
      opacity: 1,
    },
  },
}))((props: styledTabProps) => <Tab disableRipple {...props} />);

// PolicyTab is main tab which wraps basic TRASA policy configs and device hygiene policy config
function PolicyTab(props: any) {
  const [tabValue, setTabValue] = React.useState(0);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setTabValue(newValue);
  };
  return (
    <div>
      <StyledTabs
        value={tabValue}
        onChange={handleChange}
        centered
        aria-label="styled tabs example"
      >
        <StyledTab label="Basic Policy" />
        <StyledTab label="Device Hygiene (Beta)" />
      </StyledTabs>
      <TabPanel value={tabValue} index={0}>
        <TrasaUAC
          dayAndTime={props.dayAndTime}
          deleteDayAndTime={props.deleteDayAndTime}
          updateDayAndTime={props.updateDayAndTime}
          tfaRequired={props.tfaRequired}
          handle2FAChange={props.handle2FAChange}
          recordSession={props.recordSession}
          handleSessionRecordChange={props.handleSessionRecordChange}
          fileTransfer={props.fileTransfer}
          handleFileTransferChange={props.handleFileTransferChange}
          ipSource={props.ipSource}
          handleIPSourceChange={props.handleIPSourceChange}
          expiry={props.expiry}
          handleExpiryChange={props.handleExpiryChange}
        />
      </TabPanel>
      <TabPanel value={tabValue} index={1}>
        <DevicePolicy
          dhBlocking={props.dhBlocking}
          handleDHBlockingChange={props.handleDHBlockingChange}
        />
      </TabPanel>
    </div>
  );
}

/// ////////////////////////////////////////////////////////////////////////////////

function PolicyName(props: any) {
  const classes = useStyles();
  return (
    <div className={classes.stepperRoot}>
      <Grid container spacing={2}>
        <br />

        <Grid item xs={4} sm={4} md={4}>
          <Typography variant="h3">Name: </Typography>
        </Grid>
        <Grid item xs={8} sm={8} md={8}>
          <TextField
            fullWidth
            onChange={props.handlePolicyNameChange()}
            name="policyName"
            defaultValue={props.policyName}
            // defaultValue = {props.update? props.policyName: ''}
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
      </Grid>
    </div>
  );
}

/// //////////////////////////////////////////////////////////////////////////////////

function TabPanel(props: any) {
  const { children, value, index, ...other } = props;

  return (
    <Typography
      component="div"
      role="tabpanel"
      hidden={value !== index}
      id={`vertical-tabpanel-${index}`}
      aria-labelledby={`vertical-tab-${index}`}
      {...other}
    >
      {value === index && <Box p={3}>{children}</Box>}
    </Typography>
  );
}

type reviewAccessProps = {
  policyName: string;
  dayAndTime: any;
  tfaRequired: boolean;
  recordSession: boolean;
  fileTransfer: boolean;
  ipSource: string;
  expiry: string;
  devicePolicy: DevicePolicyProps;
};
/// //////////////////////////////////////////////////////////////////////////////////////////////
/// //////////////////////////////////////////////////////////////////////////////////////////////
export function ReviewAccess(props: reviewAccessProps) {
  const classes = useStyles();
  const [tabValue, setTabValue] = React.useState(0);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setTabValue(newValue);
  };

  useEffect(() => {
    console.log('device Policy: ', props.devicePolicy);
  }, [props]);

  return (
    <div>
      <StyledTabs
        value={tabValue}
        onChange={handleChange}
        centered
        aria-label="styled tabs example"
      >
        <StyledTab label="Basic Policy" />
        <StyledTab label="Device Hygiene (Beta)" />
      </StyledTabs>
      {/* Basic Policy */}
      <TabPanel value={tabValue} index={0}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">Policy Name:</Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography className={classes.textFieldInputBig} variant="h4">
                  {props.policyName}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">2FA: </Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.tfaRequired ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">Session Recording: </Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.recordSession ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">File Transfers: </Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.fileTransfer ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">IP Source: </Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.ipSource}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                {/* {JSON.stringify(props.dayAndTime)}  */}
                <Typography variant="h3">dayAndTime:</Typography>
              </Grid>
              <Grid item xs={10}>
                <Table className={classes.table}>
                  <TableHead>
                    <TableRow>
                      <StyledTableCell>SN</StyledTableCell>
                      <StyledTableCell>Days</StyledTableCell>
                      <StyledTableCell>From Time</StyledTableCell>
                      <StyledTableCell>To Time</StyledTableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {props.dayAndTime.map((perm: any, index: number) => (
                      <TableRow key={index}>
                        <StyledTableCell>{index}</StyledTableCell>
                        <StyledTableCell component="th" scope="row">
                          {`${perm.days} , `}
                        </StyledTableCell>
                        <StyledTableCell>{perm.fromTime}</StyledTableCell>
                        <StyledTableCell>{perm.toTime}</StyledTableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>

                <br />
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={2}>
                <Typography variant="h3">Policy Expiry: </Typography>
              </Grid>
              <Grid item xs={10}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.expiry}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </TabPanel>

      {/* Device Hygiene */}
      <TabPanel value={tabValue} index={1}>
        <Grid container spacing={2}>
          <Typography variant="h4">All policies are blocking*</Typography>
          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Untrusted devices:</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography className={classes.textFieldInputBig} variant="h4">
                  {props.devicePolicy.blockUntrustedDevices ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Autologin Enabled:</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography className={classes.textFieldInputBig} variant="h4">
                  {props.devicePolicy.blockAutologinEnabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Idle screen lock disabled: </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockIdleScreenLockDisabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Remote login Enabled: </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockRemoteLoginEnabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Jailbroken device (Mobile device): </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockJailBroken ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Debugging enabled (Mobile device):</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockDebuggingEnabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Emulated device (Mobile device):</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockEmulated ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Disk not encrypted (Workstation): </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockEncryptionNotSet ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Firewall disabled: </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockFirewallDisabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h3">Antivirus disabled (windows only) </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h4" className={classes.textFieldInputBig}>
                  {props.devicePolicy.blockAntivirusDisabled ? 'enabled' : 'disabled'}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </TabPanel>
    </div>
  );
}
