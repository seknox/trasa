import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import axios from 'axios';
import React, { useState } from 'react';
import ProgressHOC from '../../../utils/Components/Progressbar';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    // flexgrow: 1,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  textField: {
    // marginLeft: 100,
    paddingLeft: theme.spacing(1),
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    width: 500,
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
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
}));

export default function OrgAccountSetting(props: any) {
  const classes = useStyles();

  const [orgSetting, setOrgSetting] = useState({
    orgName: '',
    domain: '',
    primaryContact: '',
    timezone: '',
  });

  function hndlValueChange(e: any) {
    setOrgSetting({ ...orgSetting, [e.target.name]: e.target.value });
  }

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/org/detail`)
      .then((r) => {
        if (r.data.status === 'success') {
          setOrgSetting(r.data.data[0]);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const [reqStatus, setReqStatus] = useState(false);

  function submitSetting() {
    setReqStatus(true);

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/org/update`, orgSetting)
      .then(() => {
        setReqStatus(false);
      })
      .catch((error) => {
        if (error) {
          setReqStatus(false);
          console.log(error);

          // commented out for local debug window.location.href = '/login'
        } else {
          setReqStatus(false);
          // Something happened in setting up the request that triggered an Error
          console.log('Error', error.message);
        }
      });
  }

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
              <b>Organization Account Setting</b>
            </Typography>
          </ExpansionPanelSummary>
        </Grid>

        <Grid item xs={12} sm={12}>
          <ExpansionPanelDetails>
            <Grid container spacing={2} direction="row" justify="center">
              <Grid item xs={12}>
                <Divider light />
              </Grid>

              <Grid item xs={12} sm={9}>
                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={3}>
                    <Typography variant="h4">Account Name : </Typography>
                  </Grid>
                  <Grid item xs={9}>
                    <TextField
                      fullWidth
                      onChange={hndlValueChange}
                      name="orgName"
                      variant="outlined"
                      size="small"
                      value={orgSetting.orgName}
                      // defaultValue={props.emailSetting.serverAddress}
                    />
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={3}>
                    <Typography variant="h4">Primary Contact : </Typography>
                  </Grid>
                  <Grid item xs={9}>
                    <TextField
                      fullWidth
                      onChange={hndlValueChange}
                      name="primaryContact"
                      value={orgSetting.primaryContact}
                      // defaultValue={props.emailSetting.serverPort}
                      // className={ classes.textField }
                      variant="outlined"
                      size="small"
                    />
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={3}>
                    <Typography variant="h4">Timezone : </Typography>
                  </Grid>
                  <Grid item xs={9}>
                    <TextField
                      fullWidth
                      onChange={hndlValueChange}
                      name="timezone"
                      value={orgSetting.timezone}
                      // defaultValue={props.emailSetting.authKey}
                      // className={ classes.textField }
                      variant="outlined"
                      size="small"
                    />
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </ExpansionPanelDetails>
        </Grid>
        {reqStatus ? (
          <div>
            <ProgressHOC /> <br />
          </div>
        ) : null}
        <ExpansionPanelActions>
          <Button variant={reqStatus ? 'text' : 'contained'} onClick={submitSetting}>
            {' '}
            Save{' '}
          </Button>
        </ExpansionPanelActions>
      </ExpansionPanel>
      <br /> <br /> <br />
    </div>
  );
}
