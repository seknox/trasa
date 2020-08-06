require('expect-puppeteer')

// import { getTotp } from '../utils/totpgen'
const getTotp = require("../utils/totpgen");


const loginData = {
  email: 'tree',
  password: 'changemenever1#',
  orgID: '',
  intent: '',
  idpName: '',
};

const TRASA_HOSTNAME="https://app.trasa"
const TOTP_SSC="ZJAEXMC2YTJFQIO7"




export const logintests = () => {

  it('should display "Dashboard Login" text on page', async () => {
    await expect(page).toMatch('Dashboard Login')
  })

  it('should fill form', async () => {
  
   await page.type('#email', loginData.email)
    await page.type('#password', loginData.password);

 
    await page.click('#submit')
    await page.screenshot({path: 'shot.png'});

    let loginResp = await page.waitForResponse('https://app.trasa/idp/login');

    expect(loginResp.status()).toBe(200)
    await expect(page).toMatch('Choose second step verification method')

    loginResp.json().then(data=>{
        expect(data.status).toBe("success")
    })


  await page.click("[id=totpButton]")
  // await page.waitForSelector("[name=totpVal]")

  await page.type("[name=totpVal]",getTotp(TOTP_SSC))
  await page.keyboard.press("Enter")

  let tfaResp = await page.waitForResponse('https://app.trasa/idp/login/tfa');


  expect(tfaResp.status()).toBe(200)

 //  await page.waitFor(8000)

  // tfaResp.json().then(data=>{
  //   // console.log('status: ', data)
  //     //expect(data.status).toBe("success")
  // })

  
  await page.screenshot({path: 'screen.png'});



 
})
}

// getTotp(TOTP_SSC)