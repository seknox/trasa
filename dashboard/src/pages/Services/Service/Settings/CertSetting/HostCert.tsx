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
import Constants from '../../../../../Constants';
import ProgressHOC from '../../../../../utils/Components/Progressbar';

const fileDownload = require('js-file-download');

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
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
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },
}));

export default function HostCerts(props: any) {
  const [loader, setLoader] = useState(false);
  const [certVal, setCert] = useState();
  // const [keyFile, setKey] = useState();
  // const [csrFile, setCsr] = useState();
  const classes = useStyles();

  const handleCertFile = (e: any) => {
    setCert(e.target.value);
  };

  const onDownload = () => {
    setLoader(true);

    const config = {
      responseType: 'blob',
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/services/hostcerts/download/${props.serviceID}`,
        config,
      )
      .then((r) => {
        setLoader(false);
        fileDownload(r.data, 'server-certs.zip', 'application/zip');
      });
  };

  const onUpload = () => {
    setLoader(true);
    const formData = {
      certVal,
      serviceID: props.serviceID,
    };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/hostcerts/update`, formData)
      .then((r) => {
        if (r.data.status === 'success') {
          setLoader(false);
        } else {
          setLoader(false);
        }
      });
  };

  return (
    <div className={classes.root}>
      <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
        <Grid item xs={12}>
          <Typography variant="h3"> Upload SSH Host certificate </Typography>
        </Grid>
        <br />
        <Grid item xs={12}>
          <Divider light />
          <br /> <br />
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Grid container spacing={3}>
                <Grid item xs={5}>
                  <Typography variant="h4"> Upload host key of upstream server : </Typography>
                  <Typography variant="h5"> (Must be known hosts public key format) </Typography>
                </Grid>
                <Grid item xs={7}>
                  <TextField
                    fullWidth
                    // label="App Name"
                    rows="10"
                    multiline
                    onChange={handleCertFile}
                    name="sslCert"
                    value={certVal}
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

            <Grid container spacing={2} direction="column" justify="flex-end" alignItems="center">
              <Grid item xs={12}>
                <br />
                <br />
                <Button variant="contained" color="secondary" type="submit" onClick={onUpload}>
                  Upload
                </Button>
                <br /> <br />
                <br />
              </Grid>
            </Grid>
          </Grid>
          <br />
          <Divider light />
          <br /> <br />
          <Grid item xs={12}>
            <Typography variant="h4"> Download host key signed by SSH Host CA : </Typography>
            <Typography variant="h5">
              {' '}
              (setting available at TRASA Certificate authority page){' '}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Grid container spacing={2} direction="column" justify="flex-end" alignItems="center">
              <Grid item xs={12}>
                <br />
                <br />
                <Button variant="contained" color="secondary" type="submit" onClick={onDownload}>
                  Generate and Download
                </Button>
                <br />
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      <br />
      {loader ? <ProgressHOC /> : ''}
    </div>
  );
}
