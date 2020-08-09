
require('expect-puppeteer')
const isReachable = require('is-reachable');

import Constants from '../Constants'

import { InitialUserLoginAndDeviceEnrol, LoginTfa } from '../src/login'
import { CreateService } from '../src/service/create'
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

     await page.goto(Constants.TRASA_DASHBOARD)
  })



describe("user crud",usercrud)
// describe('InitialLoginAndEnroDevice', InitialUserLoginAndDeviceEnrol)
describe('Login', LoginTfa)
describe('Service CRUD', CreateService)
// describe("SSH",sshtest)
