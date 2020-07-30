import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

const useStyles = makeStyles((theme) => ({
  divider: {
    maxWidth: 400
  }
}));

type osSecurityProps = {
  blockJailBroken: boolean;
  changeBlockJailBroken: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;
  blockDebuggingEnabled: boolean;
  changeBlockDebuggingEnabled: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;
  blockEmulated: boolean
  changeBlockEmulated : (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;

  blockCriticalAutoUpdateDisabled: boolean
  changeBlockCriticalAutoUpdateDisabled : (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;


}

export default function OsSecurity(props: osSecurityProps) {
  const classes = useStyles();
  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Typography component="span" variant="h4">
          {' '}
          Mobile{' '}
        </Typography>
        {/* <Typography component="span" variant="h5"> (This will also appy to your 2FA mobile devices)</Typography> */}
        <Divider light className={classes.divider} />
      </Grid>
      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block jailbroken device: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.blockJailBroken}
              onChange={props.changeBlockJailBroken}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block debugging enabled device: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.blockDebuggingEnabled}
              onChange={props.changeBlockDebuggingEnabled}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block emulated device: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.blockEmulated}
              onChange={props.changeBlockEmulated}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>
      {/* <Grid item xs={12}>
        <Typography component="span" variant="h4">
          {' '}
          Workstation{' '}
        </Typography>
        <Divider light className={classes.divider} />

      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block emulated device: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
                checked={props.blockCriticalAutoUpdateDisabled}
                onChange={props.changeBlockCriticalAutoUpdateDisabled}
                color="secondary"
            />
          </Grid>
        </Grid>
      </Grid> */}



    </Grid>
  );
}
