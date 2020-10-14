import React, {Component} from 'react';
import {
  Body,
  Button,
  Container,
  Content,
  Header,
  Icon,
  Left,
  List,
  ListItem,
  Right,
  Separator,
  Switch,
  Text,
} from 'native-base';
import {Image, StatusBar} from 'react-native'
//import Icon from 'react-native-vector-icons/dist/MaterialIcons';
import {Linking, ImageBackground} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import RNFS from 'react-native-fs';
//import { DocumentPicker, DocumentPickerUtil } from 'react-native-document-picker';
import DocumentPicker from 'react-native-document-picker';
import Permissions from 'react-native-permissions';
import {CommonActions} from '@react-navigation/native';
import {DrawerContentScrollView} from '@react-navigation/drawer';

import * as SecureStore from 'expo-secure-store';
//const SeknoxIcon= createIconSetFromFontello(seknoxIconConig);
import Prompt from 'react-native-input-prompt';
import Constants from '../../Extra/Constants';
//Drawer menu
export default class Settings extends Component {
  constructor() {
    super();
    this.state = {
      cloudSync: false,
      floCode: false,
      trasaURL: '',
      promptVisible: false,
    };
  }

  async componentDidMount() {
    let pinEnabled = await AsyncStorage.getItem('TRASA_PIN_ENABLED');

    if (pinEnabled === 'TRUE') {
      this.setState({pinEnabled: true});
    } else {
      this.setState({pinEnabled: false});
    }

    let url = await AsyncStorage.getItem('TRASA_URL');
    //console.log(url,"__________")
    this.setState({trasaURL: url});
  }

  changePinEnabledStatus = async (status) => {
    let pin = await SecureStore.getItemAsync('APP_PIN');
    //console.log(pin)
    if (pin) {
      this.setState({pinEnabled: status});
      await AsyncStorage.setItem(
        'TRASA_PIN_ENABLED',
        status ? 'TRUE' : 'FALSE',
      );
    } else {
      alert('Pin not set');
    }
  };

  changePin = () => {
    this.props.navigation.navigate('ChangePin');
  };

  //Check storage permission
  checkReadPermission = async () => {
    //return true
    const response = await Permissions.check(
      Permissions.PERMISSIONS.ANDROID.READ_EXTERNAL_STORAGE,
    );
    // Response is one of: 'authorized', 'denied', 'restricted', or 'undetermined'
    if (response !== Permissions.RESULTS.GRANTED) {
      const res = await Permissions.request(
        Permissions.PERMISSIONS.ANDROID.READ_EXTERNAL_STORAGE,
      );
      if (res === Permissions.RESULTS.GRANTED) {
        return true;
      } else {
        return false;
      }
    } else {
      return true;
    }
  };

  checkWritePermission = async () => {
    //return true
    const response = await Permissions.check(
      Permissions.PERMISSIONS.ANDROID.WRITE_EXTERNAL_STORAGE,
    );
    // Response is one of: 'authorized', 'denied', 'restricted', or 'undetermined'
    if (response !== Permissions.RESULTS.GRANTED) {
      const res = await Permissions.request(
        Permissions.PERMISSIONS.ANDROID.WRITE_EXTERNAL_STORAGE,
      );
      if (res === Permissions.RESULTS.GRANTED) {
        return true;
      } else {
        return false;
      }
    } else {
      return true;
    }
  };

  importCodes = async () => {
    if (await this.checkReadPermission()) {
      DocumentPicker.pick({
        type: [DocumentPicker.types.allFiles],
      })
        .then((res) => {
          RNFS.readFile(res.uri, 'utf8')
            .then((a) => {
              //Combine codes from file and asyncstorage
              let jsonList = JSON.parse(a);
              if (
                jsonList[0].secret !== undefined &&
                jsonList[0].secret !== null
              ) {
                AsyncStorage.getItem('TOTPlist').then((b) => {
                  if (b !== null) {
                    b = JSON.parse(b);
                    jsonList = jsonList.concat(b);
                  }
                  let strList = JSON.stringify(jsonList);
                  AsyncStorage.setItem('TOTPlist', strList).then(() => {
                    alert('Imported');

                    //this.props.navigation.setParams({totpList: jsonList})
                    const setParamsAction = CommonActions.setParams({
                      params: {totpList: jsonList},
                      key: 'Home',
                    });
                    this.props.navigation.dispatch(setParamsAction);
                    //this.onFloatCode(this.props.floatCodes) //Tiny Hack to remount the component
                  });
                });
              }
            })
            .catch((reason) => {
              alert(reason);
            });
        })
        .catch((e) => {
          alert('Invalid File');
        });
    }
  };
  exportCodes = async () => {
    if (await this.checkWritePermission() && await this.checkReadPermission()) {
      //console.log(RNFS.ExternalStorageDirectoryPath+"/Trasa/",RNFS.exists(RNFS.ExternalStorageDirectoryPath +"/Trasa/"))
      RNFS.mkdir(RNFS.ExternalStorageDirectoryPath + '/Trasa/')
        .then((result) => {
          //console.log('result', result);
          const path =
              RNFS.ExternalStorageDirectoryPath +
              '/Trasa/trasaTotp' +
              new Date().toISOString().replace(/:/g,"-") +
              '.trasa';


          AsyncStorage.getItem('TOTPlist').then((list) => {
            RNFS.writeFile(path, list, 'utf8')
                .then(() => {
                  alert('Exported to ' + path);
                })
                .catch((reason) => {
                  console.error(reason)
                  alert(reason);
                });
          });

        })
        .catch((err) => {
          console.log('err', err);
        });

      //var path = RNFS.DocumentDirectoryPath +"/"+Date.now().toString()+ 'test.txt';
    }
  };

