
require('expect-puppeteer')
const isReachable = require('is-reachable');

import Constants from '../Constants'

import { InitialUserLoginAndDeviceEnrol, LoginTfa } from '../src/login'
import { CreateService } from '../src/service/create'
import { logintests } from '../src/login'
import usercrud from "../src/user";
// import {sshtest} from "../src/ssh";

 jest.setTimeout(30000)

beforeAll(async () => {
  let up = false
  while (up !== true) {
  await isReachable('app.trasa:443').then(u => up = u)
      await page.waitFor(200)
  }

     await page.goto(Constants.TRASA_DASHBOARD)
  })



// describe('InitialLoginAndEnroDevice', InitialUserLoginAndDeviceEnrol)
describe('Login', LoginTfa)
describe('Service CRUD', CreateService)
// describe("user crud",usercrud)

// describe("SSH",sshtest)
