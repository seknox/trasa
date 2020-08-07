import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Drawer from '@material-ui/core/Drawer';
import { makeStyles } from '@material-ui/core/styles';
import Save from '@material-ui/icons/PersonAdd';
import React, { useState } from 'react';
import UserCrud from './crud';

const useStyles = makeStyles((theme) => ({
  content: {
    width: '100%',
    flexGrow: 1,
    // backgroundColor: theme.palette.background.default,
    backgroundColor: 'white',
    padding: 24,
    height: 'calc(100% - 56px)',
    // height: '100%',
    marginTop: 56,
    [theme.breakpoints.up('sm')]: {
      height: 'calc(100% - 64px)',
      marginTop: 64,
    },
  },
  drawer: {
    background: 'white',
    width: 456,
  },
  // handles apppbar flex
  verticalBar: {
    display: 'flex',
    // height: '100%',
    flexDirection: 'column',
  },
  Up: {
    alignItems: 'flext-start',
    // justifyContent: 'center'
  },
  Down: {
    alignItems: 'flext-end',
    //  justifyContent: 'center'
  },

  drawerContent: {
    marginLeft: 30,
  },
}));

type UserCreateUpdateDrawerProps = {
  updateUserTable: (n: any) => void;
};

type Anchor = 'top' | 'left' | 'bottom' | 'right';

export default function UserCreateUpdateDrawer(props: UserCreateUpdateDrawerProps) {
  const classes = useStyles();
  const user = {
    ID: '',
    firstName: '',
    middleName: '',
    lastName: '',
    email: '',
    password: '',
    userRole: '',
    userName: '',
    userName2: '',
    cpassword: '',
    status: false,
    CreatedAt: 0,
  };
  // const [open, setopen] = useState(false);
  const [right, setright] = useState(false);
  const [Link, setlink] = useState('');
  const [dlgOpen, setdlgOpen] = useState(false);

  const openDrawer = (anchor: Anchor, open: boolean) => (
    event: React.KeyboardEvent | React.MouseEvent,
  ) => {
    if (
      event.type === 'keydown' &&
      ((event as React.KeyboardEvent).key === 'Tab' ||
        (event as React.KeyboardEvent).key === 'Shift')
    ) {
      return;
    }

    setright(open);
  };

  const closeDrawer = () => {
    setright(false);
  };

  const handleDlgClose = () => {
    setdlgOpen(false);
  };

  const handleVerifyLinkChange = (val: any) => {
    setlink(val);
    setdlgOpen(true);
  };

  return (
    <div>
      <Button name={"createUserBtn"} variant="contained" size="small" onClick={openDrawer('right', true)}>
        <Save />
        Create User
      </Button>
      <br />
      <br />
      <Drawer anchor="right" open={right} onClose={openDrawer('right', false)}>
        <div className={classes.drawer}>
          <div className={classes.drawerContent}>
            <UserCrud
              userData={user}
              update={false}
              updateUserTable={props.updateUserTable}
              handleDrawerClose={closeDrawer}
              handleVerifyLinkChange={handleVerifyLinkChange}
            />
          </div>
        </div>
      </Drawer>

      <Dialog
        open={dlgOpen}
        fullWidth
        onClose={handleDlgClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">Verification Link</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            Send this link to user for signup. <br />
            Link: {Link}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button variant="contained" onClick={handleDlgClose} color="primary" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
