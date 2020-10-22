import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import { SetPasswordComponent } from './PasswordSetupReset';


jest.setTimeout(30000);

configure({ adapter: new Adapter() });

const cPassData = {
    password: 'cH@ng3meNever1#a',
    cpassword: 'cH@ng3meNever1#a',
};


describe('Test change password component', () => {


    test('submit button should be enabled by default', () => {
        const page = mount(<SetPasswordComponent update={false} token={'12345'}  />);


        expect(page.find('super')).toBeTruthy();

        const button = page.find('button').first()
        expect(button.props().disabled).toBe(false);


    });






});
