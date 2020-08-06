const getToken = require("../utils/totpgen");
require('expect-puppeteer')


const TRASA_HOSTNAME="https://app.trasa"
const TOTP_SEC="AV2COXZHVG4OAFSF"
const SSH_USERNAME="bhrg3se"
const SSH_PASSWORD="testpass"

export const sshtest= () => {

    beforeAll(async () => {
         // browser = await puppeteer.launch();
         // page = await browser.newPage();
        await page.goto('http://localhost:3000')

    })


    //console.log(getToken(process.env.TOTP_SEC))

    it('should connect ssh successfully', async function () {

        await page.waitForSelector('[name=email]');
        await page.type("[name=email]","root")
        await page.type("[name=password]","")


        await page.waitForSelector('[type=submit]');

         await page.click("[type=submit]")
       // await page.keyboard.press('Enter');


        let resp=await page.waitForResponse(response => {
           // console.log(response)
            return response.url() === TRASA_HOSTNAME + '/idp/login'
        });//




            expect(resp.status()).toBe(200)

            resp.json().then(data=>{
                expect(data.status).toBe("success")

            })


        await page.click("[id=totpButton]")
        await page.waitForSelector("[name=totpVal]")
        console.info(3)
        await page.type("[name=totpVal]",getToken(TOTP_SEC))
        await page.keyboard.press("Enter")
        resp=await page.waitForResponse(resp=> resp.url()===TRASA_HOSTNAME+"/idp/login/tfa")

        expect(resp.status()).toBe(200)
        let data=await resp.json()
        expect(data.status).toBe( "success")


        await page.waitForRequest(TRASA_HOSTNAME+'/api/v1/my');




        //Go to my page
        await page.goto(TRASA_HOSTNAME+"/my",{waitUntil: 'load'})


        //Search for auth app.
        await page.type("[placeholder=\"Search Authapps by name or hostname\"]","aws")



        //click on username
        await page.waitForSelector("[name='127.0.0.1']")
        await page.click("[name='127.0.0.1']")


        console.log("app selected")
        await page.waitForSelector("[name='"+SSH_USERNAME+"']")


        //This promise will wait for new popup
        const newPagePromise = new Promise(x => browser.once('targetcreated', target => x(target.page())));
        await page.click("[name='"+SSH_USERNAME+"']")

        const popup = await newPagePromise;

        console.log("pop up opened")
        //wait for password field to appear
        await popup.waitForSelector("[type=password]")
        await popup.type("[type=password]",SSH_PASSWORD)
        await popup.keyboard.press("Enter")

        await popup.waitForSelector("#totpButton")
        await popup.click("#totpButton")

        await popup.waitForSelector("[name=totpVal]")
        // await popup.type("[name=totpVal]",getToken(TOTP_SEC))
        await popup.type("[name=totpVal]",totp(TOTP_SEC))


        await popup.keyboard.press("Enter")

        console.log("totp entered in console")

        await delay(10000)

        console.log("after delay")


        await popup.type(".xterm-helper-textarea","ls")


        // await popup.keyboard.type("ls")
        await popup.keyboard.press("Enter")



        for (let i=0;i<10;i++){
            console.log("typing ls")
            await popup.keyboard.type("ls")
            await popup.keyboard.press("Enter")
            await delay(1000)
        }


        await page.screenshot({path: '/Users/bhrg3se/seknox/code/trasa/trasa-oss/tests/'+process.pid+'.png'});


    });


    afterAll(() => {
        return  browser.close();
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
    //     console.log(e)
    //     return false
    // }




    return true
}
