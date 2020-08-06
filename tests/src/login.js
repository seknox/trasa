require('expect-puppeteer')

const loginData = {
  email: 'tree',
  password: 'changemenever1#',
  orgID: '',
  intent: '',
  idpName: '',
};

function delay(time) {
  return new Promise(function(resolve) { 
      setTimeout(resolve, time)
  });
}

export const logintests = () => {
  beforeAll(async () => {
    await page.goto('http://localhost:3000')
  })

  it('should display "Dashboard Login" text on page', async () => {
    await expect(page).toMatch('Dashboard Login')
  })

  it('should fill form', async () => {
    // await expect(page).toFillForm('form[name="loginform"]', {
    //   email: loginData.email,
    //   password: loginData.password,
    // })
 
    await page.waitForSelector('#email');
    await page.type('#email', loginData.email)//);
   // await page.waitFor(1000)
    await page.type('#password', loginData.password);
    await page.keyboard.press('Enter');

    // await page.focus("#email")
    // await page.keyboard.type('test')
    // await page.focus("#password")
    // await page.keyboard.type('jeskl')

 
    await page.click('#submit')
    await page.waitFor(4000)

    await page.screenshot({path: 'screenshot.png'});
  })

 
}



