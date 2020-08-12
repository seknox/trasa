require('expect-puppeteer')
const axios=require('axios')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const CreateService = () => {
    let getAllServicesPromise
  beforeAll(async () => {
      getAllServicesPromise=page.waitForResponse(r=>r.url().includes('/api/v1/services/all'))
    await page.goto(Constants.TRASA_DASHBOARD+'/services')
    })

  it('Should Create a new service '+Constants.TRASA_DASHBOARD+'/services', async () => {
    //await page.goto(Constants.TRASA_DASHBOARD+'/services', {timeout: 0});
    // await page.waitForNavigation({waitUntil: 'domcontentloaded'})
    await expect(page).toMatch('HTTPs applications')

      await getAllServicesPromise

    await page.click('#create-new-service-button')
    await expect(page).toMatch('Integrate New Service')

    await page.waitForSelector('#serviceName')
    await page.type('#serviceName', ServicesMock[0].name)
     await page.click('#serviceType');
     await page.click('#ssh')
     await page.type('#hostname', ServicesMock[0].hostname)


      const respPromise=page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/services/create');
     const navPromise=page.waitForNavigation()
     await page.click('#submit')

     let resp = await respPromise;

    expect(resp.status()).toBe(200)

    //page.waitFor(30000)
    await navPromise

    //let serviceID = ''

    // resp.json().then(d=>{
    //     expect(d.status).toBe("success")
    //     console.log('dddddddddddddd: ', d)
    //    // serviceID = d.data[0].ID
    // })

    // console.log('serviceID was::: ', serviceID)

    // resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/services'+serviceID );

    // await expect(page).toMatch(ServicesMock[0].name)

    await page.screenshot({path: 'src/service/create.png'})

  })




}

// getTotp(TOTP_SSC)

