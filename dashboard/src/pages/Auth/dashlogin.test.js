import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import Index from './index';
import Constants from '../../Constants';

configure({ adapter: new Adapter() });

const loginData = {
  email: 'tree',
  password: 'changemenever1#',
  orgID: '',
  intent: '',
  idpName: '',
};


describe('Testing Login Flow with different intents (connects to TRASA server)...', () => {
  test('Should redirect to /overview page when tfa returns success', () => {
    const login = mount(<Index userData={loginData} intent='AUTH_REQ_DASH_LOGIN' autofillEmail={false} title="Dashboard Login"
    showForgetPass proxyDomain='trasa.seknox.com' />);

    const email = login.find({ name: 'email' }).find('input');
    expect(email.props().value).toEqual(loginData.email);

    const password = login.find({ name: 'password' }).find('input');
    expect(password.props().value).toEqual(loginData.password);


    
  });

  test('Should initate access proxy flow when intent is AUTH_HTTP_ACCESS_PROXY', () => {
 
  });


});
