require('expect-puppeteer')

import { logintests } from '../src/login'
import {sshtest} from "../src/ssh";


// describe('Login', logintests)
describe("SSH",sshtest)
