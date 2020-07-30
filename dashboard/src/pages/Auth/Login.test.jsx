import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import Login from './index';

configure({ adapter: new Adapter() });

const loginData = {
  email: 'sakshyam@seknox.com',
  password: 'test',
  orgID: '',
  intent: '',
  idpName: '',
  userName: '',
};

describe('Testing <Login /> Component...', () => {
  test('Should use props to fill email and password', () => {
    const login = mount(<Login userData={loginData} />);

    const email = login.find({ name: 'email' }).find('input');
    const password = login.find({ name: 'password' }).find('input');
    expect(email.props().value).toEqual('sakshyam@seknox.com');
    expect(password.props().value).toEqual('test');
  });

  test('Should not submit empty fields', () => {});

  test('Should show loader bar untill request is resolved', () => {});

  test('Should hide forget password button for showForgetPass props', () => {});

  test('Should hide email textField for autofillEmail props ', () => {
    const login = mount(<Login userData={loginData} autofillEmail />);

    const email = login.find({ name: 'email' }).find('input');
    expect(email).toEqual({});
  });

  test('Should submit form with username and password', () => {});
});
