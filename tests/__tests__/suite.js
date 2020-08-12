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
import {DeleteService} from "../src/service/delete";
import {UserAccessMap} from "../src/accessmap/useraccessmap";
import {AddUserToGroup} from "../src/user/add_user_to_group";
import {Jpt} from "../src/usercrud";
import {UserGroupAccessMap} from "../src/accessmap/usergroup_accessmap";
import {CreateServiceGroup} from "../src/service/create_group";
import {AddServiceToGroup} from "../src/service/add_service_to_group";
import {UserGroupServiceGroupAccessMap} from "../src/accessmap/usergroup_appgroup_map";
import {RenameServiceGroup} from "../src/service/rename_service_group";
import {DeleteServiceGroup} from "../src/service/delete_service_group";
import {RenameAccessMap} from "../src/accessmap/rename_accessmap_privilege";

// import {sshtest} from "../src/ssh";

 jest.setTimeout(60000)

beforeAll(async () => {
    await page.setDefaultTimeout(10000);
  let up = false
  while (up !== true) {
  let u=await isReachable('app.trasa:443')
      up=u
   }

    // await page.goto(Constants.TRASA_DASHBOARD)
  })
describe('InitialLoginAndEnroDevice', InitialUserLoginAndDeviceEnrol)

// describe('Login', LoginTfa)
describe("policy create",CreatePolicy)

describe("user create",CreateUser)
describe("user group create",CreateUserGroup)
describe("add user to group",AddUserToGroup)

//
describe('Service create', CreateService)
describe("service group create",CreateServiceGroup)
// describe("add service to group",AddServiceToGroup)

describe('User access map',UserAccessMap)
describe('User group access map',UserGroupAccessMap)
// describe('User group service group access map',UserGroupServiceGroupAccessMap)

describe('rename access map privilege',RenameAccessMap)



describe('Rename service group',RenameServiceGroup)
describe('Delete service group',DeleteServiceGroup)


// describe("user update",UpdateUser)
// describe("user delete",DeleteUser)
// describe("policy update",UpdatePolicy)
//
//  describe('Service update', UpdateService)
//  describe('Service delete', DeleteService)

// describe("SSH",sshtest)
