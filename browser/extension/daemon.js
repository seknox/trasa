var totpCode = '';

async function cancelRequestsWithoutTRASAHeader(requestDetails) {
  let hostHost = parseUri(requestDetails.url).host;
  let gwHosts = await getHosts();
  if (gwHosts.includes(hostHost)) {
    let sessionData = await getSessionID();
    let sessionValue = sessionData.sessionStore[hostHost];

    if (sessionValue === undefined) {
      // if (hostPath === '/favicon.ico'){
      //   return {cancel: true};
      // }

      let t = new Date();

      //console.log('sending init: ', requestDetails.url)
      let resp = await initHttpRequest(hostHost, '');
      if (resp != false) {
        return { cancel: false };
      }

      return { cancel: true };
    } else {
      //console.log("allowing request for(after): ", requestDetails.url);
      //console.log("allowing request for(after): ", sessionValue);
      return { cancel: false };
    }
  }
}

// addTRASAHeader adds trasa headers to outgoing requests if domain matches root domain.
// if the header value is undefined, we initiate 2FA request which also starts new session.
// It is the responsibility of extension to continuously monitor open tabs and look for domains to create
// session pool. if the tab is closed, extension shold send logout session request and delete the local session
// value from storage.
async function addTRASAHeader(e) {
  // Loop through all available request headers
  for (var header of e.requestHeaders) {
    if (header.name.toLowerCase() === 'host') {
      let hostHost = parseUri(header.value).host;
      for (var originHeader of e.requestHeaders) {
        if (originHeader.name.toLowerCase() === 'origin') {
          hostHost = parseUri(originHeader.value).host;
        }
      }

      // we only add session token if the origin matches hostlists returned from getHosts.
      // If not it is meant for domains not controlled by TRASA.
      let gwHosts = await getHosts();

      //console.log('req host: ', hostHost, ' available hosts: ', gwHosts )
      if (gwHosts.includes(hostHost)) {
        let sessionData = await getSessionID();
        let sessionValue = sessionData.sessionStore[hostHost];

        // check if sessionID is undefined. If true we should initiate initHttpSession request
        if (sessionValue) {
          let session = sessionValue.split(':');
          e.requestHeaders.push({ name: 'TRASA-X-SESSION', value: session[0] });
          e.requestHeaders.push({ name: 'TRASA-X-CSRF', value: session[1] });
        }
      }
    }
  }
  return { requestHeaders: e.requestHeaders };
}

// we modify our requests here.
browser.webRequest.onBeforeSendHeaders.addListener(addTRASAHeader, { urls: ['<all_urls>'] }, [
  'blocking',
  'requestHeaders',
]);

// //listen for requests untill trasa session values are set. we cancel every request untill it has sessions.
// browser.webRequest.onBeforeRequest.addListener(
//   cancelRequestsWithoutTRASAHeader,
//   { urls: ["<all_urls>"] },
//   ["blocking"]
// );

// Get TRASA headers from response
async function getTRASAHeader(e) {
  for (var header of e.responseHeaders) {
    // console.log('######### ' + e.url)
    let host = parseUri(e.url).host;
    // console.log('host value ' + host)

    // If header contains trasa-auth, previosuly we used to send inithttpRequest here as well.
    // this case only happens when if user's stored session values are already expired.
    // So now we only delete user's stored session value.
    if (header.name.toLowerCase() === 'trasa-auth') {
      let sessionData = await getSessionID();
      // console.log('failed host: ', host)
      // delete sessionData.sessionStore[host]
      // browser.storage.local.set({'sessionStore': sessionData.sessionStore})
      // let sv = await getSessionID()
      // hostlist.splice( hostlist.indexOf(host), 1 );
    }
  }
  return { responseHeaders: e.responseHeaders };
}

// Listen for onHeaderReceived for the target page.
// Set "blocking" and "responseHeaders".
// s

///////////////////////////////////////////////////////////////////////////////
////// Communicate with content script   //////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

browser.runtime.onMessage.addListener(csMsgListener);

