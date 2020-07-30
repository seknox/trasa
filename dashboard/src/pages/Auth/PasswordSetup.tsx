import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import React, { useEffect, useState } from 'react';
// Form validation
import zxcvbn from 'zxcvbn';
import PassStrength from '../../utils/Components/PasswordStrengthMeter';
import Progress from '../../utils/Components/Progressbar';
import { SetupOrChangePass, VerifyToken } from './api/auth';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: '5%',
    marginLeft: '30%',
  },
  checkInbox: {
    fontSize: 15,
    color: 'green',
  },
  // //card
  card: {
    minWidth: 275,
    // marginLeft: '50%',
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    // padding: '10px 100px',
    // width: 'calc(100% - 4px)',
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
  // const [userData, setUserData] = useState<any>();

  const token = props.match.params;

  useEffect(() => {
    VerifyToken(token, setShowSetPasswordCard);
  }, []);

  return (
    <div>
      <div>
        {/* Herein starts inner comopnents. we add create user component and user view component */}

        {showSetPasswordCard ? (
          <SetPasswordComponent update={false} token={token} />
        ) : (
          ' Aww Snap. Good guess but your token is invalid '
        )}
      </div>
    </div>
  );
}

type SetPasswordComponentProps = {
  token: string;
  update: boolean;
  closeUpdatePassDlg?: () => void;
};

// set password
export function SetPasswordComponent(props: SetPasswordComponentProps) {
  const classes = useStyles();
  const [data, setData] = useState({ password: '', cpassword: '', token: '' });
  const [loading, setLoading] = useState(true);
  const [progress, setProgress] = useState(true);
  const [zscore, setZScore] = useState<zxcvbn.ZXCVBNScore>();

  const { token } = props;

  const zxcvbnscore = () => {
    const score = zxcvbn(data.cpassword);
    const sc = score.score; // -2
    setZScore(sc);
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    // const email = e.target.value;
    // setState({data});
    setData({ ...data, [event.target.name]: event.target.value });

    // setState({country: event.target.value});
    if (event.target.name === 'cpassword') {
      zxcvbnscore();
    }
  };

  const handleSubmit = (event: React.FormEvent<{}>) => {
    setProgress(true);
    setLoading(true);
    const reqData = data;
    event.preventDefault();
    reqData.token = token;
    SetupOrChangePass(props.update, reqData, props.closeUpdatePassDlg);
  };

  return (
    <div>
      <div className={classes.root}>
        <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
          <Grid item xs={6}>
            <Card className={classes.card}>
              <CardContent>
                <Typography variant="h2" component="h2" className={classes.setPasswordtext}>
                  Set your super secret password. <br />
                </Typography>

                {!loading ? (
                  <Typography variant="h3" component="h2" className={classes.checkInbox}>
                    Good to Login! <br />
                  </Typography>
                ) : (
                  ''
                )}

                <br />
                <form onSubmit={handleSubmit}>
                  <Grid container spacing={2} alignItems="center" direction="row" justify="center">
                    <Grid item xs={6}>
                      <TextField
                        fullWidth
                        label="Password"
                        onChange={handleChange}
                        name="password"
                        type="password"
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
                        {/* <TextField  className={classes.selectCustom} autoComplete="off" type="password" onChange={e => setState({ password: e.target.value })} /> */}
                        <PassStrength password={data.password} />
                      </div>
                    </Grid>
                    <Grid item xs={6}>
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
                        {/* <TextField  className={classes.selectCustom} autoComplete="off" type="password" onChange={e => setState({ password: e.target.value })} /> */}
                        <PassStrength password={data.cpassword} />
                      </div>
                    </Grid>
                  </Grid>

                  <Grid
                    container
                    spacing={2}
                    alignItems="center"
                    direction="row"
                    justify="flex-end"
                  >
                    <Grid item xs={6}>
                      <Button
                        disabled={zscore ? zscore <= 2 : undefined}
                        variant="contained"
                        color="primary"
                        type="submit"
                      >
                        Submit
                      </Button>
                    </Grid>
                  </Grid>

                  {progress ? <Progress /> : ''}
                </form>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </div>
    </div>
  );
}
