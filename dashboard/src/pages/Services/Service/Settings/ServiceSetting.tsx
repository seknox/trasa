import Button from '@material-ui/core/Button';
import green from '@material-ui/core/colors/green';
import purple from '@material-ui/core/colors/purple';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import { ServiceDetailProps } from '../../../../types/services';
import ProgressHOC from '../../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  // card
  card: {
    minWidth: 275,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    marginBottom: 16,
    // fontSize: 54,
    color: '#311B92', // theme.palette.text.secondary,
    // backgroundColor: '#1A237E',
  },
  pos: {
    marginBottom: 12,
    color: theme.palette.text.secondary,
  },

  // form
  formControl: {
    margin: theme.spacing(1),
  },
  inputLabelFocused: {
    color: purple[500],
  },
  inputInkbar: {
    '&:after': {
      backgroundColor: purple[500],
    },
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12,
  },
  successText: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },

  fab: {
    margin: theme.spacing(2),
    background: '#000080',
  },
  fab2: {
    margin: theme.spacing(2),
  },
  users: {
    margin: theme.spacing(2),
  },
  dividerInset: {
    margin: `5px 0 0 ${theme.spacing(9)}px`,
  },
  credbuttons: {
    marginTop: 20,
    background: 'navy',
    color: 'white',
  },
  fabButton: {
    margin: theme.spacing(2),
    background:
      'linear-gradient(to right, #021B79, #0575E6)' /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */,
  },
  Warning: {
    color: 'maroon',
    fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  WarningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(1),
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
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

  submitButton: {
    float: 'right',
  },
}));

type ServicesettingProps = {
  newApp: boolean;
  serviceDetail: ServiceDetailProps;
};

