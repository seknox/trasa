import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Divider from '@material-ui/core/Divider';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import axios from 'axios';
import React from 'react';
import Constants from '../../../Constants';
import { ConfirmDialogue, NotifDialogue } from '../../../utils/Components/Confirms';
import ProgressHOC from '../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  // card
  card: {
    minWidth: 275,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    marginBottom: 16,
    color: '#311B92', // theme.palette.text.secondary,
  },
  pos: {
    marginBottom: 12,
    color: theme.palette.text.secondary,
  },

  // form
  formControl: {
    margin: theme.spacing(1),
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
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  selectCustom: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 31,
    // marginTop: 5,
    // padding: '10px 100px',
    width: 100,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  smallTextBox: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    // height: 31,
    // marginTop: 5,
    // padding: '10px 100px',
    width: 30,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

export default function PasswordSetting() {
  const classes = useStyles();
  // const [open, setopen] = React.useState(false);
  const [passPolicy, setpassPolicy] = React.useState({
    minimumChars: 8,
    expiry: '',
    enforceStrongPass: false,
  });
  const [passPolicyStatus, setpassPolicyStatus] = React.useState(false);
  const [loader, setLoader] = React.useState(false);

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/all`)
      .then((response) => {
        if (response.status === 403) {
          window.location.href = '/login';
        }
        //  this.setState({appData: response.data})
        const resp = response.data.data[0];
        const policy = JSON.parse(response.data.data[0].passPolicy.settingValue);
        setpassPolicy(policy);
        setpassPolicyStatus(resp.passPolicy.status);
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  const changePolicy = (name: string) => (event: any) => {
    switch (name) {
      case 'minimumChars':
        setpassPolicy({ ...passPolicy, minimumChars: event.target.value });
        // passPolicy[name] = event.target.value;
        break;
      case 'expiry':
        setpassPolicy({ ...passPolicy, expiry: event.target.value });
        break;
      case 'enforceStrongPass':
        setpassPolicy({ ...passPolicy, enforceStrongPass: event.target.value });
        break;
      default:
        break;
    }
  };

  const submitPasswordPolicy = () => {
    setLoader(true);

    const policyState = passPolicy;
    // change this to boolean. Stupid dom casts boolean to string when used as value in checkbox
    // policyState.enforceStrongPass = (policyState.enforceStrongPass === 'true')
    // change this integer
    // TODO @sshahcodes check if minimum chars string cause problem in backend
    const req = { policy: policyState, enable: !passPolicyStatus };
    setpassPolicyStatus(!passPolicyStatus);

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/passwordpolicy/update`, req)
      .then((response) => {
        if (response.data.status === 'success') {
          setLoader(false);
        }
      })
      .catch((error) => {
        console.log(error);
        setLoader(false);
      });
  };

  return (
    <div className={classes.root}>
      <ExpansionPanel>
        <Grid item xs={12} sm={12}>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1a-content"
            id="panel1a-header"
          >
            <Typography component="h4" variant="h3">
              <b>Password security Policy</b>
            </Typography>
          </ExpansionPanelSummary>
        </Grid>
        <Grid item xs={12} sm={12}>
          <ExpansionPanelDetails>
            <Settings
              passPolicy={passPolicy}
              updatePolicySetting={submitPasswordPolicy}
              passPolicyStatus={passPolicyStatus}
              changePolicy={changePolicy}
            />
          </ExpansionPanelDetails>
          {loader && <ProgressHOC />}
        </Grid>
      </ExpansionPanel>
      <br /> <br /> <br />
    </div>
  );
}

function Settings(props: any) {
  const classes = useStyles();
  const { passPolicy, passPolicyStatus } = props;
  const [confirmDlgState, setconfirmDlgState] = React.useState(false);
  const [notifDlgState, setnotifDlgState] = React.useState(false);
  const [loader, setLoader] = React.useState(false);

  const enforceNow = () => {
    setLoader(true);
    setconfirmDlgState(!confirmDlgState);

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/globalsettings/passwordpolicynow/enforce`)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success') {
          setnotifDlgState(!notifDlgState);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <Grid container spacing={2} direction="row" justify="center">
      <Divider light />
      <br />

      <Grid item xs={12} sm={4}>
        <div className={classes.settingHeader}>
          {' '}
          Enforce Strong Password :{' '}
          <span>
            <Checkbox
              // defaultChecked
              checked={!!passPolicy.enforceStrongPass}
              // value={"enforceStrongPass"}
              value={passPolicy.enforceStrongPass}
              name={passPolicy.enforceStrongPass}
              color="primary"
              onChange={props.changePolicy('enforceStrongPass')}
              inputProps={{ 'aria-label': 'secondary checkbox' }}
            />
          </span>{' '}
        </div>
      </Grid>

      <Grid item xs={12} sm={4}>
        <div className={classes.settingHeader}>
          Minimum Characters :{' '}
          <span>
            <TextField
              onChange={props.changePolicy('minimumChars')}
              name="minimumChars"
              value={passPolicy.minimumChars}
              defaultValue={passPolicy.minimumChars}
              InputProps={{
                disableUnderline: true,
                autoComplete: 'off',
                classes: {
                  input: classes.smallTextBox,
                },
              }}
            />
          </span>
        </div>
      </Grid>

      <Grid item xs={12} sm={4}>
        <div className={classes.settingHeader}>
          {' '}
          Expiration :{' '}
          <span>
            <select
              className={classes.selectCustom}
              name="expiry"
              onChange={props.changePolicy('expiry')}
            >
              {expiry.map((name, i) => (
                <option key={i} value={name} selected={passPolicy.expiry === name}>
                  {name}
                </option>
              ))}
            </select>
          </span>{' '}
        </div>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2} direction="column" alignItems="center" justify="center">
          <Divider variant="middle" style={{ color: 'navy' }} />
          <Grid item xs={6}>
            <FormControl fullWidth>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Switch
                      checked={!!passPolicyStatus}
                      onChange={props.updatePolicySetting}
                      name="settingState"
                      defaultValue={passPolicyStatus}
                      value="passPolicyStatus"
                      color="primary"
                    />
                  }
                  label={
                    passPolicyStatus ? (
                      <div className={classes.settingSHeader}> policy is enabled </div>
                    ) : (
                      <div className={classes.settingSHeader}>policy is disabled </div>
                    )
                  }
                />
              </FormGroup>
            </FormControl>
          </Grid>
          <Grid item xs={6}>
            <Button
              variant="contained"
              color="secondary"
              disabled={!!passPolicyStatus}
              onClick={passPolicyStatus ? () => setconfirmDlgState(!confirmDlgState) : () => 0}
            >
              Enforce password change now
            </Button>
          </Grid>
        </Grid>
      </Grid>

      <NotifDialogue
        open={notifDlgState}
        close={() => setnotifDlgState(!notifDlgState)}
        confirmMessage="Policy is enforced. Users will be required to change password in next login attempt."
      />
      <ConfirmDialogue
        open={confirmDlgState}
        close={() => setconfirmDlgState(!confirmDlgState)}
        confirmFunc={enforceNow}
        confirmMessage="Enforcing will require every users to change password in next login."
      />

      {loader ? <ProgressHOC /> : null}
    </Grid>
  );
}

// const expiry = [{"Never": "never"}, {"30 days": "30"}, {"60 days": "60"}, {"90 days": "90"}, {"180 days": "180"}, {"365 days": "360"}]
const expiry = ['Never', '30 days', '60 days', '90 days', '180 days', '365 days'];
