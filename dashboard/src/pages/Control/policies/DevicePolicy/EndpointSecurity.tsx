import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

const useStyles = makeStyles((theme) => ({}));

type endpointSecurityprops = {
  // TODO @types
  dhBlocking: any;
  handleDHBlockingChange: any;
};

export default function EndpointSecurity(props: endpointSecurityprops) {
  const classes = useStyles();
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if device disk is not encrypted: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockEncryptionNotSet}
              onChange={props.handleDHBlockingChange('blockEncryptionNotSet')}
              name="blockEncryptionNotSet"
              color="secondary"
              value={props.dhBlocking.blockEncryptionNotSet}
            />
          </Grid>
        </Grid>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if firewall not enabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockFirewallDisabled}
              onChange={props.handleDHBlockingChange('blockFirewallDisabled')}
              name="blockFirewallDisabled"
              color="secondary"
              value={props.dhBlocking.blockFirewallDisabled}
            />
          </Grid>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if antivirus not enabled (Windows Only): </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.dhBlocking.blockAntivirusDisabled}
              onChange={props.handleDHBlockingChange('blockAntivirusDisabled')}
              name="blockAntivirusDisabled"
              color="secondary"
              value={props.dhBlocking.blockAntivirusDisabled}
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}
