import React from 'react';
import ReactDOM from 'react-dom';
import LoggedIn from './loggedin';
import LoginPage from './login';
import { Grid } from '@material-ui/core';
import TRASALogo from './assets/trasa.png';

//var browser = require("webextension-polyfill");

function Home() {
  const [loggedIn, setLoggedIn] = React.useState(false);

  React.useEffect(() => {
    browser.storage.local.get().then(function (token) {
      if (token.extID) {
        setLoggedIn(true);
      }
    });
  }, []);

  function setLoginTrue() {
    setLoggedIn(true);
  }

  return (
    <Grid container spacing={2} direction="column" justify="center" alignItems="center">
      <Grid item xs={12}>
        <br />
        <img src={TRASALogo} alt="trasa-logo" width={150} />
      </Grid>
      <Grid item xs={12}>
        {loggedIn ? <LoggedIn /> : <LoginPage setLoginTrue={setLoginTrue} />}
      </Grid>
      <Grid item xs={12}>
        v1.0
      </Grid>
    </Grid>
  );
}

ReactDOM.render(<Home />, document.getElementById('app'));
