import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import Chip from '@material-ui/core/Chip';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import Divider from '@material-ui/core/Divider';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import InputLabel from '@material-ui/core/InputLabel';
import LinearProgress from '@material-ui/core/LinearProgress';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import Stepper from '@material-ui/core/Stepper';
import { withStyles } from '@material-ui/core/styles';
import Switch from '@material-ui/core/Switch';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import AccountCircleIcon from '@material-ui/icons/AccountCircle';
import PolicyIcon from '@material-ui/icons/Assignment';
import Send from '@material-ui/icons/Send';
// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import classNames from 'classnames';
import Downshift from 'downshift';
import deburr from 'lodash/deburr';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import Constants from '../../../../Constants';

const styles = (theme) => ({
  container: {
    // display: 'flex',
    // flexWrap: 'wrap',
  },
  paper: {
    // display: 'flex',
    // justifyContent: 'center',
    // flexWrap: 'wrap',
    // padding: theme.spacing(1) / 2,
  },
  root: {
    flexGrow: 1,
    // marginLeft: '4%',
    // flexDirection:'column'
  },
  stepperRoot: {
    flexgrow: 1,
    // marginBotton: '5%',
  },
  expiry: {
    marginLeft: 10,
    width: 100,
  },
  useradd: {
    flexGrow: 1,
  },
  dialogue: {
    // width: '500%',
    // maxWidth: '500%',
  },
  formControl: {
    // margin: theme.spacing(1),
    // minWidth: 120,
    maxWidth: 500,
  },
  stepperContent: {
    flexGrow: 1,
    marginTop: 50,
    marginBottom: 10,
  },
  chips: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  chip: {
    margin: theme.spacing(1) / 4,
  },

  backButton: {
    marginRight: theme.spacing(1),
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },

  timePickerTextField: {
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
    width: 100,
  },

  timePickerContainer: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  button: {
    margin: theme.spacing(1),
    marginTop: '5%',
  },
  rightIcon: {
    marginLeft: theme.spacing(1),
  },
  snackbar: {
    position: 'absolute',
    background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
    borderRadius: 3,
    border: 0,
    color: 'white',
    height: 48,
    padding: '0 30px',
    boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .30)',
  },
  checkInbox: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  stepperButtonLeft: {
    // marginTop:'15%',
    marginLeft: theme.spacing(1),
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },

  stepperButtonRight: {
    marginTop: '15%',
    marginRight: '10%',
  },
  newPaper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),

    color: theme.palette.text.secondary,
  },
  table: {
    marginLeft: 50,
  },
  manageUser: {
    padding: theme.spacing(2),
  },
  menuHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
  paperAuto: {
    position: 'absolute',
    zIndex: 1,
    marginTop: theme.spacing(1),
    left: 0,
    right: 0,
    minHeight: '20px',
  },
  chip: {
    margin: `${theme.spacing(1) / 2}px ${theme.spacing(1) / 4}px`,
  },
  inputRoot: {
    flexWrap: 'wrap',
  },
  inputInput: {
    width: 'auto',
    flexGrow: 1,
  },
  divider: {
    height: theme.spacing(2),
  },
});

function getSteps() {
  return ['Select user groups', 'Assign Policy', 'Review'];
}

