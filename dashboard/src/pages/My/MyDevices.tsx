import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import ListItemText from '@material-ui/core/ListItemText';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import React, { useState } from 'react';
import DialogueWrapper from '../../utils/Components/DialogueWrapComponent';
import LoginBox from '../Auth';
import EnrolDeviceComponent from '../Device/EnrolMobileDevice';
import FidoU2f from '../Device/Fido';
import SignU2f from '../Device/Fido/Sign';
import UserDevices from '../Users/UserDevices';

export default function MyDevices(props: any) {
  const [enrolMobileDeviceDlg, setEnrolMobileDeviceDlg] = useState(false);
  const [enrolYubikeyDeviceDlg, setEnrolYubikeyDeviceDlg] = useState(false);

  const changeEnrolMobileDeviceDlgDlg = () => {
    setEnrolMobileDeviceDlg(!enrolMobileDeviceDlg);
  };

  const closeEnrolMobileDeviceDlgDlg = () => {
    setEnrolMobileDeviceDlg(false);
  };

  function changeEnrolYubikeyDeviceDlg() {
    setEnrolYubikeyDeviceDlg(!enrolYubikeyDeviceDlg);
  }

  function closeEnrolYubikeyDeviceDlg() {
    setEnrolYubikeyDeviceDlg(false);
  }

  const [hasAuthenticated, setHasAuthentecated] = useState(false);

  function changeHasAuthenticated() {
    setHasAuthentecated(!hasAuthenticated);
  }

  const [enrolDeviceDetail, setEnrolDeviceDetail] = useState({});

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        {/* <Button variant="contained" color="secondary" onClick={changeViewDlgState}>
                        Enrol 2fa Device
                      </Button> */}
        <EnrolDeviceManu
          changeEnrolMobileDeviceDlgDlg={changeEnrolMobileDeviceDlgDlg}
          changeEnrolYubikeyDeviceDlg={changeEnrolYubikeyDeviceDlg}
        />
        <br />
        <br />
        <Divider light />
      </Grid>
      <Grid item xs={12}>
        <UserDevices renderFor="myRoute" />
      </Grid>

      <div>
        <br />
        <br />
      </div>

      {/* <Grid item xs={12}>
                    <FidoU2f />
                    </Grid> */}

      <DialogueWrapper
        open={enrolMobileDeviceDlg}
        handleClose={closeEnrolMobileDeviceDlgDlg}
        title="Enrol 2FA device"
        maxWidth="md"
        fullScreen
      >
        <EnrolMobileDeviceDlg
          hasAuthenticated={hasAuthenticated}
          userData={props.userData}
          changeHasAuthenticated={changeHasAuthenticated}
          setEnrolDeviceDetail={setEnrolDeviceDetail}
          enrolDeviceDetail={enrolDeviceDetail}
        />
      </DialogueWrapper>

      <DialogueWrapper
        fullScreen
        open={enrolYubikeyDeviceDlg}
        handleClose={closeEnrolYubikeyDeviceDlg}
        title="Enrol 2FA device"
        maxWidth="md"
      >
        <EnrolYubikey hasAuthenticated={hasAuthenticated} userData={props.userData} />
      </DialogueWrapper>
    </Grid>
  );
}

function EnrolMobileDeviceDlg(props: any) {
  return (
    <Grid container spacing={2} direction="column">
      <Grid item xs={12}>
        {props.hasAuthenticated === true ? (
          <EnrolDeviceComponent
            enrolDeviceDetail={props.enrolDeviceDetail}
            testEnrol={() => 'test'}
          />
        ) : (
          <LoginBox
            autofillEmail
            intent="AUTH_REQ_ENROL_DEVICE"
            showForgetPass={false}
            title="Authenticate again to proceed"
            userData={props.userData}
            proxyDomain=""
            changeHasAuthenticated={props.changeHasAuthenticated}
            setData={props.setEnrolDeviceDetail}
          />
        )}
      </Grid>
    </Grid>
  );
}

