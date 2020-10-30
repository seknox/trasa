
import { configure, mount } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import UserDevices from './index';
import Button from '@material-ui/core/Button';

import axios from 'axios';
import {act} from "react-dom/test-utils";
import {Switch} from "@material-ui/core";
import IconButton from "@material-ui/core/IconButton";
import { DeleteConfirmDialogue } from '../../../utils/Components/Confirms';
import Constants from "../../../Constants";


jest.mock('axios');

configure({ adapter: new Adapter() });


const mockResponseData={"status":"success","reason":"devices fetched.","intent":"SingleUserDevices","data":[{"mobile":[{"userID":"","orgID":"","deviceID":"8196ac93-7314-4a83-aa7e-7695be13b314","machineID":"","deviceType":"mobile","fcmToken":"","publicKey":"","deviceFinger":"","trusted":false,"deviceHygiene":{"deviceInfo":{"deviceName":"Bhargab’s iPhone","deviceVersion":"","machineID":"sa4d5sa5d3s54a3d54as3d45as35d","brand":"Apple","manufacturer":"Apple","deviceModel":"unknown"},"deviceOS":{"osName":"iOS","osVersion":"13.4.1","kernelType":"","kernelVersion":"","readableVersion":"","latestSecurityPatch":"unknown","autoUpdate":true,"pendingUpdates":null,"jailBroken":false,"debugModeEnabled":false,"isEmulator":false},"loginSecurity":{"autologinEnabled":false,"loginMethod":"","passwordLastUpdated":"","tfaConfigured":false,"idleDeviceScreenLockTime":"","idleDeviceScreenLock":false,"remoteLoginEnabled":false},"networkInfo":{"hostname":"unknown","domainControl":false,"domainName":"","interfaceName":"wifi","ipAddress":"192.168.100.183","macAddress":"00:00:00:00:00:00","wirelessNetwork":true,"openWifiConn":false,"networkName":"","networkSecurity":""},"endpointSecurity":{"epsConfigured":false,"epsVendorName":"","epsVersion":"","epsMeta":"","firewallEnabled":false,"firewallPolicy":"","deviceEncryptionEnabled":false,"deviceEncryptionMeta":""},"lastCheckedTime":0},"addedAt":1604038231}],"workstation":[{"userID":"","orgID":"","deviceID":"bb9f780e-8c63-46ad-8821-db3dfe9b721c","machineID":"","deviceType":"workstation","fcmToken":"","publicKey":"","deviceFinger":"","trusted":false,"deviceHygiene":{"deviceInfo":{"deviceName":"Mac OS X","deviceVersion":"10.15.6","machineID":"BHUBUWSBHJWDBJHWDBJWHDB","brand":"Apple","manufacturer":"Apple","deviceModel":""},"deviceOS":{"osName":"Mac OS X","osVersion":"10.15.6","kernelType":"darwin","kernelVersion":"19.6.0","readableVersion":"","latestSecurityPatch":"","autoUpdate":true,"pendingUpdates":null,"jailBroken":false,"debugModeEnabled":false,"isEmulator":false},"loginSecurity":{"autologinEnabled":false,"loginMethod":"","passwordLastUpdated":"Wed Dec 22 10:43:49 +0545 2020","tfaConfigured":false,"idleDeviceScreenLockTime":"","idleDeviceScreenLock":true,"remoteLoginEnabled":true},"networkInfo":{"hostname":"Bhargabs-MacBook-Pro.local","domainControl":false,"domainName":"","interfaceName":"en5,en0,awdl0,llw0,utun0,utun1,utun2","ipAddress":"192.168.100.171","macAddress":"AC:ASH:ASDJSDNHSJD","wirelessNetwork":false,"openWifiConn":false,"networkName":"","networkSecurity":""},"endpointSecurity":{"epsConfigured":true,"epsVendorName":"","epsVersion":"","epsMeta":"","firewallEnabled":false,"firewallPolicy":"","deviceEncryptionEnabled":false,"deviceEncryptionMeta":""},"lastCheckedTime":0},"addedAt":1602219626}],"browser":[],"hToken":[]}]}

const waitForComponentToPaint = async (wrapper) => {
    await act(async () => {
        await new Promise(resolve => setTimeout(resolve, 0));
        wrapper.update();
    });
};

describe('Testing <UserDevices /> Component...', () => {
    test('Should have trust option when rendered for users page', async () => {

        axios.get.mockImplementationOnce(() => Promise.resolve({data:mockResponseData}));

        const userDevices = mount(<UserDevices userID={""} renderFor={"userRoute"}/>);

        await waitForComponentToPaint(userDevices)

        expect(axios.get).toHaveBeenCalledWith(
            `${Constants.TRASA_HOSTNAME}/api/v1/user/devices/all/`,
        );

        const viewDetailButton = await userDevices.find(Button);
        expect(viewDetailButton.length).toBe(2);



        viewDetailButton.at(0).simulate('click')
        // await waitForComponentToPaint(userDevices)

        let tswitch =userDevices.find(Switch)
        expect(tswitch.length).toBe(1);





    });


    test('Should not have trust option when rendered for my page', async () => {

        axios.get.mockImplementationOnce(() => Promise.resolve({data:mockResponseData}));

        const userDevices = mount(<UserDevices userID={""} renderFor={"myRoute"}/>);

        await waitForComponentToPaint(userDevices)

        expect(axios.get).toHaveBeenCalledWith(
            `${Constants.TRASA_HOSTNAME}/api/v1/my/devices`,
        );

        const viewDetailButton = await userDevices.find(Button);
        expect(viewDetailButton.length).toBe(2);
        console.log(viewDetailButton)
        viewDetailButton.at(0).simulate('click')
        // await waitForComponentToPaint(userDevices)

        let tswitch =userDevices.find(Switch)
        expect(tswitch.length).toBe(0);

        let deleteBtn =  userDevices.find(IconButton)
        expect(deleteBtn.length).toBe(1);

        deleteBtn.simulate('click')

        let deleteConfirmationDiag = userDevices.find(DeleteConfirmDialogue)
        // expect(deleteConfirmationDiag.length).toBe(1);
        expect(deleteConfirmationDiag.at(0).prop("open")).toBe(true)

        let confirmBtn=await deleteConfirmationDiag.at(0).find(Button)
        expect(confirmBtn.length).toBe(2);


        //TODO check API url too

        // axios.get.mockImplementationOnce(() => Promise.resolve({data: {status:"success"}}));


        await confirmBtn.at(1).simulate('click')

        // expect(axios.get).toHaveBeenCalledWith(
        //     `${Constants.TRASA_HOSTNAME}/api/v1/my/devices/delete/8196ac93-7314-4a83-aa7e-7695be13b314`,
        // );



    });



});
