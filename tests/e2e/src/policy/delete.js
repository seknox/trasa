//MUIDataTableSelectCell-0

import Constants from '../../Constants'

require('expect-puppeteer')
const policyData=require('../../mock_data/policies')



//create user
export const DeletePolicy = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/control',{waitUntil:'load'})
    })

    it('Should update a policy '+Constants.TRASA_DASHBOARD+'/control', async () => {
        const testPolicy=policyData.policies[0]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")


        await expect(page).toMatch('Policies')

        await page.waitForSelector('[id="new policy"]',{timeout:5000})



        await page.click('#MUIDataTableSelectCell-0')

        const respPromise=page.waitForResponse(r=>r.url().includes('policy/delete'))

        await page.click('button[title=Delete]')



        let resp = await respPromise

        await expect(resp.status()).toBe(200)

        // const respData=await resp.json()
        //
        // await expect(respData.status).toBe('success')




        // await page.waitFor(20000)
        //  await page.waitFor(10000)


    })




}



