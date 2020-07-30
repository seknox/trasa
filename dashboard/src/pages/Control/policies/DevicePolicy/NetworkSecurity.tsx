import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import React from 'react';

const useStyles = makeStyles((theme) => ({}));

type networkSecProps = {
  blockOpenWifiConn: boolean;
  changeBlockOpenWifiConn: (event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => void;
}
export default function NetworkSecurity(props: networkSecProps) {
  const classes = useStyles();
  return (
    <Grid container spacing={2}>
      <Grid item xs={12} sm={12} md={12}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={12} md={5}>
            <Typography variant="h4">Block if device is connected to open wifi </Typography>
          </Grid>
          <Grid item xs={12} sm={6} md={2} lg={2}>
            <Switch
              checked={props.blockOpenWifiConn}
              onChange={props.changeBlockOpenWifiConn}
              color="secondary"
            />
          </Grid>
        </Grid>
      </Grid>

      {/* <Grid item xs={12} sm={12} md={12}>
                <Grid container spacing={2}>
                <Grid item xs={12} sm={12} md={5}>
                <Typography variant="h4">Only allow from these mac address: </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={2} lg={2} >
                    <Switch
                        // checked= {props.fileTransfer}
                        // onChange={props.handleFileTransferChange}
                        color="secondary"
                    />                            
                </Grid>
                </Grid>
                </Grid> */}
    </Grid>
  );
}