class ServicegroupUserGroupMap extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
      activeStep: 0,
      progress: false,
      loading: true,
      error: false,
      inputValue: '',
      userGInput: '',
      policyInput: '',
      selectedItem: [],
      selectedPolicies: [],
      selectedUserGroups: [],
      getUserGroupsSt: false,
      policyState: false,
      userName: '',
    };
  }

  getPolicies = (val) => {
    this.setState({ selectedPolicies: val });
  };

  getUserGroups = (val) => {
    this.setState({ selectedUserGroups: val });
  };

  handleChangeUsername = (e) => {
    this.setState({ userName: e.target.value });
  };

  handleKeyDown = (type) => (event) => {
    if (type == 'usergroup') {
      const { userGInput, selectedUserGroups } = this.state;
      // const { inputValue, selectedItem } = this.state;
      if (selectedUserGroups.length && !userGInput.length && event.key === 'Backspace') {
        this.setState({
          selectedUserGroups: selectedUserGroups.slice(0, selectedUserGroups.length - 1),
        });
      }
    } else {
      const { policyInput, selectedPolicies } = this.state;
      // const { inputValue, selectedItem } = this.state;
      if (selectedPolicies.length && !policyInput.length && event.key === 'Backspace') {
        this.setState({ selectedPolicies: selectedPolicies.slice(0, selectedPolicies.length - 1) });
      }
    }
  };

  handleInputChange = (type) => (event) => {
    if (type == 'usergroup') {
      this.setState({ userGInput: event.target.value });
    } else {
      this.setState({ policyInput: event.target.value });
    }
  };

  handleChange = (type) => (item) => {
    if (type == 'usergroup') {
      let { selectedUserGroups } = this.state;
      if (selectedUserGroups.indexOf(item) === -1) {
        selectedUserGroups = [...selectedUserGroups, item];
      }
      this.setState({ userGInput: '', selectedUserGroups });
    } else {
      let { selectedPolicies } = this.state;
      if (selectedPolicies.length < 1) {
        if (selectedPolicies.indexOf(item) === -1) {
          selectedPolicies = [...selectedPolicies, item];
        }
        this.setState({ policyInput: '', selectedPolicies });
      }
    }
  };

  handleDelete = (item, type) => () => {
    // console.log('************* ', item, type)
    if (type == 'usergroup') {
      this.setState((state) => {
        const selectedUserGroups = [...state.selectedUserGroups];
        selectedUserGroups.splice(selectedUserGroups.indexOf(item), 1);
        return { selectedUserGroups };
      });
    } else {
      this.setState((state) => {
        const selectedPolicies = [...state.selectedPolicies];
        selectedPolicies.splice(selectedPolicies.indexOf(item), 1);
        return { selectedPolicies };
      });
    }
  };

  submitGroupMapRequest = () => {
    this.setState({ progress: true });

    const { selectedUserGroups, selectedPolicies, userName } = this.state;

    const usergroups = selectedUserGroups.map((group) => {
      return group.groupID;
    });

    const policyID = selectedPolicies.map((policy) => {
      return policy.policyID;
    });

    let req;
    let url = '';
    if (this.props.isGroup) {
      const ServicegroupID = this.props.groupMeta.groupID;
      req = {
        ServicegroupID,
        userGroupID: usergroups,
        policyID,
        userName,
        ServicegroupType: 'servicegroup',
      };
      url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/mapgroups/create`;
    } else {
      const { serviceID } = this.props;
      req = {
        ServicegroupID: ID,
        userGroupID: usergroups,
        policyID,
        userName,
        ServicegroupType: 'service',
      };
      url = `${Constants.TRASA_HOSTNAME}/api/v1/services/add/usergroup`;
    }

    this.setState({ error: false, loading: true });

    axios.post(url, req).then((response) => {
      this.setState({ progress: false });
      if (response.data.status === 'success') {
        this.setState({ loading: false });

        window.location.reload();
      } else {
        this.setState({ error: true });
      }
      console.log(response.data);
    });
  };

  handleClickOpen = () => {
    this.setState({ open: true });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  handleNext = () => {
    const { activeStep } = this.state;
    this.setState({
      activeStep: activeStep + 1,
    });
  };

  handleBack = () => {
    const { activeStep } = this.state;
    this.setState({
      activeStep: activeStep - 1,
    });
  };

  handleReset = () => {
    this.setState({
      activeStep: 0,
    });
  };

  getStepContent = (stepIndex) => {
    switch (stepIndex) {
      case 0:
        return (
          <AddUserGroupsHOC
            userGroups={this.props.userGroups}
            inputValue={this.state.inputValue}
            selectedItem={this.state.selectedItem}
            handleKeyDown={this.handleKeyDown}
            handleInputChange={this.handleInputChange}
            handleChange={this.handleChange}
            handleDelete={this.handleDelete}
            getUserGroups={this.getUserGroups}
            getUserGroupsSt={this.state.getUserGroupsSt}
            selectedUserGroups={this.state.selectedUserGroups}
            userGInput={this.state.userGInput}
            handleChangeUsername={this.handleChangeUsername}
            userName={this.state.userName}
          />
        );
      case 1:
        return (
          <AddPolicyHOC
            policies={this.props.policies}
            policyState={this.state.policyState}
            getPolicies={this.getPolicies}
            selectedPolicies={this.state.selectedPolicies}
            policyInput={this.state.policyInput}
            handleKeyDown={this.handleKeyDown}
            handleInputChange={this.handleInputChange}
            handleChange={this.handleChange}
            handleDelete={this.handleDelete}
          />
        );
      default:
        return (
          <ReviewAccessHOC
            selectedPolicies={this.state.selectedPolicies}
            selectedUserGroups={this.state.selectedUserGroups}
            submitGroupMapRequest={this.submitGroupMapRequest}
          />
        );
    }
  };

  render() {
    const { classes } = this.props;
    const steps = getSteps();
    const { progress, loading, error, activeStep } = this.state;
    return (
      <div>
        <Dialog
          fullWidth
          maxWidth="lg"
          open={this.props.open}
          onClose={this.props.handleClose}
          aria-labelledby="form-dialog-title"
        >
          <DialogContent>
            <DialogContentText>
              <Typography variant="h2"> Assign UserGroup to {this.props.serviceName}</Typography>
              <Typography variant="h4">
                Users wont be able to access protected hosts unless assigned here.
              </Typography>
            </DialogContentText>

            <Divider light />
            <Grid container>
              <div className={classes.stepperContent}>{this.getStepContent(activeStep)}</div>
            </Grid>

            <Divider light />
          </DialogContent>

          <div className={classes.stepperRoot}>
            <Grid container>
              <Grid item xs={1} sm={1} md={1}>
                <Button
                  className={classes.stepperButtonLeft}
                  disabled={activeStep === 0}
                  onClick={this.handleBack}
                  variant="contained"
                  color="primary"
                >
                  Back
                </Button>
              </Grid>

              <Grid item xs={10} sm={10} md={10}>
                {this.state.activeStep === steps.length ? (
                  <div>
                    <Typography className={classes.instructions}>
                      All steps completed - you're finished
                    </Typography>
                    <Button onClick={this.handleReset}>Reset</Button>
                  </div>
                ) : (
                  <Stepper activeStep={activeStep} alternativeLabel>
                    {steps.map((label) => {
                      return (
                        <Step key={label}>
                          <StepLabel>{label}</StepLabel>
                        </Step>
                      );
                    })}
                  </Stepper>
                )}
                {progress ? <ProgressHOC /> : ''}
              </Grid>

              <Grid item xs={1} sm={1} md={1}>
                <Button
                  className={classes.stepperButtonRight}
                  variant="contained"
                  color="secondary"
                  onClick={this.handleNext}
                >
                  {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
                </Button>
              </Grid>
            </Grid>
          </div>
        </Dialog>
      </div>
    );
  }
}

ServicegroupUserGroupMap.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ServicegroupUserGroupMap);

/// ////////////////////////////////////////////////////////////////////////////

function renderInput(inputProps) {
  const { InputProps, classes, ref, ...other } = inputProps;

  return (
    <TextField
      InputProps={{
        inputRef: ref,
        classes: {
          root: classes.inputRoot,
          input: classes.inputInput,
        },
        ...InputProps,
      }}
      {...other}
    />
  );
}

function renderSuggestion({ suggestion, index, itemProps, highlightedIndex, selectedItem }) {
  const isHighlighted = highlightedIndex === index;
  const isSelected = (selectedItem || '').indexOf(suggestion.groupID) > -1;

  return (
    <MenuItem
      {...itemProps}
      key={suggestion.groupID}
      selected={isHighlighted}
      component="div"
      style={{
        fontWeight: isSelected ? 500 : 400,
      }}
    >
      {suggestion.groupName}
    </MenuItem>
  );
}
renderSuggestion.propTypes = {
  highlightedIndex: PropTypes.number,
  index: PropTypes.number,
  itemProps: PropTypes.object,
  selectedItem: PropTypes.string,
  suggestion: PropTypes.shape({ label: PropTypes.string }).isRequired,
};

function getSuggestions(value, ServicesThatCanBeAdded) {
  const inputValue = deburr(value.trim()).toLowerCase();
  const inputLength = inputValue.length;
  let count = 0;

  ServicesThatCanBeAdded = ServicesThatCanBeAdded || [];

  return inputLength === 0
    ? ServicesThatCanBeAdded
    : ServicesThatCanBeAdded.filter((suggestion) => {
        const keep =
          count < 5 && suggestion.groupName.slice(0, inputLength).toLowerCase() === inputValue;

        if (keep) {
          count += 1;
        }

        return keep;
      });
}

class AddUserGroups extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      error: false,
      submit: false,
      checkedA: true,
    };
  }

  render() {
    const { classes } = this.props;
    const { userGInput, selectedUserGroups } = this.props;
    return (
      <div className={classes.root}>
        <Downshift
          id="downshift-multiple"
          inputValue={userGInput}
          onChange={this.props.handleChange('usergroup')}
          selectedItem={selectedUserGroups}
        >
          {({
            getInputProps,
            getItemProps,
            isOpen,
            inputValue: userGInput,
            selectedItem: selectedUserGroups,
            highlightedIndex,
          }) => (
            <div className={classes.container}>
              {renderInput({
                fullWidth: true,
                classes,
                InputProps: getInputProps({
                  startAdornment: selectedUserGroups.map((item) => (
                    <Chip
                      key={item.groupID}
                      tabIndex={-1}
                      label={item.groupName}
                      className={classes.chip}
                      onDelete={this.props.handleDelete(item.groupID, 'usergroup')}
                    />
                  )),
                  onChange: this.props.handleInputChange('usergroup'),
                  onKeyDown: this.props.handleKeyDown('usergroup'),
                  placeholder: 'Select multiple user groups',
                }),
                label: 'Label',
              })}
              {isOpen ? (
                <Paper className={classes.paper} square>
                  <Grid container spacing={2}>
                    <Grid item>
                      {getSuggestions(userGInput, this.props.userGroups).map((suggestion, index) =>
                        renderSuggestion({
                          suggestion,
                          index,
                          itemProps: getItemProps({ item: suggestion }),
                          highlightedIndex,
                          selectedItem: selectedUserGroups,
                        }),
                      )}
                    </Grid>
                  </Grid>
                </Paper>
              ) : null}
            </div>
          )}
        </Downshift>
        <Divider />
        <br />

        {this.props.handleChangeUsername ? (
          <Grid container spacing={2}>
            <Grid item xs={2}>
              UserName:
            </Grid>
            <Grid item xs={8}>
              <TextField value={this.props.userName} onChange={this.props.handleChangeUsername} />
            </Grid>
          </Grid>
        ) : (
          ''
        )}
      </div>
    );
  }
}

AddUserGroups.propTypes = {
  classes: PropTypes.object.isRequired,
};

export const AddUserGroupsHOC = withStyles(styles)(AddUserGroups);

function renderInputP(inputProps) {
  const { InputProps, classes, ref, ...other } = inputProps;

  return (
    <TextField
      InputProps={{
        inputRef: ref,
        classes: {
          root: classes.inputRoot,
          input: classes.inputInput,
        },
        ...InputProps,
      }}
      {...other}
    />
  );
}

function renderSuggestionP({ suggestion, index, itemProps, highlightedIndex, selectedItem }) {
  const isHighlighted = highlightedIndex === index;
  const isSelected = (selectedItem || '').indexOf(suggestion.groupID) > -1;

  return (
    <MenuItem
      {...itemProps}
      key={suggestion.policyID}
      selected={isHighlighted}
      component="div"
      style={{
        fontWeight: isSelected ? 500 : 400,
      }}
    >
      {suggestion.policyName}
    </MenuItem>
  );
}
renderSuggestionP.propTypes = {
  highlightedIndex: PropTypes.number,
  index: PropTypes.number,
  itemProps: PropTypes.object,
  selectedItem: PropTypes.string,
  suggestion: PropTypes.shape({ label: PropTypes.string }).isRequired,
};

function getSuggestionsP(value, ServicesThatCanBeAdded) {
  const inputValue = deburr(value.trim()).toLowerCase();
  const inputLength = inputValue.length;
  let count = 0;

  return inputLength === 0
    ? ServicesThatCanBeAdded
    : ServicesThatCanBeAdded.filter((suggestion) => {
        const keep =
          count < 5 && suggestion.policyName.slice(0, inputLength).toLowerCase() === inputValue;

        if (keep) {
          count += 1;
        }

        return keep;
      });
}

class AddPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      error: false,
      submit: false,
      inputValue: '',
      selectedItem: this.props.selectedUserGroups,
      selectedItems: [],
    };
  }

  handleKeyDown = (event) => {
    const { inputValue, selectedItem } = this.state;
    if (selectedItem.length && !inputValue.length && event.key === 'Backspace') {
      this.setState({
        selectedItem: selectedItem.slice(0, selectedItem.length - 1),
      });
    }
  };

  handleInputChange = (event) => {
    this.setState({ inputValue: event.target.value });
  };

  handleChange = (item) => {
    let { selectedItem } = this.state;

    if (selectedItem.indexOf(item) === -1) {
      selectedItem = [...selectedItem, item];
    }

    this.setState({ inputValue: '', selectedItem });
  };

  handleDelete = (item) => () => {
    this.setState((state) => {
      const selectedItem = [...state.selectedItem];
      selectedItem.splice(selectedItem.indexOf(item), 1);
      return { selectedItem };
    });
  };

  render() {
    const { classes } = this.props;
    //  const { loading, error,inputValue, selectedItem  } = this.state;
    const { policyInput, selectedPolicies } = this.props;
    return (
      <div className={classes.root}>
        <Downshift
          id="downshift-multiple"
          inputValue={policyInput}
          onChange={this.props.handleChange('policy')}
          selectedItem={selectedPolicies}
        >
          {({
            getInputProps,
            getItemProps,
            isOpen,
            inputValue: policyInput,
            selectedItem: selectedPolicies,
            highlightedIndex,
          }) => (
            <div className={classes.container}>
              {renderInputP({
                fullWidth: true,
                classes,
                InputProps: getInputProps({
                  startAdornment: selectedPolicies.map((item) => (
                    <Chip
                      key={item.policyID}
                      tabIndex={-1}
                      label={item.policyName}
                      className={classes.chip}
                      onDelete={this.props.handleDelete(item.policyID)}
                    />
                  )),
                  onChange: this.props.handleInputChange('policy'),
                  onKeyDown: this.props.handleKeyDown('policy'),
                  placeholder: 'Select policy',
                }),
                label: 'Label',
              })}
              {isOpen ? (
                <Paper className={classes.paper} square>
                  <Grid container spacing={2}>
                    <Grid item>
                      {getSuggestionsP(policyInput, this.props.policies).map((suggestion, index) =>
                        renderSuggestionP({
                          suggestion,
                          index,
                          itemProps: getItemProps({ item: suggestion }),
                          highlightedIndex,
                          selectedItem: selectedPolicies,
                        }),
                      )}
                    </Grid>
                  </Grid>
                </Paper>
              ) : null}
            </div>
          )}
        </Downshift>
        {/* {this.props.policyState? this.props.getPolicies(this.state.selectedItem): ''} */}
      </div>
    );
  }
}

AddPolicy.propTypes = {
  classes: PropTypes.object.isRequired,
};

export const AddPolicyHOC = withStyles(styles)(AddPolicy);
/// ////////////////////////////////////////////////////////////////////////////////

/// //////////////////////////////////////////////////////////////////////////////////////////////

class ReviewAccess extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userData: [],
    };
  }

  render() {
    const { classes } = this.props;
    // const { userData } = this.state;
    return (
      <div>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <InputLabel htmlFor="select-multiple-users-chip">Selected groups</InputLabel>
            <br />
            {this.props.selectedUserGroups.map((group) => (
              <Chip
                key={group.groupID}
                label={group.groupName}
                value={group.groupName}
                avatar={
                  <Avatar>
                    <AccountCircleIcon />
                  </Avatar>
                }
              />
            ))}
          </Grid>
          <br />

          <Grid item xs={12}>
            <Divider />
            <InputLabel htmlFor="select-multiple-users-chip">Selected policies</InputLabel>
            <br />
            {this.props.selectedPolicies.map((policy) => (
              <Chip
                key={policy.policyID}
                label={policy.policyName}
                value={policy.policyName}
                avatar={<PolicyIcon className={classNames(classes.leftIcon, classes.iconSmall)} />}
              />
            ))}

            <br />
          </Grid>

          <Grid container spacing={2} alignItems="flex-end" direction="row" justify="center">
            <Button
              className={classes.button}
              variant="contained"
              color="secondary"
              onClick={this.props.submitGroupMapRequest}
            >
              Submit
              <Send className={classes.rightIcon} />
            </Button>
          </Grid>
        </Grid>
      </div>
    );
  }
}

ReviewAccess.propTypes = {
  classes: PropTypes.object.isRequired,
};

const ReviewAccessHOC = withStyles(styles)(ReviewAccess);

function Progress(props) {
  const { classes } = props;
  return (
    <div className={classes.root}>
      <LinearProgress />
    </div>
  );
}

Progress.propTypes = {
  classes: PropTypes.object.isRequired,
};

const ProgressHOC = withStyles(styles)(Progress);

class AppUserUsername extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userData: [],
    };
  }

  permissionVal = (permission) => {
    // const { permission } = this.props.PermissionsAssigned
    const val = JSON.stringify({ permission });
    return val;
  };

  isInvalidUsername = (k) => {
    // if Appusers already contain the username, its invalid
    // if(!this.props.selectedUsers[k]) return true
    console.log(this.props.addedUsers[k].username);
    // return true
    return (
      this.props.AppUsers.map((u) => u.username).indexOf(this.props.addedUsers[k].username) >= 0
    );
  };

  // <ReviewAccessHOC Users={this.state.addedUsers} Permissions={this.state.permissions} Request={this.addUserReq}/>;
  render() {
    const { classes } = this.props;
    // const { userData } = this.state;
    console.log(this.props.AppUsers.map((u) => u.username));

    return (
      <div>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Paper className={classes.paper}>
              <div className={classes.manageUser}>
                <b> Selected User(s)</b>
                <br />
                <br />
                {this.props.Users.map((user, k) => {
                  // let isInvalid=this.props.AppUsers.map(u=>u.username).indexOf(user.userName)>-1 //|| (()=>{let temp=this.props.Users.map(u=>u.userName);return temp.indexOf(user.userName)!=temp.lastIndexOf(user.userName)})()
                  return (
                    <div>
                      <Grid container spacing={2} justify="space-between">
                        <Grid item xs={6}>
                          {`${k + 1}. ${user.email}`}
                        </Grid>
                        <Grid item xs={3}>
                          <TextField
                            value={user.userName}
                            onChange={(e) => {
                              this.props.handleChangeUsername(e.target.value, k);
                            }}
                            error={user.isInvalid}
                            label={user.isInvalid ? 'Username Taken' : 'Username'}
                          />
                        </Grid>
                        <Grid item xs={3}>
                          <FormControlLabel
                            control={
                              <Switch
                                checked={user.is2FAEnabled}
                                onChange={(e) => {
                                  this.props.handle2FAChange(e.target.checked, k);
                                }}
                                // value={user.is2FAEnabled}
                                color="secondary"
                              />
                            }
                            label="Enable 2FA"
                          />
                        </Grid>
                      </Grid>
                      <Divider light />
                    </div>
                  );
                })}
              </div>
            </Paper>
          </Grid>
        </Grid>
      </div>
    );
  }
}

AppUserUsername.propTypes = {
  classes: PropTypes.object.isRequired,
};

const AppUserUsernameHOC = withStyles(styles)(AppUserUsername);
