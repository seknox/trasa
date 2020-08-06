require('expect-puppeteer')

import { getTotp } from '../utils/totpgen'

const loginData = {
  email: 'tree',
  password: 'changemenever1#',
  orgID: '',
  intent: '',
  idpName: '',
};

function delay(time) {
  return new Promise(function(resolve) { 
      setTimeout(resolve, time)
  });
}



export const logintests = () => {

  it('should display "Dashboard Login" text on page', async () => {
    await expect(page).toMatch('Dashboard Login')
  })

  it('should fill form', async () => {
    // await expect(page).toFillForm('form[name="loginform"]', {
    //   email: loginData.email,
    //   password: loginData.password,
    // })
 
  
   await page.type('#email', loginData.email)
    await page.type('#password', loginData.password);

 
    await page.click('#submit')
    await page.waitFor(2000)

    await page.screenshot({path: 'screenshot.png'});
    await expect(page).toMatch('Choose second step verification method')
  })



 
}



