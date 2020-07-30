import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { makeStyles, Theme } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import React, { useState } from 'react';
import Totp from '../../assets/totp.svg';
import TrasaLogo from '../../assets/trasa-ni.svg';
import U2f from '../../assets/u2f.svg';
import U2fY from '../../assets/yubi2.svg';
import LinearProgress from '../../utils/Components/Progressbar';
import SignU2f from './SignYoubikeyU2f';

// STYLES
const useStyles = makeStyles((theme: Theme) => ({
  tfaButton: {
    boxShadow: '0 0 10px 0 rgba(0,0,0,0.1)',
    border: '1px solid white',
    color: '#051384',
    fontSize: '34px',
    padding: '0 30px',
    '&:hover': {
      background: 'white',
      transform: 'translateY(-3px)',
      boxShadow: '0 4px 20px 0 #0b1728ff', // rgba(0,0,0,0.12)',
    },
  },

  padMiddle: {
    textAlign: 'center',
    // marginLeft: '50%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    minWidth: 750,
    minHeight: 200,
    paddingRight: 50,
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  paperSmall: {
    backgroundColor: 'transparent',
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  typo: {
    color: '#1A237E',
    // fontSize: '20px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  heading: {
    color: '#0b1728ff',
    fontSize: '20px',
    fontFamily: 'Open Sans, Rajdhani',
  },

  cprightText: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    // marginLeft: '10%'
  },
}));

type TfaProp = {
  sendTfa: (event: React.FormEvent<Element>, tfaMethod: string, totpCode: string) => void;
  loader: boolean;
};

// TODO pending Yubikey signer check.
export default function TfaPage(props: TfaProp) {
  const [totpDlgState, settotpDlgState] = useState(false);
  const [totpVal, setTotpVal] = useState('');

  function getTotpval(e: React.ChangeEvent<HTMLInputElement>) {
    setTotpVal(e.target.value);
  }

  const classes = useStyles();
  return (
    <div>
      <Paper className={classes.paper}>
        <img src={TrasaLogo} height={100} width={200} alt="LOGO" />
        <div className={classes.heading}> Choose second step verification method </div>

        <br />
        <Divider light />
        <br />
        <Grid container spacing={2} direction="row" alignItems="center" justify="center">
          <Grid item xs={12} sm={6} md={6} lg={4}>
            <Paper className={classes.paperSmall} elevation={0}>
              <Typography component="h3" variant="h4" className={classes.typo}>
                TOTP
              </Typography>
              <br />
              <Button
                variant="outlined"
                className={classes.tfaButton}
                onClick={() => settotpDlgState(true)}
                id="totpButton"
              >
                <img src={Totp} height={100} width={150} alt="totp icon" />
              </Button>
            </Paper>
          </Grid>

          <Grid item xs={12} sm={6} md={6} lg={4}>
            <Paper className={classes.paperSmall} elevation={0}>
              <Typography component="h3" variant="h4" className={classes.typo}>
                TRASA U2F
              </Typography>
              <br />
              <Button
                variant="outlined"
                className={classes.tfaButton}
                name="u2fButton"
                id="u2fButton"
                onClick={(e) => {
                  props.sendTfa(e, 'U2F', '');
                }}
              >
                <img src={U2f} height={100} width={150} alt="u2f icon" />
              </Button>
            </Paper>
          </Grid>

          <Grid item xs={12} sm={6} md={6} lg={4}>
            <Paper className={classes.paperSmall} elevation={0}>
              <Typography component="h3" variant="h4" className={classes.typo}>
                Yubikey U2F
              </Typography>
              <br />
              <Button
                variant="outlined"
                className={classes.tfaButton}
                onClick={(e) => {
                  props.sendTfa(e, 'U2FY', '');
                }}
                id="u2fyButton"
              >
                <img src={U2fY} alt="yubi icon" />
              </Button>
            </Paper>
          </Grid>

          <Grid item xs={12} />
        </Grid>
        {props.loader && <LinearProgress />}
        <br />
        <div className={classes.padMiddle}>
          <div className={classes.cprightText}>
            {' '}
            Trasa Dashboard v2.5. Seknox Cybersecurity (P) Ltd. - 2019{' '}
          </div>{' '}
        </div>
        <TotpDlg
          loader={props.loader}
          totpVal={totpVal}
          getTotpVal={getTotpval}
          sendTfa={props.sendTfa}
          open={totpDlgState}
          handleClose={() => settotpDlgState(false)}
        />
      </Paper>
    </div>
  );
}

type SignDeviceDialogueProps = {
  token: string;
  open: boolean;
  intent: string;

  close: () => void;
};

export const SignDeviceDialogue = (props: SignDeviceDialogueProps) => {
  return (
    <div>
      <Dialog
        open={props.open}
        onClose={props.close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        maxWidth="sm"
        scroll="paper"
      >
        <DialogContent>
          <SignU2f token={props.token} intent={props.intent} />
        </DialogContent>
        <DialogActions>
          <Button onClick={props.close} color="primary" variant="contained" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

type TotpProp = {
  open: boolean;
  handleClose: () => void;
  sendTfa: (event: React.FormEvent<Element>, tfaMethod: string, totpCode: string) => void;
  loader: boolean;
  getTotpVal: (e: React.ChangeEvent<HTMLInputElement>) => void;
  totpVal: string;
};

function TotpDlg(props: TotpProp) {
  return (
    <Dialog open={props.open} onClose={props.handleClose} aria-labelledby="form-dialog-title">
      <DialogTitle id="form-dialog-title">Enter 6 digit code</DialogTitle>
      <DialogContent>
        <DialogContentText>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              props.sendTfa(e, 'totp', props.totpVal);
            }}
          >
            <TextField
              fullWidth
              autoFocus
              onChange={props.getTotpVal}
              name="totpVal"
              variant="outlined"
              size="small"
            />
          </form>
        </DialogContentText>
        {props.loader && <LinearProgress />}
      </DialogContent>
      <DialogActions>
        <Button
          size="small"
          variant="contained"
          color="primary"
          onClick={(e) => {
            props.sendTfa(e, 'totp', props.totpVal);
          }}
          id="submitTotpButton"
        >
          Submit
        </Button>
        <Button size="small" onClick={props.handleClose} color="primary">
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  );
}
