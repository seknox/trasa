import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useState } from 'react';
import Constants from '../../../../../Constants';
import ProgressHOC from '../../../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    '& .MuiTextField-root': {
      margin: theme.spacing(1),
      width: 200,
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
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },
}));

export default function SSLCerts(props: any) {
  const [certs, setCerts] = useState({ sslCert: '', sslKey: '', caCert: '', serviceID: '' });
  // const [sslCert, setsslCert] = useState('')
  // const [sslKey, setsslKey] = useState('')
  // const [caCert, setcaCert] = useState('')
  const [loader, setLoader] = useState(false);
  const classes = useStyles();

  const handleChange = (name: any) => (event: any) => {
    const val = event.target.value;
    setCerts({ ...certs, [name]: val });
  };

  const handleSubmit = (event: any) => {
    setLoader(true);

    const req = { ...certs, serviceID: props.serviceID };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/sslcerts/update/${props.serviceID}`, req)
      .then((response) => {
        setLoader(false);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Typography variant="h3">Upload your certificate</Typography>
      </Grid>
      <Grid item xs={3}>
        <Typography variant="h4">SSL Cert :</Typography>
      </Grid>
      <Grid item xs={9}>
        <TextField
          fullWidth
          // label="Service name"
          rows="10"
          multiline
          onChange={handleChange('sslCert')}
          name="sslCert"
          value={certs.sslCert}
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
      <Grid item xs={3}>
        <Typography variant="h4">SSL Key :</Typography>
      </Grid>
      <Grid item xs={9}>
        <TextField
          fullWidth
          // label="Service name"
          rows="10"
          multiline
          onChange={handleChange('sslKey')}
          name="sslKey"
          value={certs.sslKey}
          InputProps={{
            // disableUnderline: true,
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
      <Grid item xs={3}>
        <Typography variant="h4">CA Cert :</Typography>
      </Grid>
      <Grid item xs={9}>
        <TextField
          fullWidth
          // label="Service name"
          multiline
          rows="10"
          id="standard-multiline-static"
          onChange={handleChange('caCert')}
          name="caCert"
          value={certs.caCert}
          InputProps={{
            // disableUnderline: true,
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
      <Grid container spacing={2} direction="column" justify="flex-end" alignItems="center">
        <Grid item xs={12}>
          <br />
          <br />
          <Button variant="contained" color="secondary" type="submit" onClick={handleSubmit}>
            Submit
          </Button>
          <br />
          {loader ? <ProgressHOC /> : ''}
        </Grid>
      </Grid>
    </Grid>
  );
}
