import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles, Theme } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { HeaderFontSize, TitleFontSize } from '../../../utils/Responsive';

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paperTrans: {
    backgroundColor: 'transparent',
    // padding: theme.spacing(1),
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
  success: {
    paddingLeft: 5,
    paddingRight: 0,
    maxWidth: 50,
    background: 'green',
    color: 'white',
  },
  failed: {
    paddingLeft: 10,
    // paddingRight: 5,
    maxWidth: 50,
    background: 'maroon',
    color: 'white',
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
  aggHeaders: {
    color: '#1b1b32',
    fontSize: HeaderFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function CAStats() {
  const [castat, setCAstat] = useState({
    sshHostCA: false,
    sshUserCA: false,
    sshSystemCA: false,
    httpCA: false,
  });

  useEffect(() => {
    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/ca`;

    axios
      .get(reqPath)
      .then((r) => {
        setCAstat(r.data.data[0]);
      })
      .catch((error) => {
        if (error.response) {
          console.log(error.response.data);
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log('Error', error.message);
        }
      });
  }, []);
  const statusDiv = (val: boolean) => {
    if (val) {
      return 'Available';
    }
    return 'n.a';
  };

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Certificate Authority </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> SSH Host CA </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {statusDiv(castat.sshHostCA)} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> SSH User CA </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {statusDiv(castat.sshUserCA)} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={4}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> HTTP CA </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {statusDiv(castat.httpCA)} </b>
            </div>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

export default CAStats;
