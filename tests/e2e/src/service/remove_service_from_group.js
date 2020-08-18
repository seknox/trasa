import Constants from '../../Constants'
import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')
const groupData=require('../../mock_data/groups')
import { ServicesMock } from '../../mock_data/services'



//create user
export const RemoveServiceFromGroup = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')
    })

    it('Should add services to a group '+Constants.TRASA_DASHBOARD+'/services', async () => {
        await expect(page).toMatch('HTTPs applications')


        await clickWithText(page,"Service Groups",'span')

        await expect(page).toMatch("Service Group")

        await page.waitForSelector('#someservicegroup')

        let respPromise= page.waitForResponse(r=>r.url().includes('/api/v1/groups/service'))
        await page.click('#someservicegroup')
        await respPromise


        await page.waitForSelector('#MUIDataTableSelectCell-0')
        await page.click('#MUIDataTableSelectCell-0')



        respPromise=page.waitForResponse(r=>r.url().includes('groups/service/update'))

        await page.click('button[title=Delete]')


        const resp=await respPromise
        await expect(resp.status()).toBe(200)




    })




}

