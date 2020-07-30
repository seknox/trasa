import Divider from '@material-ui/core/Divider';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import ProgressHOC from '../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
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

export default function SSHCertSetting(props: any) {
  const [mandatoryCert, setMandatoryCert] = useState(false);

  const [reqStatus, setReqStatus] = useState(false);
  const classes = useStyles();

  function handleConfigChange(e: any) {
    setMandatoryCert(e.target.checked);
    submitSetting(e.target.checked);
  }

  function submitSetting(val: any) {
    setReqStatus(true);

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/sshcert/update`, {
        mandatoryCert: val,
      })
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

  useEffect(() => {
    setMandatoryCert(props.status);
  }, [props.status]);

  return (
    <div className={classes.root}>
      {/* <Paper className={classes.paper}>  */}
      {/* <Grid container spacing={2} direction="row"  justify="center"> */}
      <ExpansionPanel>
        <Grid item xs={12} sm={12}>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1a-content"
            id="panel1a-header"
          >
            <Typography component="h4" variant="h3">
              <b>TRASA SSH Auth </b>
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
                  <Grid item xs={6}>
                    <Typography variant="h4">Mandatory Certificate Auth : </Typography>
                  </Grid>
                  <Grid item xs={3}>
                    <FormControl fullWidth>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={!!mandatoryCert}
                              onChange={handleConfigChange}
                              name="mandatoryCert"
                              // defaultValue={mandatoryCert}
                              value="mandatoryCert"
                              color="primary"
                            />
                          }
                          label={mandatoryCert ? <div>enabled </div> : <div>disabled </div>}
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>
              </Grid>
              {/* </Grid> */}
            </Grid>
          </ExpansionPanelDetails>
        </Grid>
        {reqStatus ? (
          <div>
            <ProgressHOC /> <br />
          </div>
        ) : null}
        <ExpansionPanelActions />
      </ExpansionPanel>
      {/* </Grid> */}
      {/* </Paper> */}
      <br /> <br /> <br />
    </div>
  );
}
