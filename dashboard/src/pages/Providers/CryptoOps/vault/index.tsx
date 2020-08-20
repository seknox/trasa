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

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },

  paperSmaller: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    marginTop: '20%',

    // marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
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
  paperLarge: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 400,
    minHeight: 300,
    // height: '100%',
    // maxWidth: '100%',
    // maxHeight: '100%',
    // marginLeft: '5%',
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperHeighted: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)'
    minWidth: 800,
    minHeight: 300,
    // height: '100%',
    // maxWidth: '100%',
    // maxHeight: '100%',
    marginTop: '5%',
    marginBotton: '5%',
    padding: theme.spacing(2),
    textAlign: 'center',
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
    fontSize: '16px',
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
}));

export default function VaultConfig() {
  const classes = useStyles();
  const [open, setopen] = useState(false);
  const [unsealKeys, setunsealKeys] = useState([]);
  const [decryptKeys, setdecryptKeys] = useState([]);
  const [encRootToken, setencRootToken] = useState('');
  // const [sealStatus, setsealStatus] = useState({});
  const [initStatus, setinitStatus] = useState({ status: false, initOn: '' });
  const [tokenStatus, settokenStatus] = useState({ sealed: false });
  const [tsxVault, settsxVault] = useState(false);
  const [loader, setLoader] = useState(false);
  const [keyval, setKeyval] = useState('');
  const [vaultReinitDlgState, setvaultReinitDlgState] = useState(false);

  useEffect(() => {
    getVaultStatus();
  }, []);

  // const handleUnsealProcessOpen = () => {
  //   setopen(true);
  // };

  const handleInitDialogueOpen = () => {
    submitVaultInitRequest();
  };

  const handleInitDialogueClose = () => {
    setopen(false);
  };

  const getVaultStatus = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/status`)
      .then((response) => {
        if (response.data.data[0].initStatus.updatedOn) {
          const date = new Date(response.data.data[0].initStatus.updatedOn * 1000);
          response.data.data[0].initStatus.initOn = date.toDateString();
        }

        setinitStatus(response.data.data[0].initStatus);
        // setsealStatus(response.data.data[0].sealStatus);
        settokenStatus(response.data.data[0].tokenStatus);
        settsxVault(response.data.data[0].tsxvault);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const submitVaultInitRequest = () => {
    setLoader(true);

    const reqData = { secretShares: 5, secretThreshold: 3 };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/init`, reqData)
      .then((response) => {
        const resp = response.data.data[0];

        if (resp.tsxvault) {
          settsxVault(true);
          setunsealKeys(resp.unsealKeys);
          setLoader(false);
          setopen(true);
        } else {
          setdecryptKeys(resp.decryptKeys);
          setencRootToken(resp.encRootToken);
          settsxVault(true);
          setunsealKeys(resp.unsealKeys);
          setLoader(false);
          setopen(true);
        }

        getVaultStatus();
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const handleUnsealKeyInputChange = (event: any) => {
    setKeyval(event.target.value);
  };

  // const submitUnsealKey = () => {
  //   setunsealProgress(true);

  //   const reqData = { key: keyval };

  //   axios
  //     .post(`${Constants.TRASA_HOSTNAME}/api/v1/crypto/vault/unseal`, reqData)
  //     .then((response) => {
  //       if (response.data.status === 'success') {
  //         setsealStatus(response.data.data[0]);
  //       }

  //       setunsealProgress(false);
  //     })
  //     .catch((error) => {
  //       console.log(error);
  //     });
  // };

  const SubmitDecryptKey = () => {
    setLoader(true);
    const reqData = { key: keyval };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/providers/vault/tsxvault/decrypt`, reqData)
      .then((response) => {
        if (response.data.status === 'success') {
          settokenStatus(response.data.data[0]);
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
                  {initStatus.status
                    ? 'TRASA vault is initialized'
                    : 'TRASA vault is not initialized'}
                </div>

                <br />
                {initStatus.status ? (
                  <div>
                    {' '}
                    <div className={classes.initStatusText}>
                      {' '}
                      vault was initialized on {initStatus.initOn} <br />{' '}
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
                  unsealKeys={unsealKeys}
                  decryptKeys={decryptKeys}
                  handleClose={handleInitDialogueClose}
                  encRootToken={encRootToken}
                  tsxVault={tsxVault}
                />
              </Grid>
              {loader ? <ProgressHOC /> : ''}
            </Grid>
          </Paper>
        </Grid>

        {/* Grid End */}
        <Grid item xs={12} sm={6} md={7}>
          <Grid container spacing={2}>
            {/* {this.state.tsxVault ? null : (
              <Grid item xs={12} sm={12} md={12}>
                <Paper className={classes.paper}>
                  <div className={classes.heading}>
                    Sealed state: {this.state.sealStatus.sealed ? 'Sealed' : 'Unsealed'}
                  </div>

                  {this.state.sealStatus.sealed ? (
                    <Unseal
                      open={this.state.open}
                      handleClose={this.handleClose}
                      sealStatus={this.state.sealStatus}
                      submitUnsealKey={this.submitUnsealKey}
                      handleUnsealKeyInputChange={this.handleUnsealKeyInputChange}
                      unsealProgress={this.state.unsealProgress}
                    />
                  ) : (
                    <Button classes={{ root: classes.initedButton, label: classes.label }}>
                      Unsealed
                    </Button>
                  )}
                </Paper>
              </Grid>
            )} */}

            <Grid item xs={12} sm={12} md={12}>
              <Paper className={classes.paper}>
                <div className={classes.heading}>
                  Token state: {tokenStatus.sealed ? 'Encrypted' : 'Retrieved'}
                </div>

                {tokenStatus.sealed ? (
                  <DecryptRoot
                    open={open}
                    // handleClose={handleClose}
                    sealStatus={tokenStatus}
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

function ShowTokens(props: any) {
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
            <Typography variant="h3"> Secret Tokens</Typography>
            <p>
              {' '}
              These Keys are only shown here once. Distribute these keys to relevent parties and
              store it secretly.
            </p>
            <p> These Keys will be required every time TRASA service is rebooted.</p>
            <br />
          </DialogContentText>

          <Divider light />
          <Grid container>
            <Grid item xs={12} sm={10} md={10}>
              <div>Decryption Keys:</div>

              {props.unsealKeys.map((v: any, i: any) => (
                <div className={classes.keys}>
                  {' '}
                  Key {i + 1} {' ) '} {'  '} {v}{' '}
                </div>
              ))}

              <br />
            </Grid>

            {/* {!props.tsxVault ? (
              <div>
                {' '}
                <Grid item xs={12} sm={10} md={10}>
                  <div>Decryption Keys:</div>
                  <div>
                    {props.decryptKeys.map((v: any, i: any) => (
                      <div className={classes.keys}>
                        {' '}
                        Key {i + 1} {' ) '} {'  '} {v}{' '}
                      </div>
                    ))}

                    <br />
                  </div>
                </Grid>
                <Grid item xs={12} sm={10} md={10}>
                  <div>Encrypted Token (store this for backup):</div>
                  <div>
                    <div className={classes.keys}> {props.encRootToken} </div>
                  </div>
                </Grid>
              </div>
            ) : null} */}
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
