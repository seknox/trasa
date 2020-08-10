require('expect-puppeteer')
const axios=require('axios')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const DeleteService = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')
    })

    it('Should delete a service '+Constants.TRASA_DASHBOARD+'/services', async () => {
        //await page.goto(Constants.TRASA_DASHBOARD+'/services', {timeout: 0});
        // await page.waitForNavigation({waitUntil: 'domcontentloaded'})
        await expect(page).toMatch('HTTPs applications')

        await page.waitForSelector('#'+ServicesMock[1].name)

        await page.click('#'+ServicesMock[1].name)
        //  await page.waitForNavigation({timeout:5000})

        await expect(page).toMatch("Configurations")

        await page.click('#deleteBtn')
        await page.click('#deleteConfirmBtn')




        let resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/services/delete');

        expect(resp.status()).toBe(200)

        //page.waitFor(30000)
        await page.waitForNavigation({waitUntil: "networkidle2"})





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