function EnrolYubikey(props: any) {
  const [registerDeviceDialogueState, setRegisterDeviceDialogueState] = useState(false);
  const [signDeviceDialogueState, setSignDeviceDialogueState] = useState(false);

  const changeRegisterDeviceDialogueState = () => {
    setRegisterDeviceDialogueState(!registerDeviceDialogueState);
  };

  const changeSignDeviceDialogueState = () => {
    setSignDeviceDialogueState(!signDeviceDialogueState);
  };

  return (
    <Grid container spacing={2} direction="column">
      <Grid item xs={12}>
        {props.hasAuthenticated === false ? (
          <LoginBox
            autofillEmail
            showForgetPass={false}
            intent="AUTH_REQ_ENROL_DEVICE"
            title="Authenticate again to proceed"
            userData={props.userData}
            proxyDomain=""
            changeHasAuthenticated={props.changeHasAuthenticated}
            setData={props.setEnrolDeviceDetail}
          />
        ) : (
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <Typography component="h1" variant="h3">
                Register your yubikey
              </Typography>
            </Grid>
            <Grid item xs={6}>
              <Typography component="span" variant="h4">
                {' '}
                Register your device:
              </Typography>
              <Button
                variant="contained"
                color="secondary"
                onClick={changeRegisterDeviceDialogueState}
              >
                Register device
              </Button>
            </Grid>
            <Grid item xs={6}>
              <Typography component="span" variant="h4">
                {' '}
                Test your device:{' '}
              </Typography>
              <Button variant="contained" color="secondary" onClick={changeSignDeviceDialogueState}>
                Test
              </Button>
            </Grid>
            <SignDeviceDialogue
              open={signDeviceDialogueState}
              close={changeSignDeviceDialogueState}
            />
            <RegisterDeviceDialogue
              open={registerDeviceDialogueState}
              close={changeRegisterDeviceDialogueState}
            />
          </Grid>
        )}
      </Grid>
    </Grid>
  );
}

const RegisterDeviceDialogue = (props: any) => {
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
          <FidoU2f />
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

const SignDeviceDialogue = (props: any) => {
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
        {/* <DialogTitle id="alert-dialog-title">{"Confirm Delete?"}</DialogTitle> */}
        <DialogContent>
          <SignU2f />
        </DialogContent>
        <DialogActions>
          {/* <Button onClick={()=>{props.close() ; props.deleteFunc(props.device.deviceID, props.device.index)}} color="primary" variant="contained">
                          Yes, Delete
                      </Button> */}
          <Button onClick={props.close} color="primary" variant="contained" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

const MMenu = (props: any) => (
  <Menu
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'right',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
);

const StyledMenu = withStyles((theme) => ({
  paper: {
    border: '1px solid #d3d4d5',
  },
}))(MMenu) as typeof MMenu;

const StyledMenuItem = withStyles((theme) => ({
  root: {
    '&:hover': {
      backgroundColor: theme.palette.secondary.main,
      '& .MuiListItemIcon-root, & .MuiListItemText-primary': {
        color: theme.palette.common.white,
      },
    },
  },
}))(MenuItem) as typeof MenuItem;

function EnrolDeviceManu(props: any) {
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <div>
      <Button
        aria-controls="customized-menu"
        aria-haspopup="true"
        variant="contained"
        color="secondary"
        onClick={handleClick}
      >
        Enrol 2FA Device
      </Button>
      <StyledMenu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleClose}>
        <StyledMenuItem onClick={props.changeEnrolMobileDeviceDlgDlg}>
          <ListItemText primary="Mobile Phone" />
        </StyledMenuItem>
        <StyledMenuItem>
          <ListItemText primary="Yubikey (n/a)" />
        </StyledMenuItem>
      </StyledMenu>
    </div>
  );
}
