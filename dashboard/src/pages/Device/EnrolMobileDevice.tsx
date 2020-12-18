import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { makeStyles } from '@material-ui/core/styles';
import QRCode from 'qrcode.react';
import React, {useState} from 'react';
import servicestoreIcon from '../../assets/loginpage/appstore.png';
import PlaystoreIcon from '../../assets/loginpage/playstore.png';
import TrasaLogo from '../../assets/trasa-ni.svg';
import Dialogue from "../../utils/Components/DialogueWrapComponent";
import {TextField} from "@material-ui/core";
import Constants from "../../Constants";
import axios from 'axios';

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
    Height: '100%',
    minHeight: 900,
    margin: '0',
  },

  card: {
    // contained: true,
    textAlign: 'center',
    minWidth: 400,
    minHeight: 300,
    backgroundColor: 'white', // '#d0d3d4', // //rgba(1,1,35,1)
  },
  title: {
    marginBottom: 16,
    fontSize: 20,
    color: '#1A237E',
  },
  input: {
    color: 'teal',
  },
  button: {
    backgroundColor: '#1A237E',
    color: 'white',
    minWidth: 7,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    minWidth: 500,
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
  tfaButton: {
    background: 'white', // '#051384',
    borderRadius: 3,
    border: '1px solid #0b1728ff',
    color: '#051384',
    fontSize: '34px',
    padding: '0 30px',
    // minWidth: 100,
    // minHeight: 70,
    // marginLeft: '10%',
    borderColor: 'grey',
    // boxShadow: '0 3px 5px 2px  #0eafb9 ',
  },
  buttonHead: {
    // background: '#051384',
    borderRadius: 3,
    border: 0,
    color: '#051384',
    fontSize: '20px',
    padding: '0 30px',
    //  minWidth: 200,
    //  minHeight: 100,
    marginLeft: '33%',
  },
  padMiddle: {
    textAlign: 'center',
    // marginLeft: '50%',
  },
  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    fontFamily: 'Open Sans, Rajdhani',
    // padding: '10px 12px',
    width: 'calc(100% - 4px)',

    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldRoot: {
    // paddingRight: theme.spacing(2),
  },
  textFieldInputBig: {
    borderRadius: 4,
    // backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    paddingLeft: theme.spacing(1),
    // marginLeft: theme.spacing(1),
    //    padding: '10px 100px',
    //     width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 14,
    marginLeft: 4,
    fontFamily: 'Open Sans, Rajdhani',
  },
  errorText: {
    color: 'white',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
    background: 'maroon',
  },
  fpBtn: {
    backgroundColor: ' #e6eaea',
    color: 'black',
    minWidth: 17,
    marginRight: '35%',
  },
  enBtn: {
    backgroundColor: '#e6eaea', // '#6085f1',
    color: 'black',
    minWidth: 300,
    // marginLeft: '10%'
  },
  cprightText: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    // marginLeft: '10%'
  },
  noteHeading: {
    fontSize: '13px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

export default function EnrolDevice(props: any) {
  const classes = useStyles();
  const [totpDialogOpen,setTotpDialogOpen] = useState(false);
  return (
    <div className={classes.root}>
      <br />
      <Card className={classes.card}>
        <div className={classes.padMiddle}>
          <img src={TrasaLogo} height={100} width={200} alt="trasa-logo" />

          <br />
          <div className={classes.heading}> 1. Get TRASA app for your mobile device. </div>
          <Button
            onClick={() =>
              window.open('  https://services.apple.com/us/app/trasa/id1411267389', '_blank')
            }
          >
            <img src={servicestoreIcon} alt="servicestoreIcon" width={150} height={40} />
          </Button>
          <Button
            onClick={() =>
              window.open('https://play.google.com/store/services/details?id=com.trasa', '_blank')
            }
          >
            <img src={PlaystoreIcon} alt="playstoreIcon" width={150} height={33} />
          </Button>
          <br />
          {/* <Divider light /> */}
          <br />
        </div>

        <CardContent>
          {/* <Divider  light />  */}
          <div className={classes.heading}>
            {' '}
            2. Once TRASA app is ready on your device, <br /> scan this QR image with the app.
          </div>
          <br />
          <div className={classes.padMiddle}>
            <QRCode
              value={`otpauth://totp/trasa:${props.enrolDeviceDetail.account}?trasaType=private&deviceID=${props.enrolDeviceDetail.deviceID}&issuer=${encodeURIComponent(props.enrolDeviceDetail.orgName)}&secret=${props.enrolDeviceDetail.totpSSC}&trasaURL=${props.enrolDeviceDetail.cloudProxyURL}`}
              size={256}
            />
            <br /> <br />
            <div className={classes.noteHeading}>
              {' '}
              Note: your device must have an internet connection to sync.
            </div>
          </div>
          <br />
          {/* <Divider  light /> */}
          <br />
          <div className={classes.heading}>3. Test your device : </div> <br />
          <Button
            style={{ background: 'navy', color: 'white' }}
            onClick={() => setTotpDialogOpen(true)}
          >
            Done
          </Button>
        </CardContent>
      </Card>
      <ConfirmTOTPDialog open={totpDialogOpen} handleClose={()=>setTotpDialogOpen(false)} deviceID={props.enrolDeviceDetail.deviceID}/>
    </div>
  );
}

type ConfirmTOTPDialogProps = {
  open:boolean,
  deviceID:string,
  handleClose: ()=>void,
}

function ConfirmTOTPDialog(props:ConfirmTOTPDialogProps) {
  const [totp,setTOTP] = useState("")
  const onSubmit = (e: any) => {
    e.preventDefault()
    axios.post(Constants.TRASA_HOSTNAME+"/idp/login/checktfa",{totpCode:totp,deviceID:props.deviceID}).then(r=>{
      if(r.data.status=="success"){
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
    })
  }
  return(

      <Dialogue open={props.open} fullScreen={false} handleClose={props.handleClose} maxWidth={"sm"} title={"Verify TOTP"}>
        <form onSubmit={onSubmit}>
          <TextField
              fullWidth
              label="TOTP"
              // defaultValue={loginData.password}
              onChange={(e)=>{setTOTP(e.target.value)}}
              id="totp"
              name="totp"
              type="number"
              value={totp}
              variant="outlined"
              size="small"
          />
          <Button  type={"submit"} onSubmit={onSubmit} variant="contained" >Verify</Button>
        </form>

      </Dialogue>
  )
}
