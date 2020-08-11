import Constants from '../../Constants'
import {clickWithText} from "../../utils/utils";

require('expect-puppeteer')
const groupData=require('../../mock_data/groups')



//create user
export const CreateServiceGroup = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services')
    })

    it('Should Create a new service group '+Constants.TRASA_DASHBOARD+'/services', async () => {
        const tesGroup=groupData.servicegroups[0]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")
        await expect(page).toMatch('HTTPs applications')


        await clickWithText(page,"Service Groups",'span')


        await clickWithText(page,"Create group",'span')


        await expect(page).toMatch('Create service Group')

        await page.type("[name=groupName]",tesGroup.name)


        // await page.waitFor(5000)

        await page.click("#createGroupSubmitBtn")


        let resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/groups/create');


        expect(resp.status()).toBe(200)


    })




}

