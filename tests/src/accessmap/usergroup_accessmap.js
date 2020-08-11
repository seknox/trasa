import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const UserGroupAccessMap = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should update a service '+Constants.TRASA_DASHBOARD+'/services', async () => {
        //await page.goto(Constants.TRASA_DASHBOARD+'/services', {timeout: 0});
        // await page.waitForNavigation({waitUntil: 'domcontentloaded'})
        await expect(page).toMatch('HTTPs applications')


        await page.waitForSelector('#'+ServicesMock[0].name)

        const policyLoaded=page.waitForResponse(r=>r.url().includes('policy/all'))
        const userGroupsLoaded=page.waitForResponse(r=>r.url().includes('accessmap/servicegroup/usergroupstoadd'))


        await page.click('#'+ServicesMock[0].name)
        //  await page.waitForNavigation({timeout:5000})

        await expect(page).toMatch("Configurations")


        await clickWithText(page,"Access Maps",'span')

        await expect(page).toMatch("User Access Map")

        await page.click('#assignUserGroupBtn')


        await policyLoaded
        await userGroupsLoaded


        await expect(page).toMatch("Select Groups")
        await expect(page).toMatch("Select Policy")
        await expect(page).toMatch("Assign Privilege")


        await clickWithText(page,"Select...",'span')
        await page.waitFor(1000)

        await clickWithText(page,"somegroup",'span')
        await clickWithText(page,"Select...",'span')

        await page.waitFor(1000)

        await clickWithText(page,"full",'span')

        await page.click('[name=privilege]')
        await page.type('[name=privilege]','admin')

        await page.click('#submitAccessMapBtn')




        //page.waitFor(30000)
        await page.waitForNavigation()

        //  await page.waitForSelector()


    })




}


