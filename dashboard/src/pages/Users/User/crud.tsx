// import Paper from '@material-ui/core/Paper';
import Button from '@material-ui/core/Button';
import green from '@material-ui/core/colors/green';
// import { blue100 } from '@material-ui/core/styles/colors';
import purple from '@material-ui/core/colors/purple';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../Constants';
import { Usertype } from '../../../types/users';
import ProgressHOC from '../../../utils/Components/Progressbar';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginTop: '10%',
    marginRight: '3%',
    marginLeft: '3%',
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
    fontSize: 14,
    color: theme.palette.text.secondary,
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
  paper: {
    padding: 16,
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  buttonSuccess: {
    backgroundColor: green[500],
    '&:hover': {
      backgroundColor: green[700],
    },
  },
  fabProgress: {
    color: green[500],
    position: 'absolute',
    top: -6,
    left: -6,
    zIndex: 1,
  },
  buttonProgress: {
    color: green[500],
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12,
  },
  successText: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  content: {
    width: '100%',
    flexGrow: 1,
    backgroundColor: theme.palette.background.default,
    padding: 24,
    height: 'calc(100% - 56px)',
    // height: '100%',
    marginTop: 56,
    [theme.breakpoints.up('sm')]: {
      height: 'calc(100% - 64px)',
      marginTop: 64,
    },
  },
  // handles apppbar flex
  verticalBar: {
    display: 'flex',
    // height: '100%',
    flexDirection: 'column',
    // backgroundColor: 'rgba(1,1,35,1)',
    // alignItems: 'flex-start',
    // justifyContent: 'space-between',
  },
  Up: {
    alignItems: 'flext-start',
    // justifyContent: 'center'
  },
  Down: {
    alignItems: 'flext-end',
    //  justifyContent: 'center'
  },

  // Right drawer
  list: {
    minWidth: 775,
  },
  fullList: {
    minWidth: 775,
  },
  heading: {
    marginTop: 40,
    marginLeft: 50,
    marginBottom: 20,
    color: 'black',
    fontSize: '24px',
    fontFamily: 'Open Sans, Rajdhani',
  },
}));

type UsercrudProps = {
  update: boolean;
  userData: Usertype;
  updateUserTable: (val: any) => void;
  handleVerifyLinkChange: (val: any) => void;
  handleDrawerClose: () => void;
};

