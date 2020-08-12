import Constants from "../../Constants";
require('expect-puppeteer')
const usersData=require('../../mock_data/users')


export const DeleteUser= ()=>{
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/users')
    })

    it('Should delete user '+Constants.TRASA_DASHBOARD+'/users', async () => {
        const testUser=usersData.users[1]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")
        await expect(page).toMatch('Users')


        //matches href=/users/user/*

        await page.waitForSelector('[id="'+usersData.users[1].email+'"]')
        await page.click('[id="'+usersData.users[1].email+'"]')
        await expect(page).toMatch('Account Overview')


        await page.click('#userMenuBtn')

        await page.waitForSelector('#userDeleteBtn')
        await page.click('#userDeleteBtn')


        await page.waitForSelector('#confirmDeleteBtn')

        const respPromise=page.waitForResponse(resp=>{
            return (resp.url().includes('/api/v1/user/delete') )
        });
        await page.click('#confirmDeleteBtn')



        let resp = await respPromise

        await expect(resp.status()).toBe(200)

        // const respData=await  resp.json()
        //
        // await expect(respData.status).toBe("success")
        // // await page.waitFor(20000)
        //  await page.waitFor(10000)
        // await expect(page).toMatch('Verification Link')
        // await page.screenshot({path: 'src/user/update.png'})


    })


}

