import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  warningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
  dialogTitle: {
    fontSize: '20px',
  },
  dialogContent: {
    fontSize: '16px',
  },
}));

type DlgProps = {
  open: boolean;
  close: () => void;
  confirmMessage: string;
  deleteFunc: () => void;
};

export const DeleteConfirmDialogue = (props: DlgProps) => {
  const { open, close, confirmMessage, deleteFunc } = props;
  const classes = useStyles();
  return (
    <div>
      <Dialog
        open={open}
        onClose={close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle className={classes.dialogTitle} disableTypography>
          Confirm Delete?
        </DialogTitle>
        <DialogContent>
          <DialogContentText className={classes.dialogContent}>{confirmMessage}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              close();
              deleteFunc();
            }}
            className={classes.warningButton}
          >
            Yes, Delete
          </Button>
          <Button onClick={close} color="primary" variant="contained" autoFocus>
            No
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

type ConfirmDlgProps = {
  open: boolean;
  close: () => void;
  confirmMessage: string;
  confirmFunc: () => void;
};

export const ConfirmDialogue = (props: ConfirmDlgProps) => {
  const classes = useStyles();
  const { open, close, confirmMessage, confirmFunc } = props;
  return (
    <div>
      <Dialog
        open={open}
        onClose={close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle className={classes.dialogTitle} disableTypography>
          Confirm Action!
        </DialogTitle>
        <DialogContent>
          <DialogContentText className={classes.dialogContent}>{confirmMessage}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            className={classes.warningButton}
            onClick={() => {
              close();
              confirmFunc();
            }}
          >
            Yes, I Agree
          </Button>
          <Button onClick={close} color="primary" variant="contained" autoFocus>
            No
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

type NotifDialogueProps = {
  open: boolean;
  close: () => void;
  confirmMessage: string;
};

export const NotifDialogue = (props: NotifDialogueProps) => {
  const classes = useStyles();
  const { open, close, confirmMessage } = props;
  return (
    <div>
      <Dialog
        open={open}
        onClose={close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle className={classes.dialogTitle} disableTypography>
          Response
        </DialogTitle>
        <DialogContent>
          <DialogContentText className={classes.dialogContent}>{confirmMessage}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={close} color="primary" variant="contained" autoFocus>
            Okay
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
