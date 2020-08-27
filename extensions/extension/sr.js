// check for active domains.
// this function is used to sync active user session.
// a list of active domains that matches with root domain is kept in an array
// this function fires in every 10 seconds, collects domains from open tabs and matches with
// stored session store. If local session store has domain not found in open tabs,
// extension initiates logout session process and removes sessionID and hostname for sessionStore.

// add tab listener
// browser.tabs.onRemoved.addListener(onTablClosed);

// async function onTablClosed() {
//   console.log("tab is closed:");
//   var querying = browser.tabs.query({});
//   querying.then(ManageHostlist, onError);
// }

window.setInterval(function () {
  // console.log('not callsed')
  var querying = browser.tabs.query({});
  querying.then(ManageHostlist, onTabsError);
}, 5000);

var srMap = new Map();
var srCounter = new Map();

var hostlist = [];

// ManageHostlist implements host session management.
// this function is used to sync active user session.
// a list of active domains that matches with root domain is kept in an array
// this function is called on tabs.onRemoved event. collects domains from open tabs and matches with
// stored session store. If local session store has domain not found in open tabs,
// extension initiates logout session process and removes sessionID and hostname for sessionStore.
async function ManageHostlist(tabs) {
  // get available domains from active tabs
  var parsedDomainFromTab = [];
  for (let tab of tabs) {
    let domain = parseUri(tab.url).host;
    parsedDomainFromTab.push(domain);
    // console.log('parsed domains from tabs: ', parsedDomainFromTab, ' domain is: ', domain)
  }

  // get authorized hosts from local storage
  let gwHosts = await getHosts();

  // find open tabs and push domain name (domains which exists in getHosts) to hostlist if
  // it does not already exist
  for (let domain of parsedDomainFromTab) {
    if (hostlist.includes(domain) === false) {
      // let gwHosts =  await getHosts()
      if (gwHosts.includes(domain)) {
        hostlist.push(domain);
      }
    }
  }

  // reverse above function now to check if contained domain hostlist
  // is not found in any of tabs. If true, this function should also remove
  // sessionID and Csrf token from local storage for this domain
  for (let hst of hostlist) {
    if (parsedDomainFromTab.includes(hst) === false) {
      let sessionData = await getSessionID();
      let sessionValue = sessionData.sessionStore[hst];
      if (sessionValue || sessionValue !== undefined) {
        let sv = sessionValue.split(':');
        var formData = new FormData();
        formData.append('sid', sv[0]);
        for (let [k, v] of srMap) {
          formData.append(k, v);
        }
        //   // empty srMap values with empty string
        for (let [k, v] of srMap) {
          srMap.set(k, '');
        }

        if (srMap.size === 0) {
          return;
        }

        let host = await getTrasacore();
        let url = host + '/api/woa/proxy/http/rmsession';
        axios
          .post(url, formData, { 'Content-Type': 'multipart/form-data' })
          .catch(function (response) {
            console.log(response);
          });

        // Delete Sr stores
        srMap.delete(sv[0]);
        srCounter.delete(sv[0]);

        // delete session data
        delete sessionData.sessionStore[hst];
        browser.storage.local.set({ sessionStore: sessionData.sessionStore }).then((s, e) => {
          // console.log(s, e)
          hostlist.splice(hostlist.indexOf(hst), 1);
          parsedDomainFromTab.splice(parsedDomainFromTab.indexOf(hst), 1);
        });
        let sid = await getSessionID();
      }

      // }
    }
  }
}

// capturing visible tab on each half a second
window.setInterval(function () {
  browser.tabs.query({ active: true }).then((tabs) => {
    let urlval = tabs[0].url;
    let domain = parseUri(urlval).host;
    //console.log("DOMAIN", domain)
    browser.storage.local.get('sessionStore').then((s) => {
      // console.log(s.sessionStore)
      // console.log(s.sessionStore !== undefined)
      if (s.sessionStore !== undefined) {
        let sessionValue = s.sessionStore[domain];
        if (sessionValue) {
          let session = sessionValue.split(':');

          if (session[2] === 'true') {
            // console.log("capture: ",session[2])
            let capturing = browser.tabs.captureVisibleTab();
            capturing.then(onCaptured, onCaptureError);
          }
        }
      }
    });
  }, onTabsError);
}, 500);

function onTabsError(e) {
  console.error('onTabsError: ', e);
}

function onCaptureError(e) {
  console.error('onCaptureError: ', e);
}
// onCaptures we store image uri in ur session map.
// session map should be picked up by another function for ssr delivery.
async function onCaptured(imageUri) {
  let currentDomain = await browser.tabs.query({ active: true }).then((tabs) => {
    let urlval = tabs[0].url;
    let domain = parseUri(urlval).host;
    return domain;
  }, onTabsError);

  let gwHosts = await getHosts();
  if (gwHosts.includes(currentDomain)) {
    let sessionData = await getSessionID();
    let sessionValue = sessionData.sessionStore[currentDomain];
    if (sessionValue) {
      // get session value
      let session = sessionValue.split(':');

      // trim image data
      let img = imageUri.split(',')[1];
      let ctr = srCounter.get(session[0]);
      if (!ctr) {
        ctr = 1;
        srCounter.set(session[0], 1);
      }
      let counterStr = [ctr, img];
      let imgWithID = counterStr.join(':');
      // get map data based on session ID:
      // console.log(srMap);
      let mapValue = srMap.get(session[0]);
      if (typeof mapValue !== 'undefined') {
        let newValue = mapValue.concat('::', imgWithID);

        // set updated map value
        srMap.set(session[0], newValue);

        //let newCounter = ctr + 1;
        srCounter.set(session[0], ctr + 1);
      }
    }
  }
}

// send all avialable sr captures to trasacore
// called each 5 seconds interval.
window.setInterval(async function () {
  var formData = new FormData();
  let extID = await getExtID();
  // formData.append("extID", extID);
  for (let [k, v] of srMap) {
    formData.append(k, v);
  }
  //   // empty srMap values with empty string
  for (let [k, v] of srMap) {
    srMap.set(k, '');
  }

  if (srMap.size === 0) {
    return;
  }

  if (hostlist.length === 0) {
    return;
  }

  let host = await getTrasacore();
  let url = host + '/api/woa/proxy/http/getsession';
  axios.post(url, formData, { 'Content-Type': 'multipart/form-data' }).catch(function (response) {
    console.error(response);
  });
}, 5000);
