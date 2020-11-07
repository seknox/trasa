/*
FirstTimePasswordSetup does two things. validate user setPassword token.
If the token is validated, present change password card or display invalid token alert.

*/

import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
// Form validation
import zxcvbn from 'zxcvbn';
import QueryString from 'query-string';
import Constants from '../../Constants';
import PassStrength from '../../utils/Components/PasswordStrengthMeter';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: '5%',
    marginLeft: '40%',
  },
  checkInbox: {
    fontSize: 15,
    color: 'green',
  },
  // //card
  card: {
    minWidth: 250,
    maxWidth: 300,

  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
    maxWidth: '200px'
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,

    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
  },
  setPasswordtext: {
    fontSize: 25,
    padding: '1px',
    color: 'rgba(1,1,35,1)',
    fontWeight: 'bold',
    fontFamily: 'Open Sans ,Rajdhani',
  },

}));

type ValidatesetPasswordTokenProps = {
  match: any;
  location: any;
};
export default function ValidatesetPasswordToken(props: ValidatesetPasswordTokenProps) {
  // const [veriryTokenStatus, setVeriryTokenStatus] = useState([]);
  const [showSetPasswordCard, setShowSetPasswordCard] = useState(false);
  const [token, setToken] = useState<string | string[] | null | undefined>('');

  const verifyToken = () => {
    const hashed = QueryString.parse(props.location.hash);
    setToken(hashed.token);

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/woa/verify/${hashed.token}`)
      .then((response) => {
        // setUserData(response.data);

        if (response.data.status === 'success') {
          setShowSetPasswordCard(true);
        }
      })
      .catch((error: any) => {
        console.error(error);
      });
  };
  useEffect(verifyToken, []);

  return (
    <div>


        {showSetPasswordCard ? (
          <SetPasswordComponent update={false} token={token} />
        ) : (
            ' Aww Snap. Good guess but your token is invalid '
          )}

    </div>
  );
}

type SetPasswordComponentProps = {
  token: string | string[] | null | undefined;
  update: boolean;
  closeUpdatePassDlg?: () => void;
};

type PassReq = {
  password: string;
  cpassword: string;
  token: string | string[] | null | undefined;
};
// set password
export function SetPasswordComponent(props: SetPasswordComponentProps) {
  const [data, setData] = useState<PassReq>({ password: '', cpassword: '', token: '' });
  // const [country, setCountry] = useState([]);
  const [loader, setLoader] = useState(false);
  // const [password, setPassword] = useState('');
  const [zscore, setZScore] = useState<zxcvbn.ZXCVBNScore>(0);



  const classes = useStyles();

  const zxcvbnscore = () => {
    const score = zxcvbn(data.cpassword);
    const sc = score.score; // -2
    setZScore(sc);
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {

    setData({ ...data, [event.target.name]: event.target.value });
    zxcvbnscore()


  };

  const handleSubmit = (event: React.FormEvent<{}>) => {
     event.preventDefault();

    if (data.password !== data.cpassword) {
      alert('Your passwords does not match')    
      return

    }

    if (zscore < 2) {
      alert('Please enter strong password (the password strength should indicate "Good")') 
      return      
     }


    setLoader(true);
    const reqData = data;
   

    let url = `${Constants.TRASA_HOSTNAME}/api/woa/setup/password`;
    reqData.token = props.token;
    if (props.update === true) {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/my/changepass`;
    }

    axios
      .post(url, reqData)
      .then((response) => {
        setLoader(false);
        // if response status is success, close loader and reqirect user to login page.
        if (response.data.status === 'success') {
          if (props.update && props.closeUpdatePassDlg) {
            props.closeUpdatePassDlg();
          }
          window.location.href = '/login';
        }
        // console.log(response.data);
      })
      .catch((e) => {
        setLoader(false);
        console.error(e);
      });
  };

  return (
    <div className={classes.root}>
      <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
        <Grid item xs={6}>
          <Card className={classes.card}>
            <CardContent>
              <Typography variant="h2" component="h2" className={classes.setPasswordtext}>
                Set your super secret password. <br />
              </Typography>

              {loader ? (
                <Typography variant="h3" component="h2" className={classes.checkInbox}>
                  Good to Login! <br />
                </Typography>
              ) : (
                  ''
                )}

              <br />
              <form onSubmit={handleSubmit}>
                <Grid container spacing={2} alignItems="center" direction="row" justify="center">
                  <Grid item xs={12} >
                    <TextField
                      fullWidth
                      label="Password"
                      onChange={handleChange}
                      name="password"
                      type="password"
                      autoFocus
                      value={data.password}
                      InputProps={{
                        disableUnderline: true,
                        classes: {
                          root: classes.textFieldRoot,
                          input: classes.textFieldInputBig,
                        },
                      }}
                      InputLabelProps={{
                        shrink: true,
                        className: classes.textFieldFormLabel,
                      }}
                    />
                    <div>
                      <PassStrength password={data.password} />
                    </div>
                  </Grid>
                  <Grid item xs={12} >
                    <TextField
                      fullWidth
                      label="Confirm Password"
                      onChange={handleChange}
                      name="cpassword"
                      type="password"
                      value={data.cpassword}
                      InputProps={{
                        disableUnderline: true,
                        classes: {
                          root: classes.textFieldRoot,
                          input: classes.textFieldInputBig,
                        },
                      }}
                      InputLabelProps={{
                        shrink: true,
                        className: classes.textFieldFormLabel,
                      }}
                    />
                    <div>
                      <PassStrength password={data.cpassword} />
                    </div>
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" direction="row" justify="flex-end">
                  <Grid item xs={6}>
                    <Button variant="contained" color="primary" name="submit" type="submit">
                      Submit
                    </Button>
                  </Grid>
                </Grid>

                {loader ? <Progress /> : ''}
              </form>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </div>
  );
}

function Progress() {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <LinearProgress />
    </div>
  );
}
