import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import VaultImg from '../../../../assets/cred-vault.svg';
import Constants from '../../../../Constants';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import DecryptRoot from './VaultDecrypt';
// import Unseal from './VaultUnseal';
import Trasalogo from '../../../../assets/trasa-ni.svg';
import HVaultImg from '../../../../assets/secretstorage/vault-enterprise.svg';
import ConfigSecretProvider from './SecretProviderConfig';
import TextField from '@material-ui/core/TextField';
import FormControl from '@material-ui/core/FormControl';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },

  initButton: {
    background: '#000080', // '#0A2053', //'#0A2053',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 38,
    padding: '0 30px',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'skyblue',
    },
    // boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
  },
  initedButton: {
    background: '#025531',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 38,
    padding: '0 30px',
  },
  label: {
    textTransform: 'capitalize',
  },
  heading: {
    color: 'black',
    fontSize: '24px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  keys: {
    color: 'navy',
    fontSize: '18px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  initStatusText: {
    color: 'grey',
    fontSize: '14px',
    fontFamily: 'Open Sans, Rajdhani',
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
  divider: {
    marginTop: 50,
    marginBottom: 50,
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
  trasaLogo: {
    padding: theme.spacing(1),
    paddingTop: 40,
    // marginLeft: 30,
    minHeight: 100,
    textAlign: 'center',
    minWidth: 100,
  },
  logo: {
    padding: theme.spacing(1),
    paddingTop: 30,
    // marginLeft: 30,
    minHeight: 100,
    textAlign: 'center',
    minWidth: 100,
  },
  selectCustom: {
    width: 300,
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function SecretStorate() {
  const classes = useStyles();

  const [initStatus, setinitStatus] = useState({ status: false, initOn: '', settingValue: '{}' });
  const [tokenStatus, settokenStatus] = useState({ sealed: false });
  const getVaultStatus = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/status`)
      .then((r) => {
        if (r.data.data[0].initStatus.updatedOn) {
          const date = new Date(r.data.data[0].initStatus.updatedOn * 1000);
          r.data.data[0].initStatus.initOn = date.toDateString();
        }

        setinitStatus(r.data.data[0].initStatus);
        settokenStatus(r.data.data[0].tokenStatus);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  useEffect(() => {
    getVaultStatus();
  }, []);

  return (
    <div>
      <Typography variant="h1">TsxVault</Typography>
      <VaultConfig
        initStatus={initStatus}
        tokenStatus={tokenStatus}
        settokenStatus={settokenStatus}
        getVaultStatus={getVaultStatus}
      />
      {/* <br /> */}
      <Divider light className={classes.divider} />
      <ConfigSecretProvider credprovProps={initStatus.settingValue} />
    </div>
  );
}

export function VaultConfig(props: any) {
  const classes = useStyles();
  const [open, setopen] = useState(false);
  const [decryptKeys, setdecryptKeys] = useState([]);
  // const [sealStatus, setsealStatus] = useState({});

  const [loader, setLoader] = useState(false);
  const [keyval, setKeyval] = useState('');
  const [vaultReinitDlgState, setvaultReinitDlgState] = useState(false);

  // const handleUnsealProcessOpen = () => {
  //   setopen(true);
  // };

  const handleInitDialogueOpen = () => {
    submitVaultInitRequest();
  };

  const handleInitDialogueClose = () => {
    setopen(false);
  };

  const submitVaultInitRequest = () => {
    setLoader(true);

    const reqData = { secretShares: 5, secretThreshold: 3 };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/init`, reqData)
      .then((r) => {
        const resp = r.data.data[0];
        setLoader(false);
        if (resp.showDlg === true) {
          setdecryptKeys(resp.decryptKeys);

          setopen(true);
        }

        props.getVaultStatus();
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const handleUnsealKeyInputChange = (event: any) => {
    setKeyval(event.target.value);
  };

  const SubmitDecryptKey = () => {
    setLoader(true);
    const reqData = { key: keyval };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/decrypt`, reqData)
      .then((response) => {
        if (response.data.status === 'success') {
          props.settokenStatus(response.data.data[0]);
        }
        setLoader(false);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const confirmReinitDialogue = () => {
    setvaultReinitDlgState(true);
  };

  const handleVaultReinitDlgClose = () => {
    setvaultReinitDlgState(false);
  };

  // send delete vault request to trasacore. If the response is success,
  // we expect vaule "submit-init". If this value is present, we invoke submitInit function.
  const reinitVault = () => {
    setLoader(true);

    axios
      .delete(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/reinit`)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success' && response.data.reason === 'submit-init') {
          handleVaultReinitDlgClose();
          submitVaultInitRequest();
        }
        handleVaultReinitDlgClose();
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div className={classes.root}>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6} md={5}>
          <Paper className={classes.paper}>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={4} md={4}>
                <img src={VaultImg} alt="vault-icon" />
              </Grid>

              <Grid item xs={12} sm={8} md={8}>
                <div className={classes.heading}>
                  {props.initStatus.status
                    ? 'TRASA vault is initialized'
                    : 'TRASA vault is not initialized'}
                </div>

                <br />
                {props.initStatus.status ? (
                  <div>
                    {' '}
                    <div className={classes.initStatusText}>
                      {' '}
                      Initialized on {props.initStatus.initOn} <br />{' '}
                    </div>{' '}
                    <br />
                    <Button classes={{ root: classes.initedButton, label: classes.label }}>
                      Active
                    </Button>{' '}
                    <Button
                      classes={{ root: classes.initButton, label: classes.label }}
                      onClick={confirmReinitDialogue}
                    >
                      Re-Initialize
                    </Button>
                  </div>
                ) : (
                  <div>
                    {' '}
                    <Button
                      classes={{ root: classes.initButton, label: classes.label }}
                      onClick={handleInitDialogueOpen}
                    >
                      Initialize
                    </Button>
                  </div>
                )}

                <VaultReinitDlg
                  open={vaultReinitDlgState}
                  close={handleVaultReinitDlgClose}
                  reinitVault={reinitVault}
                />

                <ShowTokens
                  open={open}
                  decryptKeys={decryptKeys}
                  handleClose={handleInitDialogueClose}
                />
              </Grid>
              {loader ? <ProgressHOC /> : ''}
            </Grid>
          </Paper>
        </Grid>

        {/* Grid End */}
        <Grid item xs={12} sm={6} md={7}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={12} md={12}>
              <Paper className={classes.paper}>
                <div className={classes.heading}>
                  Token state: {props.tokenStatus.sealed ? 'Encrypted' : 'Retrieved'}
                </div>

                {props.tokenStatus.sealed ? (
                  <DecryptRoot
                    open={open}
                    sealStatus={props.tokenStatus}
                    SubmitDecryptKey={SubmitDecryptKey}
                    handleUnsealKeyInputChange={handleUnsealKeyInputChange}
                    loader={loader}
                  />
                ) : (
                  <Button classes={{ root: classes.initedButton, label: classes.label }}>
                    Token Retrieved
                  </Button>
                )}
              </Paper>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

export type ResourceStatsFilterProps = {
  entityType: string;
  entityID: string;
};

export type showTokenProps = {
  open: boolean;
  handleClose: () => void;
  decryptKeys: Array<string>;
};

export function ShowTokens(props: showTokenProps) {
  const classes = useStyles();
  return (
    <div>
      <Dialog
        disableBackdropClick
        disableEscapeKeyDown
        fullWidth
        maxWidth="lg"
        open={props.open}
        onClose={props.handleClose}
        aria-labelledby="form-dialog-title"
      >
        <DialogContent>
          <DialogContentText>
            <Typography variant="h3"> Decryption Keys</Typography>
            <Typography variant="h4">
              {' '}
              These Keys are only shown here once. Store them securely. <br />
              <b>Pro tip: </b>It is better to distribute these keys to multiple trusted team
              members. <br />
              If TRASA server is restarted, these keys will be required to retrieve the master key.
            </Typography>
            <br />
          </DialogContentText>

          <Divider light />
          <Grid container>
            <Grid item xs={12} sm={10} md={10}>
              <div>Decryption Keys:</div>

              {props.decryptKeys.map((v: any, i: any) => (
                <div className={classes.keys} key={i}>
                  {' '}
                  Key {i + 1} {' ) '} {'  '} {v}{' '}
                </div>
              ))}

              <br />
            </Grid>
          </Grid>

          <Divider light />

          <div>
            <br />
            <Grid container>
              <Grid item xs={10} sm={10} md={10} />

              <Grid item xs={1} sm={1} md={1}>
                <Button variant="contained" color="secondary" onClick={props.handleClose}>
                  CLose
                </Button>
              </Grid>
            </Grid>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}

const VaultReinitDlg = (props: any) => {
  const classes = useStyles();
  return (
    <div>
      <Dialog
        open={props.open}
        onClose={props.close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {' '}
          <div className={classes.Warning}> !!! WARNING !!! </div>
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Re initializing vault will delete every credentials stored in vault causing data loss
            and may lock you out from servers. Make sure you know what you are doing.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              props.close();
              props.reinitVault();
            }}
            className={classes.WarningButton}
          >
            Yes, Remove Everything from Vault.
          </Button>
          <Button onClick={props.close} className={classes.initButton} autoFocus>
            No
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
