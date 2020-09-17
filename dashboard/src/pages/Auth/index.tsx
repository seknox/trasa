import Card from '@material-ui/core/Card';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Typography } from '@material-ui/core';
import Constants from '../../Constants';
import { UserData } from '../../types/auth';
import EnrolMobileDevice from '../Device/EnrolMobileDevice';
import LoginBox from './Login';
import OrgSelect from './OrgSelect';
import TfaBox from './Tfa';

import DialogueWrapper from '../../utils/Components/DialogueWrapComponent';

const useStyles = makeStyles(() => ({
  card: {
    // contained: true,
    textAlign: 'center',
    minWidth: 400,
    minHeight: 300,
    backgroundColor: 'white', // '#d0d3d4', // //rgba(1,1,35,1)
  },
  padMiddle: {
    textAlign: 'center',
    // marginLeft: '50%',
  },
}));

type LoginIndexProps = {
  userData: UserData;
  autofillEmail: boolean;
  intent:
  | 'AUTH_REQ_DASH_LOGIN'
  | 'AUTH_REQ_CHANGE_PASS'
  | 'AUTH_REQ_ENROL_DEVICE'
  | 'AUTH_HTTP_ACCESS_PROXY';
  title: string;
  showForgetPass: boolean;
  tfaRequired?: boolean;
  proxyDomain: string | string[] | null | undefined;
  // loader: boolean;
  setData: (value: any) => void;

  changeHasAuthenticated: () => void | null;
};

const loginResps = {
  AUTH_RESP_NOTIF_LICENSE: 'AUTH_RESP_NOTIF_LICENSE',
  AUTH_RESP_ENROL_DEVICE: 'AUTH_RESP_ENROL_DEVICE',
  AUTH_RESP_TFA_REQUIRED: 'AUTH_RESP_TFA_REQUIRED',
  AUTH_RESP_SELECT_ORG: 'AUTH_RESP_SELECT_ORG',
  AUTH_RESP_CHANGE_PASS: 'AUTH_RESP_CHANGE_PASS',
  AUTH_RESP_RESET_PASS: 'AUTH_RESP_RESET_PASS',
  AUTH_RESP_TFA_DH_REQUIRED: 'AUTH_RESP_TFA_DH_REQUIRED',
};

// message extension to set session details
function msgContentScriptWithSessionDetails(d: any) {
  window.postMessage(
    {
      direction: 'tsxdashboard',
      message: { type: 'setSession', data: d },
    },
    Constants.TRASA_HOSTNAME,
  );
}

