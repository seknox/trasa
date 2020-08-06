require('expect-puppeteer')
const isReachable = require('is-reachable');

import { logintests } from '../src/login'
// import {sshtest} from "../src/ssh";

jest.setTimeout(30000)
beforeAll(async () => {


  let up = false 
  while (up !== true) {
  await isReachable('app.trasa:443').then(u => up = u)
  }


  
    //const page = await browser.newPage()
    await page.goto('https://app.trasa')
  })

describe('Login', logintests)
// describe("SSH",sshtest)
