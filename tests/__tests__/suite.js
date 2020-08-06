require('expect-puppeteer')

import { logintests } from '../src/login'
import {sshtest} from "../src/ssh";

jest.setTimeout(30000)
beforeAll(async () => {
    //const page = await browser.newPage()
    await page.goto('http://localhost:3000')
  })

describe('Login', logintests)
describe("SSH",sshtest)
