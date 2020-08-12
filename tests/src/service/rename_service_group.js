

import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')


// CreateService
export const RenameServiceGroup = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should rename group '+Constants.TRASA_DASHBOARD+'/services', async () => {

        await expect(page).toMatch('HTTPs applications')


        await clickWithText(page,"Service Groups",'span')



        await expect(page).toMatch("Service Group")

        await page.waitForSelector('#someservicegroup')
        await page.click('#someservicegroup')

        await page.click('#editGroupBtn')

        await page.waitForSelector('[name=groupName]')
        await page.type('[name=groupName]','somegroup1')

        const navigationPromise=page.waitForNavigation()
        const respPromise = page.waitForResponse(r=>r.url().includes('api/v1/groups/update'))

        await page.click('#createGroupSubmitBtn')

        const resp= await respPromise
        await expect(resp.status()).toBe(200)
        await navigationPromise

        //  await page.waitForSelector()


    })




}


