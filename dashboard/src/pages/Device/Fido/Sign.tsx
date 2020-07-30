import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React from 'react';
import U2fReIcon from '../../../assets/u2fregs.svg';
import Constants from '../../../Constants';
import { StartTimeout, StopTimeout } from './index';

const useStyles = makeStyles((theme) => ({
  root: {
    //  padding: theme.spacing(2)
  },

  content: {
    padding: theme.spacing(2),
  },

  paper: {
    textAlign: 'center',
    padding: theme.spacing(2),
  },
}));

export default function MyDevices() {
  const classes = useStyles();
  // const [serviceID, setserviceID] = React.useState('');
  const [progressStatus, setprogressStatus] = React.useState('Getting data from server...');
  const [error, setError] = React.useState(false);
  const [timeout, settimeout] = React.useState(0);
  // const [registerRequest, setregisterRequest] = React.useState();
  // const [registration, setregistration] = React.useState({ response: { errorCode: '' } });

  const RegisterRequestResponse = (regData: any) => {
    setError(false);
    setprogressStatus('Registering your device...');

    const reqData = regData;
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/my/regresp`, reqData)
      .then((r) => {
        if (r.data.status === 'success') {
          setError(true);
          setprogressStatus(
            'Your device is registered. You can check by logging out and logging in again by selecting yubico as 2nd factor method.',
          );
        } else {
          setError(true);
          setprogressStatus(
            'Oops! We encountered some error. Try again after few seconds. If the problem persists, contact your administrator.',
          );
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const U2fRegister = (serviceID: string, registerRequest: any) => {
    StartTimeout(timeout, Date.now(), 30000, settimeout);
    setprogressStatus('Press the golden disk');
    (window as any).u2f.register(serviceID, [registerRequest], [], (deviceResponse: any) => {
      StopTimeout(settimeout);

      if (deviceResponse.errorCode) {
        setError(true);
        setprogressStatus(
          'Oops! We encountered some error. Try again after few seconds. If the problem persists, contact your administrator.',
        );
        // registration.response = { errorCode: deviceResponse.errorCode };
      } else {
        RegisterRequestResponse(deviceResponse);
      }
      // this.setState({registration})
    });
  };

  const RegisterRequest = () => {
    setError(false);
    setprogressStatus('Getting data from server...');

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/my/regreq`)
      .then((response) => {
        // setregisterRequest(response.data.data[0].registerRequests[0]);
        // setserviceID(response.data.data[0].serviceID);
        U2fRegister(response.data.data[0].serviceID, response.data.data[0].registerRequests[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  React.useEffect(() => {
    RegisterRequest();
  }, []);

  return (
    <div className={classes.root}>
      <div className={classes.content}>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Grid container spacing={2} direction="column" alignItems="center" justify="center">
              <Grid item xs={12}>
                <Typography component="h3" variant="h1">
                  <b>Authorize with your security key.</b>
                </Typography>
              </Grid>

              <Grid item xs={12}>
                <Typography component="h3" variant="h3">
                  Insert your device into a USB port. If the key has blinking light, press the gold
                  disk.
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
              <Typography component="h3" variant="h1" style={{ color: error ? 'red' : 'navy' }}>
                {progressStatus}
              </Typography>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <LinearBuffer />
          </Grid>
        </Grid>
      </div>
    </div>
  );
}

function LinearBuffer() {
  const classes = useStyles();
  const [completed, setCompleted] = React.useState(0);
  const [buffer, setBuffer] = React.useState(10);

  const progress = React.useRef(() => {});
  React.useEffect(() => {
    progress.current = () => {
      if (completed > 100) {
        setCompleted(0);
        setBuffer(10);
      } else {
        const diff = Math.random() * 10;
        const diff2 = Math.random() * 10;
        setCompleted(completed + diff);
        setBuffer(completed + diff + diff2);
      }
    };
  }, []);

  React.useEffect(() => {
    function tick() {
      progress.current();
    }
    const timer = setInterval(tick, 500);

    return () => {
      clearInterval(timer);
    };
  }, []);

  return (
    <div className={classes.root}>
      <LinearProgress variant="buffer" value={completed} valueBuffer={buffer} />
      {/* <LinearProgress variant="buffer" value={completed} valueBuffer={buffer} color="secondary" /> */}
    </div>
  );
}
