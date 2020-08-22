import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')

import Constants from '../../Constants' //= require('../../Constants')

import { ServicesMock } from '../../mock_data/services'

// CreateService
export const RenameAccessMap = () => {

    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')

    })

    it('Should add a user access map '+Constants.TRASA_DASHBOARD+'/services', async () => {
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



        await policyLoaded
        await usersLoaded



        await page.click('#editUserPrivilegeBtn')
        await page.waitForSelector('#username')
        await page.type('#username','admin2')
        let saved =page.waitForResponse(r=>r.url().includes('accessmap/service/user/update'))
        await page.click('#savePrivilegeBtn')

        let resp= await saved
        await expect(resp.status()).toBe(200)
        let respData = await resp.json()
        await expect(respData.status).toBe('success')

        await page.click('#editGroupPrivilegeBtn')
        await page.waitForSelector('#username')
        await page.type('#username','admin2')

        saved =page.waitForResponse(r=>r.url().includes('servicegroup/usergroup/update'))
        await page.click('#savePrivilegeBtn')

        resp= await saved
        await expect(resp.status()).toBe(200)
        respData = await resp.json()
        await expect(respData.status).toBe('success')



    })




}


