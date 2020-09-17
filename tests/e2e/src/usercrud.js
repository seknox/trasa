import {clickWithText} from "../utils/utils";

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

const TRASA_HOSTNAME="https://localhost"
let TOTP_SSC=""




export const SomeTest = () => {



  it('should fill form', async () => {

   await  page.goto('https://www.google.com')
      clickWithText(page,"Gmail")

      await page.waitFor(5000)
      await expect(page).toMatch("Google")

})
}

// getTotp(TOTP_SSC)
