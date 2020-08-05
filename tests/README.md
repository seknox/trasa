# Integration test for TRASA



This is a Jest and Puppeteer based project for integration (or more of end-to-end) testing TRASA.

We render full dashboard inside puppeteer and run test suites which will 

```
|--> Test user interactions with dashboard. This in turn will

  |--> Test API served from TRASA server. This will cover

    |--> All TRASA functions that are required while serving API.
```

**This project is intended for integration test.** 

**All unit tests reside in respective project.**


## Prerequisite:

TRASA server running on `https://app.trasa`

## TODO

- Auto spin TRASA server.
- Setup with CI server.