import Constants from "../../Constants";
require('expect-puppeteer')
const usersData=require('../../mock_data/users')


export const UpdateUser= ()=>{
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/users')
    })

    it('Should update user '+Constants.TRASA_DASHBOARD+'/users', async () => {
        const testUser=usersData.users[1]
        // await page.goto(Constants.TRASA_DASHBOARD+"/users")
        await expect(page).toMatch('Users')


        //matches href=/users/user/*

        await page.waitForSelector('[id="'+usersData.users[0].email+'"]')
        await page.click('[id="'+usersData.users[0].email+'"]')
        await expect(page).toMatch('Account Overview')


        await page.click('#userMenuBtn')

        await page.waitForSelector('#userUpdateBtn')

        await page.click('#userUpdateBtn')

        await page.waitForSelector("[name=firstName]")




        await page.click("[name=firstName]",{ clickCount: 3 })
        await page.type("[name=firstName]",testUser.firstName)

        await page.click("[name=lastName]",{ clickCount: 3 })
        await page.type("[name=lastName]",testUser.lastName)
        await page.click("[name=email]",{ clickCount: 3 })
        await page.type("[name=email]",testUser.email)
        await page.click("[name=userName]",{ clickCount: 3 })
        await page.type("[name=userName]",testUser.userName)

        //select org role
        await page.click("#userRole")

        await page.click("#orgAdmin")

        // await page.waitFor(5000)

        await page.click("#submit")


        let resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/user/update');

       await expect(resp.status()).toBe(200)

        // const respData=await  resp.json()
        //
        // await expect(respData.status).toBe("success")

        // await page.screenshot({path: 'src/user/update.png'})


    })


}

