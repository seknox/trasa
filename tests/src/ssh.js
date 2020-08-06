const getToken = require("../utils/totpgen");
require('expect-puppeteer')


const TRASA_HOSTNAME="https://app.trasa"
const TOTP_SEC="AV2COXZHVG4OAFSF"

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
        await page.type("[name=password]","Anyth1n&123")

        console.info(1)
        await page.waitForSelector('[type=submit]');

         await page.click("[type=submit]")
       // await page.keyboard.press('Enter');


        let resp=await page.waitForResponse(response => {
           // console.log(response)
            return response.url() === TRASA_HOSTNAME + '/idp/login'
        });//

        console.info(2)


        it('should return success when logged in',async function () {
            assert.equal(resp.status(), 200, "statusCode is not 200")
            let data=await resp.json()
            assert.equal(data.status, "success", "status in trasa response is not success")


        });

        await page.click("[id=totpButton]")
        await page.waitForSelector("[name=totpVal]")

        await page.type("[name=totpVal]",getToken(TOTP_SEC))
        await page.keyboard.press("Enter")
        resp=await page.waitForResponse(resp=> resp.url()===TRASA_HOSTNAME+"/idp/login/tfa")

        await page.screenshot({path: '/Users/bhrg3se/seknox/code/trasa/trasa-oss/tests/'+process.pid+'.png'});

        assert.equal(resp.status(), 200, "statusCode is not 200")
        let data=await resp.json()
        assert.equal(data.status, "success", "status in trasa response is not success")


        await page.waitForRequest(TRASA_HOSTNAME+'/api/v1/my');

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
