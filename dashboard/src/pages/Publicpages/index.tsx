import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Activate from './Activate';
// import Enroll2faDevice from './enroll2fadevice';
// import ForgotPassword from './ForgotPassword';
import VerifyToken from './PasswordSetupReset';
// import SignUp from './signup';

const PublicView = () => {
  return (
    <Switch>
      <Route exact path="/woa/verify" component={VerifyToken} />
      {/* <Route path='/woa/enrol/device' component={Enroll2faDevice} /> */}
      {/* <Route path='/woa/forgotpass' component={ForgotPassword} /> */}
      <Route path="/woa/activate" component={Activate} />
      {/* <Route path="/woa/signup" component={SignUp} /> */}
    </Switch>
  );
};

export default withRouter(PublicView);
