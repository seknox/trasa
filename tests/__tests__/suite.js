require('expect-puppeteer')

import { logintests } from '../src/login'

jest.setTimeout(30000)
beforeAll(async () => {
    //const page = await browser.newPage()
    await page.goto('http://localhost:3000')
  })
  
describe('Login', logintests)