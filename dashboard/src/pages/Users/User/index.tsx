// import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Route, RouteComponentProps, Switch, withRouter } from 'react-router-dom';
import Constants from '../../../Constants';
import Layout from '../../../Layout/DashboardBase';
import Headers from '../../../Layout/HeaderAndContent';
import { Usertype } from '../../../types/users';
import AccessStats from '../AccessStats';
import UserDevices from '../UserDevices';
import UserResources from '../UserResources';
import UserOverview from './UserOverview';

type Userpageindexprops = {
  userID: string;
};

function UserPage(props: Userpageindexprops) {
  const [userData, setuserData] = useState<Usertype>({
    ID: '',
    firstName: '',
    middleName: '',
    lastName: '',
    email: '',
    password: '',
    userRole: '',
    userName: '',
    userName2: '',
    cpassword: '',
    status: false,
    CreatedAt: 0,
  });
  const [UsernameForPage, setUsernameForPage] = useState('');
  const [userAccessMaps, setuserAccessMaps] = useState([]);
  const [userDevices, setuserDevices] = useState<any>([]);
  const [userGroups, setuserGroups] = useState([]);

  useEffect(() => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/user/${props.userID}`)
      .then((response) => {
        setuserData(response.data.data[0].user);
        setUsernameForPage(
          `${response.data.data[0].user.firstName} ${response.data.data[0].user.lastName}`,
        );
        setuserAccessMaps(response.data.data[0].userAccessMaps);
        setuserDevices(response.data.data[0].userDevices);
        setuserGroups(response.data.data[0].userGroups);
      })
      .catch((error) => {
        console.log(error);
      });
  }, [props.userID]);

  const getRoute = (staticVal: any, dynamicVal: any) => {
    return staticVal + dynamicVal;
  };

  return (
    <Layout>
      <Headers
        pageName={[
          { name: 'Users', route: '/users' },
          // { name: 'All Users', route: '/users' },
          { name: UsernameForPage, route: getRoute('/users/user/', props.userID) },
        ]}
        tabHeaders={['Account Overview', 'User Devices', 'Groups & Services', 'Access Stats']}
        Components={[
          <UserOverview
            userData={userData}
            userDevices={userDevices}
            userGroups={userGroups}
            userAccessMaps={userAccessMaps}
          />,

          <UserDevices userID={props.userID} />,

          <UserResources
            userID={props.userID}
            userGroups={userGroups}
            userAccessMaps={userAccessMaps}
          />,

          <AccessStats userID={props.userID} />,
        ]}
      />
    </Layout>
  );
}

type TParams = { userid: string };

const UserPageIndex = ({ match }: RouteComponentProps<TParams>) => {
  return (
    <div>
      <UserPage userID={match.params.userid} />
    </div>
  );
};

const UserPageRoute = () => {
  return (
    <Switch>
      <Route exact path="/users/user/:userid" component={UserPageIndex} />
    </Switch>
  );
};

export default withRouter(UserPageRoute);
