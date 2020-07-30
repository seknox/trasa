import { Grid, TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import ProgressHOC from '../../../utils/Components/Progressbar';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  content: {
    width: '100%',
    flexGrow: 1,
    // backgroundColor: theme.palette.background.default,
    padding: 24,
    height: 'calc(100% - 56px)',
    // height: '100%',
    marginTop: 26,
    // [theme.breakpoints.up('sm')]: {
    //   height: 'calc(100% - 64px)',
    // },
  },

  paper: {
    padding: theme.spacing(2),
    // minWidth: 800,
  },

  textFieldInputBig: {
    borderRadius: 4,
    padding: theme.spacing(1),
    // backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '14px',
    fontWeight: 100,
    fontFamily: 'Open Sans, Rajdhani',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  selectCustom: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 17,
    // marginTop: 5,
    // padding: '10px 100px',
    // width: 100,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  fab: {
    fontSize: '10px',
  },
}));

export default function InternalHosts() {
  const classes = useStyles();
  const [reqStatus, setReqStatus] = useState(false);

  const [internalHosts, setHosts] = useState();

  useEffect(() => {
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/gateway/internalhosts`, config)
      .then((response) => {
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          setHosts(data);
        }
      });
  }, []);

  function hostChange(e: any) {
    setHosts(e.target.value);
  }

  const submitInternalHosts = () => {
    setReqStatus(true);
    const req = { internalHosts };
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/gateway/updateinternalhosts`, req, config)
      .then((r) => {
        setReqStatus(false);
      });
  };

  return (
    <div className={classes.root}>
      {/* <Grid container spacing={2} direction="row"  justify="center"> */}
      <ExpansionPanel>
        <Grid item xs={12} sm={12}>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1a-content"
            id="panel1a-header"
          >
            <Typography component="h4" variant="h3">
              <b>HTTP Internal Hosts</b>
            </Typography>
          </ExpansionPanelSummary>
        </Grid>
        <Grid item xs={12}>
          <Divider light />
        </Grid>

        <Grid item xs={12} sm={12}>
          <ExpansionPanelDetails>
            <Grid container spacing={2} alignItems="center" justify="center">
              <Grid item xs={1}>
                <Typography variant="h4">Hosts: </Typography>
              </Grid>
              <Grid item xs={11}>
                <TextField
                  fullWidth
                  onChange={hostChange}
                  name="internalHosts"
                  value={internalHosts}
                  defaultValue={internalHosts}
                />
              </Grid>
            </Grid>
          </ExpansionPanelDetails>
        </Grid>
        {reqStatus ? (
          <div>
            <ProgressHOC /> <br />
          </div>
        ) : null}
        <ExpansionPanelActions>
          <Button variant={reqStatus ? 'text' : 'contained'} onClick={submitInternalHosts}>
            {' '}
            Save{' '}
          </Button>
        </ExpansionPanelActions>
      </ExpansionPanel>
      {/* </Grid> */}
      <br /> <br /> <br />
    </div>
  );
}
