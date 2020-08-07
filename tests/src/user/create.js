import {TRASA_HOSTNAME} from "../../Constants";

require('expect-puppeteer')
const usersData=require('../../mock_data/users')


//create user
export default async ()=>{
    const testUser=usersData.users[0]
    await page.goto(TRASA_HOSTNAME+"/users",{waitUntil:"load"})

    await page.waitFor(3000)
    await page.screenshot({path: 'screen.png'});


    await expect(page).toMatch("Users")
    await expect(page).toMatch("Create User")
    await page.click('[name=createUserBtn]')

    await page.waitForSelector("[name=firstName]",{timeout:1000})

    await page.type("[name=firstName]",testUser.firstName)
    await page.type("[name=lastName]",testUser.lastName)
    await page.type("[name=email]",testUser.email)
    await page.type("[name=userName]",testUser.userName)

    //select org role
    await page.click("[name=userRole]")
    await page.waitForSelector("[data-value=orgAdmin]",{timeout:1000})
    await page.click("[data-value=orgAdmin]")

    // //select status
    // await page.click("[name=status]")
    // await page.waitForSelector("[data-value=active]",{timeout:1000})
    // await page.click("[data-value=active]")


    await page.click("[name=createUserSubmitBtn]")


    const resp=await page.waitForResponse(resp=> {
        if (resp.url() === TRASA_API + "/api/v1/user/create"){
            // //console.log(resp)
            return true
        }
    })


    const respData= await resp.json()



    await expect(page).toMatch('Send this link to user for signup')




    //{"POST":{"scheme":"https","host":"app.trasa","filename":"/api/v1/user/create","remote":{"Address":"127.0.0.1:443"}}}


    //Send this link to user for signup.



}
