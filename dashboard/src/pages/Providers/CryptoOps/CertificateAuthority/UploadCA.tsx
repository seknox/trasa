import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useState } from 'react';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import Constants from '../../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    maxWidth: 600,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  notifSnack: {
    width: '100%',
    '& > * + *': {
      marginTop: theme.spacing(2),
    },
  },
  inputButton: {
    background: '#0575E6' /* fallback for old browsers */,
    color: 'white',
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

export default function TlsConfig(props: any) {
  const [loader, setLoader] = useState(false);
  const [statusMsg, setstatusMsg] = useState('');
  const [respStatus, setRespStatus] = useState(false);
  const [success, setSuccess] = useState(false);
  const [certFile, setCert] = useState();
  const [keyFile, setKey] = useState();
  const [csrFile, setCsr] = useState();
  const classes = useStyles();

  const handleClose = (event: any, reason: any) => {
    if (reason === 'clickaway') {
      return;
    }
    setRespStatus(false);
  };

  const handleCertFile = (e: any) => {
    setCert(e.target.value);
  };

  const handleKeyFile = (e: any) => {
    setKey(e.target.value);
  };

  const handleCsrFile = (e: any) => {
    setCsr(e.target.value);
  };

  const onUpload = () => {
    setLoader(true);
    setRespStatus(false);
    setstatusMsg('');
    const formData = {
      certVal: certFile,
      keyVal: keyFile,
      csrVal: csrFile,
    };

    axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/ca/tsxca/upload`, formData).then((r) => {
      if (r.data.status === 'success') {
        setLoader(false);
        setRespStatus(true);
        setstatusMsg(r.data.reason);
        setSuccess(true);
      } else {
        setLoader(false);
        setRespStatus(true);
        setstatusMsg(r.data.reason);
        setSuccess(false);
      }
    });
  };

  return (
    <div className={classes.root}>
      <Grid container spacing={2} alignItems="center" direction="column" justify="center">
        <Grid item xs={12}>
          {/* <Typography variant="h2"> Tls Cert </Typography> */}
          <Typography variant="h4"> (* All files must be in pem format) </Typography>
        </Grid>
        <br />
        <Grid item xs={12}>
          <Divider light />
          <br /> <br />
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={3}>
                  <Typography variant="h4"> Paste Cert file: </Typography>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    // label="Service name"
                    rows="10"
                    multiline
                    onChange={handleCertFile}
                    name="sslCert"
                    value={certFile}
                  />
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={3}>
                  <Typography variant="h4"> Paste Key file : </Typography>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    // label="Service name"
                    rows="10"
                    multiline
                    onChange={handleKeyFile}
                    name="sslCert"
                    value={keyFile}
                  />
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={3}>
                  <Typography variant="h4"> Paste CSR file : </Typography>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    fullWidth
                    // label="Service name"
                    rows="10"
                    multiline
                    onChange={handleCsrFile}
                    name="sslCert"
                    value={csrFile}
                  />
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={6}>
              <Button variant="contained" color="secondary" onClick={onUpload}>
                Submit
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <br />
      {loader ? <ProgressHOC /> : ''}
    </div>
  );
}
