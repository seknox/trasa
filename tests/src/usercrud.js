require('expect-puppeteer')
const axios=require('axios')
// import { getTotp } from '../utils/totpgen'
const getTotp = require("../utils/totpgen");


const loginData = {
  email: 'root',
  password: 'changeme',
  orgID: '',
  intent: '',
  idpName: '',
};

const TRASA_HOSTNAME="https://app.trasa"
let TOTP_SSC=""




export const logintests = () => {

  it('should display "Dashboard Login" text on page', async () => {
    await expect(page).toMatch('Dashboard Login')
  })

  it('should fill form', async () => {

   await page.type('#email', loginData.email)
    await page.type('#password', loginData.password);


    await page.click('#submit')
    // await page.screenshot({path: 'shot.png'});

    let loginResp = await page.waitForResponse(TRASA_HOSTNAME+'/idp/login');

    expect(loginResp.status()).toBe(200)

      await expect(page).toMatch('Get TRASA app for your mobile device.')


      const deviceData = await loginResp.json()
      TOTP_SSC=deviceData.data[0].totpSSC

      await axios({
          method: 'post',
          url: deviceData.data[0].cloudProxyURL + '/api/v1/passmydevicedetail',
          data: {
              deviceId: deviceData.data[0].deviceID,
              fcmToken: 'askdnsanduasdosuajasdguyagsdygsac7sacsyaubchjasb',
              publicKey: '',
              deviceFinger: JSON.stringify({"deviceInfo": {"brand": "samsung", "deviceModel": "SM-A305F", "deviceName": "Galaxy A30", "deviceVersion": "", "machineID": "81836ac40d412942", "manufacturer": "samsung"}, "deviceOS": {"autoUpdate": true, "debugModeEnabled": false, "isEmulator": false, "jailBroken": false, "kernelType": "", "kernelVersion": "", "latestSecurityPatch": "2020-06-01", "osName": "Android", "osVersion": "10", "pendingUpdates": null, "readableVersion": ""}, "endpointSecurity": {"deviceEncryptionEnabled": false, "deviceEncryptionMeta": "", "epsConfigured": false, "epsMeta": "", "epsVendorName": "", "epsVersion": "", "firewallEnabled": false, "firewallPolicy": ""}, "lastCheckedTime": 0, "loginSecurity": {"autologinEnabled": false, "idleDeviceScreenLock": false, "idleDeviceScreenLockTime": "", "loginMethod": "", "passwordLastUpdated": "", "remoteLoginEnabled": false, "tfaConfigured": false}, "networkInfo": {"domainControl": false, "domainName": "", "hostname": "SWDI7613", "interfaceName": "wifi", "ipAddress": "192.168.100.125", "macAddress": "FC:AA:B6:57:71:8B", "networkName": "", "networkSecurity": "", "openWifiConn": false, "wirelessNetwork": true}}),
          },
      }).then(r=>{
          expect(r.data.status).toBe("success")
      })
      await page.goto("https://app.trasa/login",{waitUntil:"load"})

      await page.type('#email', loginData.email)
      await page.type('#password', loginData.password);


      await page.click('#submit')

    await expect(page).toMatch('Choose second step verification method')

    loginResp.json().then(data=>{
        expect(data.status).toBe("success")
    })


  await page.click("[id=totpButton]")
  // await page.waitForSelector("[name=totpVal]")

  await page.type("[name=totpVal]",getTotp(TOTP_SSC))
  await page.keyboard.press("Enter")

  let tfaResp = await page.waitForResponse(TRASA_HOSTNAME+'/idp/login/tfa');


  expect(tfaResp.status()).toBe(200)

 //  await page.waitFor(8000)

  // tfaResp.json().then(data=>{
  //   // console.log('status: ', data)
  //     //expect(data.status).toBe("success")
  // })


  //await page.screenshot({path: 'screen.png'});




})
}

// getTotp(TOTP_SSC)
