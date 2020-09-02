import React, {Component} from 'react';
import {Alert, StyleSheet, View} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import axios from 'axios';
import Constants from './Extra/Constants';
// import SplashScreen from 'react-native-splash-screen';
import SplashScreen from "react-native-bootsplash";
import messaging from '@react-native-firebase/messaging';
import {storeNotification} from './Overview/NotificationStorage';
import { RSA } from 'react-native-rsa-native';
import * as SecureStore from 'expo-secure-store';
import {CommonActions} from '@react-navigation/native';
import getDeviceInfo from './Device/Device';
import PinLock from './Auth/PinLock';

export default class SplashScreenPage extends Component {
  constructor(props) {
    super(props);
    this.loaded = false;
    this.state = {unlocked: false};
  }

  showTwoFA(notification) {
    this.goTo('TwoFA', {
      appName: notification.appName,
      ipAddr: notification.ipAddr,
      challenge: notification.challenge,
      orgName: notification.orgName,
      time: notification.time,
    });
  }

  async generateToken() {

   // console.info("isDeviceRegisteredForRemoteMessages",isDeviceRegisteredForRemoteMessages)
    messaging()
      .getToken()
      .then((token) => {
        if (token) {
          console.info(token);
          AsyncStorage.setItem('FCMToken', token);
        } else {
          console.error('FCM token not generated');
          this.generateToken();
        }
      })
      .catch((reason) => console.log(reason));
  }

  async initializeFirebase() {
    //await messaging().registerDeviceForRemoteMessages();
  // console.info("2222222")

    //this.generateToken();
    // const channel = new firebase.notifications.Android.Channel('2fa', '2fa', firebase.notifications.Android.Importance.Max)
    //     .setDescription('My apps test channel');
    //
    // // Create the channel
    // firebase.notifications().android.createChannel(channel);

    const authStatus = await messaging().requestPermission();
    const enabled =
      authStatus === messaging.AuthorizationStatus.AUTHORIZED ||
      authStatus === messaging.AuthorizationStatus.PROVISIONAL;

    if (!enabled) {
      console.log('This app needs notification permission to function properly');
    }

    await this.generateToken()

    //});

    this.notificationListener = messaging().onMessage((notification) => {
      if(!notification?.data){
        return
      }
      storeNotification(notification);
      if (
        notification.data.type !== 'adhoc' &&
        notification.data.type !== 'REQUEST'
      ) {
        this.showTwoFA(notification.data);
      }
    });

    this.notificationOpenedListener = messaging().onNotificationOpenedApp(
      (notification) => {

        if(!notification?.data){
          return
        }
        storeNotification(notification);
        if (
          notification.data?.type !== 'GRANT' &&
          notification.data?.type !== 'REQUEST'
        ) {
          this.showTwoFA(notification.data);
        }
      },
    );

    /*    this.onTokenRefreshListener = firebase.messaging().onTokenRefresh(token => {
                //console.log("TOKEN (refreshUnsubscribe)", token);
                AsyncStorage.setItem("FCMToken",token);
                AsyncStorage.setItem("FCMTokenSynced","FALSE");
                //  this.fcmUpdate(token);

            });*/
  }

  /*
        componentWillUnmount() {
            this.notificationOpenedListener();
            this.notificationListener();
        }
    */


  async generateKeypair(){
    const check=await SecureStore.getItemAsync("PRIVATE_KEY")
   // console.log("check if already generated",check)
    if(check){
      return
    }

   // console.log("generating key pair")
    const keys=await RSA.generateKeys(4096)
   // console.log("key pair generated")

    await SecureStore.setItemAsync("PRIVATE_KEY",keys.private)
    await SecureStore.setItemAsync("PUBLIC_KEY",keys.public)
//    console.log("key pair saved")

  }


  async componentDidMount() {
    this.generateKeypair()


    Constants.hostname = await AsyncStorage.getItem('TRASA_URL');
    if(!Constants.hostname){
      Constants.hostname="https://api.trasa.seknox.com"
    }
    // Constants.hostname= "http://192.168.100.91"

    let pinEnabled = await AsyncStorage.getItem('TRASA_PIN_ENABLED');

    if (pinEnabled === 'TRUE') {
      await this.setState({unlocked: false});
    } else {
      await this.setState({unlocked: true});
    }

    await this.initializeFirebase();

    const notification = await messaging().getInitialNotification();


    if (notification?.data) {
      // Get information about the notification that was opened
      //console.log(notification,notificationOpen.action)
      storeNotification(notification);
      if (

        notification.data.type !== 'adhoc' &&
        notification.data.type !== 'REQUEST'
      ) {
        SplashScreen.hide();
        this.showTwoFA(notification.data);
      }else {
        this.goTo('Home');
      }


    } else {
      SplashScreen.hide();

      this.goTo('Home');
    }

   // console.info('close');
    SplashScreen.hide();

    /*let deviceId = RNSecureKeyStore.get("deviceId").then(did=>{
            if(!did)
            {
                this.props.navigation.navigate("Auth")
            }
            else
            {
                this.props.navigation.navigate("Overview")
            }

            SplashScreen.close({
                animationType: SplashScreen.animationType.scale,
                duration: 800,
                delay: 0,
            })



        }).catch(reason=>{
            this.props.navigation.navigate("Auth")
            SplashScreen.close({
                animationType: SplashScreen.animationType.scale,
                duration: 0,
                delay: 0,
            })



        })


    */

    //SplashScreen.close(SplashScreen.animationType.scale, 850, 500)
  }

  // closeSplashScreen() {
  //   SmartSplashScreen.close({
  //     animationType: SmartSplashScreen.animationType.scale,
  //     duration: 850,
  //     delay: 0,
  //   });
  // }

  fcmUpdate(token) {
    SecureStore.getItemAsync('deviceId').then((deviceId) => {
      axios({
        method: 'post',
        //TODO maybe chage this
        url: Constants.hostname + '/api/v1/remote/auth/fcmupdate',
        data: {
          token: token,
          deviceId: deviceId,
        },
      })
        .then((value) => {
          if (value.data.success === 'true') {
            AsyncStorage.setItem('FCMTokenSync', 'TRUE');
          }
        })
        .catch((reason) => Alert.alert('FCM sync error ' + reason));
    });
  }

  syncAccount(list) {
    SecureStore.getItemAsync('deviceId').then((deviceId) => {
      axios({
        method: 'post',
        //TODO maybe change this
        url: Constants.hostname + '/api/v1/remote/auth/synctoken',
        data: {
          list: list,
          deviceId: deviceId,
        },
      })
        .then((value) => {
          if (value.data.success === 'true') {
            AsyncStorage.setItem('accSynced', 'TRUE');
          }
        })
        .catch((reason) => Alert.alert('Account sync error ' + reason));
    });
  }

  goTo = (page, params) => {
    this.setState({page: page, params: params});
    if (this.state.unlocked) {
      this.props.navigation.navigate(page, params);
    }
  };

  unlock = () => {
    this.setState({unlocked: true});
    this.props.navigation.navigate(this.state.page, this.state.params);
  };

  render() {
    return (
      <View style={styles.container}>
        {this.state.unlocked ? null : <PinLock unlock={this.unlock} />}
      </View>
    );
  }
}

const styles = StyleSheet.create({
  icon: {
    width: 24,
    height: 24,
  },
  container: {
    flex: 1,
    //     backgroundColor:'#0e0343', //'rgba(1,1,35,1)',
    alignItems: 'flex-start',
    justifyContent: 'center',
  },
  button: {
    flex: 1,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-around',
  },
});

/*


*/
