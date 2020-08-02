import Button from '@material-ui/core/Button';
import green from '@material-ui/core/colors/green';
import purple from '@material-ui/core/colors/purple';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import mixpanel from 'mixpanel-browser';
import React, { useState } from 'react';
import Constants from '../../../../Constants';
// Form validation
import ProgressHOC from '../../../../utils/Components/Progressbar';
import { ServiceDetailProps } from '../../../../types/services';

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

  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
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
  selectInput: {
    borderRadius: 4,
    position: 'relative',
    backgroundColor: theme.palette.background.paper,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    width: 'auto',
    padding: '10px 26px 10px 12px',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
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

type ProxySettingProps = {
  proxyConfig: ProxyConfigProps;
  serviceDetail: ServiceDetailProps;
};

type ProxyConfigProps = {
  routeRule: string;
  passHostHeader: boolean;
  upstreamServer: string;
  strictTLSValidation: boolean;
};

export default function ProxySetting(props: ProxySettingProps) {
  const classes = useStyles();
  const [proxy, setProxy] = useState(props.proxyConfig);
  const [loader, setLoader] = useState(false);

  const handleChange = (name: any) => (event: any) => {
    let val = event.target.value;
    if (name === 'passHostHeader' || name === 'strictTLSValidation') {
      val = event.target.checked;
    }

    setProxy({ ...proxy, [name]: val });
  };

  const handleSubmit = () => {
    setLoader(true);
    mixpanel.track('http proxy');

    const { serviceDetail } = props;
    const req = {
      name: serviceDetail.serviceName,
      serviceID: serviceDetail.ID,
      proxy: {
        // routeRule: proxy.routeRule,
        passHostHeader: proxy.passHostHeader,
        upstreamServer: proxy.upstreamServer,
        strictTLSValidation: proxy.strictTLSValidation,
      },
    };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/httpproxy/update`, req)
      .then((r) => {
        setLoader(false);
        if (r.data.status === 'success') {
          // update parent proxy data here
          window.location.reload();
        }
      })
      .catch(() => {
        setLoader(false);
      });
  };

  return (
    <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
      <Grid item xs={12}>
        <Typography variant="h2"> Proxy settings </Typography>
        <Typography variant="h4"> (only for HTTPs applications) </Typography>
      </Grid>
      <Grid item xs={12}>
        <br />

        {/* <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Route rule : </div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <TextField
              fullWidth
              // label="App Name"
              onChange={handleChange('routeRule')}
              name="routeRule"
              value={proxy.routeRule}
              defaultValue={proxy.routeRule}
              InputProps={{
                disableUnderline: true,
                classes: {
                  root: classes.textFieldRoot,
                  input: classes.textFieldInputBig,
                },
              }}
              InputLabelProps={{
                shrink: true,
                className: classes.textFieldFormLabel,
              }}
            />
          </Grid>
        </Grid> */}

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Upstream Server :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <TextField
              fullWidth
              //  label="App Admin Username"
              onChange={handleChange('upstreamServer')}
              name="upstreamServer"
              value={proxy.upstreamServer}
              defaultValue={proxy.upstreamServer}
              // value={this.state.upstreamServers}
              // defaultValue={proxyConfig.upstreamServers}
              InputProps={{
                disableUnderline: true,
                classes: {
                  root: classes.textFieldRoot,
                  input: classes.textFieldInputBig,
                },
              }}
              InputLabelProps={{
                shrink: true,
                className: classes.textFieldFormLabel,
              }}
            />
          </Grid>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Pass Host Header :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <FormControl fullWidth>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Switch
                      checked={proxy.passHostHeader}
                      onChange={handleChange('passHostHeader')}
                      name="passHostHeader"
                      value={proxy.passHostHeader}
                      // defaultValue={proxy.passHostHeader}
                      color="primary"
                    />
                  }
                  label={
                    proxy.passHostHeader ? (
                      <div className={classes.settingSHeader}>true </div>
                    ) : (
                      <div className={classes.settingSHeader}>false </div>
                    )
                  }
                />
              </FormGroup>
            </FormControl>
          </Grid>
        </Grid>

        <Grid container spacing={2}>
          <Grid item xs={5} sm={5} md={5}>
            <div className={classes.settingHeader}>Strict TLS validation :</div>
          </Grid>
          <Grid item xs={7} sm={7} md={7}>
            <FormControl fullWidth>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Switch
                      checked={proxy.strictTLSValidation}
                      onChange={handleChange('strictTLSValidation')}
                      name="strictTLSValidation"
                      value={proxy.strictTLSValidation}
                      // defaultValue={proxy.passHostHeader}
                      color="primary"
                    />
                  }
                  label={
                    proxy.strictTLSValidation ? (
                      <div className={classes.settingSHeader}>true </div>
                    ) : (
                      <div className={classes.settingSHeader}>false </div>
                    )
                  }
                />
              </FormGroup>
            </FormControl>
          </Grid>
        </Grid>

        <br />
        <br />
      </Grid>
      <Grid container alignItems="flex-start" justify="flex-end" direction="row">
        <Grid item xs={12}>
          <Button variant="contained" onClick={handleSubmit} className={classes.submitButton}>
            Submit
          </Button>
        </Grid>

        <br />
        {loader ? <ProgressHOC /> : ''}
      </Grid>
    </Grid>
  );
}
