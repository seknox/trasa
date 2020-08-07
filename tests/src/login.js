require('expect-puppeteer')
const axios=require('axios')
// import { getTotp } from '../utils/totpgen'
const getTotp = require("../utils/totpgen");

const rootUserTfaDevice = require("../mock_data/devices")

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

      page.waitFor(5000)
      await axios({
          method: 'post',
          url: deviceData.data[0].cloudProxyURL + '/api/v1/passmydevicedetail',
          data: {
              deviceId: deviceData.data[0].deviceID,
              fcmToken: 'askdnsanduasdosuajasdguyagsdygsac7sacsyaubchjasb',
              publicKey: '',
              deviceFinger: JSON.stringify(rootUserTfaDevice),
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
