const getToken = require("../utils/totpgen");
const puppeteer = require('puppeteer');
const assert = require('assert')
 // let mocha = require('mocha');
// let describe = jest.describe;


// var totp = require('totp-generator');

let browser
let page



describe("Login into dashboard using totp",async () => {

    console.log(getToken(process.env.TOTP_SEC))
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(process.env.TRASA_HOSTNAME);


    await page.type("[name=email]",process.env.TRASA_ID)
    await page.type("[name=password]",process.env.TRASA_PASSWORD)

    await page.click("[name=submit]")
    let resp=await page.waitForResponse(response => response.url() === process.env.TRASA_HOSTNAME+'/idp/login');//



    it('should return success when logged in',async function () {
        assert.equal(resp.status(), 200, "statusCode is not 200")
        let data=await resp.json()
        assert.equal(data.status, "success", "status in trasa response is not success")


    });

    await page.click("[type=button]")
    await page.waitForSelector("[name=totpVal]")

    await page.type("[name=totpVal]",getToken(process.env.TOTP_SEC))
    await page.keyboard.press("Enter")
    resp=await page.waitForResponse(resp=> resp.url()===process.env.TRASA_HOSTNAME+"/idp/login/tfa")

    assert.equal(resp.status(), 200, "statusCode is not 200")
    data=await resp.json()
    assert.equal(data.status, "success", "status in trasa response is not success")


    await page.waitForRequest(process.env.TRASA_HOSTNAME+'/api/v1/my');

    await browser.close();
});










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
