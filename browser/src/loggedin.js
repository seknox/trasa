import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Shield from '@material-ui/icons/Security';
import axios from 'axios';
import React, { useState } from 'react';
import useStyles from './styles';
import ProgressHOC from './utils';

function LoggedIn() {
  const [loading, setLoader] = useState(false);
  const [failed, setFailed] = useState(false);
  const [error, setError] = useState('');

  function getExt() {
    setLoader(true);
    browser.storage.local.get((data) => {
      let hostName = data.trasaCore;
      let extID = data.extID;
      let req = { extID: extID };

      let url = hostName + '/auth/device/ext/sync';
      // console.log ('url: ', url)
      axios
        .post(url, req)
        .then((response) => {
          if (response.data.status === 'success') {
            // console.log(response.data)
            // console.log('sync: ', response.data.data[0])
            var sessionStore = new Object();
            sessionStore['testhost'] = 'testsession';

            browser.storage.local.set({
              hosts: response.data.data[0].hosts,
              wsPath: response.data.data[0].wsPath,
              trasaCore: hostName,
              rootDomain: response.data.data[0].rootDomain,
              ssoDomain: response.data.data[0].ssoDomain,
              trasaDACom: response.data.data[0],
              sessionStore: sessionStore,
              loggedIn: true,
            });
            setLoader(false);
          } else {
            setFailed(true);
            setError(response.data.reason);
            setLoader(false);
          }
        })
        .catch((error) => {
          console.log(error);
          setError('could not connect to trasa server');
          setLoader(false);
        });
    });
  }

  function deregister() {
    setLoader(true);
    browser.storage.local.set({
      hosts: '',
      trasaCore: '',
      extID: '',
      loggedIn: false,
      sessionStore: {},
      trasaDACom: '',
    });
    setLoader(false);
  }

  const classes = useStyles();

  return (
    <Grid container spacing={2} direction="column" alignItems="center" justify="center">
      <Grid item xs={12}>
        <Shield fontSize="large" className={classes.icon} />
        <br />
        <div>Device cyber hygiene managed by TRASA</div>
        <br />

        {failed ? (
          <div>
            <div className={classes.errorText}>
              {' '}
              {error.length > 0 ? error : 'We encountered some error. contact your administrator'}
            </div>{' '}
            <br />
          </div>
        ) : (
          ''
        )}
        <br />
        {loading ? (
          <div>
            {' '}
            <ProgressHOC /> <br />
          </div>
        ) : (
          ''
        )}
      </Grid>

      <Grid container direction="row" justify="center" spacing={2}>
        <Grid item xs={3}>
          <Button variant="contained" className={classes.button} size="small" onClick={getExt}>
            Sync
          </Button>
          <br />
        </Grid>

        <Grid item xs={3}>
          <Button variant="contained" className={classes.button} size="small" onClick={deregister}>
            deregister
          </Button>
        </Grid>

        <br />
      </Grid>
      <Grid item xs={12}>
        v30.03.2020
      </Grid>
    </Grid>
  );
}

export default LoggedIn;
