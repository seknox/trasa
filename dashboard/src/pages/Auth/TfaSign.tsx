import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React from 'react';
import U2fReIcon from '../../assets/u2fregs.svg';
import Constants from '../../Constants';
import { SetPasswordComponent } from './PasswordSetup';
import { StartTimeout, StopTimeout } from '../Device/Fido';

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
    '& > * + *': {
      marginTop: theme.spacing(2),
    },
  },
  content: {
    padding: theme.spacing(2),
  },

  paper: {
    textAlign: 'center',
    padding: theme.spacing(2),
  },
}));

export default function TfaSign(props: any) {
  const classes = useStyles();
  // const [serviceID, setserviceID] = React.useState('');
  const [progressStatus, setprogressStatus] = React.useState('Getting data from server...');
  const [error, setError] = React.useState(false);
  const [timeout, settimeout] = React.useState(0);
  // const [registerRequest, setregisterRequest] = React.useState();
  // const [challenge, setchallenge] = React.useState('');
  // const [changePassToken, setchangePassToken] = React.useState('');
  // const [changePassDlgState, setchangePassDlgState] = React.useState(false);
  // const [signDeviceDialogueState, setsignDeviceDialogueState] = React.useState(false);
  // const [progress, setprogress] = React.useState(false);

  const U2fSign = (serviceID: string, signChall: string, registeredKey: any) => {
    StartTimeout(timeout, Date.now(), 30000, settimeout);
    setError(false);
    setprogressStatus('Press the golden disk');

    (window as any).u2f.sign(serviceID, signChall, [registeredKey], (deviceResponse: any) => {
      // console.log('Response', deviceResponse)
      StopTimeout(settimeout);
      if (deviceResponse.errorCode) {
        setError(true);
        setprogressStatus(
          'Oops! We encountered some error. Try again after few seconds. If the problem persists, contact your administrator.',
        );

        console.log('ErrorCode:', deviceResponse.errorCode);
        // registration['response'] = {errorCode: deviceResponse.errorCode}
      }
    });
  };

  const SignRequest = () => {
    setError(false);
    setprogressStatus('Processing tfa request...');

    const req = { tfaMethod: 'u2fy', token: props.token, intent: props.intent };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/idp/login/tfa`, req)
      .then((response) => {
        // this.setState({registeredKeys: response.data.data[0].registeredKeys[0], serviceID:response.data.data[0].serviceID,  })
        U2fSign(
          response.data.data[0].serviceID,
          response.data.data[0].challenge,
          response.data.data[0].registeredKeys[0],
        );
      })
      .catch((err) => {
        console.log('Error', err);
        window.location.href = '/login';
      });
  };

  React.useEffect(() => {
    SignRequest();
  }, []);

  return (
    <div className={classes.root}>
      <div className={classes.content}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Grid container spacing={2} direction="column" alignItems="center" justify="center">
              <Grid item xs={12}>
                <Typography component="h3" variant="h3">
                  <b>Authorize with your security key.</b>
                </Typography>
              </Grid>

              <Grid item xs={12}>
                <Typography component="h3" variant="h5">
                  Insert your device into a USB port and press the gold disk.
                </Typography>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12}>
            <img src={U2fReIcon} alt="u2f Registration icon" />
            {/* <Divider variant="middle" />
              <br /><br /><br /> */}
          </Grid>

          <Grid container spacing={2} direction="column" alignItems="center" justify="center">
            <Grid item xs={12}>
              <Typography component="h3" variant="h6" style={{ color: error ? 'red' : 'navy' }}>
                {progressStatus}
              </Typography>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            {/* {progress ? <LinearBuffer /> : null} */}
          </Grid>
        </Grid>
        <ChangePassDialogue token="" open={false} close={() => 0} />
      </div>
    </div>
  );
}

const ChangePassDialogue = (props: any) => {
  return (
    <div>
      <Dialog
        open={props.open}
        onClose={props.close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        // minWidth="xl"
        fullScreen
        scroll="paper"
      >
        <DialogTitle id="alert-dialog-title">Password change required by policy</DialogTitle>
        <DialogContent>
          <SetPasswordComponent token={props.token} update={false} />
        </DialogContent>
      </Dialog>
    </div>
  );
};
