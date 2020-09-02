import React, { Component } from 'react';
// import {createDrawerNavigator} from 'react-navigation';
import { createDrawerNavigator } from '@react-navigation/drawer';
import { NavigationContainer } from '@react-navigation/native';
import Settings from './Settings/Settings';
import ScanPage, { CodeInput } from './TotpScan';
import HelpPage from './InfoPages/HelpPage';
import AboutSeknox from './InfoPages/AboutSeknox';
import AboutTrasa from './InfoPages/AboutTrasa';
import TwoFA from '../Extra/2FA';
import OrgPage from './OrgPage';
import { Home } from './Home';
import { Notifications, ViewNotification } from './Notifications';
import SplashScreen from '../SplashScreen';
import ChangePin from './Settings/ChangePin';
// import {NativeRouter, Route, Link} from 'react-router-native';
//
// import getTheme from '../../native-base-theme/components';
// import platform from "../../native-base-theme/variables/platform";
// import {StyleProvider} from "native-base";
// import commonColor from "../../native-base-theme/variables/commonColor";
// Overview page

// const Router = () => {
//   return (
//     <NativeRouter>
//       <Route exact path="/" component={SplashScreen} />
//       <Route path="/HelpPage" component={HelpPage} />
//       <Route path="/AboutSeknox" component={AboutSeknox} />
//       <Route path="/AboutTrasa" component={AboutTrasa} />
//       <Route path="/Home" component={TotpPage} />
//       <Route path="/OrgPage" component={OrgPage} />
//       <Route path="/Notifications" component={Notifications} />
//       <Route path="/ViewNotification" component={ViewNotification} />
//       <Route path="/TwoFA" component={TwoFA} />
//       <Route path="/ScanPage" component={ScanPage} />
//       <Route path="/CodeInput" component={CodeInput} />
//       <Route path="/ChangePin" component={ChangePin} />
//     </NativeRouter>
//   );
// };

// Org specific page
const Drawer = createDrawerNavigator();

function RootNavigator() {
  // https://github.com/react-navigation/react-navigation/issues/7561#issuecomment-636566410
  // this is for frawer flicker issue
  const showDrawer = false;

  return (
    <NavigationContainer>
      <Drawer.Navigator
        drawerContent={(props) => <Settings {...props} />}
        initialRouteName="SplashScreen"
        drawerStyle={{ width: !showDrawer ? null : 280 }}
        //  hideStatusBar={true}
        // drawerType={"front"}
        backBehavior="history"
      >
        <Drawer.Screen
          name="SplashScreen"
          component={SplashScreen}
          options={{ headerShown: false }}
        />
        <Drawer.Screen name="Help" component={HelpPage} options={{ headerShown: false }} />
        <Drawer.Screen
          name="AboutSeknox"
          component={AboutSeknox}
          options={{ headerShown: false }}
        />
        <Drawer.Screen name="AboutTrasa" component={AboutTrasa} options={{ headerShown: false }} />
        <Drawer.Screen name="Home" component={Home} options={{ headerShown: false }} />
        <Drawer.Screen name="OrgPage" component={OrgPage} options={{ headerShown: false }} />
        <Drawer.Screen
          name="Notifications"
          component={Notifications}
          options={{ headerShown: false }}
        />
        <Drawer.Screen
          name="ViewNotification"
          component={ViewNotification}
          options={{ headerShown: false }}
        />
        <Drawer.Screen name="TwoFA" component={TwoFA} options={{ headerShown: false }} />
        <Drawer.Screen name="ScanPage" component={ScanPage} options={{ headerShown: false }} />
        <Drawer.Screen name="CodeInput" component={CodeInput} options={{ headerShown: false }} />
        <Drawer.Screen name="ChangePin" component={ChangePin} options={{ headerShown: false }} />
      </Drawer.Navigator>
    </NavigationContainer>
  );
}

// const RootDrawer = createDrawerNavigator(
//   {
//     SplashScreen: {
//       screen: SplashScreen,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     Help: {
//       screen: HelpPage,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     AboutSeknox: {
//       screen: AboutSeknox,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     AboutTrasa: {
//       screen: AboutTrasa,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     Home: {
//       screen: TotpPage,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     OrgPage: {
//       screen: OrgPage,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     Notifications: {
//       screen: Notifications,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     ViewNotification: {
//       screen: ViewNotification,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     TwoFA: {
//       screen: TwoFA,
//       navigationOptions: {
//         header: null,
//       },
//     },
//
//     ScanPage: {
//       screen: ScanPage,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     CodeInput: {
//       screen: CodeInput,
//       navigationOptions: {
//         header: null,
//       },
//     },
//     ChangePin: {
//       screen: ChangePin,
//       navigationOptions: {
//         header: null,
//       },
//     },
//   },
//   {
//     contentComponent: Settings,
//     initialRouteName: 'SplashScreen',
//     backBehavior: 'history',
//   },
// );

export default RootNavigator;
