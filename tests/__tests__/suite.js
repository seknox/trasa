
require('expect-puppeteer')
const isReachable = require('is-reachable');

import { logintests } from '../src/login'
import usercrud from "../src/user";
// import {sshtest} from "../src/ssh";

 jest.setTimeout(60000)

beforeAll(async () => {
  let up = false
  while (up !== true) {
  await isReachable('app.trasa:443').then(u => up = u)
      await page.waitFor(200)
  }

     await page.goto('https://app.trasa')
  })



describe('Login', logintests)
describe("user crud",usercrud)
// describe("SSH",sshtest)
