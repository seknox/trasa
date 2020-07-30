import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import Index from './index';

configure({ adapter: new Adapter() });

const loginData = {
  email: 'sakshyam@seknox.com',
  password: 'test',
  orgID: '',
  intent: '',
  idpName: '',
};

describe('Testing <IndexLogin /> Component...', () => {
  test('Should show login page by default', () => {
    const login = mount(<Index userData={loginData} />);

    const email = login.find({ name: 'email' }).find('input');
    expect(email.props().value).toEqual('sakshyam@seknox.com');
  });

  test('Should show tfa page if tfaRequired state is true', () => {
    const index = mount(<Index userData={loginData} tfaRequired />);
    expect(index.find({ name: 'email' }).find('input').exists()).toBeFalsy();
  });

  test('Should show enrol device page if enrolDevice state is true', () => {});
});
