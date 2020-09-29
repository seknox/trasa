// getTrasacore retreives trasacore address from local storage.
async function getTrasacore() {
  let val = await browser.storage.local.get();
  if (val.trasaCore) {
    return val.trasaCore;
  } else {
    return 'https://trasa.seknox.com';
  }
}

// getExtID retreives extension token from local storage
async function getExtID() {
  let val = await browser.storage.local.get();
  return val.extID;
}

// getHosts retreives hosts array from local storage
async function getHosts() {
  let val = await browser.storage.local.get();
  if (val.hosts) {
    return val.hosts;
  } else {
    let emptyArr = ['test'];
    return emptyArr;
  }
}

// getSessionID takes hostname as argument and returns its corresponding sessionID
async function getSessionID() {
  //let p = await browser.storage.local.get()
  //console.log(p)
  let val = await browser.storage.local.get('sessionStore');
  //console.log("## sesion store: ", val.sessionStore)

  return val;
}

// get root domain retreives root domain from local storage.
async function getRootDomain() {
  let val = await browser.storage.local.get(); //.then(val => {
  //	console.log('rootDomain '+ val.rootDomain)
  return val.rootDomain;
}

// getSsoDomainFromOrigin retreives single sign on domain from origin header.
// is the sso matches, trasa extension auto fills the password for user.
function getSsoDomainFromOrigin(domain) {
  let split = domain.split('.');
  let parsed = split.slice(Math.max(split.length - 4, 0));
  let res = parsed.join('.');
  return res;
}

// getSsoDomain retreives single sign on domain from local storage.
async function getSsoDomain() {
  let val = await browser.storage.local.get(); //.then(val => {
  //  console.log('session: ', val )
  return val.ssoDomain;
}

// getTrasauser retreives trasauser from local storage
async function getTrasauser() {
  let val = await browser.storage.local.get(); //.then(val => {
  return val.trasauser;
}

// getRootDomainFromOrigin retreives root domain from Origin header value.
// (i.e. trasa.io from app.trasa.io)
function getRootDomainFromOrigin(domain) {
  let split = domain.split('.');
  let parsed = split.slice(Math.max(split.length - 3, 0));
  let res = parsed.join('.');
  return res;
}

async function getWsHost() {
  let val = await browser.storage.local.get();
  if (val.wsPath) {
    return val.wsPath;
  } else {
    return 'wss://app.trasa.io';
  }
}

async function isLoggedIn() {
  let val = await browser.storage.local.get();
  if (val.loggedIn === undefined || typeof val.loggedIn === 'undefined') {
    return false;
  } else {
    return val.loggedIn;
  }
}

async function trasaDaCom() {
  let val = await browser.storage.local.get();
  if (
    val.trasaDACom === false ||
    val.trasaDACom === undefined ||
    typeof val.trasaDACom === 'undefined'
  ) {
    return false;
  } else {
    return true;
  }
}
