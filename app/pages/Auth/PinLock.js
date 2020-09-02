import Icon from 'react-native-vector-icons/Ionicons';
import React, {useEffect, useRef, useState} from 'react';
import {ImageBackground, SafeAreaView, StatusBar, Text} from 'react-native';
import ReactNativePinView from 'react-native-pin-view';
import * as SecureStore from 'expo-secure-store';

const PinLock = (props) => {
  const pinView = useRef(null);
  const [showRemoveButton, setShowRemoveButton] = useState(false);
  const [enteredPin, setEnteredPin] = useState('');
  const [showCompletedButton, setShowCompletedButton] = useState(false);
  useEffect(() => {
    return cleanUp;
  }, []);

  const cleanUp = () => {
    pinView.current.clearAll();
  };

  const checkPin = async (value) => {
    let pin = await SecureStore.getItemAsync('APP_PIN');
    //console.log(pin);
    if (value == pin) {
      props.unlock();
    } else {
      alert('Invalid Pin');
      pinView.current.clearAll();
    }
  };
  return (
    <>
      <StatusBar barStyle="light-content" />
      <SafeAreaView
        style={{
          flex: 1,
          backgroundColor: 'rgba(0,0,0,0.5)',
          justifyContent: 'center',
          alignItems: 'center',
        }}>
        <Text
          style={{
            paddingTop: 24,
            paddingBottom: 48,
            color: 'rgba(255,255,255,0.7)',
            fontSize: 48,
          }}>
          TRASA
        </Text>
        <ReactNativePinView
          inputSize={32}
          ref={pinView}
          pinLength={4}
          buttonSize={60}
          onValueChange={(value) => {
            setEnteredPin(value);
            if (value.length === 4) {
              checkPin(value);
            } else {
              setShowCompletedButton(false);
            }
          }}
          buttonAreaStyle={{
            marginTop: 24,
          }}
          inputAreaStyle={{
            marginBottom: 24,
          }}
          inputViewEmptyStyle={{
            backgroundColor: 'transparent',
            borderWidth: 1,
            borderColor: '#FFF',
          }}
          inputViewFilledStyle={{
            backgroundColor: '#FFF',
          }}
          buttonViewStyle={{
            borderWidth: 1,
            borderColor: '#FFF',
          }}
          buttonTextStyle={{
            color: '#FFF',
          }}
          onButtonPress={(key) => {
            if (key === 'custom_left') {
              pinView.current.clear();
            }
          }}
          customLeftButton={
            <Icon name={'ios-backspace'} size={36} color={'#FFF'} />
          }
        />
      </SafeAreaView>
    </>
  );
};
export default PinLock;
