import {UpdatePolicy} from "../src/policy/update";

require('expect-puppeteer')
const isReachable = require('is-reachable');

import Constants from '../Constants'

import { InitialUserLoginAndDeviceEnrol, LoginTfa } from '../src/login'
import { CreateService } from '../src/service/create'
import { logintests } from '../src/login'
import {CreateUser} from "../src/user/create";
import {UpdateUser} from "../src/user/update";
import {DeleteUser} from "../src/user/delete";
import {CreateUserGroup} from "../src/user/creategroup";
import {CreatePolicy} from "../src/policy/create";
import {UpdateService} from "../src/service/update";

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
// describe("user create",CreateUser)
// describe("user update",UpdateUser)
// describe("user delete",DeleteUser)
// describe("user group create",CreateUserGroup)
// describe("policy create",CreatePolicy)
// describe("policy create",UpdatePolicy)
// describe('InitialLoginAndEnroDevice', InitialUserLoginAndDeviceEnrol)

 describe('Service CRUD', CreateService)
 describe('Service CRUD', UpdateService)
// describe("SSH",sshtest)
