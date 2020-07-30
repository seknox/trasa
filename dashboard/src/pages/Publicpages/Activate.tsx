// import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
// import Input, { InputLabel, InputAdornment } from '@material-ui/core/Input';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import CopyIcon from '@material-ui/icons/FileCopy';
import axios from 'axios';
import React from 'react';
import CopyToClipboard from 'react-copy-to-clipboard';
import TrasaLogo from '../../assets/trasa-white.svg';
import ProgressHOC from '../../utils/Components/Progressbar';
import Constants from '../../Constants';

// STYLES
const useStyles = makeStyles((theme) => ({
  root: {
    // position: 'absolute',
    backgroundColor: 'rgba(1,1,35,1)', // '#0b1728ff',
    width: '100%',
    height: '100%',
    padding: theme.spacing(4),
    margin: '0',
  },
  transformMiddle: {
    // position: 'absolute',
    marginTop: '10%',
    // right: '50%',
    //  transform: 'translate(50%,-50%)',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(4),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperActivate: {
    backgroundColor: '#fdfdfd',
    marginTop: '5%',
    padding: theme.spacing(4),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  transPaper: {
    backgroundColor: 'transparent',
    padding: theme.spacing(4),
    // textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  ActivateHeader: {
    color: 'white',
    fontSize: '20px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  divMin: {
    //  minWidth: 700,
  },
  heading: {
    color: 'white',
    fontSize: '20px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(2),
  },
  subHeading: {
    color: 'white',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(1),
  },
  icon: {
    verticalAlign: 'bottom',
    height: 20,
    width: 20,
  },
  details: {
    alignItems: 'center',
  },
  column: {
    flexBasis: '33.33%',
  },
  activateHeading: {
    color: 'black',
    fontSize: '25px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(2),
  },
  menuKey: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    // padding: theme.spacing(2),
  },
  submitButton: {
    marginLeft: '50%',
    color: 'white',
    background: 'navy',
    borderRadius: 3,
    '&:hover': {
      marginLeft: '50%',
      color: 'white',
      background: 'navy',
      borderRadius: 3,
    },
  },
  copyButton: {
    //  marginLeft: '50%',
    color: 'white',
    background: 'navy',
    borderRadius: 3,
    '&:hover': {
      color: 'white',
      background: 'navy',
      borderRadius: 3,
    },
  },
  errorText: {
    color: 'white',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    background: 'maroon',
  },
}));

export default function Activate() {
  const classes = useStyles();
  const [machineID, setmachineID] = React.useState('');
  // const [msg, setmsg] = React.useState('');
  const [signedLicense, setsignedLicense] = React.useState('');
  const [clipboard, setClipboard] = React.useState('');
  const [loader, setLoader] = React.useState(false);
  // const [activationFailed, setactivationFailed] = React.useState(false);

  const handleChange = (e: any) => {
    setsignedLicense(e.target.value);
  };

  const submitLicense = (event: any) => {
    setLoader(true);
    event.preventDefault();

    const req = { signedLicense };

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/license/activate`, req)
      .then((response) => {
        setLoader(false);
        if (response.data.status === 'success') {
          window.location.href = '/login';
        } else {
          setLoader(false);
        }
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  React.useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/lisc`)
      .then((response) => {
        if (response.data.intent === 'NotifLicense') {
          setmachineID(response.data.data[0]);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <div>
      <div className={classes.root}>
        <div className={classes.transformMiddle}>
          <Grid container>
            <Grid item xs={12} sm={6}>
              <div className={classes.divMin}>
                <Paper className={classes.transPaper}>
                  <img src={TrasaLogo} height={100} width={300} alt="trasa-logo" />
                  <div className={classes.heading}>
                    {' '}
                    Experience the most <b> powerfull access control platform </b> on the planet{' '}
                  </div>
                  <div className={classes.subHeading}> 1. Copy your license key </div>
                  <div className={classes.subHeading}> 2. Uplaod the key to license signer </div>
                  <div className={classes.subHeading}> 3. Download signed license file </div>
                  <div className={classes.subHeading}>
                    {' '}
                    4. Upload the signed license file here to activate.
                  </div>
                </Paper>
              </div>
            </Grid>
            <Grid item xs={12} sm={4}>
              <Paper className={classes.paperActivate}>
                <div className={classes.column}>
                  <div className={classes.activateHeading}>Activate</div>
                  {loader && <ProgressHOC />}
                  <Divider light />
                  <br />
                  <br />
                  <Grid container spacing={2}>
                    <Grid item xs={3} sm={3} md={3}>
                      <div className={classes.menuKey}> Step 1: </div>
                    </Grid>
                    <Grid item xs={9} sm={9} md={9}>
                      <CopyToClipboard text={clipboard}>
                        <Button
                          className={classes.copyButton}
                          size="small"
                          onClick={() => {
                            setClipboard(machineID);
                          }}
                          aria-label="copy"
                        >
                          Copy License Key : <CopyIcon fontSize="small" />
                        </Button>
                      </CopyToClipboard>
                    </Grid>

                    <Grid item xs={3} sm={3} md={3}>
                      <div className={classes.menuKey}> Step 2: </div>
                    </Grid>
                    <Grid item xs={9} sm={9} md={9}>
                      <div className={classes.menuKey}> Paste license file below and submit: </div>
                      <TextField
                        id="outlined-multiline-static"
                        // label="Multiline"
                        multiline
                        autoFocus
                        fullWidth
                        rows="4"
                        defaultValue=""
                        margin="normal"
                        variant="outlined"
                        onChange={handleChange}
                      />
                    </Grid>
                    {/* <div className={classes.submitButton}> */}
                    <Button
                      onClick={submitLicense}
                      aria-label="copy"
                      className={classes.submitButton}
                    >
                      Submit
                    </Button>
                    {/* </div> */}
                  </Grid>
                </div>
              </Paper>
            </Grid>
          </Grid>
        </div>
      </div>
    </div>
  );
}
