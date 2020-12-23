import { Grid, TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import React from 'react';
import { default as ProgressHOC } from '../../../../utils/Components/Progressbar';
import axios from 'axios';
import Constants from '../../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  content: {
    width: '100%',
    flexGrow: 1,
    padding: 24,
    height: 'calc(100% - 56px)',
    marginTop: 26,
  },

  paper: {
    padding: theme.spacing(5),
    width: 700,
  },

  maxWidth: {
    width: 300,
  },
  provName: {
    color: 'black',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function SecretProvider(props: any) {
  const classes = useStyles();

  const [loader, setLoader] = React.useState(false);

  const [credprovProp, setCredprovProp] = React.useState({
    providerName: '',
    providerAddr: '',
    providerAccessToken: '',
  });

  React.useEffect(() => {
    let v = JSON.parse(props.credprovProps);
    setCredprovProp(v);
  }, [props.credprovProps]);

  function credprovnameChange(e: React.ChangeEvent<{ value: unknown }>) {
    setCredprovProp({ ...credprovProp, providerName: e.target.value as string });
  }

  // change cred prov property
  function credprovPropChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.name === 'providerAddr') {
      setCredprovProp({ ...credprovProp, providerAddr: e.target.value });
      return;
    }

    if (e.target.name === 'providerAccessToken') {
      setCredprovProp({ ...credprovProp, providerAccessToken: e.target.value });
      return;
    }
  }

  function submitReq(e: any) {
    e.preventDefault();

    setLoader(true);
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/credprov`, credprovProp)
      .then(() => {
        setLoader(false);
      })
      .catch((error) => {
        console.error(error);
        setLoader(false);
      });
  }

  return (
    <div>
      <Typography variant="h1">Service secret storage providers</Typography>
      <br />

      <Grid container spacing={2} direction="row" justify="center" alignItems="center">
        <Grid item xs={3}>
          <Typography variant="h3">Currenlty active provider: </Typography>
        </Grid>
        <Grid item xs={9}>
          <FormControl>
            <Select
              name="providerName"
              defaultValue={credprovProp.providerName}
              value={credprovProp.providerName}
              onChange={credprovnameChange}
              variant="outlined"
              inputProps={{
                classes: {
                  root: classes.maxWidth,
                },
              }}
            >
              <MenuItem value="CREDPROV_TSXVAULT">
                <p className={classes.provName}>TsxVault </p>
              </MenuItem>
              <MenuItem value="CREDPROV_HCVAULT">
                <p className={classes.provName}>Hashicorp Vault</p>
              </MenuItem>
            </Select>
          </FormControl>
        </Grid>

        <Grid container spacing={6}>
          <Divider light />
          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <Typography variant="h4">Provider Address: </Typography>
              </Grid>
              <Grid item xs={9}>
                <TextField
                  name="providerAddr"
                  variant="outlined"
                  className={classes.maxWidth}
                  onChange={credprovPropChange}
                  value={credprovProp.providerAddr}
                />
              </Grid>
            </Grid>

            <Grid container spacing={2}>
              <Grid item xs={3}>
                <Typography variant="h4">Provider Access Token: </Typography>
              </Grid>
              <Grid item xs={9}>
                <TextField
                  name="providerAccessToken"
                  variant="outlined"
                  type="password"
                  className={classes.maxWidth}
                  onChange={credprovPropChange}
                  value={credprovProp.providerAccessToken}
                />
              </Grid>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid container spacing={2} direction="column" justify="center" alignItems="center">
              <Grid item xs={3}>
                <Button variant="contained" onClick={submitReq}>
                  Submit
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      <br />
      <br />
      <br />
      {loader && <ProgressHOC />}
    </div>
  );
}
