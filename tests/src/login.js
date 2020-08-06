require('expect-puppeteer')

export const logintests = () => {
  beforeAll(async () => {
    await page.goto('https://app.trasa')
  })

  it('should display "Dashboard Login" text on page', async () => {
    await expect(page).toMatch('Dashboard Login')
  })
}



