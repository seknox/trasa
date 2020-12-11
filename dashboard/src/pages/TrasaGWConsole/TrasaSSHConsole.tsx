import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, {FormEvent, useEffect, useRef} from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import Constants from '../../Constants';
import '../../utils/styles/xterm.css';
import TfaComponent from '../Auth/Tfa'
import QueryString from 'query-string';
// import TfaDialog from './2faDialog';
// const term = new Terminal();

export function SSHLiveSession(props: any) {
  const [isCredDialogOpen, setCredDialogOpen] = React.useState(false);
  const [is2FADialogOpen, set2FADialogOpen] = React.useState(false);
  const [password, setPassword] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [isDeviceHygeneRequired, setDeviceHygeneRequired] = React.useState(false);

  const termRef = useRef(new Terminal());

  const askPassword = () => {
    // createGuacConnection(prompt("Enter Password"),false)
    setCredDialogOpen(true);
  };

  //
  // TODO define device hygiene type
  const connect = (e:FormEvent<Element> | null, tfaMethod: string,totpCode:string,) => {
    const term = termRef.current;

    set2FADialogOpen(false);




    let wskt: WebSocket;
    if (props.connID) {
      wskt = new WebSocket(`${Constants.TRASA_GUAC_HOSTNAME_WEBSOCKET}/accessproxy/ssh/join`,"xterm");
    } else {
      wskt = new WebSocket(`${Constants.TRASA_GUAC_HOSTNAME_WEBSOCKET}/accessproxy/ssh/connect`,"xterm");
    }

    wskt.onclose = (e) => {
      console.log('Socket is closed. Reconn will be attempted in 1 second.', e.reason);
      // setTimeout(() => {
      //     connect();
      // }, 1000);
    };

    wskt.onerror = (err) => {
      console.error('Socket encountered error: ', err, 'Closing socket');
      wskt.close();
    };

    wskt.onopen = () => {
      const newReq = {
        connID: props.connID,
        //session: localStorage.getItem('X-SESSION'),
        csrf: localStorage.getItem('X-CSRF'),
        serviceID: props.serviceID,
        privilege: props.username,
        totpCode: totpCode,
        tfaMethod,
        password,
        width: window.innerWidth - 10,
        height: window.innerHeight - 20,
        email,
        hostname: props.hostname,
      };

      wskt.send(JSON.stringify(newReq));
      term.onData((data) => {
        //  console.log(data)
        wskt.send(data);
        //  term.write(data);
      });
      // setInterval(() => {
      //         wskt.send("pong");
      // }, 5000);
    };
    wskt.onmessage = (evt) => {
      //  console.log('--------------- ', evt.data)
      // alert("Message is received...");
      term.write(evt.data);
    };
  };

  useEffect(() => {


    let askPassCheck = false;
    axios
      .get(
        `${Constants.TRASA_HOSTNAME}/api/v1/my/authmeta/${props.serviceID}/${encodeURIComponent(props.username)}`,
      )
      .then((resp) => {
        setEmail(resp.data.data[0].trasaID);
        setDeviceHygeneRequired(resp.data.data[0].isDeviceHygeneRequired);
        askPassCheck = resp.data.data[0].isPasswordRequired;

        const term = termRef.current;
        const container = document.getElementById('xterm');

        const fitAddon = new FitAddon();

        term.loadAddon(fitAddon);
        if (container) {
          term.open(container);
        }

        fitAddon.fit();

        termRef.current = term;

        if (askPassCheck && !props.connID) {
          askPassword();
        } else {
          setCredDialogOpen(false);
          if(props.tfaRequired=="true"){
            set2FADialogOpen(true);
          }else {
            connect( null,"","");
          }
        }
      })
      .catch((e) => {
        console.log(e)
        askPassword()
        setDeviceHygeneRequired(false);
        askPassCheck = true;
      });
  }, [props.serviceID, props.username]);

  const onPassChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const onPasswordSubmit = (e:React.FormEvent<any>) => {
    e.preventDefault()
    setCredDialogOpen(false);
    if(props.tfaRequired=="true"){
      set2FADialogOpen(true);
    }else {
      connect( null,"","");
    }
    // set2FADialogOpen(true);
  };

  // message extension to fetch device hygiene
  const messageContentScript = () => {
    console.debug("sending msg to extension")
    window.postMessage(
      {
        direction: 'tsxdashboard',
        message: { type: 'exportDeviceHygiene', data: '' },
      },
      Constants.TRASA_HOSTNAME,
    );
  };


  return (
    <div>
      {/* <input onChange={(e)=>{setState({sessionID:e.target.value})}}/> */}
      {/* <button onClick={connect}>connect</button> */}
      <div id="xterm" style={{ position: 'absolute', bottom: 0, right: 0, left: 0, top: 0 }} />

      <Dialog
    //  fullWidth
      maxWidth='lg'
          open={is2FADialogOpen}
      >
        <TfaComponent

            sendTfa={connect}
            loader={false}
        />
      </Dialog>


      <Dialog
        open={isCredDialogOpen}
        onBackdropClick={() => {
          setCredDialogOpen(false);
        }}
      >
        <DialogTitle>Password not set</DialogTitle>
        <DialogContent>
          <DialogContentText>Enter remote host password</DialogContentText>
          <form onSubmit={onPasswordSubmit}>
            <TextField
              autoFocus
              value={password}
              onChange={onPassChange}
              margin="dense"
              id="pass"
              label="Password"
              type="password"
              fullWidth
            />
          </form>
        </DialogContent>
        <DialogActions>
          <Button size="small" variant="contained" color="primary" onClick={onPasswordSubmit}>
            Submit
          </Button>
          <Button
            size="small"
            onClick={() => {
              setCredDialogOpen(false);
            }}
            color="primary"
          >
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}


export default function TsshConsoleMain(props: any){

  const [data, setData] = React.useState<any>({serviceID:'',username: '', email: '', connID: '', hostname: '' })
  React.useEffect(() => {
    const hashed = QueryString.parse(props.location.hash);

    if (hashed) {
      setData({serviceID: hashed.serviceID,username: hashed.username, email: hashed.email, connID: hashed.connID, hostname: hashed.hostname})
    }
  }, [props.location.hash])


  return (<div>
        {data.username ?
            <SSHLiveSession
                serviceID={data.serviceID}
                username={data.username}
                email={data.email}
                connID={data.connID}
                hostname={data.hostname}/> : <div></div>}
  </div>

  )
}
