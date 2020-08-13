require('expect-puppeteer')
const axios=require('axios')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const UpdateService = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')
    })

    it('Should update a service '+Constants.TRASA_DASHBOARD+'/services', async () => {
        //await page.goto(Constants.TRASA_DASHBOARD+'/services', {timeout: 0});
        // await page.waitForNavigation({waitUntil: 'domcontentloaded'})
        await expect(page).toMatch('HTTPs applications')


        await page.waitForSelector('#'+ServicesMock[0].name)

        let respPromise=page.waitForResponse(r=>r.url().includes('/api/v1/service'))

        await page.click('#'+ServicesMock[0].name)
      //  await page.waitForNavigation({timeout:5000})

        await expect(page).toMatch("Configurations")
        await respPromise

        await page.click('#configEditBtn')


        await page.click('#serviceName', {clickCount:3})
        await page.type('#serviceName', ServicesMock[1].name)
        await page.click('#serviceType');
        await page.click('#rdp')
        await page.click('#hostname',{clickCount:3})
        await page.type('#hostname', ServicesMock[1].hostname)

         respPromise =page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/services/update');
        const navPromise = page.waitForNavigation()
        await page.click('#submit')

        let resp = await respPromise
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

        //await page.screenshot({path: 'src/service/create.png'})

    })




}

// getTotp(TOTP_SSC)

