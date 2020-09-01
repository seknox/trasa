import Icon from 'react-native-vector-icons/Ionicons';
import React, { useEffect, useRef, useState } from 'react';
import { SafeAreaView, StatusBar, Text } from 'react-native';
import ReactNativePinView from 'react-native-pin-view';
import * as SecureStore from 'expo-secure-store';

const ChangePin = (props) => {
  const pinView = useRef(null);
  const [showRemoveButton, setShowRemoveButton] = useState(false);
  const [enteredPin, setEnteredPin] = useState('');
  const [instruction, setInstruction] = useState('Enter Old Pin');
  const [step, setStep] = useState(1);
  const [newPin, setNewPin] = useState('');

  const [showCompletedButton, setShowCompletedButton] = useState(false);

  const cleanUp = () => {
    pinView.current.clearAll();
    setStep(0);
    setNewPin('');
    setEnteredPin('');
    setInstruction('Enter Old Pin');
    pinView.current.clearAll();
  };

  useEffect(async () => {
    cleanUp();
    const pin = await SecureStore.getItemAsync('APP_PIN');
    //  console.log(pin)
    if (!pin) {
      setStep(1);
      setInstruction('Enter New Pin');
    }
    return cleanUp;
  }, []);

  const checkPin = async (value) => {
    switch (step) {
      case 0:
        const pin = await SecureStore.getItemAsync('APP_PIN');
        if (value === pin) {
          setNewPin(value);
          setStep(1);
          setInstruction('Enter New Pin');
        } else {
          setNewPin(value);
          alert('Invalid Pin');
        }
        setEnteredPin('');
        pinView.current.clearAll();
        break;

      case 1:
        setNewPin(value);
        setStep(2);
        setEnteredPin('');
        setInstruction('Confirm New Pin');
        pinView.current.clearAll();
        break;
      case 2:
        if (newPin === value) {
          await SecureStore.setItemAsync('APP_PIN', newPin);
          cleanUp();
          alert('Pin Changed');
          props.navigation.goBack();
        } else {
          setStep(1);
          setNewPin('');
          setEnteredPin('');
          setInstruction('Try Again');
        }
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
        }}
      >
        <Text
          style={{
            paddingTop: 24,
            paddingBottom: 48,
            color: 'rgba(255,255,255,0.7)',
            fontSize: 48,
          }}
        >
          TRASA
        </Text>
        <Text>{instruction}</Text>
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
            showRemoveButton ? <Icon name="ios-backspace" size={36} color="#FFF" /> : undefined
          }
          customRightButton={
            showCompletedButton ? <Icon name="ios-unlock" size={36} color="#FFF" /> : undefined
          }
        />
      </SafeAreaView>
    </>
  );
};
export default ChangePin;
