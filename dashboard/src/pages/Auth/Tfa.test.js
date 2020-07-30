import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import Tfa from './index';

configure({ adapter: new Adapter() });

const loginData = {
  email: 'sakshyam@seknox.com',
  password: 'test',
  orgID: '',
  intent: '',
  idpName: '',
};

describe('Testing <Tfa /> Component...', () => {
  test('Should prompt for totp dialog', () => {});
  test('should show loader in u2f request', () => {});
  test('Should show loader in totp request', () => {});
});
