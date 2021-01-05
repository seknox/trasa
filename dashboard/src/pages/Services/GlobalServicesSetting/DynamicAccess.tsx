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
import DynamicAccessRulesTable from "./DynamicAccessRulesTable";
import {Option} from "react-multi-select-component/dist/lib/interfaces";

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
  const [allRules, setAllRules] = React.useState<any>([]);


  const handleGroupsSelect = (groups: any) => {

    if (groups.length > 1) {
      const last = groups[groups.length - 1];
      groups = [last];
    }

    setState({ ...state, userGroups: groups });
  };

  const handlePolicySelect = (pol: any) => {
    if (pol.length > 1) {
      const last = pol[pol.length - 1];
      pol = [last];
    }
    setState({ ...state, policy: pol });
  };



  React.useEffect(() => {
    axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/dynamic`).then((r) => {

      const data = r?.data?.data?.[0];
      const dataArr = data.map(function (n: any) {

        return [
          n.groupName,
          n.policyName,
          n.ruleID,
        ];
      });

      setAllRules(dataArr);
    });
    axios.get(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/dynamic/usergroups`).then((r) => {
      const groups = r.data?.data?.[0].map((g: any) => ({
        label: g,
        id: g,
        value: g,
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
      setState({
        ...state,
        status: props.settings?.status,
      });

  }, [props.settings]);

  const [reqStatus, setReqStatus] = useState(false);
  const classes = useStyles();

  function changeSetting(event: any) {
    setState({ ...state, [event.target.name]: event.target.checked });

    const data = {
      status: event.target.checked ,
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

  function addRule(event:any) {
    const data = {
      groupName: state.userGroups?.[0].id,
      policyID: state.policy?.[0]?.id,
    }
    axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/dynamic/create`,data).then((r) => {
      if(r.data.status=="success"){
        const data = r?.data?.data?.[0];
        const dataArr = data.map(function (n: any) {

          return [
            n.groupName,
            n.policyName,
            n.ruleID,
          ];
        });
        setAllRules(dataArr)
      }
    })
    }

    function deleteRule(ruleID:string) {
    const data = {
      ruleID:ruleID
    }

      axios.post(`${Constants.TRASA_HOSTNAME}/api/v1/accessmap/dynamic/delete`,data).then((r) => {
        if(r.data.status=="success"){
          const data = r?.data?.data?.[0];
          const dataArr = data.map(function (n: any) {

            return [
              n.groupName,
              n.policyName,
              n.ruleID,
            ];
          });
          setAllRules(dataArr)
        }
      })

    }

  const valueRenderer=(label: string)=>(selected:Option[], options:Option[]) =>{
    if(selected && selected.length==1){
      return selected[0].label
    }else {
      return label
    }
  }

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
                              onChange={changeSetting}
                              name="status"
                            />
                          }
                          label={state.status ? <div>enabled </div> : <div>disabled </div>}
                        />
                      </FormGroup>
                    </FormControl>
                  </Grid>
                </Grid>


                <Grid container direction={"row"}>

                    <Grid item xs={3}>
                      <MultiSelect
                          labelledBy="Group"
                          options={allGroups}
                          value={state.userGroups}
                          valueRenderer={valueRenderer("Select Group...")}
                          onChange={handleGroupsSelect}
                      />
                    </Grid>


                    <Grid item xs={3}>
                      <MultiSelect
                          labelledBy="Policy"
                          options={allPolicies}
                          value={state.policy}
                          valueRenderer={valueRenderer("Select Policy...")}
                          onChange={handlePolicySelect}
                      />
                    </Grid>

                  <Grid container xs={3} >
                    <Button onClick={addRule} variant="contained" color="secondary">
                      Add Rule
                    </Button>
                  </Grid>

                </Grid>


                <DynamicAccessRulesTable
                    data={allRules}
                    deleteRule={deleteRule}
                />
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
