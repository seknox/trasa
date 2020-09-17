import axios from 'axios';
import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Constants from '../../Constants';
import Layout from '../../Layout/DashboardBase';
import Headers from '../../Layout/HeaderAndContent';
import GroupPage from './Group';
import UserGroups from './Group/GroupTable';
import UserPage from './User';
import UsersList from './User/userListTable';

function AllUsers() {
  const [allUsers, setallUsers] = React.useState([]);
  const [groups, setgroups] = React.useState([]);

  function getAllUsers() {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/all`)
      .then((response) => {
        let dataArr = [];
        const data = response.data.data[0];
        dataArr = data.map(function (n: any) {
          const date = new Date(n.CreatedAt * 1000);
          return [
            n.email.toString(),
            n.firstName,
            n.lastName,
            n.userName.toString(),
            n.userRole.toString(),
            date.toDateString(),
            n.idpName,
            n.status,
            n.ID,
          ];
        });

        setallUsers(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  }

  function getUserGroups() {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/user`)
      .then((response) => {
        //  /   this.setState({groups: response.data})
        let dataArr = [];
        const data = response.data.data[0];

        dataArr = data.map(function (n: any) {
          // console.log('************ ', n)
          const cdate = new Date(n.createdAt * 1000).toDateString();
          const udate = new Date(n.updatedAt * 1000).toDateString();
          return [
            n.groupName,
            n.memberCount,
            n.status ? 'Active' : 'Disabled',
            cdate,
            udate,
            n.groupID,
          ];
        });
        setgroups(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  }

  React.useEffect(() => {
    getAllUsers();
    getUserGroups();
  }, []);

  return (
    <div>
      <Layout>
        <Headers
          pageName={[{ name: 'Users', route: '/users' }]}
          tabHeaders={['All Users', 'User Groups']}
          Components={[
            <UsersList allUsers={allUsers} groupID="" isGroup={false} deleteUsers={() => {}} />,
            <UserGroups allUsers={allUsers} groups={groups} />,
          ]}
        />
      </Layout>
    </div>
  );
}

const UserView = () => {
  return (
    <Switch>
      <Route exact path="/users" component={AllUsers} />
      {/* <Route exact ='/users/allusers' component={AllUsers} /> */}
      <Route path="/users/user/:userid" component={UserPage} />
      {/* <Route exact path="/users/groups" component={AllUsers} /> */}
      <Route path="/users/groups/group/:groupid" component={GroupPage} />
      {/* <Route path='/users/:groupid/user/:userid' component={AppUserSetting} /> */}
    </Switch>
  );
};

export default withRouter(UserView);
