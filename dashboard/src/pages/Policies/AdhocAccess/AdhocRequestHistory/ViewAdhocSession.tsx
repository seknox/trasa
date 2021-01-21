import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import 'date-fns';
import Moment from 'moment';
import React from 'react';

type viewAdhocSessionProps = {
  reqText: string;
  sessionID: string;
  serviceType: string;
  orgID: string;
  handleRequestDialogueClose: () => void;
  viewAdhocSessionDlgState: boolean;
  localeAuthorizedOn: string;
  localeAuthorizedPeriod: string;
  timestamp: string;
};

export default function ViewAdhocSession(props: viewAdhocSessionProps) {
  return (
    <div>
      <Dialog
        onClose={props.handleRequestDialogueClose}
        aria-labelledby="customized-dialog-title"
        open={props.viewAdhocSessionDlgState}
        fullWidth
        maxWidth="sm"
      >
        <DialogTitle id="customized-dialog-title">Adhoc Session</DialogTitle>
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={12} md={12}>
              <Typography variant="h3">Requested reason for access:</Typography>{' '}
              <Typography variant="h4">{props.reqText}</Typography>
            </Grid>
            <Grid item xs={12} sm={12} md={12}>
              <Typography variant="h3">Sanctioned on: </Typography>{' '}
              <Typography variant="h4">{props.localeAuthorizedOn}</Typography>
            </Grid>
            <Grid item xs={12} sm={12} md={12}>
              <Typography variant="h3">Authorized Till on:</Typography>{' '}
              <Typography variant="h4">{props.localeAuthorizedPeriod}</Typography>
            </Grid>
            <Grid item xs={12} sm={12} md={12}>
              <Typography variant="h3">User Sessions: </Typography>
              <ViewSession
                sessionID={props.sessionID}
                serviceType={props.serviceType}
                orgID={props.orgID}
                timestamp={props.timestamp}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={props.handleRequestDialogueClose} variant="contained" color="secondary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}

const ViewSession = (props: any) => {
  const sessionArr = props.sessionID.split(',');

  const d = Moment.unix(props.timestamp);

  return (
    <div>
      <Grid container spacing={2}>
        {sessionArr.map((session: any, i: number) =>
          props.serviceType === 'ssh' ||
          props.serviceType === 'ssh' ||
          props.serviceType === 'rdp' ||
          props.serviceType === 'vnc' ||
          props.serviceType === 'http' ? (
            <Grid item xs={3}>
              {` ( ${i + 1} ) `}
              <Button
                style={{ color: 'white', background: 'navy' }}
                onClick={() => {
                  // window.location.href="/overview/session/"+sessionID
                  if (props.serviceType === 'ssh') {
                    window.open(
                      `/monitor/sessions/view#type=ssh&year=${d.year()}&month=${
                        d.months() + 1
                      }&day=${d.date()}&sessionID=${session}`,
                      '_blank',
                    );
                  } else if (props.serviceType === 'guac-ssh') {
                    window.open(
                      `/monitor/sessions/view#type=guac-ssh&year=${d.year()}&month=${
                        d.months() + 1
                      }&day=${d.date()}&sessionID=${session}`,
                      '_blank',
                    );
                  } else if (props.serviceType === 'guac-rdp') {
                    window.open(
                      `/monitor/sessions/view#type=guac-rdp&year=${d.year()}&month=${
                        d.months() + 1
                      }&day=${d.date()}&sessionID=${session}`,
                      '_blank',
                    );
                  } else if (props.serviceType === 'rdp') {
                    window.open(
                      `/monitor/sessions/view#type=rdp&year=${d.year()}&month=${
                        d.months() + 1
                      }&day=${d.date()}&sessionID=${session}`,
                      '_blank',
                    );
                  } else if (props.serviceType === 'http') {
                    window.open(
                      `/monitor/sessions/view#type=http&year=${d.year()}&month=${
                        d.months() + 1
                      }&day=${d.date()}&sessionID=${session}`,
                      '_blank',
                    );
                  }
                }}
              >
                View Session
              </Button>{' '}
            </Grid>
          ) : (
            <Grid item xs={3}>
              {` ( ${i + 1} ) `}
              <Button style={{ color: 'black' }}> Not Recorded </Button>{' '}
            </Grid>
          ),
        )}
      </Grid>
    </div>
  );
};
