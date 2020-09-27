// postMessage listener
window.addEventListener('message', function (event) {
  if (event.source == window && event.data && event.data.direction == 'tsxdashboard') {
    // console.log('received: ',event.origin )
    browser.runtime.sendMessage({
      origin: event.origin,
      type: event.data.message.type,
      data: event.data.message.data,
    });
    browser.runtime.onMessage.addListener((msg) => {
      // console.log('sending response: ', msg)
      if (msg) {
        window.postMessage(
          {
            direction: 'trasaExt',
            message: msg.data,
          },
          msg.domain,
        );
      }
    });
  }
});
