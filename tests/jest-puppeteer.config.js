module.exports = {
    launch: {
      // dumpio: true,
      ignoreHTTPSErrors: true,
       slowMo: 150,
      headless: false,
     // devtools: true,
      defaultViewport: {
        width: 1300,
        height: 800,
      }
    },
     browser: 'chromium',
    browserContext: 'incognito',

  }