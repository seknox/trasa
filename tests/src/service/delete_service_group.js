

import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')


// CreateService
export const DeleteServiceGroup = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should delete group '+Constants.TRASA_DASHBOARD+'/services', async () => {

        await expect(page).toMatch('HTTPs applications')


        await clickWithText(page,"Service Groups",'span')



        await expect(page).toMatch("Service Group")

        await page.waitForSelector('#somegroup1')
        await page.click('#somegroup1')




        const navigationPromise=page.waitForNavigation()
        const respPromise = page.waitForResponse(r=>r.url().includes('api/v1/groups/delete'))

        await page.click('button[title=delete]')

        const resp= await respPromise
        await expect(resp.status()).toBe(200)
        await navigationPromise

        //  await page.waitForSelector()


    })




}


