import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

const useStyles = makeStyles((theme) => ({}));

type endpointSecurityprops = {
  blockEncryptionNotSet: boolean;
  changeBlockEncryptionNotSet: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;

  blockFirewallDisabled: boolean;
  changeBlockFirewallDisabled: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;

  blockAntivirusDisabled: boolean;
  changeBlockAntivirusDisabled: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;


}

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
              checked={props.blockEncryptionNotSet}
              onChange={props.changeBlockEncryptionNotSet}
              color="secondary"
            />
          </Grid>
        </Grid>
          <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if firewall not enabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
                checked={props.blockFirewallDisabled}
                onChange={props.changeBlockFirewallDisabled}
                color="secondary"
            />

          </Grid>
        </Grid>

  <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if antivirus not enabled (Windows Only): </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
                checked={props.blockAntivirusDisabled}
                onChange={props.changeBlockAntivirusDisabled}
                color="secondary"
            />

          </Grid>
        </Grid>



      </Grid>
    </Grid>
  );
}
