import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles, Theme } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { HeaderFontSize, TitleFontSize, ValueTextFontSize } from '../../../utils/Responsive';

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
  aggHeadersBig: {
    color: '#000080', // '#1b1b32',
    fontSize: TitleFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  headerMain: {
    color: '#1b1b32',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  aggHeaders: {
    color: '#1b1b32',
    fontSize: HeaderFontSize(),
    fontFamily: 'Open Sans, Rajdhani',
  },
  valueFont: {
    color: '#000080', // '#1b1b32',
    fontSize: ValueTextFontSize(), // '14px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

function BackupStats() {


  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Backup </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={12}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> enterprise only </b>
            </div>
          </Paper>
        </Grid>



      </Grid>
    </div>
  );
}

export default BackupStats;
