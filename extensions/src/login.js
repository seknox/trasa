//import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import InputAdornment from '@material-ui/core/InputAdornment';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Lock from '@material-ui/icons/Lock';
import World from '@material-ui/icons/Public';
import axios from 'axios';
import React, { useState, useEffect } from 'react';
import useStyles from './styles';
import ProgressHOC from './utils';

var browser = require('webextension-polyfill');

export default function LoginPage(props) {
  const [data, setData] = useState({
    trasaID: '',
    password: '',
    tfaMethod: 'u2f',
    totpCode: '',
    trasacore: '',
  });
  const [loading, setLoader] = useState(false);
  const [failed, setFailed] = useState(false);
  const [error, setError] = useState('');
  const [respData, setRespData] = useState({});
  const [showOrgs, setShowOrgs] = useState(false);
  const [orgs, setOrgs] = useState([
    { orgID: 'testID', orgName: 'testName' },
    { orgID: 'bestID', orgName: 'bestName' },
  ]);
  const classes = useStyles();

  function handleChange(event) {
    setData({ ...data, [event.target.name]: event.target.value });
  }

  useEffect(() => {
    if (respData.exToken) {
      sendBsrDetails();
    }
  }, [respData]);

  async function handleSubmit() {
    setLoader(true);

    data['UA'] = navigator.userAgent;

    function getChromeVersion() {
      var raw = navigator.userAgent.match(/Chrom(e|ium)\/([0-9]+)\./);
      return raw ? raw[2] : false;
    }

    let browserName = 'chrome';
    let browserVersion = '58';
    let build = '0';
    if (getChromeVersion() !== false) {
      browserName = 'chrome';
      browserVersion = getChromeVersion();
    } else {
      let brsr = await browser.runtime.getBrowserInfo();
      browserName = brsr.name;
      browserVersion = brsr.version;
      build = brsr.buildID;
    }

    let hostName = data['trasacore'];
    if (data['totpCode'].length === 0) {
      data['tfaMethod'] = 'u2f';
    } else {
      data['tfaMethod'] = 'totp';
    }

    let brsrDetails = {
      name: browserName,
      version: browserVersion,
      userAgent: navigator.userAgent,
      build: build,
    };

    let exts = await browser.management.getAll();
    let enrolDeviceData = {
      hostName: hostName,
      authData: data,
      deviceBrowser: brsrDetails,
      browserExtensions: exts,
    };

    // native message trasaExtNative to register device.
    let message = { intent: 'enrolDevice', data: enrolDeviceData };

    // chrome.runtime.sendNativeMessage('trasaextnative',
    // { text: "Hello" },
    // function(response) {
    //   console.log("Received " + response);
    // });

    // if (chrome.runtime.lastError) {
    //   console.error(chrome.runtime.lastError);
    // }

    var sending = browser.runtime.sendNativeMessage('trasaextnative', message);
    let resp = await sending.then(rcvd, onError);

    if (resp.status === true) {
      let v = JSON.parse(resp.data);
      browser.storage.local.set({
        hosts: v.hosts,
        trasaCore: hostName,
        extID: v.extID,
        loggedIn: true,
        trasaDACom: v.trasaDACom,
      });

      var sessionStore = new Object();
      sessionStore['testhost'] = 'testsession';
      browser.storage.local.set({ sessionStore: sessionStore });
      props.setLoginTrue();

      setLoader(false);
    } else {
      setError(resp.intent);
      setFailed(true);
      setLoader(false);
    }
  }

  function onError(e) {
    if (browser.runtime.lastError) {
      console.error(browser.runtime.lastError);
    } else {
      console.log(e);
    }
  }

  function rcvd(msg) {
    return msg;
  }

  return (
    <div className={classes.paper}>
      <Grid container spacing={2}>
        {showOrgs ? (
          <OrgSelect
            orgs={orgs}
            submitLoginRequest={handleSubmit}
            loading={loading}
            error={error}
            failed={failed}
          />
        ) : (
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Typography className={classes.title}>Register your browser </Typography>
              <br />
              {loading ? (
                <div>
                  {' '}
                  <ProgressHOC /> <br />
                </div>
              ) : (
                ''
              )}

              {failed ? (
                <div>
                  <div className={classes.errorText}>
                    {' '}
                    {'We encountered some error. contact your administrator'}
                  </div>{' '}
                  <br />
                </div>
              ) : (
                ''
              )}
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Enter TRASA hostname"
                variant="outlined"
                autoFocus
                onChange={handleChange}
                name="trasacore"
                value={data.trasacore}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <World color="primary" />
                    </InputAdornment>
                  ),
                }}
              ></TextField>
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Your Email or Username"
                variant="outlined"
                onChange={handleChange}
                name="trasaID"
                value={data.trasaID}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <AccountCircle color="primary" />
                    </InputAdornment>
                  ),
                }}
              ></TextField>
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="TOTP code (leave empty for push U2F)"
                variant="outlined"
                onChange={handleChange}
                name="totpCode"
                value={data.totpCode}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Lock color="primary" />
                    </InputAdornment>
                  ),
                }}
              ></TextField>
            </Grid>

            <Grid item xs={12}>
              {loading ? (
                <Button variant="contained" disabled className={classes.button}>
                  Authorize
                </Button>
              ) : (
                <Button
                  variant="contained"
                  className={classes.button}
                  onClick={() => {
                    handleSubmit('');
                  }}
                >
                  Authorize
                </Button>
              )}
            </Grid>
          </Grid>
        )}
      </Grid>
    </div>
  );
}

function OrgSelect(props) {
  const classes = useStyles();
  return (
    <Grid container spacing={2} direction="column" alignItems="center" justify="center">
      <Grid item xs={12}>
        <Typography className={classes.title}>Select your organization </Typography>
        <br />
        {props.loading ? (
          <div>
            {' '}
            <ProgressHOC /> <br />
          </div>
        ) : (
          ''
        )}

        {props.failed ? (
          <div>
            <div className={classes.errorText}>
              {' '}
              {props.error.length > 0
                ? error
                : 'We encountered some error. contact your administrator'}
            </div>{' '}
            <br />
          </div>
        ) : (
          ''
        )}
      </Grid>

      {/* <Grid item xs={12}>
        <Grid container spacing={2}> */}
      {props.orgs.map((org) => (
        <Grid item xs={12}>
          <Button
            variant="contained"
            className={classes.button}
            size="small"
            onClick={() => {
              props.submitLoginRequest(org.ID);
            }}
          >
            {org.orgName}
          </Button>
        </Grid>
      ))}
      {/* </Grid>
        </Grid> */}
    </Grid>
  );
}
