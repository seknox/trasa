import { Grid, TextField } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import Checkbox from '@material-ui/core/Checkbox';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Fab from '@material-ui/core/Fab';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import CopyIcon from '@material-ui/icons/FileCopy';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import CopyToClipboard from 'react-copy-to-clipboard';
import Constants from '../../../../Constants';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import ScrollableDlg from '../../../../utils/Components/ScrollableDlg';
import Progress from '../../../../utils/Components/Progressbar';

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
    padding: theme.spacing(5),
    minWidth: 800,
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
  fab: {
    fontSize: '10px',
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },
}));

type Anchor = 'top' | 'left' | 'bottom' | 'right';

export default function IdpConfigDrawer(props: any) {
  const classes = useStyles();

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

  function returnIdp(idpDetail: any) {
    switch (idpDetail.idpType) {
      case 'saml':
        return <SAMLIdP idpDetail={idpDetail} />;
      case 'ldap':
        return <LdapIdp idpDetail={idpDetail} />;
      default:
        return '';
    }
  }

  return (
    <div>
      <Button fullWidth variant="contained" size="small" onClick={toggleDrawer('right', true)}>
        Configure
      </Button>
      <Drawer anchor="right" open={state.right} onClose={toggleDrawer('right', false)}>
        <Paper className={classes.paper}>{returnIdp(props.idpDetail)}</Paper>
      </Drawer>
    </div>
  );
}