// LoginHandler wraps user login components  - (1) login (2) tfa (3) enrol device; and manages authentication flow based on response send by trasacore.
export default function AuthIndex(props: LoginIndexProps) {
  const classes = useStyles();
  const [loginData, setLoginData] = useState({
    email: '',
    password: '',
    userName: '',
    orgID: '',
    intent: '',
    idpName: '',
  });

  const [dhRequired, setdhRequired] = useState(false);
  useEffect(() => {
    if (props.userData.email !== '') {
      setLoginData(props.userData);
    }
  }, [props.userData]);

  const [enrolDeviceDetail, setenrolDeviceDetail] = useState({});
  const [, setchangePassToken] = useState('');
  const [tfaToken, setTfaToken] = useState('');
  const [orgs, setOrgs] = useState([]);

  const [, setchangePassDlgState] = useState(false);
  // const [showTfaComponent, setShowTfaComponent] = useState(false);
  const [loader, setLoader] = useState(false);
  const [tfaRequired, setTfaRequired] = useState<boolean | undefined>(false);
  const [tfaIntent, setTfaIntent] = useState('');

  useEffect(() => {
    setTfaRequired(props.tfaRequired);
  }, [props.tfaRequired]);

  useEffect(() => {
    if (props.intent === 'AUTH_HTTP_ACCESS_PROXY') {
      setTfaRequired(true);
    }
  }, [props.intent]);

  const [enrolDevice, setEnrolDevice] = useState(false);
  // const [signDeviceDialogueState, setsignDeviceDialogueState] = useState(false);

  function handleLoginDataChange(event: React.ChangeEvent<HTMLInputElement>) {
    setLoginData({ ...loginData, [event.target.name]: event.target.value });
  }

  // function checkVersion() {
  //   axios.get(Constants.TRASA_HOSTNAME+"/api/v1/check").then(r=>{
  //     if(r.data.status==="success"){
  //       localStorage.setItem("isEE","TRUE")
  //     }
  //   })
  // }

  function sendLoginRequest(
    e: React.FormEvent<Element>,
    intent: string,
    idpName: string,
    orgID: string,
  ) {
    setLoader(true);
    e.preventDefault();

    // checkVersion()

    if (intent === '') {
      loginData.intent = 'AUTH_REQ_DASH_LOGIN';
    } else {
      loginData.intent = intent;
    }

    loginData.orgID = orgID || '';
    loginData.intent = 'AUTH_REQ_DASH_LOGIN';
    if (props.autofillEmail === true) {
      if (props.userData.idpName !== 'freeipa') {
        loginData.email = props.userData.userName;
      } else {
        loginData.email = props.userData.email;
      }
      loginData.intent = props.intent;
      loginData.orgID = props.userData.orgID;
    }

    loginData.idpName = idpName;

    let reqUrl = `${Constants.TRASA_HOSTNAME}/idp/login`;

    if (intent === 'AUTH_REQ_FORGOT_PASS') {
      setTfaIntent('AUTH_REQ_FORGOT_PASS');
      reqUrl = `${Constants.TRASA_HOSTNAME}/api/woa/forgotpass`;
    }

    axios
      .post(reqUrl, loginData)
      .then((response) => {
        setLoader(false);
        if (response.data.intent === loginResps.AUTH_RESP_NOTIF_LICENSE) {
          window.location.href = '/woa/activate';
        }
        if (response.data.intent === loginResps.AUTH_RESP_TFA_REQUIRED) {
          setTfaRequired(true);
          setTfaToken(response.data.data[0]);
        }

        if (response.data.intent === loginResps.AUTH_RESP_TFA_DH_REQUIRED) {
          setTfaRequired(true);
          setTfaToken(response.data.data[0]);
          setdhRequired(true);
          // prepare and fetch device hygiene here?
        }
        if (response.data.intent === loginResps.AUTH_RESP_ENROL_DEVICE) {
          setEnrolDevice(true);
          setenrolDeviceDetail(response.data.data[0]);
        }

        if (response.data.status === loginResps.AUTH_RESP_SELECT_ORG) {
          setOrgs(response.data.data[0]);
        }
      })

      .catch((error: any) => {
        setLoader(false);
        console.log(error);
      });
  }

  function accessProxyTfa(tfaMethod: string, totpCode: string, extToken: string) {
    if (tfaMethod === 'U2FY') {
      // setsignDeviceDialogueState(true);
    } else {
      const tfaReq = {
        tfaMethod,
        totpCode,
        extToken,
        hostName: props.proxyDomain,
      };

      const url = `${Constants.TRASA_HOSTNAME}/auth/accessproxy/http`;

      axios
        .post(url, tfaReq)
        .then(async (r) => {
          setLoader(false);
          if (r.data.status === 'success') {
            await msgContentScriptWithSessionDetails({
              domain: props.proxyDomain,
              sessionID: r.data.data?.[0]?.sessionID,
              csrfToken: r.data.data?.[0]?.csrfToken,
              sessionRecord: r.data.data?.[0]?.sessionRecord,
            });

            // function to update device hygiene state and remoe listener
            const updateState = async (event: any) => {
              if (
                event.source === window &&
                event.data.direction &&
                event.data.direction === 'trasaExt'
              ) {
                try {
                  const confirm = await event.data.message;
                  console.log('confirm: ', confirm);
                  if (confirm && confirm === 'done') {
                    window.removeEventListener('message', updateState);
                    // reload here
                    window.location.assign(`https://${props.proxyDomain}`);
                  }
                } catch (err) {
                  console.error(err);
                  alert('Access proxy requires TRASA browser extension installed in your browser.');
                  // alert here that trasa could not contact extension.
                }
              }
            };

            window.addEventListener('message', updateState);
          }
        })

        .catch((error) => {
          setLoader(false);
          console.log(error);
        });
    }
  }

  const [viewDlg, setviewDlg] = useState(false);
  const changeViewDlgState = () => {
    setviewDlg(!viewDlg);
  };
  const closeViewDlg = () => {
    setviewDlg(false);
  };

  function loginTfa(
    intent: string,
    tfaMethod: string,
    totpCode: string,
    extID: string,
    clientPubKey: string,
    dh: any,
  ) {
    if (tfaMethod === 'U2FY') {
      // setsignDeviceDialogueState(true);
    } else {
      if (intent === '') {
        if (tfaIntent === 'AUTH_REQ_FORGOT_PASS') {
          intent = tfaIntent;
        } else {
          intent = props.intent;
        }
      }

      const tfaReq = {
        tfaMethod,
        totpCode,
        token: tfaToken,
        intent,
        extID,
        deviceHygiene: dh,
        clientPubKey,
      };

      const url = `${Constants.TRASA_HOSTNAME}/idp/login/tfa`;

      axios
        .post(url, tfaReq)
        .then((r) => {
          setLoader(false);
          if (r.data.status === 'success') {
            switch (true) {
              case r.data.intent === 'AUTH_RESP_RESET_PASS':
                setchangePassToken(r.data.data[0]);
                setchangePassDlgState(true);
                break;
              case r.data.intent === 'AUTH_RESP_FORGOT_PASS':
                changeViewDlgState();

                break;

              case r.data.intent === 'AUTH_RESP_ENROL_DEVICE':
                props.setData(r.data.data[0]);
                props.changeHasAuthenticated();
                break;
              case r.data.intent === 'AUTH_RESP_CHANGE_PASS':
                props.setData(r.data.data[0]);
                props.changeHasAuthenticated();
                break;
              default:
                localStorage.setItem('X-CSRF', r.data.data[0].CSRFToken);
                if (r.data.data[0].user.userRole === 'orgAdmin') {
                  const urlParams = new URLSearchParams(window.location.search);
                  const nextUrl = urlParams.get('next');
                  if (nextUrl) {
                    window.location.href = nextUrl;
                  } else {
                    window.location.href = '/overview';
                  }
                } else {
                  window.location.href = '/my';
                }
            }
          }
        })

        .catch((error) => {
          setLoader(false);
          console.log(error);
        });
    }
  }

  // message extension to fetch device hygiene
  function messageContentScriptForExtToken() {
    window.postMessage(
      {
        direction: 'tsxdashboard',
        message: { type: 'exportExtToken', data: '' },
      },
      Constants.TRASA_HOSTNAME,
    );
  }

  async function sendTfa(e: React.FormEvent<Element>, tfaMethod: string, totpCode: string) {
    e.preventDefault();

    if (props.intent === 'AUTH_HTTP_ACCESS_PROXY') {
      setLoader(true);
      await messageContentScriptForExtToken();

      // function to update device hygiene state and remoe listener
      const updateState = async (event: any) => {
        if (
          event.source === window &&
          event.data.direction &&
          event.data.direction === 'trasaExt'
        ) {
          try {
            const extToken = await event.data.message;

            if (extToken) {
              // remove listener here
              window.removeEventListener('message', updateState);
              accessProxyTfa(tfaMethod, totpCode, extToken);
            }
          } catch (err) {
            console.error(err);
            alert('Access proxy requires TRASA browser extension installed in your browser.');
            // alert here that trasa could not contact extension.
          }
        }
      };

      window.addEventListener('message', updateState);
    } else if (dhRequired) {
      setLoader(true);
      const messageContentScript = () => {
        console.debug('sending msg to extension');
        console.log('tfaToken: ', tfaToken);
        window.postMessage(
          {
            direction: 'tsxdashboard',
            message: { type: 'exportDeviceHygiene', data: tfaToken },
          },
          Constants.TRASA_HOSTNAME,
        );
      };

      messageContentScript();
      // function to update device hygiene state and remoe listener
      const updateState = (event: any) => {
        if (
          event.source === window &&
          event.data.direction &&
          event.data.direction === 'trasaExt'
        ) {
          // use event.data.message to receive object and use it.

          // remove listener here
          window.removeEventListener('message', updateState);
          try {
            console.log(event.data);
            const deviceHyg = event.data && event.data.message;
            // console.log('received::::::::; ', event.data);
            //  connect(totpCode, tfaMethod, {}, deviceHyg);
            loginTfa(
              '',
              tfaMethod,
              totpCode,
              deviceHyg.extID,
              deviceHyg.dh.clientPubKey,
              deviceHyg.dh.encryptedDH,
            );
          } catch (e) {
            console.log('failed::::::::; ', e);
            // connect(totpCode, tfaMethod, {}, {});
          }
        }
      };

      window.addEventListener('message', updateState);
    } else {
      setLoader(true);
      loginTfa('', tfaMethod, totpCode, '', '', '');
    }
  }

  function returnComponent() {
    switch (true) {
      // return TFA page
      case tfaRequired:
        return (
          <Card className={classes.card}>
            <TfaBox sendTfa={sendTfa} loader={loader} />
          </Card>
        );

      // return enrolDevice
      case enrolDevice:
        return (
          <Card className={classes.card}>
            <EnrolMobileDevice
              enrolDeviceDetail={enrolDeviceDetail}
              sendLoginRequest={sendLoginRequest}
            />
          </Card>
        );

      // return Login Page
      default:
        return (
          <LoginBox
            handleLoginDataChange={handleLoginDataChange}
            sendLoginRequest={sendLoginRequest}
            loginData={loginData}
            autofillEmail={props.autofillEmail}
            userData={props.userData}
            intent={props.intent}
            title={props.title}
            showForgetPass={props.showForgetPass}
            loader={loader}
          />
        );
    }
  }
  return (
    <div>
      {/* <OrgSelect orgs={orgs} submitLoginRequest={sendLoginRequest} /> */}
      {returnComponent()}
      <DialogueWrapper
        open={viewDlg}
        handleClose={closeViewDlg}
        title=""
        maxWidth="md"
        fullScreen={false}
      >
        <Typography variant="h5">We have sent you a reset link. Check your email.</Typography>
      </DialogueWrapper>
    </div>
  );
}
