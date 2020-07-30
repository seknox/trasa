import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import React from 'react';

const useStyles = makeStyles((theme) => ({
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
    textAlign: 'center',
  },
  heading: {
    color: 'black',
    fontSize: '16px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingHeader: {
    color: 'black',
    fontSize: '11px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  dlgContent: {
    padding: theme.spacing(4),
    //  minHeight: '1000px',
    maxWidth: '400px',
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
    color: 'black',
    fontSize: '11px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingHeaderL: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    margintTop: '100px',
    padding: '10px',
  },
  settingHeaderWP: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  radio: {
    '&$checked': {
      color: '#4B8DF8',
    },
    fontSize: '9px',
    fontFamily: 'Open Sans, Rajdhani',
  },

  checked: {},
}));

export default function Newconn(props: any) {
  const classes = useStyles();

  const [data, setData] = React.useState({ appType: '', username: '', rdpProto: '', hostname: '' });

  const handleChange = () => (event: any) => {
    setData({ ...data, [event.target.name]: event.target.value });
  };

  const submitConnectionRequest = () => {
    if (data.appType === 'rdp') {
      window.open(`/my/service/connectrdp#username=${encodeURIComponent(data.username)}&serviceID=${data.hostname}&hostname=${data.hostname}`);

    } else if (data.appType === 'ssh') {
      window.open(`/my/service/connectssh#username=${encodeURIComponent(data.username)}&serviceID=${data.hostname}&hostname=${data.hostname}`);

    }
  };

  return (
    <div>
      <Dialog
        // fullWidth={true} maxWidth={'sm'}
        open={props.open}
        onClose={props.handleNewconDlgState}
        //   aria-labelledby="alert-dialog-title"
        //   aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {' '}
          <div className={classes.heading}> Enter your connection properties </div>
        </DialogTitle>

        <div className={classes.dlgContent}>
          <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
            <TextField
              fullWidth
              onChange={handleChange()}
              name="hostname"
              value={data.hostname}
              label="hostname"
            />

            <TextField
              fullWidth
              onChange={handleChange()}
              name="username"
              label="username"
              value={data.username}
            />

            <br />

            <Grid container spacing={2}>
              <Grid item xs={5}>
                <div className={classes.settingHeaderL}> Connection Type : </div>
              </Grid>
              <Grid item xs={7}>
                <div>
                  <FormControl>
                    <RadioGroup aria-label="position" name="appType" onChange={handleChange()} row>
                      <FormControlLabel
                        value="ssh"
                        control={<Radio color="primary" checked={data.appType === 'ssh'} />}
                        label="SSH"
                        labelPlacement="start"
                        classes={{ root: classes.radio }}
                      />
                      <FormControlLabel
                        value="rdp"
                        control={<Radio color="primary" checked={data.appType === 'rdp'} />}
                        label="RDP"
                        labelPlacement="start"
                      />
                    </RadioGroup>
                  </FormControl>
                </div>
              </Grid>
            </Grid>

            {data.appType === 'rdp' ? (
              <Grid container spacing={2}>
                <Grid item xs={4}>
                  <div className={classes.settingHeaderL}> RDP Type : </div>
                </Grid>
                <Grid item xs={8}>
                  <div>
                    <FormControl>
                      <RadioGroup
                        aria-label="position"
                        name="rdpProto"
                        value={data.rdpProto}
                        onChange={handleChange()}
                        row
                      >
                        <FormControlLabel
                          value="nla"
                          control={<Radio color="primary" />}
                          label="NLA"
                          labelPlacement="start"
                        />
                        <FormControlLabel
                          value="tls"
                          control={<Radio color="primary" />}
                          label="TLS"
                          labelPlacement="start"
                        />
                        <FormControlLabel
                          value="rdp"
                          control={<Radio color="primary" />}
                          label="RDP"
                          labelPlacement="start"
                        />
                      </RadioGroup>
                    </FormControl>
                  </div>
                </Grid>
              </Grid>
            ) : null}
            <br />
          </Grid>
        </div>
        <DialogActions>
          <Button onClick={props.close} color="primary" autoFocus>
            Close
          </Button>

          <Button variant="contained" color="primary" onClick={submitConnectionRequest}>
            Connect
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
