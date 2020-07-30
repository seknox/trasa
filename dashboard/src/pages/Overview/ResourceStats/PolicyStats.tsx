import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { HeaderFontSize, TitleFontSize } from '../../../utils/Responsive';

const useStyles = makeStyles((theme) => ({
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

function PolicyStat() {
  const [policyStat, setPolicyStat] = useState({ total: 0, expired: 0 });

  useEffect(() => {
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/policies`;

    axios
      .get(reqPath, config)
      .then((r) => {
        setPolicyStat(r.data.data[0]);
        // console.log('uuuuuuuuuu ', r.data.data[0])
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

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Policies </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={6}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Total </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {policyStat.total} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={6}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Expired </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {policyStat.expired} </b>
            </div>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

export default PolicyStat;
