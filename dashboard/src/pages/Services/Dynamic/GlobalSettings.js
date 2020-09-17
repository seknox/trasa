import { withStyles } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import ProgressHOC from '../../../../components/utils/Progressbar';
import Constants from '../../../../Constants';
import { AddPolicyHOC as AddPolicy, AddUserGroupsHOC as AddUserGroup } from './oldComponents';

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    marginTop: '5%',
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
  textFieldInput: {
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    fontSize: 16,
    padding: '10px 12px',
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
  successText: {
    fontSize: 15,
    color: 'green',
  },
  errorText: {
    fontSize: 15,
    color: 'red',
  },
  settingHeader: {
    color: 'black',
    fontSize: '12px',
    fontFamily: 'Open Sans, Rajdhani',
  },
});

class GlobalSettings extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      enabled: false,
      settings: {},
      selectedUserGroups: [],
      selectedItem: {},
      userGroups: [],
      userGInput: '',
      inputValue: '',
      getUserGroupsSt: '',
      selectedPolicies: [],
      policyInput: '',
      policies: [],
    };
  }

  handleChange = (name) => (event) => {
    const val = event.target.checked;
    const temp = this.state.settings;
    temp[name] = val;
    this.setState({ settings: temp });
  };

  handleSubmit = (event) => {
    this.setState({ progress: true });
    const data = this.state.settings;
    data.userGroups = this.state.selectedUserGroups;
    data.policyID =
      this.state.selectedPolicies.length == 1 ? this.state.selectedPolicies[0].policyID : null;

    axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/services/globalsettings`, data);
  };

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

  handleChangeG = (type) => (item) => {
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

  componentDidMount() {

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/services/globalsettings`)
      .then((response) => {
        //   this.setState({progress: false,loading: false})
        console.log(response.data.data);
        if (response.data.status === 'success' && response.data.data) {
          const resp = response.data.data[0];
          const groupObj = resp.userGroups || []; // [{'groupName': resp.userGroupName, 'groupID': resp.userGroupID }]
          const policyObj = [{ policyName: resp.policyName, policyID: resp.policyID }];
          this.setState({
            settings: resp,
            selectedUserGroups: groupObj,
            selectedPolicies: policyObj,
          });
          console.log(policyObj);
        }
      });

    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/all`)
      .then((response) => {
        // this.setState({allUsers: response.data})
        const resp = response.data.data[0];

        this.setState({ policies: resp });
      })
      .catch((error) => {
        if (error.response.status === 403) {
          window.location.href = '/login';
        }
      });

    const url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/Services/usergroupstoadd`;
    axios
      .get(url)
      .then((response) => {
        const resp = response.data.data[0];
        this.setState({ userGroups: resp });
      })
      .catch((error) => {
        if (error.response.status === 403) {
          window.location.href = '/login';
        }
      });
  }

  render() {
    const { classes } = this.props;
    const { progress, loading, error } = this.state;

    return (
      <div>
        <div className={classes.root}>
          <Grid container spacing={2} alignItems="center" direction="row" justify="flex-start">
            <Grid item xs={12}>
              {loading ? (
                <Typography variant="headline" component="h2" className={classes.successText}>
                  successfully updated  service . reload the page <br />
                </Typography>
              ) : (
                ''
              )}

              {error ? (
                <Typography variant="headline" component="h2" className={classes.errorText}>
                  There were errors validating your request. Please verify your information <br />
                </Typography>
              ) : (
                ''
              )}
              <br />

              <form onSubmit={this.handleSubmit}>
                <Grid container spacing={2}>
                  <Grid item xs={5} sm={5} md={5}>
                    <div className={classes.settingHeader}>Dynamic  service  </div>
                  </Grid>
                  <Grid item xs={7} sm={7} md={7}>
                    <FormControl fullWidth>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Switch
                              enabled={this.state.settings.status}
                              checked={!!this.state.settings.status}
                              onChange={this.handleChange('status')}
                              name="status"
                              defaultValue={this.state.settings.status}
                              value="status"
                              color="primary"
                            />
                          }
                          label={
                            this.state.settings.status ? (
                              <div className={classes.settingSHeader}>enabled </div>
                            ) : (
                              <div className={classes.settingSHeader}>disabled </div>
                            )
                          }
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>

                <Grid container spacing={2}>
                  <Grid item xs={5} sm={5} md={5}>
                    <div className={classes.settingHeader}>Dynamic  service  Video Log</div>
                  </Grid>
                  <Grid item xs={7} sm={7} md={7}>
                    <FormControl fullWidth>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={!!this.state.settings.videoRecord}
                              enabled={this.state.enabled}
                              onChange={this.handleChange('videoRecord')}
                              name="videoRecord"
                              defaultValue={this.state.settings.videoRecord}
                              value="videoRecord"
                              color="primary"
                            />
                          }
                          label={
                            this.state.settings.videoRecord ? (
                              <div className={classes.settingSHeader}>enabled </div>
                            ) : (
                              <div className={classes.settingSHeader}>disabled </div>
                            )
                          }
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>

                <Grid container spacing={2}>
                  <Grid item xs={5} sm={5} md={5}>
                    <div className={classes.settingHeader}>
                      Create Service profile dynamically on access{' '}
                    </div>
                  </Grid>
                  <Grid item xs={7} sm={7} md={7}>
                    <FormControl fullWidth>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={!!this.state.settings.autocreateApp}
                              enabled={this.state.enabled}
                              onChange={this.handleChange('autocreateApp')}
                              name="autocreateApp"
                              defaultValue={this.state.settings.autocreateApp}
                              value="autocreateApp"
                              color="primary"
                            />
                          }
                          label={
                            this.state.settings.autocreateApp ? (
                              <div className={classes.settingSHeader}>enabled </div>
                            ) : (
                              <div className={classes.settingSHeader}>disabled </div>
                            )
                          }
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>

                <Grid item xs={6}>
                  <AddUserGroup
                    userGroups={this.state.userGroups}
                    inputValue={this.state.inputValue}
                    selectedItem={this.state.selectedItem}
                    handleKeyDown={this.handleKeyDown}
                    handleInputChange={this.handleInputChange}
                    handleChange={this.handleChangeG}
                    handleDelete={this.handleDelete}
                    getUserGroups={this.getUserGroups}
                    getUserGroupsSt={this.state.getUserGroupsSt}
                    selectedUserGroups={this.state.selectedUserGroups}
                    userGInput={this.state.userGInput}
                    userName={this.state.userName}
                  />
                  <AddPolicy
                    policies={this.state.policies}
                    policyState={this.state.policyState}
                    getPolicies={this.getPolicies}
                    selectedPolicies={this.state.selectedPolicies}
                    policyInput={this.state.policyInput}
                    handleKeyDown={this.handleKeyDown}
                    handleInputChange={this.handleInputChange}
                    handleChange={this.handleChangeG}
                    handleDelete={this.handleDelete}
                  />
                  {/* <AddUserGroup inputValue={this.state.inputValue} handleChange={()=>{}} selectedUserGroups={this.state.selectedUserGroups} userGInput={this.state.userGInput} handleInputChange={()=>{}} handleKeyDown={()=>{}} selectedItem={this.state.selectedItem}/> */}
                </Grid>

                <Grid item xs={12}>
                  <div className={classes.root}>
                    <div className={classes.wrapper}>
                      <Button variant="contained" color="primary" type="submit">
                        Submit
                      </Button>
                    </div>
                  </div>
                  <div style={{ color: 'red' }}>{this.state.errorText}</div>
                </Grid>
              </form>
            </Grid>
          </Grid>
          <br />
          {progress ? <ProgressHOC /> : ''}
        </div>
        <Divider light />
      </div>
    );
  }
}

GlobalSettings.propTypes = {
  classes: PropTypes.object.isRequired,
};

// export default App;
const GlobalSettingsHOC = withStyles(styles)(GlobalSettings);
export default GlobalSettingsHOC;
