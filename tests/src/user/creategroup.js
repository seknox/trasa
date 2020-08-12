import Constants from '../../Constants'

require('expect-puppeteer')
const groupData=require('../../mock_data/groups')



//create user
export const CreateUserGroup = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/users#User Groups')
    })

    it('Should Create a new user group '+Constants.TRASA_DASHBOARD+'/users#User Groups', async () => {
        const tesGroup=groupData.usergroups[0]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")
        await expect(page).toMatch('Users')

        await page.click('#createGroupBtn')
        await expect(page).toMatch('Create User Group')

        await page.type("[name=groupName]",tesGroup.name)


        // await page.waitFor(5000)

        const respPromise =page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/groups/create');

        await page.click("#createGroupSubmitBtn")


        let resp = await respPromise


         expect(resp.status()).toBe(200)

        //TODO test response status

        // const respData=await  resp.json()

      //  await expect(respData.status).toBe("success")

    })




}

