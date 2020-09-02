//'use strict';
import React, {Component} from 'react';
import {ActivityIndicator, StyleSheet, View} from 'react-native';
import {RNCamera} from 'react-native-camera';
import parse from 'url-parse';
import {
  Body,
  Button,
  Container,
  Form,
  H2,
  Header,
  Icon,
  Input,
  Item,
  Left,
  Right,
  Title,
    Text,
} from 'native-base';
import messaging from '@react-native-firebase/messaging';
import axios from 'axios';
import Constants from '../Extra/Constants';
//import RNSecureKeyStore from "react-native-secure-key-store";
import * as SecureStore from 'expo-secure-store';
//import AsyncStorage from "@react-native-community/async-storage";
//import QRCodeScanner from 'react-native-qrcode-scanner';
import BarcodeMask from 'react-native-barcode-mask';
//import DeviceInfo from 'react-native-device-info';
import DeviceInfo from '../Device/Device';
import AsyncStorage from '@react-native-community/async-storage';
import {CommonActions} from '@react-navigation/native';

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'column',
    //    backgroundColor: 'black',
  },
  preview: {
    flex: 1,
    justifyContent: 'flex-end',
    alignItems: 'center',
  },
  camera: {
    flex: 1,
    justifyContent: 'flex-end',
    alignItems: 'center',
  },
  capture: {
    flex: 0,
    //     backgroundColor: '#fff',
    borderRadius: 5,
    padding: 15,
    paddingHorizontal: 20,
    alignSelf: 'center',
    margin: 20,
  },
});
class TotpScan extends Component {
  constructor(props) {
    super(props);
    this.state = {loading: false,done:false};
  }

  componentWillUnmount() {
    this.setState({done:false});
  }

  update = async (
    item, //Callback function when barcode is read from ScanPage
  ) => {
    let totpList = JSON.parse(await AsyncStorage.getItem('TOTPlist'));
    totpList = totpList || [];
    var newItem = true;
    totpList.forEach((it) => {
      if (it.secret === item.secret) {
        newItem = false;
      }
    });
    if (newItem) {
      totpList = totpList.concat(item);

      AsyncStorage.setItem('TOTPlist', JSON.stringify(totpList)).then(a=>{
        this.setState({done:false});
        this.props.navigation.navigate('Home',{updated:totpList});
      });
    }
    this.setState({done:false});


  };

