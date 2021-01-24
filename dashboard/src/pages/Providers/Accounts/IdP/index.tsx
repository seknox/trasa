import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import AzureAD from '../../../../assets/idp/azuread.png';
import FreeIPA from '../../../../assets/idp/freeipa.png';
import GSuite from '../../../../assets/idp/gsuite.png';
import JumpCloud from '../../../../assets/idp/jumpcloud.png';
import Okta from '../../../../assets/idp/okta.png';
import Onelogin from '../../../../assets/idp/onelogin.png';
// import servicesetting from './AssignUserToApp';
// import Constants from '../../../../Constants';
import OpenLdapIcon from '../../../../assets/idp/openldap.png';
import Trasalogo from '../../../../assets/trasa-ni.svg';
import Constants from '../../../../Constants';
import CreateIdpDrawer from './CreateIdpDrawer';
import ConfigureIdpDrawer from './IdpConfDrawer';

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(5),
  },

  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    // textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  idpPaper: {
    // backgroundColor:  '#fdfdfd',
    padding: theme.spacing(2),
    textAlign: 'center',
    minWidth: 100,
    maxWidth: 200,
    minHeight: 190,
    color: theme.palette.text.secondary,
  },
  servicesDemiter: {
    marginBottom: 10,
  },
  logoPaper: {
    padding: theme.spacing(1),
    paddingTop: 40,
    // marginLeft: 30,
    minHeight: 100,
    textAlign: 'center',
    minWidth: 100,
  },
}));

export default function IdpPage() {
  const classes = useStyles();

  const [idps, setIdps] = useState([]);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/providers/uidp/all`)

      .then((response) => {
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          setIdps(data);
          // console.log('ppppppppppp: ', data)
        }
      });
  }, []);

  return (
    <div className={classes.root}>
      <CreateIdpDrawer setIdps={setIdps} idps={idps} />
      <br />
      <br />
      <br />
      <Grid container spacing={2} direction="row">
        <Grid item xs={12}>
          <div className={classes.servicesDemiter}>
            <Typography variant="h2">Configured Identity Provider(s)</Typography>
            <Divider light />{' '}
          </div>
        </Grid>

        <Grid item xs={3}>
          <Paper className={classes.idpPaper}>
            <Grid container spacing={2} direction="column" alignItems="center" justify="center">
              <Grid item xs={12}>
                <Paper className={classes.logoPaper}>
                  <img src={Trasalogo} alt="Trasa" />
                </Paper>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="h3">Built-in, default</Typography>
              </Grid>
            </Grid>
          </Paper>
        </Grid>

        {idps.map((v, k) => (
          <Grid item xs={3}>
            <RenderIdP idps={v} k={k} />
          </Grid>
        ))}
      </Grid>
    </div>
  );
}

// const idps = ['Okta', 'JumpCloud', 'GSuite', 'OpenLdapIcon', 'AzureAD', 'FreeIPA'];

function RenderIdP(props: any) {
  const classes = useStyles();

  function returnIcon(name: string, type: string) {
    // console.log(val)
    switch (true) {
      case name === 'okta':
        return <img src={Okta} alt={name} height={70} />;
      case name === 'onelogin':
        return (
          <Paper className={classes.logoPaper}>
            <img src={Onelogin} alt={name} height={20} />
          </Paper>
        );

      case name === 'freeipa':
        return <img src={FreeIPA} alt={name} height={70} />;
      case name === 'gsuite':
        return <img src={GSuite} alt={name} height={70} />;
      case name === 'jumpcloud':
        return <img src={JumpCloud} alt={name} height={70} />;
      case name === 'azuread':
        return <img src={AzureAD} alt={name} height={70} />;
      case name === 'ad':
        return <img src={AzureAD} alt={name} height={70} />;
      case type === 'ldap':
        return <img src={OpenLdapIcon} alt={name} height={70} />;
      default:
        return '';
    }
  }

  return (
    <Paper className={classes.idpPaper}>
      <Grid container direction="column" alignItems="center" justify="center" key={props.k}>
        <Grid item xs={12}>
          {returnIcon(props.idps.idpName, props.idps.idpType)}
        </Grid>
        <Grid item xs={12}>
          <Typography variant="h3">
            {props.idps.idpName === 'ad' ? 'Active Directory' : props.idps.idpName}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <ConfigureIdpDrawer idpDetail={props.idps} />
        </Grid>
      </Grid>
    </Paper>
  );
}
