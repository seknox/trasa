require('expect-puppeteer')
const axios=require('axios')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

export const CreateService = () => {

  it('Should load service page on url address: '+Constants.TRASA_DASHBOARD+'/services', async () => {
    await page.goto(Constants.TRASA_DASHBOARD+'/services');
    await expect(page).toMatch('HTTPs applications')
  })

  it('Should create a new service', async () => {
    await page.goto(Constants.TRASA_DASHBOARD+'/services');
   await page.waitForSelector('[name=create-new-service-button]')
    await page.click('[name=create-new-service-button]')
    await expect(page).toMatch('Integrate New Service')

    await page.type('#serviceName', ServicesMock[0].name)
     await page.click('#serviceType');
     await page.click('#ssh')

    await page.screenshot({path: 'src/service/create.png'})

  })


}

// getTotp(TOTP_SSC)

