import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import MultiSelect from 'react-multi-select-component';
import Constants from '../../../Constants';
import ProgressBar from '../../../utils/Components/Progressbar';
import { GroupType } from '../../../types/groups';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    height: 250,
  },
  container: {
    flexGrow: 1,
    position: 'relative',
  },
  paper: {
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
  submitRight: {
    marginLeft: '35%',
  },
  dialogueRoot: {
    minHeight: 500,
    // overflow: 'auto',
  },
  selectZindex: {
    // overflow: 'visible',
    // zIndex: 99999,
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
    height: 25,
    fontFamily: 'Open Sans, Rajdhani',
    padding: 'theme.spacing(2)',
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
    fontSize: '15px',
    height: 23,
    fontFamily: 'Open Sans, Rajdhani',
    padding: theme.spacing(1),
    transition: theme.transitions.create(['border-color', 'box-shadow']),
    '&:focus': {
      borderColor: '#80bdff',
      boxShadow: '0 0 0 0.2rem rgba(0,123,255,.25)',
    },
  },
}));

type AssignUsergroupToServicegroupProps = {
  renderFor: 'assignUserToApp' | 'assignUsergroupToServicegroup' | 'assignUserGroup';
  ID: string;
  open: boolean;
  handleClose: () => void;
  serviceName: string;
  groupMeta: GroupType;
  assignuser: boolean;
  userGroups: GroupType[];
  policies: any;
};
export default function AssignUsergroupToServicegroup(props: AssignUsergroupToServicegroupProps) {
  const classes = useStyles();

  const [loader, setLoader] = useState(false);
  const [selectedUsergroups, setSelectedUsergroups] = useState([]);
  const [selectedPolicies, setSelectedPolicies] = useState([]);
  const [privilege, setPrivilege] = useState('');

  const changePrivilege = (e: any) => {
    setPrivilege(e.target.value);
  };

  const selectPolicy = (policies: any) => {
    if (policies.length > 1) {
      const last = policies[policies.length - 1];
      policies = [last];
    }
    setSelectedPolicies(policies);
  };

  const submitGroupMapRequest = () => {
    setLoader(true);

    const usergroups = selectedUsergroups.map((group: any) => {
      if (props.renderFor === 'assignUserToApp') {
        return group.id;
      }
      return group.value;
    });

    const policyID = selectedPolicies.map((policy: any) => {
      return policy.value;
    });

    let req;
    let url = '';
    if (props.renderFor === 'assignUsergroupToServicegroup') {
      const serviceGroupID = props.groupMeta.groupID;
      req = {
        serviceGroupID,
        userGroupID: usergroups,
        policyID,
        privilege,
        mapType: 'servicegroup',
      };
      url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroup/create`;
    } else if (props.renderFor === 'assignUserGroup') {
      const { ID } = props;
      req = {
        serviceGroupID: ID,
        userGroupID: usergroups,
        policyID,
        privilege,
        mapType: 'service',
      };
      url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroup/create`;
    } else {
      const { ID } = props;
      req = { serviceID: ID, users: usergroups, privilege, policyID };
      url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/service/user/create`;
    }

    axios.post(url, req).then((response) => {
      setLoader(false);
      if (response.data.status === 'success') {
        window.location.reload();
      } else {
        console.log('err:: ');
      }
    });
  };

  return (
    <Dialog
      fullWidth
      maxWidth="lg"
      open={props.open}
      onClose={props.handleClose}
      aria-labelledby="form-dialog-title"
      classes={{
        paper: classes.dialogueRoot,
      }}
    >
      <DialogTitle>
        {/* TODO - serviceName not available here... */}
        <Typography component="span" variant="h2">
          {' '}
          Assign user group to {props.serviceName}
        </Typography>
      </DialogTitle>

      <DialogContent>
        <Divider light />

        <Grid container spacing={2} direction="row">
          <Grid item xs={5}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="h3">
                  {' '}
                  {props.assignuser ? 'Select Users' : 'Select Groups'}{' '}
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <AddUsergroup
                  usergroupsThatCanBeAdded={props.userGroups}
                  selectedUsergroups={selectedUsergroups}
                  setSelectedUsergroups={setSelectedUsergroups}
                  renderFor={props.renderFor}
                />
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={4}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="h3"> Select Policy</Typography>
              </Grid>
              <Grid item xs={12}>
                <AddPolicy
                  policies={props.policies}
                  selectedPolicies={selectedPolicies}
                  setSelectedPolicies={selectPolicy}
                />
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={3}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="h3"> Assign Privilege </Typography>
              </Grid>
              <Grid item xs={12}>
                <AddPrivilege privilege={privilege} changePrivilege={changePrivilege} />
              </Grid>
            </Grid>
          </Grid>

          <Grid container spacing={2} alignItems="flex-end" direction="row" justify="center">
            <Button variant="contained" color="secondary" onClick={submitGroupMapRequest}>
              Submit
            </Button>
          </Grid>
          <Grid item xs>
            {loader ? <ProgressBar /> : null}
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button variant="outlined" color="primary" onClick={props.handleClose}>
          Close
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export const AddUsergroup = (props: any) => {
  const [listVals, setListVals] = useState([{ label: '', value: '', id: '' }]);

  useEffect(() => {
    const vals = props.usergroupsThatCanBeAdded.map((v: any) => {
      if (props.renderFor === 'assignUserToApp') {
        return { label: `${v.firstName} ${v.lastName} (${v.email})`, value: v.email, id: v.ID };
      }
      return { label: `${v.groupName} `, value: v.groupID };
    });

    setListVals(vals);
  }, [props.usergroupsThatCanBeAdded, props.renderFor]);

  return (
    <div>
      <MultiSelect
        options={listVals}
        value={props.selectedUsergroups}
        onChange={props.setSelectedUsergroups}
        labelledBy="Select user Groups"
      />

      <br />
      <br />
    </div>
  );
};

export const AddPolicy = (props: any) => {
  const [listVals, setListVals] = useState([{ label: '', value: '', id: '' }]);

  useEffect(() => {
    const vals = props.policies.map((v: any) => {
      return { label: `${v.policyName} `, value: v.policyID };
    });

    setListVals(vals);
  }, [props.policies]);

  return (
    <div>
      <MultiSelect
        options={listVals}
        value={props.selectedPolicies}
        onChange={props.setSelectedPolicies}
        labelledBy="Select user Groups"
      />

      <br />
      <br />
    </div>
  );
};

function AddPrivilege(props: any) {
  const classes = useStyles();

  return (
    <TextField
      fullWidth
      // label="Service name"
      onChange={props.changePrivilege}
      name="privilege"
      value={props.privilege}
      InputProps={{
        disableUnderline: true,
        classes: {
          // root: classes.textFieldRoot,
          input: classes.textFieldInputBig,
        },
      }}
      InputLabelProps={{
        shrink: true,
      }}
    />
  );
}