function SAMLIdP(props: any) {
  const [loader, setLoader] = React.useState(false);
  const [idp, setIdpVal] = React.useState({
    idpID: '',
    idpName: '',
    idpType: '',
    idpMeta: '',
    endpoint: '',
    isEnabled: true,
  });

  function idpChange(e: React.ChangeEvent<HTMLInputElement>) {
    setIdpVal({ ...idp, [e.target.name]: e.target.value });
  }

  const submitIdp = () => {
    setLoader(true);

    idp.idpName = props.idpDetail.idpName;
    idp.idpID = props.idpDetail.idpID;
    idp.idpType = props.idpDetail.idpType;
    if (idp.idpMeta === '') {
      idp.idpMeta = props.idpDetail.idpMeta;
    }
    if (idp.endpoint === '') {
      idp.endpoint = props.idpDetail.endpoint;
    }
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/idp/external/update`, idp)
      .then((r) => {
        setLoader(false);

        if (r.data.status === 'success') {
          setLoader(false);
        }
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  return (
    <div>
      <SAMLConfig idp={props.idpDetail} setIdpVal={idpChange} submitIdp={submitIdp} />
      <SCIMConfig idp={props.idpDetail} setIdpVal={idpChange} submitIdp={submitIdp} />
      <br />
      {loader && <ProgressHOC />}
      <br />
      <DisableIdp idp={props.idpDetail} />
    </div>
  );
}

function SAMLConfig(props: any) {
  const classes = useStyles();
  const [clipText, setClipText] = useState('');

  return (
    <div>
      <Grid container spacing={6}>
        <Grid item xs={12}>
          <Typography variant="h2">SAML configuration</Typography>
          <Divider light />
        </Grid>

        <Divider light />
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">SAML Identity Provider :</Typography>
            </Grid>
            <Grid item xs={9}>
              <TextField
                name="idpName"
                fullWidth
                value={props.idp.idpName}
                disabled
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    // root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">Embed Link: </Typography>
            </Grid>
            <Grid item xs={9}>
              <TextField
                name="endpoint"
                // label={<Typography variant="h5">paste SAML embed link here</Typography>}
                fullWidth
                onChange={props.setIdpVal}
                defaultValue={props.idp.endpoint}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">Audience URI: </Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField
                disabled
                name="audienceURI"
                fullWidth
                value={props.idp.audienceURI}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
            <Grid item xs={1}>
              <CopyToClipboard text={clipText}>
                <Fab
                  size="small"
                  onClick={() => setClipText(props.idp.audienceURI)}
                  aria-label="copy"
                  className={classes.fab}
                  component="button"
                >
                  <CopyIcon />
                </Fab>
              </CopyToClipboard>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">Redirect Url: </Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField disabled name="redirectURL" fullWidth value={props.idp.redirectURL} />
            </Grid>
            <Grid item xs={1}>
              <CopyToClipboard text={clipText}>
                <Fab
                  size="small"
                  onClick={() => setClipText(props.idp.redirectURL)}
                  aria-label="copy"
                  className={classes.fab}
                  component="button"
                >
                  <CopyIcon />
                </Fab>
              </CopyToClipboard>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">IDP Metadata: </Typography>
            </Grid>
            <Grid item xs={9}>
              <TextField
                multiline
                fullWidth
                rows="8"
                name="idpMeta"
                defaultValue={props.idp.idpMeta}
                onChange={props.setIdpVal}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2} direction="column" justify="center" alignItems="center">
            <Grid item xs={3}>
              <Button variant="contained" onClick={props.submitIdp}>
                Update
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function SCIMConfig(props: any) {
  const classes = useStyles();
  const [clipText, setClipText] = useState('');
  const [key, setkey] = useState(props.idp.apiKey);
  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
  });

  const generateKey = () => {
    updateActionStatus({ ...actionStatus, respStatus: false, statusMsg: '', loader: true });
    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };

    axios
      .post(
        `${Constants.TRASA_HOSTNAME}/api/v1/idp/external/generatescimtoken/${props.idp.idpID}`,
        '',
        config,
      )
      .then((r) => {
        updateActionStatus({ ...actionStatus, loader: false });
        if (r.data.status === 'success') {
          setkey(r.data.data[0]);
          updateActionStatus({
            ...actionStatus,
            success: true,
            respStatus: true,
            statusMsg: r.data.reason,
          });
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
        console.log(error);
      });
  };

  return (
    <div>
      <Grid container spacing={6}>
        <Grid item xs={12}>
          <Typography variant="h2">SCIM configuration</Typography>
          <Divider light />
        </Grid>

        <Divider light />
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">SCIM endpoint :</Typography>
            </Grid>
            <Grid item xs={9}>
              <TextField
                name="scimEndpoint"
                fullWidth
                defaultValue={props.idp.scimEndpoint}
                disabled
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">SCIM api key :</Typography>
            </Grid>
            <Grid item xs={6}>
              <TextField
                name="apiKey"
                fullWidth
                defaultValue={key}
                type="password"
                value={key}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
            <Grid item xs={1}>
              <CopyToClipboard text={clipText}>
                <Fab
                  size="small"
                  onClick={() => setClipText(key)}
                  aria-label="copy"
                  className={classes.fab}
                  component="button"
                >
                  <CopyIcon />
                </Fab>
              </CopyToClipboard>
            </Grid>
            <Grid item xs={2}>
              <CopyToClipboard text={clipText}>
                <Button
                  variant="contained"
                  onClick={generateKey}
                  aria-label="copy"
                  className={classes.fab}
                >
                  Generate New Key
                </Button>
              </CopyToClipboard>
            </Grid>
          </Grid>
        </Grid>
        {actionStatus.loader ? <ProgressHOC /> : ''}
      </Grid>
    </div>
  );
}

/// /////////////////////////////////////////////////////////////
/// ////////////////       LDAP IDP        //////////////////////
/// /////////////////////////////////////////////////////////////
function LdapIdp(props: any) {
  const [loader, setLoader] = React.useState(false);

  const [idp, setIdpVal] = React.useState({
    idpID: '',
    idpName: '',
    idpType: '',
    idpMeta: '',
    clientID: '',
    audienceURI: '',
    endpoint: '',
    isEnabled: true,
  });

  function idpChange(e: React.ChangeEvent<HTMLInputElement>) {
    setIdpVal({ ...idp, [e.target.name]: e.target.value });
  }

  const submitIdp = () => {
    setLoader(true);
    idp.idpName = props.idpDetail.idpName;
    idp.idpID = props.idpDetail.idpID;
    idp.idpType = props.idpDetail.idpType;
    if (idp.idpMeta === '') {
      idp.idpMeta = props.idpDetail.idpMeta;
    }
    if (idp.endpoint === '') {
      idp.endpoint = props.idpDetail.endpoint;
    }
    if (idp.clientID === '') {
      idp.clientID = props.idpDetail.clientID;
    }
    if (idp.audienceURI === '') {
      idp.audienceURI = props.idpDetail.audienceURI;
    }
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/idp/external/update`, idp)
      .then(() => {
        setLoader(false);
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  return (
    <div>
      <LdapConfig idp={props.idpDetail} setIdpVal={idpChange} submitIdp={submitIdp} />
      {loader ? <Progress /> : ''}
      <LdapSync idpDetail={props.idpDetail} setIdpVal={idpChange} submitIdp={submitIdp} />
      <br />

      <br />
      <MigrateUserDlg idp={props.idpDetail} />
      <br />
      <DisableIdp idp={props.idpDetail} />
    </div>
  );
}

function LdapConfig(props: any) {
  const classes = useStyles();

  return (
    <div>
      <Grid container spacing={6}>
        <Grid item xs={12}>
          <Typography variant="h2">LDAP configuration</Typography>
          <Divider light />
        </Grid>

        <Divider light />
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">LDAP Identity Provider :</Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField
                name="idpName"
                fullWidth
                value={props.idp.idpName}
                disabled
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">Server Domain: </Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField
                name="endpoint"
                // label={<Typography variant="h5">paste SAML embed link here</Typography>}
                fullWidth
                onChange={props.setIdpVal}
                defaultValue={props.idp.endpoint}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">LDAP DN (search base): </Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField
                name="idpMeta"
                fullWidth
                onChange={props.setIdpVal}
                defaultValue={props.idp.idpMeta}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          {/* <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">User group: </Typography>
            </Grid>
            <Grid item xs={8}>
              <TextField
                name="audienceURI"
                fullWidth
                onChange={props.setIdpVal}
                defaultValue={props.idp.audienceURI}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid> */}

          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">Service account username: </Typography>
            </Grid>
            <Grid item xs={5}>
              <TextField
                fullWidth
                name="clientID"
                defaultValue={props.idp.clientID}
                onChange={props.setIdpVal}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h4">Service account password: </Typography>
            </Grid>
            <Grid item xs={5}>
              <TextField
                fullWidth
                name="clientSecret"
                type="password"
                defaultValue={props.idp.clientSecret}
                onChange={props.setIdpVal}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2} direction="column" justify="center" alignItems="center">
            <Grid item xs={3}>
              <Button variant="contained" onClick={props.submitIdp}>
                Update
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function LdapSync(props: any) {
  const classes = useStyles();
  const [respData, setRespData] = useState('');
  const [idp, setIdpVal] = React.useState({
    idpID: '',
    idpName: '',
    idpType: '',
    idpMeta: '',
    clientID: '',
    audienceURI: '',
    endpoint: '',
    isEnabled: true,
  });

  function idpChange(e: React.ChangeEvent<HTMLInputElement>) {
    setIdpVal({ ...idp, [e.target.name]: e.target.value });
  }
  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
    openDlg: false,
  });

  const importLdapUsers = () => {
    updateActionStatus({ ...actionStatus, respStatus: false, statusMsg: '', loader: true });

    idp.idpName = props.idpDetail.idpName;
    idp.idpID = props.idpDetail.idpID;
    idp.idpType = props.idpDetail.idpType;
    if (idp.idpMeta === '') {
      idp.idpMeta = props.idpDetail.idpMeta;
    }
    if (idp.endpoint === '') {
      idp.endpoint = props.idpDetail.endpoint;
    }
    if (idp.clientID === '') {
      idp.clientID = props.idpDetail.clientID;
    }
    if (idp.audienceURI === '') {
      idp.audienceURI = props.idpDetail.audienceURI;
    }

    const respStr = (val: any) => {
      const v =
        `Total imported users: ${val.totalUsers}\n` +
        `. Total import failed : ${val.failedCount}\n` +
        `. Following users were not imported: \n${val.failedUsers}`;

      return v;
    };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/idp/external/ldap/importusers`, idp)
      .then((r) => {
        updateActionStatus({ ...actionStatus, loader: false });
        if (r.data.status === 'success') {
          setRespData(respStr(r.data.data[0]));
          updateActionStatus({
            ...actionStatus,
            openDlg: true,
            success: true,
            respStatus: true,
            statusMsg: r.data.reason,
          });
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
        console.log(error);
      });
  };

  const closeDlg = () => {
    updateActionStatus({ ...actionStatus, openDlg: false });
  };

  return (
    <div>
      <Grid container spacing={6}>
        <Grid item xs={12}>
          <Typography variant="h2">Import users</Typography>
          <Divider light />
        </Grid>

        <Divider light />
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={3}>
              <Typography variant="h4">Import user(s) :</Typography>
            </Grid>
            <Grid item xs={6}>
              <TextField
                name="audienceURI"
                fullWidth
                defaultValue=""
                onChange={idpChange}
                // type="password"
                // value={key}
                InputProps={{
                  disableUnderline: true,
                  classes: {
                    root: classes.textFieldRoot,
                    input: classes.textFieldInputBig,
                  },
                }}
              />
            </Grid>
            <Grid item xs={3}>
              <Button
                variant="contained"
                onClick={importLdapUsers}
                aria-label="copy"
                className={classes.fab}
              >
                Import
              </Button>
            </Grid>
          </Grid>

          <ScrollableDlg
            open={actionStatus.openDlg}
            data={respData}
            handleClose={closeDlg}
            maxWidth="lg"
            title=""
            children=""
          />
        </Grid>
        {actionStatus.loader ? <ProgressHOC /> : ''}
      </Grid>
    </div>
  );
}

function MigrateUserDlg(props: any) {
  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
    openDlg: false,
  });
  const closeDlg = () => {
    updateActionStatus({ ...actionStatus, openDlg: false });
  };

  return (
    <Grid container spacing={2} justify="center" alignItems="center">
      <Grid item xs={12}>
        <br />
        <Typography variant="h2">{`Convert trasa users to/from ${props.idp.idpName}`}</Typography>
        <Divider light />
      </Grid>
      <Button
        variant="contained"
        onClick={() => updateActionStatus({ ...actionStatus, openDlg: true })}
      >
        Click to Open Transfer List
      </Button>
      <ScrollableDlg
        open={actionStatus.openDlg}
        data=""
        title={`Convert trasa users to/from ${props.idp.idpName}`}
        handleClose={closeDlg}
        maxWidth="lg"
      >
        {' '}
        <MigrateUsers idp={props.idp} />
      </ScrollableDlg>
    </Grid>
  );
}

function not(a: number[], b: number[]) {
  return a.filter((value) => b.indexOf(value) === -1);
}

function intersection(a: number[], b: number[]) {
  return a.filter((value) => b.indexOf(value) !== -1);
}

function union(a: number[], b: number[]) {
  return [...a, ...not(b, a)];
}

function MigrateUsers(props: any) {
  const [checked, setChecked] = React.useState<number[]>([]);

  const [trasaUsers, setTrasaUsers] = useState<any>([]);
  const [extIdpUsers, setExtIdpUsers] = useState<any>([]);
  const [updatedTrasaUsers, updateTrasaUsers] = useState<any>([]);
  const [updatedExtIdpUsers, updatesetExtIdpUsers] = useState<any>([]);

  const leftChecked = intersection(checked, trasaUsers);
  const rightChecked = intersection(checked, extIdpUsers);

  const [actionStatus, updateActionStatus] = useState({
    respStatus: false,
    success: false,
    loader: false,
    statusMsg: '',
    openDlg: false,
  });

  const handleToggle = (value: number) => () => {
    const currentIndex = checked.indexOf(value);
    const newChecked = [...checked];

    if (currentIndex === -1) {
      newChecked.push(value);
    } else {
      newChecked.splice(currentIndex, 1);
    }

    setChecked(newChecked);
  };

  const numberOfChecked = (items: number[]) => intersection(checked, items).length;

  const handleToggleAll = (items: number[]) => () => {
    if (numberOfChecked(items) === items.length) {
      setChecked(not(checked, items));
    } else {
      setChecked(union(checked, items));
    }
  };

  const handleCheckedRight = () => {
    setExtIdpUsers(extIdpUsers.concat(leftChecked));
    updatesetExtIdpUsers(updatedExtIdpUsers.concat(leftChecked));
    setTrasaUsers(not(trasaUsers, leftChecked));
    setChecked(not(checked, leftChecked));
  };

  const handleCheckedLeft = () => {
    setTrasaUsers(trasaUsers.concat(rightChecked));
    updateTrasaUsers(updatedTrasaUsers.concat(rightChecked));
    setExtIdpUsers(not(extIdpUsers, rightChecked));
    setChecked(not(checked, rightChecked));
  };

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/idp/users/all/trasa`)

      .then((response) => {
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          const ddataArr = data.map(function (n: any) {
            return [n.firstName, n.lastName, n.email.toString(), n.idpName, n.ID];
          });
          setTrasaUsers(ddataArr);
        }
      });

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/idp/users/all/${props.idp.idpName}`)

      .then((response) => {
        if (response.data.status === 'success') {
          const data = response.data.data[0];
          const ddataArr = data.map(function (n: any) {
            return [n.firstName, n.lastName, n.email.toString(), n.idpName, n.ID];
          });
          setExtIdpUsers(ddataArr);
        }
      });
  }, [props.idp.idpName]);

  const submitTrasaUsers = (idpName: string) => {
    updateActionStatus({ ...actionStatus, loader: true });
    const req = { idpName: '', userList: [] };

    req.idpName = idpName;
    if (idpName === 'trasa') {
      req.userList = updatedTrasaUsers;
    } else {
      req.userList = updatedExtIdpUsers;
    }

    const config = {
      headers: {
        'X-SESSION': localStorage.getItem('X-SESSION'),
        'X-CSRF': localStorage.getItem('X-CSRF'),
      },
    };
    axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/idp/users/transfer`, req, config).then(() => {
      updateActionStatus({ ...actionStatus, loader: false });
    });
  };

  const customList = (title: React.ReactNode, items: number[]) => (
    <Card>
      <CardHeader
        avatar={
          <Checkbox
            onClick={handleToggleAll(items)}
            checked={numberOfChecked(items) === items.length && items.length !== 0}
            indeterminate={numberOfChecked(items) !== items.length && numberOfChecked(items) !== 0}
            disabled={items.length === 0}
            inputProps={{ 'aria-label': 'all items selected' }}
          />
        }
        title={title}
        subheader={`${numberOfChecked(items)}/${items.length} selected`}
      />
      <Divider />
      <List dense component="div" role="list">
        {items.length > 0 ? (
          items.map((value: any) => {
            const labelId = `transfer-list-all-item-${value}-label`;

            return (
              <ListItem key={value} role="listitem" button onClick={handleToggle(value)}>
                <ListItemIcon>
                  <Checkbox
                    checked={checked.indexOf(value) !== -1}
                    tabIndex={-1}
                    disableRipple
                    inputProps={{ 'aria-labelledby': labelId }}
                  />
                </ListItemIcon>
                <ListItemText id={labelId} primary={`${value[0]} ${value[1]} [${value[2]}]`} />
              </ListItem>
            );
          })
        ) : (
          <div>
            <br />
            <br />
            <ProgressHOC />
          </div>
        )}
      </List>
    </Card>
  );

  return (
    <Grid container spacing={2} direction="row" justify="center" alignItems="center">
      <Grid item>
        {customList('TRASA users', trasaUsers)}
        <br />
        <Button variant="contained" onClick={() => submitTrasaUsers('trasa')}>
          Update TRASA users
        </Button>
      </Grid>
      <Grid item>
        <Grid container direction="column" alignItems="center">
          <Button
            variant="outlined"
            size="small"
            onClick={handleCheckedRight}
            disabled={leftChecked.length === 0}
            aria-label="move selected right"
          >
            &gt;
          </Button>
          <Button
            variant="outlined"
            size="small"
            onClick={handleCheckedLeft}
            disabled={rightChecked.length === 0}
            aria-label="move selected left"
          >
            &lt;
          </Button>
        </Grid>
      </Grid>
      <Grid item>
        {customList(`${props.idp.idpName} users`, extIdpUsers)}
        <br />
        <Button variant="contained" onClick={() => submitTrasaUsers(props.idp.idpName)}>
          Update {props.idp.idpName} users
        </Button>
      </Grid>
      <Grid item xs={12}>
        {actionStatus.loader ? (
          <div>
            <br />
            <ProgressHOC />
          </div>
        ) : (
          ''
        )}
      </Grid>
    </Grid>
  );
}