export default function Servicesetting(props: ServicesettingProps) {
  const classes = useStyles();
  const [serviceDetail, setData] = useState<ServiceDetailProps>({
    ID: '',
    serviceName: '',
    serviceType: '',
    rdpProtocol: '',
    remoteserviceName: '',
    hostname: '',
    adhoc: false,
    passthru: false,
    nativeLog: false,
  });
  const [loader, setLoader] = useState(false);

  useEffect(() => {
    setData(props.serviceDetail);
  }, [props.serviceDetail]);

  const handleChange = (name: any) => (event: any) => {
    let val = event.target.value;
    if (
      name === 'passthru' ||
      name === 'adhoc' ||
      name === 'sessionRecord' ||
      name === 'nativeLog'
    ) {
      val = event.target.checked;
    }
    setData({ ...serviceDetail, [name]: val });
  };

  function handleSubmit() {
    setLoader(true);

    let url = `${Constants.TRASA_HOSTNAME}/api/v1/services/update`;
    if (props.newApp === true) {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/services/create`;
    }

    axios
      .post(url, serviceDetail)
      .then((r) => {
        setLoader(false);
        if (r.data.status === 'success') {
          if (props.newApp === true) {
            window.location.href = `/services/service/${r.data.data[0].ID}`;
          } else {
            // update parent service config state sere
            window.location.reload();
          }
        }
      })
      .catch(() => {
        setLoader(false);
      });
  }

  return (
    <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
      <Grid item xs={12}>
        <Typography variant="h2">
          {props.newApp ? 'Integrate New Service' : `Edit ${serviceDetail.serviceName}`}{' '}
        </Typography>
      </Grid>

      <Grid item xs={12}>
        <br />

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Service name :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <TextField
              fullWidth
              // label="Service name"
              onChange={handleChange('serviceName')}
              name="serviceName"
              id="serviceName"
              value={serviceDetail.serviceName}
              // defaultValue={props.serviceDetail.serviceName}
              variant="outlined"
              size="small"
            />
          </Grid>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Application Type :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <FormControl fullWidth variant="outlined" size="small">
              <Select
                value={serviceDetail.serviceType}
                // defaultValue={serviceDetail.serviceType}
                name="serviceType"
                onChange={handleChange('serviceType')}
                inputProps={{
                  name: 'serviceType',
                  id: 'serviceType',
                  // classes: {
                  //   root: classes.selectCustom,
                  // },
                }}
              >
                <MenuItem value="ssh" id="ssh">
                  <div className={classes.settingSHeader}>SSH </div>
                </MenuItem>
                <MenuItem value="http" id="http">
                  <div className={classes.settingSHeader}>HTTP (web) </div>
                </MenuItem>

                <MenuItem value="rdp" id="rdp">
                  <div className={classes.settingSHeader}>RDP </div>
                </MenuItem>
                <MenuItem value="radius" id="radius">
                  <div className={classes.settingSHeader}>Radius </div>
                </MenuItem>
                <MenuItem value="db" id="db">
                  <div className={classes.settingSHeader}>Database </div>
                </MenuItem>
              </Select>
            </FormControl>
          </Grid>
        </Grid>

        {props.serviceDetail.serviceType === 'rdp' ? (
          <Grid container spacing={2}>
            <Grid item xs={5} sm={5} md={5}>
              <div className={classes.settingHeader}>Security Protocol : </div>
            </Grid>
            <Grid item xs={7} sm={7} md={7}>
              <FormControl fullWidth variant="outlined" size="small">
                <Select
                  // label="Security Protocol"
                  value={serviceDetail.rdpProtocol}
                  // defaultValue={props.serviceDetail.rdpProtocol}
                  name="rdpProtocol"
                  onChange={handleChange('rdpProtocol')}
                  inputProps={{
                    name: 'rdpProtocol',
                    id: 'rdpProtocol',
                  }}
                >
                  <MenuItem value="nla">
                    <div className={classes.settingSHeader}>NLA </div>
                  </MenuItem>
                  <MenuItem value="tls">
                    <div className={classes.settingSHeader}>TLS </div>
                  </MenuItem>
                  <MenuItem value="rdp">
                    <div className={classes.settingSHeader}>RDP </div>
                  </MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={5} sm={5} md={5}>
              <div className={classes.settingHeader}>Remote Service name : </div>
            </Grid>
            <Grid item xs={7} sm={7} md={7}>
              <TextField
                fullWidth
                onChange={handleChange('remoteserviceName')}
                name="remoteserviceName"
                value={serviceDetail.remoteserviceName}
                variant="outlined"
                size="small"
              />
            </Grid>
          </Grid>
        ) : null}

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            {props.serviceDetail.serviceType === 'http' ? (
              <div className={classes.settingHeader}>Domain name which points to TRASA : </div>
            ) : (
              <div className={classes.settingHeader}>Hostname : </div>
            )}
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <TextField
              fullWidth
              onChange={handleChange('hostname')}
              name="hostname"
              id="hostname"
              value={serviceDetail.hostname}
              // defaultValue={props.serviceDetail.hostname}
              // validators={['required']}
              // errorMessages={['this field is required']}
              variant="outlined"
              size="small"
            />
          </Grid>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Adhoc Access :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <FormControl fullWidth>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Switch
                      checked={!!serviceDetail.adhoc}
                      onChange={handleChange('adhoc')}
                      name="adhoc"
                      // defaultValue={props.serviceDetail.adhoc}
                      value="adhoc"
                      color="primary"
                    />
                  }
                  label={
                    serviceDetail.adhoc ? (
                      <div className={classes.settingSHeader}>enabled </div>
                    ) : (
                      <div className={classes.settingSHeader}>disabled </div>
                    )
                  }
                />
              </FormGroup>
            </FormControl>
          </Grid>
        </Grid>

        {props.newApp ? null : (
          <div>
            <Grid container spacing={2}>
              <Grid item xs={5} sm={5} md={5}>
                <div className={classes.settingHeader}>Passthru Authentication :</div>
              </Grid>
              <Grid item xs={7} sm={7} md={7}>
                <FormControl fullWidth>
                  <FormGroup>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={!!serviceDetail.passthru}
                          onChange={handleChange('passthru')}
                          name="passthru"
                          // defaultValue={props.serviceDetail.passthru}
                          value="passthru"
                          color="primary"
                        />
                      }
                      label={
                        serviceDetail.passthru ? (
                          <div className={classes.settingSHeader}>enabled </div>
                        ) : (
                          <div className={classes.settingSHeader}>disabled </div>
                        )
                      }
                    />
                  </FormGroup>
                </FormControl>
              </Grid>
            </Grid>

            {serviceDetail.serviceType === 'ssh' || 'rdp' ? (
              <Grid container spacing={2}>
                <Grid item xs={5} sm={5} md={5}>
                  <div className={classes.settingHeader}>Native Agent Logs :</div>
                </Grid>
                <Grid item xs={7} sm={7} md={7}>
                  <FormControl fullWidth>
                    <FormGroup>
                      <FormControlLabel
                        control={
                          <Switch
                            checked={!!serviceDetail.nativeLog}
                            onChange={handleChange('nativeLog')}
                            name="nativeLog"
                            // defaultValue={props.serviceDetail.nativeLog}
                            value="nativeLog"
                            color="primary"
                          />
                        }
                        label={
                          serviceDetail.nativeLog ? (
                            <div className={classes.settingSHeader}>enabled </div>
                          ) : (
                            <div className={classes.settingSHeader}>disabled </div>
                          )
                        }
                      />
                    </FormGroup>
                  </FormControl>
                </Grid>
              </Grid>
            ) : null}
          </div>
        )}

        <br />
        <br />
      </Grid>

      <Grid container alignItems="flex-start" justify="flex-end" direction="row">
        <Grid item xs={12}>
          <Button
            name="submit"
            id="submit"
            variant="contained"
            onClick={handleSubmit}
            className={classes.submitButton}
          >
            Submit
          </Button>
        </Grid>

        <br />
        {loader ? <ProgressHOC /> : ''}
      </Grid>
    </Grid>
  );
}