  barCodeRead = (d) => {

    try {
      if ( !d.data) {
        return;
      }
      this.setState({done:true});
      // const data = d.barcodes[0].data;
      const data = d.data;

      // const { dispatch } = this.props;
      const {protocol, host, pathname, query} = parse(data, true);

      const name = decodeURIComponent(pathname.replace(/^\//, ''));
      //const secret = decode(query.secret);
      const secret = query.secret || '';
      const issuer = query.issuer || '';
      const type = query.trasaType || 'public';

      if (protocol === 'otpauth:' && host === 'totp') {
        //Normal TOTP add

        this.update({
          account: name,
          secret: secret,
          issuer: issuer,
          type: type,
        });
       // this.props.navigation.navigate('Home');
      } else if (protocol === 'mobileauth:') {
        // this.setState({done:true});
        //Device Enroll
        this.enrollDevice(query,secret,name,issuer,type)
      }else {
        this.setState({done:false});

      }
    } catch (e) {
      console.log(e);
      this.setState({done:false});

      //console.log(this.props);
      this.setState({loading: false});
      alert('Invalid Code');
    }
  };





  getPublicKey=async ()=>{
    const publicKey=await SecureStore.getItemAsync("PUBLIC_KEY")
    if(!publicKey){

    }
  }


  enrollDevice=async (query,secret,name,issuer,type)=>{
    const urlfromStorage=await AsyncStorage.getItem('TRASA_URL');
    //console.log("urlfromStorage",urlfromStorage)
    let trasaURL = query.trasaURL || urlfromStorage || 'https://trasa.seknox.com';
    Constants.hostname=trasaURL
    this.setState({loading: true});
    const deviceID = query.deviceID || '';
    //AsyncStorage.getItem("FCMToken").then(fcmToken=>{

    messaging()
        .getToken()
        .then((fcmToken) => {
          if (!fcmToken) {
            alert('FCM Token not generated');
          } else {
            DeviceInfo().then((info) => {
              SecureStore.getItemAsync("PUBLIC_KEY").then(publicKey=>{
                axios({
                  method: 'post',
                  url: Constants.hostname + '/api/v1/passmydevicedetail',
                  data: {
                    deviceId: deviceID,
                    fcmToken: fcmToken,
                    publicKey: publicKey || '',
                    deviceFinger: JSON.stringify(info),
                  },
                })
                    .then((response) => {
                      // console.log(response)
                      if (response.data.status === 'success') {
                        Constants.hostname = trasaURL;
                        AsyncStorage.setItem('TRASA_URL', trasaURL);
                        SecureStore.setItemAsync('deviceId', deviceID)
                            .then(() => {
                              this.setState({loading: false});
                              this.update({
                                account: name,
                                secret: secret,
                                issuer: issuer,
                                type: type,
                              });

                            })
                            .catch((reason) => {
                              this.setState({done:false});
                              this.setState({loading: false});
                              alert(reason);
                            });
                      } else {
                        this.setState({loading: false});
                        this.setState({done:false});
                        this.props.navigation.navigate('Home');
                      }
                    })
                    .catch((reason) => {
                      this.setState({done:false});
                      console.error(reason)
                      this.setState({loading: false});
                      alert(reason);
                    });
              });

            })







          }
        });
  }



  render() {
    // this.flag = true;
    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  this.props.navigation.goBack();
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title style={{color:'white'}}>{this.state.done?"Saving":"Scanning"}</Title>
          </Body>
          <Right />
        </Header>

        {this.state.loading ? (
          <View>
            <H2>Please Wait</H2>
            <ActivityIndicator size="large" color="#0000ff" animating={true} />
          </View>
        ) : (
          <RNCamera
            style={styles.camera}
            onBarCodeRead={!this.state.done?this.barCodeRead:null}
          //  onGoogleVisionBarcodesDetected={this.barCodeRead}
            autoFocusPointOfInterest={{x: 0.5, y: 0.5}}
            captureAudio={false}
            googleVisionBarcodeType={
              RNCamera.Constants.GoogleVisionBarcodeDetection.BarcodeType
                .QR_CODE
            }
            //barCodeTypes={[RNCamera.Constants.BarCodeType.qr]}

            androidCameraPermissionOptions={{
              title: 'Permission to use camera',
              message: 'We need your permission to use your camera',
              buttonPositive: 'Ok',
              buttonNegative: 'Cancel',
            }}
            androidRecordAudioPermissionOptions={{
              title: 'Permission to use audio recording',
              message: 'We need your permission to use your audio',
              buttonPositive: 'Ok',
              buttonNegative: 'Cancel',
            }}>
            <BarcodeMask />
          </RNCamera>
        )}
      </Container>
    );
  }
}

export default TotpScan;

//Manual code input page
export class CodeInput extends Component {
  constructor(props) {
    super(props);
    this.state = {issuer: '', account: '', secret: ''};
  }

  update = async (
    item, //Callback function when barcode is read from ScanPage
  ) => {
    let totpList = JSON.parse(await AsyncStorage.getItem('TOTPlist'));
    totpList = totpList || [];
    var newItem = true;
    totpList.forEach((it) => {
      if (it.secret === item.secret) {
        newItem = false;
      }
    });
    if (newItem) {
      totpList = totpList.concat(item);
      AsyncStorage.setItem('TOTPlist', JSON.stringify(totpList))
        .then((r) => {
          this.setState({done:false});
          this.props.navigation.navigate('Home', {
            updated: this.state.name + this.state.issuer,
          });
        })
        .catch((error) => {
          this.setState({done:false});
          console.log(error);
          alert(error);
        });
    }
  };

  onAdd = () => {
    this.update({
      name: this.state.account,
      secret: this.state.secret,
      issuer: this.state.issuer,
    });

  };

  render() {
    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  this.props.navigation.navigate('Home');
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title style={{color:'white'}}>Scan</Title>
          </Body>
          <Right />
        </Header>
        <View>
          <Form>
            <Item fixedLabel>
              <Input
                placeholder={'Issuer'}
                onChangeText={(value) => {
                  this.setState({issuer: value});
                }}
              />
            </Item>
            <Item fixedLabel>
              <Input
                placeholder={'Account'}
                onChangeText={(value) => {
                  this.setState({account: value});
                }}
              />
            </Item>
            <Item fixedLabel>
              <Input
                autoCapitalize={'characters'}
                placeholder={'Secret'}
                onChangeText={(value) => {
                  this.setState({secret: value});
                }}
              />
            </Item>
  <Button primary large onPress={this.onAdd} full>
    <Text style={{color:'white',fontFamily:'Rajdhani-SemiBold'}}> Add </Text>
  </Button>



          </Form>
        </View>
      </Container>
    );
  }
}
