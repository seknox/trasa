import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import { SetPasswordComponent } from './PasswordSetupReset';
import Constants from '../../Constants';
import axios from 'axios';


configure({ adapter: new Adapter() });

const cPassData = {
    password: 'cH@ng3meNever1#a',
    cpassword: 'cH@ng3meNever1#a',
};


describe('Test change password component', () => {


    test('submit button should be disabled by default', () => {
        const page = mount(<SetPasswordComponent update={false} token={'12345'}  />);


        expect(page.find('super')).toBeTruthy();

        const button = page.find('button').first()
        expect(button.props().disabled).toBe(true);


    });

    test('Should show submit button if password and changepassword matches && zxcvbn score is equal or greater than 2', () => {
        const page = mount(<SetPasswordComponent update={false} token={'12345'}  />);

        const input1 = page.find({ name: 'password' }).find('input')
        input1.props().value = cPassData.password
        input1.props().onChange({target: {name: 'password', value: cPassData.password}})
        console.log('input1 props: ', input1.props())

  
        const input2 = page.find({ name: 'cpassword' }).find('input');
        input2.props().value = cPassData.cpassword
        input2.props().onChange({target: {name: 'cpassword', value: cPassData.cpassword}})
        console.log('input2 props: ', input2.props())

        
        expect(input1.props().value).toEqual(input2.props().value);

        const button = page.update().find({ name: 'submit' }).first()

          console.log('button props: ', button.props())
          expect(button.props().disabled).toBe(false);

    });

   
});
