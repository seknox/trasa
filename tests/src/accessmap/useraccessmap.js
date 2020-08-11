import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const UserAccessMap = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should update a service '+Constants.TRASA_DASHBOARD+'/services', async () => {
        //await page.goto(Constants.TRASA_DASHBOARD+'/services', {timeout: 0});
        // await page.waitForNavigation({waitUntil: 'domcontentloaded'})
        await expect(page).toMatch('HTTPs applications')


        await page.waitForSelector('#'+ServicesMock[0].name)

        const policyLoaded=page.waitForResponse(r=>r.url().includes('policy/all'))
        const usersLoaded=page.waitForResponse(r=>r.url().includes('/v1/user/all'))


        await page.click('#'+ServicesMock[0].name)
        //  await page.waitForNavigation({timeout:5000})

        await expect(page).toMatch("Configurations")


       await clickWithText(page,"Access Maps",'span')

        await expect(page).toMatch("User Access Map")

        await page.click('#assignUserBtn')


        await policyLoaded
        await usersLoaded


        await expect(page).toMatch("Select Users")
        await expect(page).toMatch("Select Policy")
        await expect(page).toMatch("Assign Privilege")


        await clickWithText(page,"Select...",'span')
        await page.waitFor(1000)

       await clickWithText(page,"auser",'span')
        await clickWithText(page,"Select...",'span')

        await page.waitFor(1000)

        await clickWithText(page,"full",'span')

        await page.click('[name=privilege]')
        await page.type('[name=privilege]','admin')

        await page.click('#submitAccessMapBtn')




        //page.waitFor(30000)
        await page.waitForNavigation()

      //  await page.waitForSelector()

        await page.screenshot({path: 'src/service/create.png'})

    })




}


