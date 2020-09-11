import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import TextField from '@material-ui/core/TextField';
import SettingIcon from '@material-ui/icons/SettingsApplicationsSharp';
import axios from 'axios';
import React from 'react';
import Constants from '../../Constants';
 import TfaComponent from '../Auth/Tfa';
import ConsoleMenu from './ConsoleMenu';
import Guacamole from './guacamole-common';
import QueryString from 'query-string';

export class TrasaGWConsole extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      conneted: false,
      isCredDialogOpen: false,
      is2FADialogOpen: false,
      password: '',
      totp: '',
      openMenu: false,
      progress: 0,
    };
    this.guactunnel = new Guacamole.Tunnel(
      `${Constants.TRASA_HOSTNAME_WEBSOCKET}/trasagw/rdp/tunnel`,
    );
    // Instantiate client, using an HTTP tunnel for communications.
    this.client = new Guacamole.Client(this.guactunnel);
  }

  async componentDidMount() {
    const  serviceID = this.props.serviceID;
    const userName = this.props.username;

    let askPassCheck = false;

    try {
      const resp = await axios.get(
        `${Constants.TRASA_HOSTNAME}/api/v1/my/authmeta/${serviceID}/${encodeURIComponent(userName)}`,
      );

      this.setState({ email: resp.data?.data?.[0]?.trasaID });
      this.setState({ isDeviceHygeneRequired: resp.data?.data?.[0]?.isDeviceHygeneRequired });
      askPassCheck = resp.data?.data?.[0]?.isPasswordRequired;
    } catch (e) {
      this.setState({ isDeviceHygeneRequired: false });
      askPassCheck = true;
    }

    //   }

    if (askPassCheck && !this.props.connID) {
      this.askPassword();
    } else {
      this.createGuacConnection();
    }
  }

  toggleDrawer = () => {
    this.setState({ openMenu: !this.state.openMenu });
  };


  createGuacConnection = async ( ) => {
    this.setState({ is2FADialogOpen: false });

     this.display = document.getElementById('displayyyyy');
     this.display.appendChild(this.client.getDisplay().getElement())



    const postData = {
      serviceID: this.props.serviceID,
      hostname: this.props.hostname,
      privilege: this.props.username,
      csrf: localStorage.getItem('X-CSRF'),
      //session: localStorage.getItem('X-SESSION'),
      email: this.props.email,
      optWidth: window.innerWidth - 10,
      optHeight: window.innerHeight - 20,
      appType: this.props.appType,
      rdpProto: this.props.rdpProto,
    };

    if (this.props.connID) {
      postData.connID = this.props.connID;
    }

    if (this.state.password) {
      postData.password = this.state.password;
    }

    this.guactunnel.onasktfa = (params) => {
      // console.log(params)
      // console.log(params.length)
      // console.log(params)
      if(params && params.length===1 && params[0]==="skip" ){
        // console.log("skipping tfa")
        this.client.sendtfa("skipped");
        this.initInput()
        return
      }
      this.show2FA();
    };

    this.initInput = () => {
      // Mouse
      this.mouse = new Guacamole.Mouse(this.client.getDisplay().getElement());

      this.mouse.onmousedown = this.mouse.onmouseup = this.mouse.onmousemove = (mouseState) => {
        this.client.sendMouseState(mouseState);
      };

      // Keyboard
      this.keyboard = new Guacamole.Keyboard(document);

      this.keyboard.onkeydown = (keysym) => {
        this.client.sendKeyEvent(1, keysym);
      };

      this.keyboard.onkeyup = (keysym) => {
        this.client.sendKeyEvent(0, keysym);
      };
      this.client.onclipboard = (stream, a) => {
        stream.onblob = (blob) => {
          console.log(atob(blob));
          this.setState({ clipboardVal: atob(blob) });
          // window.lastClipboardData =atob(blob)
          stream.sendAck();
        };
        // console.log(clip,a)
      };
    };

    this.guactunnel.onerror = (error) => {
      alert(this.getStatusMessage(error));
    };

    // Error handler
    this.client.onerror = (error) => {
      console.log(error, 'client');
      if (error.message === 'InvalidCreds') {
        alert('Invalid Password');
      } else if (error.message === 'TfaRequired') {
        alert('TFA REQUIRED!!!!!');
      } else {
        alert(this.getStatusMessage(error));
      }
    };

    // Connect
    try {
      this.client.connect(JSON.stringify(postData));
    } catch (e) {
      alert(e.message);
    }

    this.client.onstatechange = () => {
      this.setState({ connected: true });
    };
    // Disconnect on close
    window.onunload = () => {
      this.client.disconnect();
    };
  };

  onTfa = (event, tfaMethod, totpCode) => {
    // const client = guacClientRef.current;
    // TODO handle U2FY

    this.setState({ is2FADialogOpen: false });
    this.initInput();
    this.client.sendtfa(totpCode);
  };

  getStatusMessage = (error) => {
    switch (error.code) {
      case 256:
        return 'Unsupported';
      case 512:
        return 'Login Failed';
      case 513:
        return 'Server Busy';
      case 514:
        return 'The upstream server is not responding';
      case 515:
        return 'The upstream server encountered an error.';
      case 516:
        return 'Session Timed Out, Please Login Again';
      case 517:
        return 'Resource Conflict';
      case 518:
        return 'Resource Closed';
      case 519:
        //TODO change message
        return 'The upstream server does not appear to exist, or cannot be reached over the network. In most cases, the upstream server is the remote desktop server.';
      case 520:
        return 'The upstream server is refusing to service connections.';
      case 521:
        return 'Session Conflict';
      case 522:
        return 'The session within the upstream server has ended because it appeared to be inactive';
      case 523:
        return 'The session within the upstream server has been forcibly closed.';
      case 768:
        return 'The parameters of the request are illegal or otherwise invalid.';
      case 769:
        return 'Username/Password invalid';
      case 771:
        return 'Permission was denied';
      case 776:
        return 'The client (usually the browser) is taking too long to respond.';
      case 781:
        return 'The client has sent more data than the protocol allows.';
      case 783:
        return 'The client has sent data of an unexpected or illegal type.';
      case 797:
        return 'The client is already using too many resources. Existing resources must be freed before further requests are allowed.';
      case 3339:
        return error.message;
      default:
        return 'Unknown Error';
    }
  };

  askPassword = () => {
    this.setState({ isCredDialogOpen: true });
  };

  show2FA = () => {
    this.setState({ is2FADialogOpen: true });
  };

  // //uploads file
  onUpload = () => {
    if (this.state.files) {
      const file = this.state.files[0];
      const formData = new FormData();
      formData.append('file', file);
      axios
        .post(`${Constants.TRASA_HOSTNAME}/api/v1/my/upload_file`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
            'X-CSRF': localStorage.getItem('X-CSRF'),
          },
          onUploadProgress: (progressEvent) => {
            const p = Math.floor((progressEvent.loaded * 100) / progressEvent.total);
            this.setState({ progress: p });
          },
        })
        .then((r) => {
          if (r.data.status === 'success') {
            alert('Uploaded');
          } else {
            alert('Upload fail');
          }
        });
    }
  };

  onClipboardChange = (e) => {
    this.setState({ clipboardVal: '' });
    const stream = this.client.createClipboardStream('text/plain');
    stream.onack = (status) => {
      if (status.isError()) {
        console.log(status);
      } else {
        stream.sendBlob(e.target.value);
        stream.onack((s) => {
          if (s.isError()) {
            console.log(s);
          } else {
            console.log(s);
          }
        });
      }

      // Notify of any errors from the Guacamole server
    };
  };

  CtrlAltDel = (e) => {
    const KEYSYM_CTRL = 65507;
    const KEYSYM_ALT = 65513;
    const KEYSYM_DELETE = 65535;
    this.client.sendKeyEvent(1, KEYSYM_CTRL);
    this.client.sendKeyEvent(1, KEYSYM_ALT);
    this.client.sendKeyEvent(1, KEYSYM_DELETE);
    this.client.sendKeyEvent(0, KEYSYM_DELETE);
    this.client.sendKeyEvent(0, KEYSYM_ALT);
    this.client.sendKeyEvent(0, KEYSYM_CTRL);
  };

  onPassChange = (e) => {
    this.setState({ password: e.target.value });
  };

  onPasswordSubmit = (e) => {
    e.preventDefault()
    this.createGuacConnection();
    this.setState({ isCredDialogOpen: false });
  };

  onFileOpen = (e) => {
    this.setState({ files: e.target.files });
  };

  render() {
    return (
      <div>
        {/* <div id={"keyboardContainer"} > <div id={"showKeyboard"} /> </div> */}
        <Grid container>
          <Grid item xs={12}>
            <div>
              <Button size="small" style={floatingButton} onClick={this.toggleDrawer}>
                <SettingIcon />
              </Button>
              <ConsoleMenu
                progress={this.state.progress}
                openMenu={this.state.openMenu}
                toggleDrawer={this.toggleDrawer}
                onClipboardChange={this.onClipboardChange}
                CtrlAltDel={this.CtrlAltDel}
                onFileOpen={this.onFileOpen}
                onUpload={this.onUpload}
                clipboardVal={this.state.clipboardVal}
                appType={this.appType}
              />
              {this.state.connected ? null : <LinearProgress />}
              <div id="displayyyyy" style={styles.display} />
            </div>
            <Dialog
                maxWidth={"lg"}
                open={this.state.is2FADialogOpen}
                closeDialog={() => {
                  this.setState({ is2FADialogOpen: false });
                }}
            >
                <TfaComponent
                    sendTfa={this.onTfa}
                />
            </Dialog>

            <Dialog
              open={this.state.isCredDialogOpen}
              onBackdropClick={() => {
                this.setState({ isCredDialogOpen: false });
              }}
            >
              <DialogTitle>Password not set</DialogTitle>
              <DialogContent>
                <DialogContentText>Enter remote host password</DialogContentText>
                <form onSubmit={this.onPasswordSubmit}>
                  <TextField
                    autoFocus
                    value={this.state.password}
                    onChange={this.onPassChange}
                    margin="dense"
                    id="pass"
                    label="Password"
                    type="password"
                    fullWidth
                  />
                </form>
              </DialogContent>
              <DialogActions>
                <Button
                  size="small"
                  variant="contained"
                  color="primary"
                  onClick={this.onPasswordSubmit}
                >
                  Submit
                </Button>
                <Button size="small" onClick={this.handleClose} color="primary">
                  Cancel
                </Button>
              </DialogActions>
            </Dialog>
          </Grid>
        </Grid>
      </div>
    );
  }
}

const styles = {
  display: {
    // align:"middle",
    margin: 'auto',
    width: '100%',
    height: '100%',
    // display: "inlineBlock",
    border: '3px', // solid green;
    padding: '10px',
  },
};

const floatingButton = {
  backgroundColor: 'navy',
  // width: '100%',
  // maxWidth: '2px',
  color: 'white',
  top: 0,
  right: '50%',
  position: 'absolute',
  zIndex: 9999,
};


export default function GuacConsoleMain(props){

  const [data, setData] = React.useState({serviceID:'',username: '', email: '', connID: '', hostname: '' })
  React.useEffect(() => {
    const hashed = QueryString.parse(props.location.hash);

    if (hashed) {
      setData({serviceID: hashed.serviceID,username: hashed.username, email: hashed.email, connID: hashed.connID, hostname: hashed.hostname})
    }
  }, [props.location.hash])


  return (
      <div>
        { data.username===""?<div></div>:( <TrasaGWConsole
            serviceID={data.serviceID}
            username={data.username}
            email={data.email}
            connID={data.connID}
            hostname={data.hostname}/>)}
      </div>
  )
}
