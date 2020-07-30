import React, { useEffect, useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import axios from 'axios';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import Grid from '@material-ui/core/Grid';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Typography from '@material-ui/core/Typography';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import { Checkbox } from '@material-ui/core';
import MultiSelect from 'react-multi-select-component';
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button';
import ProgressHOC from '../../../utils/Components/Progressbar';
import Constants from '../../../Constants';

const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1,
  },
  paper: {
    backgroundColor: '#fdfdfd',
    padding: theme.spacing(2),
    color: theme.palette.text.secondary,
  },
  textField: {
    // marginLeft: 100,
    paddingLeft: theme.spacing(1),
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    color: 'black',
    fontSize: '12px',
    width: 500,
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
  selectCustom: {
    fontSize: 15,
    borderRadius: 4,
    backgroundColor: theme.palette.common.white,
    border: '1px solid #ced4da',
    height: 31,
    // marginTop: 5,
    // padding: '10px 100px',
    width: 100,
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

export default function DynamicAccess(props: any) {
  const [state, setState] = React.useState<any>({
    status: true,
    userGroups: [],
    policy: [],
    admins: [],
  });

  const [allGroups, setAllGroups] = React.useState<any>([]);
  const [allPolicies, setAllPolicies] = React.useState<any>([]);

  const handleChange = (event: any) => {
    setState({ ...state, [event.target.name]: event.target.checked });
  };

  const handleGroupsSelect = (groups: any) => {
    setState({ ...state, userGroups: groups });
  };

  const handlepolicySelect = (groups: any) => {
    if (groups.length > 1) {
      const last = groups[groups.length - 1];
      groups = [last];
    }
    setState({ ...state, policy: groups });
  };

  React.useEffect(() => {
    axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/user`).then((r) => {
      const groups = r.data?.data?.[0].map((g: any) => ({
        label: g.groupName,
        id: g.groupID,
        value: g.groupID,
      }));
      setAllGroups(groups);
    });
    axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/all`).then((r) => {
      const policies = r.data?.data?.[0].map((g: any) => ({
        label: g.policyName,
        id: g.policyID,
        value: g.policyID,
      }));
      setAllPolicies(policies);
    });
  }, []);

  React.useEffect(() => {
    try {
      const settVal = JSON.parse(props.settings?.settingValue);

      const groups = [];
      for (let i = 0; i < allGroups.length; i++) {
        for (let j = 0; j < settVal.userGroups.length; j++) {
          if (allGroups[i].id == settVal.userGroups[j]) {
            groups.push(allGroups[i]);
          }
        }
      }

      let policy = {};
      for (let i = 0; i < allPolicies.length; i++) {
        if (allPolicies[i].id == settVal.policyID) {
          policy = allPolicies[i];
        }
      }

      setState({
        ...state,
        userGroups: groups,
        policy: [policy],
        status: props.settings?.status,
      });
      // setState({ ...state, ['policy']: [policy] });
    } catch (e) {
      console.error(e);
    }
  }, [props.settings, allPolicies, allGroups]);

  const [reqStatus, setReqStatus] = useState(false);
  const classes = useStyles();

  function submitSetting(val: any) {
    const data = {
      userGroups: state.userGroups.map((g: any) => g.id),
      policyID: state.policy?.[0]?.id,
      status: state.status,
    };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/system/settings/dynamicaccess/update`, data)
      .then((response) => {
        setReqStatus(false);
      })
      .catch((error) => {
        console.log(error);
      });
  }
  // TODO
  // useEffect(() => {
  //     setEnableDynamicAccess(props.status);
  // }, [props.status]);

  return (
    <div className={classes.root}>
      {/* <Paper className={classes.paper}>  */}
      {/* <Grid container spacing={2} direction="row"  justify="center"> */}
      <ExpansionPanel>
        <Grid item xs={12} sm={12}>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1a-content"
            id="panel1a-header"
          >
            <Typography component="h4" variant="h3">
              <b>Dynamic Access </b>
            </Typography>
          </ExpansionPanelSummary>
        </Grid>

        <Grid item xs={12}>
          <Typography component="h5" variant="h5">
            NOTE: Dynamic access enables users to access services which are not yet added to TRASA.
          </Typography>
        </Grid>

        <Grid item xs={12} sm={12}>
          <ExpansionPanelDetails>
            <Grid container spacing={2} direction="row" justify="center">
              <Grid item xs={12}>
                <Divider light />
              </Grid>

              <Grid item xs={12} sm={9}>
                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={6}>
                    <Typography variant="h4">Enable Dynamic Access : </Typography>
                  </Grid>
                  <Grid item xs={3}>
                    <FormControl fullWidth>
                      <FormGroup>
                        <FormControlLabel
                          control={
                            <Checkbox
                              checked={state.status}
                              onChange={handleChange}
                              name="status"
                            />
                          }
                          label={state.status ? <div>enabled </div> : <div>disabled </div>}
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={6}>
                    <Typography variant="h4">User groups allowed Dynamic Access : </Typography>
                  </Grid>
                  <Grid item xs={3}>
                    <MultiSelect
                      labelledBy="Select user groups"
                      options={allGroups}
                      value={state.userGroups}
                      onChange={handleGroupsSelect}
                    />
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" justify="center">
                  <Grid item xs={6}>
                    <Typography variant="h4">Policy : </Typography>
                  </Grid>
                  <Grid item xs={3}>
                    <MultiSelect
                      labelledBy="Select policy"
                      options={allPolicies}
                      value={state.policy}
                      onChange={handlepolicySelect}
                    />
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="center" justify="center">
                  <Button onClick={submitSetting} variant="contained" color="secondary">
                    Submit
                  </Button>
                </Grid>
              </Grid>
              {/* </Grid> */}
            </Grid>
          </ExpansionPanelDetails>
        </Grid>
        {reqStatus ? (
          <div>
            <ProgressHOC /> <br />
          </div>
        ) : null}
        <ExpansionPanelActions />
      </ExpansionPanel>
      {/* </Grid> */}
      {/* </Paper> */}
      <br /> <br /> <br />
    </div>
  );
}
