import Constants from '../../Constants'

require('expect-puppeteer')
const policyData=require('../../mock_data/policies')



//create user
export const CreatePolicy = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/control',{waitUntil:'load'})
    })

    it('Should Create a new policy '+Constants.TRASA_DASHBOARD+'/control', async () => {
        const testPolicy=policyData.policies[0]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")


        await expect(page).toMatch('Policies')

        await page.click('#createPolicyBtn')
        await expect(page).toMatch('Create New Policy')

        await page.type("[name=policyName]",testPolicy.name)

        await page.click('#nextBtn')

        await expect(page).toMatch('Mandatory 2FA:')
        await page.click('[type=checkbox]')




        await page.click("#sessionRecordingTab")
        await expect(page).toMatch('Record Session: ')
        await page.click('[type=checkbox]')


        await page.click("#dayTimePolicyTab")
        await expect(page).toMatch('Days:')
        await expect(page).toMatch('Time:')


        await page.click('#daysPolicyDropDown')
        await page.click('[data-value=Sunday]')
        await page.click('[data-value=Monday]')
        await page.click('[data-value=Tuesday]')
        await page.click('[data-value=Wednesday]')
        await page.click('[data-value=Thursday]')
        await page.click('[data-value=Friday]')
        await page.click('[data-value=Saturday]')

        await page.keyboard.press("Escape")

        await page.focus('#FROM')
        await page.type('#FROM','0100am')
await page.focus('#TO')
        await page.type('#TO','1159pm')


        await page.click('#addBtn')

        await page.click('#nextBtn')
        //await page.click('#nextBtn')








        const navPromise=page.waitForResponse(r=>r.url().includes('policy/create'))
        await page.click("#submitBtn")


        let resp = await navPromise

        await expect(resp.status()).toBe(200)

       // await page.waitForNavigation({ waitUntil: 'load' })

        await page.waitForSelector('[id="'+testPolicy.name+'"]',{timeout:5000})


        await expect(page).toMatch(testPolicy.name)


        // await page.waitFor(20000)
        //  await page.waitFor(10000)


    })




}


