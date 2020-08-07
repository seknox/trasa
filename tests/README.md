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
- Run test on already running trasa instance with `yarn test`
- Create fresh test infrastructure and test with `yarn run spinNTest`

## TODO

- Setup with CI server (Github action? ).


## Test Spec

### 1. Initial login flow

- [x] root user initial device enrol
- [x] root user complete login

### 2. Verify all intents login and tfa handlers support

- [ ] dashlogin
- [ ] access-proxy
- [ ] forget password
- [ ] change password
- [ ] enrol device

### 3. User CRUD

- [ ] create user
- [ ] update user
- [ ] delete user
- [ ] create user group
- [ ] update user group
- [ ] add users to group
- [ ] remove user from group
- [ ] delete group

### 4. Service CRUD

- [ ] create service
- [ ] update service
- [ ] delete service
- [ ] create service group
- [ ] update service group
- [ ] add service to group
- [ ] remove service from group
- [ ] delete group

### 5. Device CRUD

- [ ] enrol mobile device
- [ ] update mobile device hygiene
- [ ] remove mobile device
- [ ] enrol workstation 
- [ ] update workstation hygiene
- [ ] remove workstation

### 6. Policy CRUD

- [ ] create policy
- [ ] update policy
- [ ] delete policy


### 7. Access MAP

- [ ] assign user to service
- [ ] edit privilege
- [ ] remove user from service
- [ ] assign usergroup to service
- [ ] edit privilege
- [ ] remove usergroup from service
- [ ] assign usergroup to servicegroup
- [ ] edit privilege
- [ ] remove usergroup from servicegroup