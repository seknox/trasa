import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')


// CreateService
export const UserGroupServiceGroupAccessMap = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should add a user group service group access map '+Constants.TRASA_DASHBOARD+'/services', async () => {

        await expect(page).toMatch('HTTPs applications')


        await clickWithText(page,"Service Groups",'span')


        const policyLoaded=page.waitForResponse(r=>r.url().includes('policy/all'))
        const userGroupsLoaded=page.waitForResponse(r=>r.url().includes('accessmap/servicegroup/usergroupstoadd'))



        await expect(page).toMatch("Service Group")

        await page.waitForSelector('#someservicegroup')
        await page.click('#someservicegroup')



        await clickWithText(page,"Assigned User Groups",'span')

        await expect(page).toMatch("Group Access Map")

        await page.click('#assignUserGroupBtn')


        await policyLoaded
        await userGroupsLoaded


        await expect(page).toMatch("Select Groups")
        await expect(page).toMatch("Select Policy")
        await expect(page).toMatch("Assign Privilege")


        await clickWithText(page,"Select...",'span')
        await page.waitFor(1000)

        await clickWithText(page,"someusergroup",'span')
        await clickWithText(page,"Select...",'span')

        await page.waitFor(1000)

        await clickWithText(page,"full-policy",'span')

        await page.click('[name=privilege]')
        await page.type('[name=privilege]','admin')

        const navigationPromise=page.waitForNavigation()
        await page.click('#submitAccessMapBtn')


        await navigationPromise

        //  await page.waitForSelector()


    })




}