export default function UserCrud(props: UsercrudProps) {
  const classes = useStyles();
  const [loader, setLoader] = useState(false);
  const [data, setData] = useState(
    props.update
      ? props.userData
      : {
        firstName: '',
        middleName: '',
        lastName: '',
        email: '',
        password: '',
        userRole: '',
        userName: '',
        userName2: '',
        cpassword: '',
        status: true,
      },
  );
  // const [passMethod, setpassMethod] = useState('selfPassSetup');

  useEffect(() => {
    if (props.update) {
      setData(props.userData);
    }
  }, [props.update, props.userData]);

  function handleChange(event: any) {
    if (event.target.name === 'status') {
      setData({ ...data, [event.target.name]: event.target.value === 'active' });
    } else {
      setData({ ...data, [event.target.name]: event.target.value });
    }
  }

  const handleSubmit = (event: any) => {
    event.preventDefault();

    const req = { user: data };

    setLoader(true);

    axios({
      method: 'post',
      url: `${Constants.TRASA_HOSTNAME}/api/v1/user/${props.update ? 'update' : 'create'}`,
      data: req,
    })
      .then((response) => {
        setLoader(false);

        if (response.data.status === 'success') {
          if (!props.update) props.updateUserTable(response.data.data[0].user);
          props.handleVerifyLinkChange(response.data.data[0].confirmLink);
          props.handleDrawerClose();
        }
      })
      .catch((error) => {
        setLoader(false);
        console.log(error);
      });
  };

  // const handlePasswordOptionChange = (val: any) => (event: any) => {
  //   setpassMethod(val);
  // };

  return (
    <div>
      <div className={classes.heading}>{props.update ? 'Update User' : 'Create User'}</div>
      <Divider light />
      <br />
      {/* loader text */}

      <form onSubmit={handleSubmit} onError={(errors) => console.log(errors)}>
        <Grid item xs={10}>
          First Name:
          <TextField
            fullWidth
            //  label="First Name"
            onChange={handleChange}
            name="firstName"
            value={data.firstName}
            variant="outlined"
            size="small"
          />
        </Grid>
        <Grid item xs={10}>
          Middle Name:
          <TextField
            fullWidth
            //    label="Middle Name"
            onChange={handleChange}
            name="middleName"
            value={data.middleName}
            variant="outlined"
            size="small"
          />
        </Grid>

        <Grid item xs={10}>
          Last Name:
          <TextField
            fullWidth
            //  label="Last Name"
            onChange={handleChange}
            name="lastName"
            value={data.lastName}
            variant="outlined"
            size="small"
          />
        </Grid>
        <Grid item xs={10}>
          Username:
          <TextField
            fullWidth
            // label="Username"
            onChange={handleChange}
            name="userName"
            value={data.userName}
            variant="outlined"
            size="small"
          />
        </Grid>

        <Grid item xs={10}>
          Email:
          <TextField
            fullWidth
            // label="Email"
            onChange={handleChange}
            name="email"
            value={data.email}
            variant="outlined"
            size="small"
          />
        </Grid>

        {/* <Grid item xs={12} >
            Password: <br />

          <FormControlLabel
                    control={
                      <Checkbox
                      checked={this.state.passMethod === "selfPassSetup"}
                      onChange={this.handlePasswordOptionChange("selfPassSetup")}
                      value="selfPassSetup"
                      color="primary"

                     />
                    }
                    label="Let user self setup their password"
                  />
              <FormControlLabel
                    control={
                      <Checkbox
                      checked={this.state.passMethod === "autoGenPass"}
                      onChange={this.handlePasswordOptionChange("autoGenPass")}
                      value="autoGenPass"
                      color="primary"

                     />
                    }
                    label="Auto generate strong password for user."
                  />


         </Grid>


      {(this.state.passMethod === "selfPassSetup" || this.state.passMethod === "autoGenPass") ? '': <div>  <Grid container spaceing={24}> <Grid item xs={6} >
         <TextValidator
         //fullWidth
                label="Password"
                onChange={this.handleChange}
                name="password"
               // disabled = {!this.props.update}
                type="password"
                value={data.password}
                validators=  {[ 'minStringLength: 8']} // {!this.props.update && [ 'minStringLength: 8']}
                errorMessages={[ 'Weak Password']}
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

         </Grid>



         <Grid item xs={6} >
         <TextValidator
               // fullWidth
                label="Confirm Password"
                onChange={this.handleChange}
                name="cpassword"
                type="password"
              //  disabled = {!this.props.update}
                value={data.cpassword}
                validators={ ['isPasswordMatch']}
                errorMessages={['password does not match']}
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


        </Grid> </Grid>    </div>} */}

        <Grid item xs={10}>
          User Role:
          <FormControl fullWidth variant="outlined" size="small">
            {/* <FormHelperText>User Role</FormHelperText> */}
            {/* <InputLabel>{!data.userRole && 'User Role'}</InputLabel> */}
            <Select
              // className={classes.selectCustom}
              name="userRole"
              id="userRole"
              value={data.userRole}
              onChange={handleChange}
              // defaultValue={'selfUser'}
              inputProps={{
                name: 'userRole',
                id: 'userRole',
              }}
            >
              <MenuItem value="orgAdmin" id="orgAdmin">
                Org Admin
              </MenuItem>
              <MenuItem value="selfUser" id="selfUser">
                Normal User
              </MenuItem>
            </Select>
            {/* <FormHelperText>{!data.userRole && 'Default option is Normal User'}</FormHelperText> */}
          </FormControl>
        </Grid>

        <Grid item xs={10}>
          Status:
          <FormControl fullWidth variant="outlined" size="small">
            {/* <FormHelperText>Status</FormHelperText> */}
            {/* <InputLabel>{!data.userRole && 'User Role'}</InputLabel> */}
            <Select
              // className={classes.selectCustom}
              value={data.status ? 'active' : 'disabled'}
              onChange={handleChange}
              // defaultValue={'active'}
              inputProps={{
                name: 'status',
                id: 'status',
              }}
            >
              <MenuItem value="active">active</MenuItem>
              <MenuItem value="disabled">disabled</MenuItem>
            </Select>
            {/* <FormHelperText>{!data.userRole && 'Default option is Normal User'}</FormHelperText> */}
          </FormControl>
        </Grid>
        <Grid container spacing={2} alignItems="center" direction="row" justify="flex-end">
          <Grid item xs={12}>
            {loader ? <ProgressHOC /> : ''}
            <div className={classes.root}>
              <div>
                {loader ? (
                  <Button
                    name="submit"
                    variant="contained"
                    id="submit"
                    disabled
                    color="secondary"
                    type="submit"
                  >
                    Submit
                  </Button>
                ) : (
                  <Button
                    name="submit"
                    id="submit"
                    variant="contained"
                    color="secondary"
                    type="submit"
                  >
                    Submit
                  </Button>
                )}
                <br />
              </div>
            </div>
          </Grid>
        </Grid>
      </form>
      {/* </Grid> */}

      <br />
    </div>
  );
}
