const getToken = require("../utils/totpgen");
require('expect-puppeteer')


const TRASA_HOSTNAME="http://localhost:3000"
const TRASA_API="https://localhost"
const TOTP_SEC="AV2COXZHVG4OAFSF"
const SSH_USERNAME="bhrg3se"
const SSH_PASSWORD="testpass"

export const sshtest= () => {

    beforeAll(async () => {
        await page.goto(TRASA_HOSTNAME)

    })



    it('should connect ssh successfully', async function () {

        await page.waitForSelector('[name=email]');
        await page.type("[name=email]","root")
        await page.type("[name=password]","changeme")

        await page.waitForSelector('[type=submit]');

         await page.click("[type=submit]")
       // await page.keyboard.press('Enter');


        let resp=await page.waitForResponse(response => {
            return response.url() === TRASA_API + '/idp/login'
        });//




            expect(resp.status()).toBe(200)

            resp.json().then(data=>{
                expect(data.status).toBe("success")

            })


        await page.click("[id=totpButton]")
        await page.waitForSelector("[name=totpVal]",{timeout:1000})
        await page.type("[name=totpVal]",getToken(TOTP_SEC))


        await page.keyboard.press("Enter")

        // await page.waitFor(2000)

        resp=await page.waitForResponse(resp=> {
            if (resp.url() === TRASA_API + "/idp/login/tfa"){
               // //console.log(resp)
                return true
            }
        })

      //  expect(resp.ok()).toBe(true)


        expect(resp.status()).toBe(200)


     // //   console.info(resp)
     //    resp.json().then(data=>{
     //        console.info(data)
     //        expect(data.data.status).toBe( "success")
     //
     //    }).catch(e=>{
     //        console.error(e)
     //        throw new Error(e)
     //    })




        //Go to my page
        await page.goto(TRASA_HOSTNAME+"/my",{waitUntil: 'load'})


        // //Search for auth app.
        // await page.type("[placeholder=\"Search services by name or hostname\"]","127.0.0.1")


        //select element based on inner text
        //
        // await page.evaluate(() => {
        //     let btns = [...document.querySelectorAll("button")];
        //     btns.forEach(function (btn) {
        //         if (btn.toString().includes("Connect") )
        //             btn.click();
        //     });
        // });
        // await page.click("button.MuiButton-root:nth-child(2)")

        //click on username


        await page.waitFor(5000)

        await page.click("[name='127.0.0.1']")


        await page.waitForSelector("[id='"+SSH_USERNAME+"']",{timeout:1000})


        //This promise will wait for new popup
        const newPagePromise = new Promise(x => browser.once('targetcreated', target => x(target.page())));
        await page.click("[id='"+SSH_USERNAME+"']")
        await page.screenshot({path: '/Users/bhrg3se/seknox/code/trasa/trasa-oss/tests/my'+process.pid+'.png'});

        const popup = await newPagePromise;

        //wait for password field to appear
        await popup.waitForSelector("[type=password]",{timeout:1000})
        await popup.type("[type=password]",SSH_PASSWORD)
        await popup.keyboard.press("Enter")

        await popup.waitForSelector("#totpButton")
        await popup.click("#totpButton")

        await popup.waitForSelector("[name=totpVal]")
        // await popup.type("[name=totpVal]",getToken(TOTP_SEC))
        await popup.type("[name=totpVal]",getToken(TOTP_SEC))


        await popup.keyboard.press("Enter")


        await popup.waitFor(10000)



        await popup.type(".xterm-helper-textarea","ls")


        await popup.keyboard.press("Enter")



        for (let i=0;i<10;i++){
            ////console.log("typing ls")
            await popup.keyboard.type("ls")
            await popup.keyboard.press("Enter")
            await popup.waitFor(1000)
        }


      //  await popup.screenshot({path: '/Users/bhrg3se/seknox/code/trasa/trasa-oss/tests/'+process.pid+'.png'});


    });


};










async function checkresp(resp) {
    if(!resp.ok()){
        return false
    }

    if(!(resp.status() === 200 )){
        return false
    }
    //
    // try {
    //     let data=await resp.data.json()
    //     if (data.status!=="success"){
    //         return false
    //     }
    // }
    // catch (e) {
    //     ////console.log(e)
    //     return false
    // }




    return true
}
