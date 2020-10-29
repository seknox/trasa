//'use strict';
import React, {Component, useEffect, useState} from 'react';
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

  camera: {
    flex: 1,
    justifyContent: 'flex-end',
    alignItems: 'center',
  },

});
function TotpScan (props){

  const [loading,setLoading] = useState(false)
  const [done,setDone] = useState(false)


  useEffect(()=>{
    setDone(false)
  },[props.route?.params?.rand])



  const update = async (
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
        props.navigation.navigate('Home',{updated:totpList});
      });
    }


  };

  const barCodeRead = (d) => {

    try {
      if ( !d.data ) {
        return;
      }
      setDone(true)

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

        update({
          account: name,
          secret: secret,
          issuer: issuer,
          type: type,
        });
       // this.props.navigation.navigate('Home');
      } else if (protocol === 'mobileauth:') {
        // this.setState({done:true});
        //Device Enroll
        enrollDevice(query,secret,name,issuer,type)
      }else {

      }
    } catch (e) {
      console.error(e);
      console.info("done:true")

      //console.log(this.props);
      alert('Invalid Code');
    }
  };





  const getPublicKey=async ()=>{
    const publicKey=await SecureStore.getItemAsync("PUBLIC_KEY")
    if(!publicKey){

    }
  }


 const  enrollDevice=async (query,secret,name,issuer,type)=>{
    const urlfromStorage=await AsyncStorage.getItem('TRASA_URL');
    //console.log("urlfromStorage",urlfromStorage)
    let trasaURL = query.trasaURL || urlfromStorage || 'https://trasa.seknox.com';
    Constants.hostname=trasaURL
   setLoading(true);
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
                              setLoading(false);
                              update({
                                account: name,
                                secret: secret,
                                issuer: issuer,
                                type: type,
                              });

                            })
                            .catch((reason) => {
                              console.error("Could not store deviceID")
                              setLoading(false);

                              alert(reason);
                            });
                      } else {
                        console.error("response is not success")
                        setLoading(false);

                        props.navigation.navigate('Home');
                      }
                    })
                    .catch((reason) => {
                      console.error("axios error",reason)
                      setLoading(false);

                      alert(reason);
                    });
              });

            })







          }
        });
  }




    return (
      <Container>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                name={'arrow-back'}
                onPress={() => {
                  props.navigation.goBack();
                }}
              />
            </Button>
          </Left>
          <Body>
            <Title style={{color:'white'}}>{done?"Saving":"Scanning"}</Title>
          </Body>
          <Right />
        </Header>

        {loading ? (
          <View>
            <H2>Please Wait</H2>
            <ActivityIndicator size="large" color="#0000ff" animating={true} />
          </View>
        ) : (
          <RNCamera
            style={styles.camera}
            onBarCodeRead={!done?barCodeRead:null}
            barCodeTypes={[RNCamera.Constants.BarCodeType.qr]}
          //  onGoogleVisionBarcodesDetected={this.barCodeRead}
            autoFocusPointOfInterest={{x: 0.5, y: 0.5}}
            captureAudio={false}
            // googleVisionBarcodeType={
            //   RNCamera.Constants.GoogleVisionBarcodeDetection.BarcodeType
            //     .QR_CODE
            // }
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
          this.props.navigation.navigate('Home', {
            updated: this.state.name + this.state.issuer,
          });
        })
        .catch((error) => {
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