/// /////////////////////////////////////////////////
/// ///////  Delete or Disable IDP  /////////////////
/// /////////////////////////////////////////////////
function DisableIdp(props: any) {
  const [loader, setLoader] = useState(false);
  const [active, setActive] = useState(props.idp.isEnabled);

  function disableIdp(state: any) {
    setLoader(true);

    const req = { idpID: props.idp.idpID, active: state };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/idp/external/activateordisable`, req)
      .then(() => {
        setLoader(false);
      })
      .catch((error) => {
        console.log(error);
        setLoader(false);
      });
  }

  function changeState(e: React.ChangeEvent<HTMLInputElement>) {
    // console.log('state: ', e.target.checked)
    setActive(e.target.checked);
    disableIdp(e.target.checked);
  }

  return (
    <Grid container spacing={2} justify="center" alignItems="center">
      <Grid item xs={12}>
        <br />
        <Typography variant="h2">{`Enable or Disable ${props.idp.idpName}`} </Typography>
        <Divider light />
      </Grid>

      <FormControl component="fieldset">
        <FormControlLabel
          // value={active}
          control={<Switch checked={active} color="secondary" onChange={changeState} />}
          label={props.idp.isEnabled ? 'active' : 'disabled'}
          labelPlacement="end"
        />
      </FormControl>

      {/* <Button variant="contained" className={classes.warningButton}>
      remove
    </Button> */}
      <Grid item xs={12}>
        {loader ? (
          <div>
            <br />
            <ProgressHOC />
          </div>
        ) : (
          ''
        )}
      </Grid>
    </Grid>
  );
}
