import DateFnsUtils from '@date-io/date-fns';
import { makeStyles } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import 'date-fns';
import { MuiPickersUtilsProvider, TimePicker } from 'material-ui-pickers';
import React, { useState } from 'react';
import ProgressHOC from '../../../../utils/Components/Progressbar';
import Constants from '../../../../Constants';

const useStyles = makeStyles((theme) => ({
  paper: {
    maxWidth: 1500,
    margin: 'auto',
    marginTop: 50,
    overflow: 'hidden',
    padding: theme.spacing(2),
  },
}));

type grantOrDenyAccessprops = {
  reqID: string;
  reqText: string;
  //serviceID: string;
  handleRequestDialogueClose: () => void;
  deleteElement: (key: string) => void;
  adhocReqDlgState: boolean;
};

export default function GrantOrDenyAccess(props: grantOrDenyAccessprops) {
  const [selectedTime, setSelectedTime] = useState(new Date());
  const [loader, setLoader] = useState(false);

  const handleSubmit = (name: string) => (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    setLoader(true);
    const req = {
      reqID: props.reqID,
      authorizedPeriod: selectedTime.getTime(),
      //  serviceID: props.serviceID,
      isAuthorized: false,
    };

    if (name !== 'reject') {
      req.isAuthorized = true;
    }

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/policy/adhoc/respond`, req)
      .then((response) => {
        setLoader(false);

        props.handleRequestDialogueClose();
        if (response.data.status === 'success') {
          props.deleteElement(props.reqID);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const handleDateChange = (time: any) => {
    setSelectedTime(time);
  };

  const classes = useStyles();
  return (
    <div>
      <Dialog
        onClose={props.handleRequestDialogueClose}
        aria-labelledby="customized-dialog-title"
        open={props.adhocReqDlgState}
        fullWidth
        maxWidth="sm"
      >
        <DialogTitle id="customized-dialog-title">
          <Typography variant="h2"> Grant or Revoke Access to this app. </Typography>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={12} md={12}>
              <Typography variant="h3"> Requested reason for access: </Typography>{' '}
              <Typography variant="h4"> {props.reqText} </Typography>
            </Grid>
            <Grid item xs={5} sm={5} md={5}>
              <Typography variant="h3"> Grant access to selected time: </Typography>
            </Grid>
            <Grid item xs={7} sm={7} md={7}>
              <MuiPickersUtilsProvider utils={DateFnsUtils}>
                <Grid container justify="space-around">
                  {/* <ThemeProvider theme={MgetMuiTheme()}> */}
                  <TimePicker
                    margin="normal"
                    label="Time picker"
                    value={selectedTime}
                    // color="primary"
                    // variant="contained"
                    onChange={handleDateChange}
                  />
                  {/* </ThemeProvider> */}
                </Grid>
              </MuiPickersUtilsProvider>
            </Grid>
          </Grid>
        </DialogContent>
        {loader ? <ProgressHOC /> : ''}
        <DialogActions>
          <Button onClick={handleSubmit('reject')} variant="contained" color="secondary">
            Reject
          </Button>
          <Button onClick={handleSubmit('grant')} variant="contained" color="secondary">
            Grant
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
