import { Grid, TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Drawer from '@material-ui/core/Drawer';
import FormControl from '@material-ui/core/FormControl';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import IdpIcon from '@material-ui/icons/SupervisorAccount';
import axios from 'axios';
import React, { useState } from 'react';
import Constants from '../../../../Constants';
import ProgressHOC from '../../../../utils/Components/Progressbar';

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
    padding: theme.spacing(4),
  },

  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  textFieldInput: {
    borderRadius: 4,
    // backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    padding: theme.spacing(1),
    // backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: '#404854',
    fontSize: '15px',
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
}));

type Anchor = 'top' | 'left' | 'bottom' | 'right';

export default function CreateIdpDrawer(props: any) {
  const classes = useStyles();
  const [idp, setIdpVal] = React.useState({
    idpName: 'okta',
    idpType: 'saml',
    idpMeta: '',
    endpoint: '',
  });
  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
  });
  // const [selectedIdp, setSelectedIdp] = React.useState('');

  // const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
  //   setSelectedIdp(event.target.value);
  // };

  const [state, setState] = React.useState({
    top: false,
    left: false,
    bottom: false,
    right: false,
  });

  const toggleDrawer = (anchor: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }

    setState({ ...state, [anchor]: open });
  };

  // const toggleDrawer = (side, open) => (event) => {
  //   if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
  //     return;
  //   }

  //   setState({ ...state, [side]: open });
  // };

  function idpChange(e: React.ChangeEvent<{ value: unknown }>) {
    setIdpVal({ ...idp, idpType: e.target.value as string });
  }

  function idpNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    setIdpVal({ ...idp, idpName: e.target.value });
  }

  function setIdpName() {
    switch (idp.idpType) {
      case 'okta':
        idp.idpName = 'okta';
        break;
      case 'freeipa':
        idp.idpName = 'freeipa';
        break;
      case 'ad':
        idp.idpName = 'ad';
        break;
      default:
        break;
    }
  }
  const submitIdp = () => {
    setIdpName();
    updateActionStatus({ ...actionStatus, respStatus: false, statusMsg: '', loader: true });

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/uidp/create`, idp)
      .then((r) => {
        updateActionStatus({ ...actionStatus, loader: false });
        if (r.data.status === 'success') {
          updateActionStatus({
            ...actionStatus,
            success: true,
            respStatus: true,
            statusMsg: r.data.reason,
          });
          // let newIdps = props.idps.push(r.data.data[0])
          props.setIdps([...props.idps, r.data.data[0]]);
        } else {
          updateActionStatus({
            ...actionStatus,
            respStatus: true,
            statusMsg: r.data.reason,
            success: false,
          });
        }
      })
      .catch((error) => {
        updateActionStatus({
          ...actionStatus,
          respStatus: true,
          statusMsg: 'something went wrong',
          success: false,
        });
      });
  };

  return (
    <div>
      <Button
        variant="contained"
        // color="primary"
        size="small"
        onClick={toggleDrawer('right', true)}
      >
        <IdpIcon />
        Create New Identity Provider
      </Button>
      <Drawer anchor="right" open={state.right} onClose={toggleDrawer('right', false)}>
        <div>
          <Paper className={classes.paper}>
            <Typography variant="h2"> Create New Identity Provider </Typography>
            <br />
            <br />
            <br />
            <Grid container spacing={2} direction="column">
              <Grid item xs={12}>
                <Typography variant="h3"> Select your Identity Provider :</Typography>
              </Grid>
              <Grid item xs={12}>
                <FormControl fullWidth>
                  <Select
                    name="idpType"
                    defaultValue={props.idpType}
                    onChange={idpChange}
                    inputProps={{
                      classes: {
                        root: classes.selectCustom,
                      },
                    }}
                  >
                    <MenuItem value="okta">
                      <div className={classes.settingSHeader}>Okta </div>
                    </MenuItem>
                    <MenuItem value="freeipa">
                      <div className={classes.settingSHeader}>FreeIPA</div>
                    </MenuItem>
                    <MenuItem value="ldap">
                      <div className={classes.settingSHeader}>Ldap</div>
                    </MenuItem>
                    <MenuItem value="ad">
                      <div className={classes.settingSHeader}>Active Directory</div>
                    </MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              {idp.idpType === 'ldap' ? (
                <div>
                  <br />
                  <br />
                  <br />
                  <Grid item xs={12}>
                    <Typography variant="h3"> Give your identity provider unique name :</Typography>
                    <br />
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      fullWidth
                      name="idpName"
                      autoFocus
                      onChange={idpNameChange}
                      InputProps={{
                        disableUnderline: true,
                        classes: {
                          root: classes.textFieldRoot,
                          input: classes.textFieldInputBig,
                        },
                      }}
                    />
                  </Grid>
                </div>
              ) : null}
              <br />
              <br />
              <br />
              <Grid item xs={12}>
                <Button fullWidth variant="contained" size="small" onClick={submitIdp}>
                  Create
                </Button>
              </Grid>
            </Grid>

            <br />
            {actionStatus.loader ? <ProgressHOC /> : ''}
          </Paper>
        </div>
      </Drawer>
    </div>
  );
}
