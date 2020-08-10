
require('expect-puppeteer')
const isReachable = require('is-reachable');

import Constants from '../Constants'

import { InitialUserLoginAndDeviceEnrol, LoginTfa } from '../src/login'
import { CreateService } from '../src/service/create'
import { logintests } from '../src/login'
import {CreateUser} from "../src/user/create";
import {UpdateUser} from "../src/user/update";
import {DeleteUser} from "../src/user/delete";
// import {sshtest} from "../src/ssh";

 jest.setTimeout(60000)

beforeAll(async () => {
  let up = false
  while (up !== true) {
  await isReachable('app.trasa:443').then(u => up = u)
  }

    // await page.goto(Constants.TRASA_DASHBOARD)
  })


describe('Login', LoginTfa)
describe("user crud",CreateUser)
describe("user crud",UpdateUser)
describe("user crud",DeleteUser)
// describe('InitialLoginAndEnroDevice', InitialUserLoginAndDeviceEnrol)

 describe('Service CRUD', CreateService)
// describe("SSH",sshtest)
