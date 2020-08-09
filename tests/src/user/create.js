import Constants from '../../Constants' 

require('expect-puppeteer')
const usersData=require('../../mock_data/users')



//create user
export const CreateUser = () => {
    beforeAll(async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/users')
        })
    
    it('Should Create a new user '+Constants.TRASA_DASHBOARD+'/users', async () => {
    const testUser=usersData.users[0]
   // await page.goto(Constants.TRASA_DASHBOARD+"/users")
    await expect(page).toMatch('Users')

    await page.click('[name=createUserBtn]')
    await expect(page).toMatch('Create User')

    await page.type("[name=firstName]",testUser.firstName)
    await page.type("[name=lastName]",testUser.lastName)
    await page.type("[name=email]",testUser.email)
    await page.type("[name=userName]",testUser.userName)

    //select org role
    await page.click("#userRole")

    await page.click("#orgAdmin")
 
   // await page.waitFor(5000)

    await page.click("#submit")
   

     let resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/user/create');

     expect(resp.status()).toBe(200)
    // await page.waitFor(20000)
  //  await page.waitFor(10000)
    await expect(page).toMatch('Verification Link')
    await page.screenshot({path: 'src/user/create.png'})
   
  
    })
  
  
  
  
  }

  

// export default async ()=>{
//     const testUser=usersData.users[0]
//     await page.goto(Constants.TRASA_DASHBOARD+"/users")

//     await expect(page).toMatch("Users")

//     await page.click('[name=createUserBtn]')
//     await expect(page).toMatch('Create User')

//     await page.type("[name=firstName]",testUser.firstName)
//     await page.type("[name=lastName]",testUser.lastName)
//     await page.type("[name=email]",testUser.email)
//     await page.type("[name=userName]",testUser.userName)

//     //select org role
//     await page.click("[name=userRole]")

//     await page.click("[data-value=orgAdmin]")



//     await page.click("[name=createUserSubmitBtn]")

//     let resp = await page.waitForResponse(Constants.TRASA_HOSTNAME+'/api/v1/user/create');

//     expect(resp.status()).toBe(200)

//     //await page.waitForNavigation({waitUntil: "networkidle2"})



//     await expect(page).toMatch('Send this link to user for signup')

//     await page.screenshot({path: 'src/user/create.png'})

//     return



// }
