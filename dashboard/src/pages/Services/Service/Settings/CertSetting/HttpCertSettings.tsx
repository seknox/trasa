import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useState } from 'react';
import Constants from '../../../../../Constants';
import ProgressHOC from '../../../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    // flexgrow: 1,
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)' // #011019
    // minWidth: 400,
    // minHeight: 300,
    padding: theme.spacing(2),
    //  textAlign: 'center',
    color: theme.palette.text.secondary,
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
  settingHeader: {},
}));

export default function TlsConfig(props: any) {
  const [loader, setLoader] = useState(false);
  const [certFile, setCert] = useState('');
  const [keyFile, setKey] = useState('');
  const classes = useStyles();
  const [configType, setConfigType] = useState('le');

  const handleCertFile = (e: any) => {
    setCert(e.target.value);
  };

  const handleKeyFile = (e: any) => {
    setKey(e.target.value);
  };

  const onUpload = () => {
    setLoader(true);

    const formData = {
      serviceID: props.serviceID,
      configType,
      certVal: '',
      keyVal: '',
    };

    if (configType === 'upload') {
      formData.certVal = certFile;
      formData.keyVal = keyFile;
    }



    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/gateway/httpproxy/configtls`, formData)
      .then((r) => {
        if (r.data.status === 'success') {
          setLoader(false);
        }
      });
  };

  const changePolicy = (name: any) => (event: any) => {
    switch (name) {
      case 'le':
        setConfigType('le');
        break;
      case 'upload':
        setConfigType('upload');
        break;
    }
  };

  const getCheckStatus = (v: any) => {
    if (v === configType) {
      return true;
    }
    return false;
  };

  return (
    <Grid container spacing={2} alignItems="center" direction="column" justify="center">
      <Grid item xs={12}>
        <Typography variant="h2"> Tls Cert </Typography>
        <Typography variant="h4"> (only for HTTPs applications) </Typography>
      </Grid>

      <Grid item xs={12}>
        <Grid container spacing={2} direction="row" alignItems="center">
          <Grid item xs={6}>
            <div className={classes.settingHeader}>
              {' '}
              Fetch from let's encrypt :{' '}
              <span>
                <Checkbox
                  // defaultChecked
                  checked={!!getCheckStatus('le')}
                  // value={"enforceStrongPass"}
                  value={configType}
                  name={configType}
                  color="primary"
                  onChange={changePolicy('le')}
                  inputProps={{ 'aria-label': 'secondary checkbox' }}
                />
              </span>{' '}
            </div>
          </Grid>
          <Grid item xs={6}>
            <div className={classes.settingHeader}>
              {' '}
              Upload certificates :{' '}
              <span>
                <Checkbox
                  // defaultChecked
                  checked={!!getCheckStatus('upload')}
                  // value={"enforceStrongPass"}
                  value={configType}
                  name={configType}
                  color="primary"
                  onChange={changePolicy('upload')}
                  inputProps={{ 'aria-label': 'secondary checkbox' }}
                />
              </span>{' '}
            </div>
          </Grid>
        </Grid>
      </Grid>

      {configType === 'upload' ? (
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={3}>
                  <Typography variant="h4"> Upload Cert file : </Typography>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    // label="Service name"
                    rows="4"
                    multiline
                    onChange={handleCertFile}
                    name="sslCert"
                    value={certFile}
                    InputProps={{
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
            </Grid>

            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={3}>
                  <Typography variant="h4"> Upload Key file : </Typography>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    // label="Service name"
                    rows="4"
                    multiline
                    onChange={handleKeyFile}
                    name="sslCert"
                    value={keyFile}
                    InputProps={{
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
            </Grid>

            <Grid container spacing={2} direction="row" justify="flex-end">
              <Grid item xs={6}>
                <Button variant="contained" onClick={onUpload}>
                  Submit
                  {/* <CloudUploadIcon/> */}
                </Button>
                <br />
                {loader ? <ProgressHOC /> : ''}
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      ) : (
        ''
      )}
    </Grid>
  );
}
