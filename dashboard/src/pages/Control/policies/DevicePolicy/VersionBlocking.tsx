import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

const useStyles = makeStyles((theme) => ({}));

type versionBlockingProps = {

}

export default function VersionBlocking(props: versionBlockingProps) {
  const classes = useStyles();
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block jailbroken or rooted mobile devices: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              // checked= {props.fileTransfer}
              // onChange={props.handleFileTransferChange}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Block device version: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              // checked= {props.fileTransfer}
              // onChange={props.handleFileTransferChange}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>

      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={4}>
            <Typography variant="h4">Device Remote Login Enabled: </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              // checked= {props.fileTransfer}
              // onChange={props.handleFileTransferChange}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}
