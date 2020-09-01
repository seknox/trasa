/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 *
 * @format
 * @flow strict-local
 */

import React from 'react';

import {StyleProvider} from 'native-base';
import commonColor from './native-base-theme/variables/commonColor';
import RootNavigator from './pages/Overview/RootNavigator';
import 'react-native-gesture-handler';
import getTheme from './native-base-theme/components';


const App: () => React$Node = () => {
  return (
    <>
      <StyleProvider style={getTheme(commonColor)}>
          <RootNavigator />
      </StyleProvider>
    </>
  );
};

export default App;
