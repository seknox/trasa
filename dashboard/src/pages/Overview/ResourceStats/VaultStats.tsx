import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { HeaderFontSize, TitleFontSize, TitleTextFontSize } from '../../../utils/Responsive';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },
  paperTrans: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(1),,
    textAlign: 'center',
  },
  paperTrans1: {
    backgroundColor: 'transparent',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(1),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 400,
    minHeight: 500,
    // padding: theme.spacing(2),
    // textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  headerMain: {
    color: '#1b1b32',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeadersBig: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  TitleText: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleTextFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    consttSize: HeaderFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function VaultStats() {
  const [vaultStatus, setVaultStatus] = useState({ tsxVault: false });

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/status`;

    axios
      .get(reqPath)
      .then((r) => {
        const d = r.data.data[0];
        try {
          const sett = JSON.parse(d.setting);

          d.setting = sett;
        } catch (e) {
          console.error(e);
          d.setting = {};
        }
        setVaultStatus(d);
        // console.log('uuuuuuuuuu ', r.data.data[0])
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Secret Storage </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={6}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> TsxVault </b>
            </div>
            <div className={classes.TitleText}>
              {' '}
              <b> {vaultStatus.tsxVault ? 'Enabled' : 'Disabled'} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Cred Store </b>
              <div className={classes.TitleText}>
                {' '}
                <b> TsxVault </b>
              </div>
            </div>
          </Paper>
        </Grid>
        {/* <Grid item xs={4}>
            <Paper className={classes.paperTrans} elevation={0}>

            <div className={classes.aggHeaders}> <b> service Creds Store </b></div>
            <div className={classes.TitleText}> <b> {vaultStatus.setting && vaultStatus.setting.credStorage}  </b></div>
            </Paper>
            </Grid> */}
      </Grid>
    </div>
  );
}
