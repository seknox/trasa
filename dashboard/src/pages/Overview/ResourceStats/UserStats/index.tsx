import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import { ResourceStatsFilterProps } from '../../../../types/analytics';
import { HeaderFontSize, TitleFontSize } from '../../../../utils/Responsive';
import UserPerIdpTypes from './UsersPerIdpType';

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
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperLarge: {
    backgroundColor: '#fdfdfd',
    minWidth: 400,
    minHeight: 500,
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
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

function UserAggregatedData(props: ResourceStatsFilterProps) {
  const [allusers, setAllusers] = useState({
    totalUsers: 0,
    totalIdps: [],
    users: 0,
    admins: 0,
    groups: 0,
    disabledUsers: 0,
  });

  useEffect(() => {
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    const reqPath = `${Constants.TRASA_HOSTNAME}/api/v1/stats/users/${props.entityType}/${props.entityID}`;

    axios
      .get(reqPath, config)
      .then((r) => {
        setAllusers(r.data.data[0]);
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
  }, [props.entityType, props.entityID]);

  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <div className={classes.headerMain}>
            <b> Users </b>
          </div>
          <Divider light />
        </Grid>

        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Total </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allusers.users} </b>
            </div>
          </Paper>
        </Grid>
        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Administrators </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allusers.admins} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Groups </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allusers.groups} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Disabled </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allusers.disabledUsers} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={2}>
          <Paper className={classes.paperTrans} elevation={0}>
            <div className={classes.aggHeaders}>
              {' '}
              <b> Idps </b>
            </div>
            <div className={classes.aggHeadersBig}>
              {' '}
              <b> {allusers.totalIdps.length} </b>
            </div>
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Divider light />
        </Grid>

        <Grid item xs={12}>
          <UserPerIdpTypes totalIdps={allusers.totalIdps} />
        </Grid>
      </Grid>
    </div>
  );
}

export default UserAggregatedData;
