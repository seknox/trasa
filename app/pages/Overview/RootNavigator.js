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


export default RootNavigator;
