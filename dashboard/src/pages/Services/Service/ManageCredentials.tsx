import Button from '@material-ui/core/Button';
import purple from '@material-ui/core/colors/purple';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Fab from '@material-ui/core/Fab';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import classNames from "classnames";
import Typography from '@material-ui/core/Typography';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
    // marginBotton: '5%',
  },
  paper: {
    backgroundColor: '#fdfdfd',
    // backgroundColor: '#00001aff', //'#E0E0E0', 'rgba(10,34,52,1)' // #011019
    // minWidth: 400,
    // minHeight: 300,
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  // card
  card: {
    minWidth: 275,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    marginBottom: 16,
    // fontSize: 54,
    color: '#311B92', // theme.palette.text.secondary,
    // backgroundColor: '#1A237E',
  },
  pos: {
    marginBottom: 12,
    color: theme.palette.text.secondary,
  },

  // form
  formControl: {
    margin: theme.spacing(1),
  },
  inputLabelFocused: {
    color: purple[500],
  },
  inputInkbar: {
    '&:after': {
      backgroundColor: purple[500],
    },
  },
  textFieldRoot: {
    padding: 1,
    'label + &': {
      marginTop: theme.spacing(3),
    },
  },

  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(1),
    width: 'calc(100% - 4px)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldInputBig: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(1),
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  textFieldFormLabel: {
    fontSize: 18,
  },

  fab: {
    margin: theme.spacing(2),
    background: '#000080',
  },
  fab2: {
    margin: theme.spacing(2),
  },
  users: {
    margin: theme.spacing(2),
  },
  dividerInset: {
    margin: `5px 0 0 ${theme.spacing(9)}px`,
  },
  credbuttons: {
    marginTop: 20,
    background: 'navy',
    color: 'white',
  },
  fabButton: {
    margin: theme.spacing(2),
    background:
      'linear-gradient(to right, #021B79, #0575E6)' /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */,
  },
  Warning: {
    color: 'maroon',
    fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  WarningButton: {
    color: 'maroon',
    //  fontSize: '44px',
    fontFamily: 'Open Sans, Rajdhani',
    '&:hover': {
      color: theme.palette.common.white,
      background: 'maroon',
    },
  },
  settingHeader: {
    color: 'black',
    fontSize: '15px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  settingSHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  selectCustom: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 17,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

export default function ManageCredential(props: any) {
  const [appData, setAppData] = useState({});
  const [managedUsers, setManagedUsers] = useState<any[][]>([]);
  const [removeDialogueOpen, setremoveDialogueOpen] = useState(false);
  const [removeCredObj, setremoveCredObj] = useState({ username: '' });
  const [loader, setLoader] = useState(false);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/services/${props.urlID}`)
      .then((response) => {
        if (response.status === 403) {
          window.location.href = '/login';
        }

        if (!response.data.data) {
          console.log(response.data.data);
          return;
        }

        const managedUsers = response.data.data[0].managedAccounts.split(',');

        const accounts = managedUsers.map(function (v: any) {
          return [v, '**********'];
        });
        accounts.shift();

        setAppData(response.data.data[0]);
        setManagedUsers(accounts);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [props.ID]);

  const onViewPassword = (appuser: string, serviceID: string, type: string) => {
    setLoader(true);
    const req = {
      username: appuser,
      serviceID,
      type,
    };

    const reqData = JSON.stringify(req);
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/creds/view`, reqData)
      .then((response) => {
        setLoader(false);
        const tempUsers = managedUsers.map((v) =>
          v[0] === appuser ? [response.data.data[0].username, response.data.data[0].credential] : v,
        );

        if (response.data.status === 'success') {
          setManagedUsers(tempUsers);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const removeCredPrompt = (appuser: string, serviceID: string, serviceType: string) => {
    const req = {
      username: appuser,
      serviceID,
      type: serviceType,
    };

    setremoveCredObj(req);
    setremoveDialogueOpen(true);
  };

  const removeCred = () => {
    const reqData = JSON.stringify(removeCredObj);

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/creds/delete`, reqData)
      .then((response) => {
        setLoader(true);

        // console.log(tempUsers)
        if (response.data.status === 'success') {
          const lmanagedUsers = managedUsers;
          const index = returnIndex(lmanagedUsers, removeCredObj.username);

          managedUsers.splice(index, 1);
          setLoader(false);
          setManagedUsers(lmanagedUsers);
          setremoveDialogueOpen(false);
        }

        // console.log(response.data);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  // returnIndex returns index of managed users array element
  function returnIndex(arr: any, val: any): number {
    for (let i = 0; i < arr.length; i++) {
      if (arr[i][0] === val) {
        return i;
      }
    }
    return 0;
  }

  function checkIfUserExists(arr: any, val: any) {
    for (let i = 0; i < arr.length; i++) {
      if (arr[i][0] === val) {
        return true;
      }
    }

    return false;
  }

  const updateManagedUsers = (username: any) => {
    const userArr = managedUsers;

    const check = checkIfUserExists(userArr, username);

    if (!check) {
      userArr.unshift([username, '**********']);

      setManagedUsers(userArr);
    }
  };

  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Grid container spacing={4} direction="row" alignItems="center" justify="center">
        <Grid item xs={12} sm={12} md={9}>
          <Paper className={classes.paper}>
            <PasswordOrKeyInputComponent
              serviceDetail={appData}
              updateManagedUsers={updateManagedUsers}
            />
            <Divider light />
            <ViewManagedCreds
              serviceDetail={appData}
              managedUsers={managedUsers}
              onViewPassword={onViewPassword}
              removeCredPrompt={removeCredPrompt}
            />
          </Paper>
        </Grid>
      </Grid>

      <Dialog open={removeDialogueOpen}>
        <DialogTitle>Are you sure?</DialogTitle>
        <DialogContent>
          <DialogActions>
            <Button onClick={removeCred} color="primary" variant="contained">
              Yes, Delete
            </Button>
            <Button
              onClick={() => {
                setremoveDialogueOpen(false);
              }}
              color="primary"
              variant="contained"
              autoFocus
            >
              No
            </Button>
          </DialogActions>
        </DialogContent>
      </Dialog>
    </div>
  );
}

function PasswordOrKeyInputComponent(props: any) {
  const [data, setData] = useState({
    type: 'password',
    username: '',
    serviceID: '',
    credential: '',
  });
  const [loader, setLoader] = useState(false);
  // const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  function handleChange(event: any) {
    setData({ ...data, [event.target.name]: event.target.value });
  }

  function handleSubmit() {
    setLoader(true);
    if (data.username.length === 0) {
      return;
    }

    data.serviceID = props.serviceDetail.ID;

    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/services/creds/store`, data)
      .then((response) => {
        setLoader(false);

        if (response.data.status === 'success') {
          props.updateManagedUsers(data.username);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const classes = useStyles();

  return (
    <div>
      <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
        <Grid item xs={12}>
          <Typography variant="h2"> Manage Credentials </Typography>
        </Grid>
      </Grid>

      <Grid container spacing={2}>
        <Grid item xs={2} sm={2} md={2}>
          {/* <Typography variant="h4">Cred type</Typography> */}
          <FormControl variant="outlined" className={classes.formControl}>
            <InputLabel id="demo-simple-select-outlined-label">Cred Type</InputLabel>
            <Select
              label="Credential Type"
              value={data.type}
              name="type"
              onChange={handleChange}
              inputProps={{
                classes: {
                  // root: classes.selectCustom
                },
              }}
            >
              <MenuItem value="password">
                <div className={classes.settingSHeader}>Password </div>
              </MenuItem>
              <MenuItem value="key">
                <div className={classes.settingSHeader}>SSH Key</div>
              </MenuItem>
            </Select>
          </FormControl>
        </Grid>

        <Grid item xs={4} sm={4} md={4}>
          <TextField
            fullWidth
            label="Username"
            onChange={handleChange}
            name="username"
            value={data.username}
            autoComplete="off"
            // type='hidden'
            // validators={['required']}
            // errorMessages={['this field is required']}
            InputProps={{
              disableUnderline: true,
              autoComplete: 'off',
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
        </Grid>
        <Grid item xs={5} sm={5} md={5}>
          <TextField
            fullWidth
            label="Credential"
            onChange={handleChange}
            name="credential"
            rows={data.type === 'key' ? '10' : '2'}
            multiline={data.type === 'key'}
            type={showPassword || data.type === 'key' ? 'text' : 'password'}
            value={data.credential}
            InputProps={{
              disableUnderline: true,
              autoComplete: 'off',
              classes: {
                root: classes.textFieldRoot,
                input: classes.textFieldInputBig,
              },
              endAdornment: (
                <InputAdornment position="start">
                  <IconButton
                    aria-label="Toggle password visibility"
                    onClick={handleClickShowPassword}
                  >
                    {showPassword ? <Visibility /> : <VisibilityOff />}
                  </IconButton>
                </InputAdornment>
              ),
            }}
            InputLabelProps={{
              shrink: true,
              className: classes.textFieldFormLabel,
            }}
          />
        </Grid>

        <Grid item xs={1} sm={1} md={1}>
          <Fab
            size="small"
            color="secondary"
            aria-label="Add"
            onClick={handleSubmit}
            className={classes.fab}
            component="button"
            type="submit"
          >
            <AddIcon />
          </Fab>
        </Grid>
      </Grid>
    </div>
  );
}

function ViewManagedCreds(props: any) {
  const [showPassword, setShowPassword] = useState(false);
  const [password, setPassword] = useState('');

  function handleChange(e: any) {
    setPassword(e.target.value);
  }

  function handleClickShowPassword() {
    setShowPassword(!showPassword);
  }

  const classes = useStyles();

  return (
    <div>
      {props.managedUsers.map((v: any) => (
        <Grid container spacing={2}>
          <Grid item xs={6} sm={6} md={6}>
            <div className={classes.users}>
              <Typography variant="h4">{v[0]}</Typography>
            </div>
          </Grid>
          <Grid item xs={3} sm={3} md={3}>
            <TextField
              fullWidth
              type={showPassword ? 'text' : 'password'}
              label="Password"
              value={v[1]}
              onChange={(e) => handleChange(e)}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="start">
                    <IconButton
                      aria-label="Toggle password visibility"
                      onClick={handleClickShowPassword}
                    >
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={1} sm={1} md={1}>
            <Button
              variant="contained"
              size="small"
              className={classes.credbuttons}
              onClick={() => {
                props.onViewPassword(v[0], props.serviceDetail.ID, 'password');
              }}
            >
              Fetch Password
            </Button>
          </Grid>

          <Grid item xs={1} sm={1} md={1}>
            <Button
              variant="contained"
              size="small"
              className={classes.credbuttons}
              onClick={() => {
                props.onViewPassword(v[0], props.serviceDetail.ID, 'key');
              }}
            >
              Fetch Key
            </Button>
          </Grid>
          <Grid item xs={1} sm={1} md={1}>
            <Fab
              size="small"
              className={classes.fab2}
              component="button"
              onClick={() => {
                props.removeCredPrompt(
                  v[0],
                  props.serviceDetail.ID,
                  props.serviceDetail.serviceType,
                );
              }}
            >
              <DeleteIcon color="error" />
            </Fab>
          </Grid>
        </Grid>
      ))}
    </div>
  );
}
