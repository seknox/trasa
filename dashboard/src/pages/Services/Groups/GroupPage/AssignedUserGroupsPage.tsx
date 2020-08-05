import Button from '@material-ui/core/Button';
import Save from '@material-ui/icons/PersonAdd';
// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import Constants from '../../../../Constants';
import NAddUsergroupToServicegroup from '../../AccessMap/NAddUsergroupToServicegroup';
import ServicegroupUsergroupMapTable from './ServicegroupUsergroupAccessMapTable';

export default function AssignedUserGroup(props: any) {
  const [assignGroupDlgState, setassignGroupDlgState] = useState(false);
  const [usergroupToAdd, setusergroupToAdd] = useState([]);

  const [policiesToAdd, setpoliciesToAdd] = useState([]);

  // const [serviceName, setserviceName] = useState('');

  const changeAssignGroupDlgState = () => {
    setassignGroupDlgState(!assignGroupDlgState);
  };

  const fetchUserGroups = () => {
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/accessmap/servicegroup/usergroupstoadd`;
    axios
      .get(url)
      .then((response) => {
        setusergroupToAdd(response.data.data[0]);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const fetchPolicies = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/all`)
      .then((response) => {
        // this.setState({allUsers: response.data})
        const resp = response.data.data[0];
        setpoliciesToAdd(resp);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    fetchUserGroups();
    fetchPolicies();
  }, []);

  return (
    <div>
      <Button
        variant="contained"
        color="secondary"
        size="small"
        onClick={changeAssignGroupDlgState}
      >
        <Save />
        Assign User Group
      </Button>
      <br />
      <br />
      <NAddUsergroupToServicegroup
        open={assignGroupDlgState}
        handleClose={changeAssignGroupDlgState}
        userGroups={usergroupToAdd}
        groupMeta={props.groupMeta}
        policies={policiesToAdd}
        renderFor="assignUsergroupToServicegroup"
        /// // TODO serviceName /////
        serviceName=""
        assignuser={false}
        ID={props.groupID}
      />
      {/* <AddUserGroupToServicegroup open={this.state.assignGroupDlgState} handleClose={this.changeAssignGroupDlgState} usergroupToAdd={this.state.usergroupToAdd} userGroups={this.state.userGroups} policiesToAdd={this.state.policiesToAdd}  policies={this.state.policies} groupMeta={this.props.groupMeta} isGroup={true} serviceName={this.state.serviceName} /> */}
      <ServicegroupUsergroupMapTable groupID={props.groupID} />
    </div>
  );
}