  render() {
    return (
      <Container>
        {/*<StatusBar barStyle = {"dark-content"} hidden = {false} backgroundColor = "#00BCD4" translucent = {false}/>*/}
        {/*/!*<ImageBackground*!/*/}
        {/*/!*  source={require('../Images/logo-blue-2.png')}*!/*/}
        {/*/!*  style={{width: null, height: 56}}*!/*/}
        {/*/!*>*!/*/}
        {/*  <Header style={{backgroundColor: 'transparent'}}>*/}
        {/*    <Image source={require('../Images/logo-blue-2.png')} style={{width:300}}/>*/}
        {/*    /!*<Body>*!/*/}
        {/*    /!*    <Title>TRASA</Title>*!/*/}
        {/*    /!*    <Subtitle>Trusted Access and Session Analytics</Subtitle>*!/*/}
        {/*    /!*</Body>*!/*/}
        {/*  </Header>*/}
        {/*/!*</ImageBackground>*!/*/}

        <DrawerContentScrollView {...this.props}>
          <Image source={require('../Images/logo-blue-2.png')} resizeMethod={"scale"} height={50} width={10} style={{height:50,width:"100%"}}/>
          <List>
            <Separator bordered>
              <Text>Preferences</Text>
            </Separator>

            <ListItem icon onPress={this.exportCodes}>
              <Left>
                <Button>
                  <Icon
                    type={'MaterialIcons'}
                    name={'unarchive'}
                    style={{fontSize: 30}}
                  />
                </Button>
              </Left>
              <Body>
                <Text>Export Secrets</Text>
              </Body>
            </ListItem>

            <ListItem icon onPress={this.importCodes}>
              <Left>
                <Button>
                  <Icon
                    type={'MaterialIcons'}
                    name={'archive'}
                    style={{fontSize: 30}}
                  />
                </Button>
              </Left>
              <Body>
                <Text>Import Secrets</Text>
              </Body>
            </ListItem>



            <ListItem icon>
              <Left>
                <Button>
                  <Icon active type={'MaterialIcons'} name="vpn-key" />
                </Button>
              </Left>
              <Body>
                <Text>Enable Pin</Text>
              </Body>
              <Right>
                <Switch
                  value={this.state.pinEnabled}
                  onValueChange={this.changePinEnabledStatus}
                />
              </Right>
            </ListItem>

            <ListItem icon onPress={this.changePin}>
              <Left>
                <Button>
                  <Icon type={'MaterialIcons'} name={'vpn-key'} />
                </Button>
              </Left>
              <Body>
                <Text>Change Pin</Text>
              </Body>
            </ListItem>

            <ListItem
              icon
              onPress={() => {
                this.setState({promptVisible: true});
              }}>
              <Left>
                <Button>
                  <Icon type={'MaterialIcons'} name={'link'} />
                </Button>
              </Left>
              <Body>
                <Text>Change Trasa URL</Text>
              </Body>
            </ListItem>

            <Separator bordered>
              <Text>Info</Text>
            </Separator>

            <ListItem
              icon
              onPress={() => {
                Linking.openURL('https://trasa.io/docs/');
                //this.props.navigation.navigate("Help")
              }}>
              <Left>
                <Button>
                  <Icon type={'Entypo'} name={'help'} />
                </Button>
              </Left>
              <Body>
                <Text>Help</Text>
              </Body>
              <Right>
                <Icon name={'arrow-forward'} />
              </Right>
            </ListItem>

            <ListItem
              icon
              onPress={() => {
                Linking.openURL('https://trasa.io/');
                //this.props.navigation.navigate("AboutTrasa")
              }}>
              <Left>
                <Button>
                  <Icon name={'information-circle'} />
                </Button>
              </Left>
              <Body>
                <Text>About Trasa</Text>
              </Body>
              <Right>
                <Icon name={'arrow-forward'} />
              </Right>
            </ListItem>

            <ListItem
              icon
              onPress={() => {
                Linking.openURL('https://seknox.com/');
                //this.props.navigation.navigate("AboutSeknox")
              }}>
              <Left>
                <Button>
                  <Icon name={'information-circle'} />
                </Button>
              </Left>
              <Body>
                <Text>About Seknox</Text>
              </Body>
              <Right>
                <Icon name={'arrow-forward'} />
              </Right>
            </ListItem>

            {/*<Separator bordered >
                            <Text>Account</Text>
                        </Separator>

                        <ListItem icon onPress={()=>{
                            RNSecureKeyStore.set("deviceId","")
                            this.props.screenProps.rootNav.navigate("Auth")
                        }}>
                            <Left>
                                <Button>
                                <Icon  type={"Ionicons"} name={"log-out"} />
                                </Button>
                            </Left>
                            <Body>
                            <Text>Logout</Text>
                            </Body>

                        </ListItem>*/}
          </List>
          <Prompt
            title={'Enter Trasa URL'}

            placeholder={this.state.trasaURL || ''}
            defaultValue={this.state.trasaURL}
            visible={this.state.promptVisible}
            onCancel={() =>
              this.setState({
                promptVisible: false,
              })
            }
            onSubmit={(value) => {
              this.setState({promptVisible: false, trasaURL: value});
              Constants.hostname = value;
              AsyncStorage.setItem('TRASA_URL', value);
            }}
          />
        </DrawerContentScrollView>
      </Container>
    );
  }



}
