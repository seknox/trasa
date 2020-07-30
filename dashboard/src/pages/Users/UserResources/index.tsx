import Grid from '@material-ui/core/Grid';
import React from 'react';
import GroupsAssignedToUser from './GroupsAssignedToUser';
import UsersAppTable from './UsersServicesTable';

export default function UserGroupsAndService(props: any) {
  return (
    <Grid container spacing={2} direction="row" justify="center">
      <Grid item xs={12} sm={12} md={5}>
        <GroupsAssignedToUser userID={props.userID} userGroups={props.userGroups} />
      </Grid>
      <Grid item xs={12} sm={12} md={7}>
        <UsersAppTable AppUsers={props.userAccessMaps} serviceID={props.userAccessMaps.serviceID} />
      </Grid>
    </Grid>
  );
}
