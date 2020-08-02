import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

type loginSecurityProps = {
  // TODO @types
  dhBlocking: any;
  handleDHBlockingChange: any;
};

export default function LoginSecurity(props: loginSecurityProps) {
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block device with autologin enabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockAutologinEnabled}
              onChange={props.handleDHBlockingChange('blockAutologinEnabled')}
              name="blockAutologinEnabled"
              color="secondary"
              value={props.dhBlocking.blockAutologinEnabled}
            />
          </Grid>
        </Grid>
      </Grid>

      {/* <Grid item xs={12} sm={12} md={12}> */}
      {/*  <Grid container spacing={2}> */}
      {/*    <Grid item xs={12} sm={12} md={5}> */}
      {/*      <Typography variant="h4">Block device if 2fa is not configured: </Typography> */}
      {/*    </Grid> */}
      {/*    <Grid item xs={12} sm={6} md={2} lg={2}> */}
      {/*      <Switch */}
      {/*        checked={props.blockTfaNotConfigured} */}
      {/*        onChange={props.changeBlockTfaNotConfigured} */}
      {/*        color="secondary" */}
      {/*      /> */}
      {/*    </Grid> */}
      {/*  </Grid> */}
      {/* </Grid> */}

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block device idle screen lock disabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockIdleScreenLockDisabled}
              onChange={props.handleDHBlockingChange('blockIdleScreenLockDisabled')}
              name="blockIdleScreenLockDisabled"
              color="secondary"
              value={props.dhBlocking.blockIdleScreenLockDisabled}
            />
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block device with remote login Enabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockRemoteLoginEnabled}
              onChange={props.handleDHBlockingChange('blockRemoteLoginEnabled')}
              name="blockRemoteLoginEnabled"
              color="secondary"
              value={props.dhBlocking.blockRemoteLoginEnabled}
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}