async function csMsgListener(message) {
  let trustedTrasaDomain = await getTrasacore();
  if (message.origin === trustedTrasaDomain) {
    switch (true) {
      case message.type === 'exportExtToken':
        browser.tabs
          .query({ active: true, currentWindow: true })
          .then(sendextToken)
          .catch(function (error) {
            console.error(`Error to fetch deviceHygiene: ${error}`);
          });
        break;
      case message.type === 'exportDeviceHygiene':
        // console.log('message data: ', message.data)
        browser.tabs
          .query({ active: true, currentWindow: true })
          .then(function (tabs) {
            sendHygiene(tabs, message.data);
          })
          .catch(function (error) {
            console.error(`Error to fetch deviceHygiene: ${error}`);
          });
        break;
      case message.type === 'setSession':
        // set session values
        srMap.set(message.data.sessionID, '');
        srCounter.set(message.data.sessionID, 1);

        let buildSession =
          message.data.sessionID + ':' + message.data.csrfToken + ':' + message.data.sessionRecord;
        sessionStore[message.data.domain] = buildSession;
        await browser.storage.local.set({ sessionStore }).then(function (e) {}, onError);

        let session = await getSessionID();

        // return message
        browser.tabs
          .query({ active: true, currentWindow: true })
          .then(setSession)
          .catch(function (error) {
            console.error(`Error to  set session: ${error}`);
          });

        break;
      default:
        console.log('default reached');
        break;
    }
  }
}

//////////////////////////////////////////////////////
// exportDeviceHygiene

async function sendHygiene(tabs, pubKey) {
  let globalConnect = await trasaDaCom();
  // console.log('global connect state: ', globalConnect)
  // console.log('global connect tabs: ', tabs)
  // console.log('global connect pubKey: ', pubKey)
  if (globalConnect === false) {
    return;
  }

  let d = await GetDeviceHygiene(pubKey);
  let trustedTrasaDomain = await getTrasacore();
  let extID = await getExtID();
  if (!extID) {
    extID = 'notRegistered';
  }
  for (let tab of tabs) {
    browser.tabs
      .sendMessage(tab.id, { data: { extID: extID, dh: d }, domain: trustedTrasaDomain })
      .catch(function (error) {
        console.error(`[Error] sendHygiene: ${error}`);
      });
  }
}

///////////////////////////////////////////////////
// exportExtToken

async function sendextToken(tabs) {
  let extID = await getExtID();
  let trustedTrasaDomain = await getTrasacore();
  for (let tab of tabs) {
    browser.tabs
      .sendMessage(tab.id, { data: extID, domain: trustedTrasaDomain })
      .catch(function (error) {
        console.error(`[Error] sendExtToken: ${error}`);
      });
  }
}

///////////////////////////////////////////////////
// setSession

async function setSession(tabs) {
  let trustedTrasaDomain = await getTrasacore();
  for (let tab of tabs) {
    browser.tabs
      .sendMessage(tab.id, { data: 'done', domain: trustedTrasaDomain })
      .catch(function (error) {
        console.error(`[Error] sendExtToken: ${error}`);
      });
  }
}

////////////////////////////////////////////////////
// http access proxy password fill

async function checkShouldFill(message) {
  // console.log('received checkShouldFill : ', message)
  if (message.shouldFillURL) {
    let domain = parseUri(message.shouldFillURL).host;
    let rootDomainFromOrigin = getSsoDomainFromOrigin(domain);
    let managed = await getSsoDomain();

    // console.log('prninting domains: ', rootDomainFromOrigin, managed)
    if (rootDomainFromOrigin == managed) {
      browser.tabs.query({ active: true, currentWindow: true }).then(shouldFill).catch(onError);
    }
  }
}

function shouldFill(tabs) {
  for (let tab of tabs) {
    browser.tabs.sendMessage(tab.id, { shouldFill: true }).catch(onError);
  }
  // browser.tabs.sendMessage(tabs[0].id, { "shouldFill": true  } ) ;
}

function onError(error) {
  console.error(`Error: ${error}`);
}

//browser.runtime.onMessage.addListener(getTrasaUser);

async function getTrasaUser(message) {
  //  console.log('received getTrasaUser: ', message)

  if (message.trasauser) {
    //   console.log('setting getTrasaUser: ', message.trasauser)
    let usersetter = await browser.storage.local.set({
      trasauser: message.trasauser,
    });
    let trasaUser = await getTrasaUser();
    //   console.log("trasauser after setting : ", trasaUser)
  }
}

//////////////////////////////////////////////
////    Prompt and init Tfa
//////////////////////////////////////////////

// async function promptAndProcessTfa(hostHost) {
//   let r = await browser.tabs.query({active: true, currentWindow: true})
//   .then(async function (tabs) {
//     for (let tab of tabs) {
//       browser.tabs.sendMessage( tab.id, { "promptTfa": true  } )
//     }
//   }).catch(onError)

//   return '.....returning.....'
//}
