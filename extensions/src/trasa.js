import React from 'react';
import ReactDOM from 'react-dom';
import LoggedIn from './loggedin';
import LoginPage from './login';

//import TRASALogo from './assets/trasa-ni.svg'

//var browser = require("webextension-polyfill");

// window.browser = (function () {
//   return window.msBrowser ||
//     window.browser ||
//     window.chrome;
// })();

class Trasa extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loggedIn: false,
    };
  }

  async componentDidMount() {
    let token = await browser.storage.local.get();
    //console.log('tokens: ', token)
    if (token.extID) {
      this.setState({ loggedIn: true });
    }
  }

  setLoginTrue = () => {
    this.setState({ loggedIn: true });
  };

  click = () => {
    this.setState({ loggedIn: !this.state.loggedIn });
  };
  render() {
    return (
      <div>
        {/* <img src={TRASALogo} alt="trasa-logo" /> */}
        {this.state.loggedIn ? <LoggedIn /> : <LoginPage setLoginTrue={this.setLoginTrue} />}
      </div>
    );
  }
}

ReactDOM.render(<Trasa />, document.getElementById('app'));
