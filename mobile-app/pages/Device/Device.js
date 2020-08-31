import React from 'react';
import DeviceInfo from 'react-native-device-info';
import JailMonkey from 'jail-monkey';
import Permissions from 'react-native-permissions';
// import AppList from 'react-native-installed-apps';
import NetInfo, { NetInfoStateType } from '@react-native-community/netinfo';
// import RNAndroidInstalledApps from 'react-native-android-installed-apps';

const getDeviceInfo = async () => {
  await Permissions.request(Permissions.PERMISSIONS.ANDROID.READ_PHONE_STATE);
  await Permissions.request(Permissions.PERMISSIONS.ANDROID.READ_PHONE_NUMBERS);
  // const info={
  //     installedApps:[],//await RNAndroidInstalledApps.getNonSystemApps(),
  //     deviceName:await DeviceInfo.getDeviceName(),
  //     brand:await DeviceInfo.getBrand(),
  //     manufacturer:await DeviceInfo.getManufacturer(),
  //     osName:await DeviceInfo.getSystemName(),
  //     osVersion:await DeviceInfo.getSystemVersion(),
  //     model:await DeviceInfo.getDeviceId(),
  //     userAgent:await DeviceInfo.getUserAgent(),
  //     isJailBroken:await JailMonkey.isJailBroken(),
  //     hooksDetected:JailMonkey.hookDetected?await JailMonkey.hookDetected():false,
  //     debugModeEnabled:JailMonkey.isDebugged?await JailMonkey.isDebugged():false,
  //     deviceID:await DeviceInfo.getUniqueId(),
  //     ipAddress:await DeviceInfo.getIpAddress(),
  //     macAddress:await DeviceInfo.getMacAddress(),
  //     readableVersion:await DeviceInfo.getReadableVersion(),
  //     securityPatch:await DeviceInfo.getSecurityPatch(),
  //     trasaAppVersion:await DeviceInfo.getVersion(),
  //     isEmulator:await DeviceInfo.isEmulator(),
  //     isPinOrFingerprintSet:await DeviceInfo.isPinOrFingerprintSet(),
  //
  //
  // }

  const deviceInfo = {
    deviceName: await DeviceInfo.getDeviceName(),
    machineID: await DeviceInfo.getUniqueId(),
    brand: await DeviceInfo.getBrand(),
    manufacturer: await DeviceInfo.getManufacturer(),
    deviceModel: await DeviceInfo.getModel(),
  };

  const deviceOS = {
    osName: await DeviceInfo.getSystemName(),
    osVersion: await DeviceInfo.getSystemVersion(),
    kernelType: '',
    kernelVersion: '',
    latestSecurityPatch: await DeviceInfo.getSecurityPatch(),
    autoUpdate: true,
    jailBroken: await JailMonkey.isJailBroken(),
    debugModeEnabled: JailMonkey.isDebugged ? await JailMonkey.isDebugged() : false,
    isEmulator: await DeviceInfo.isEmulator(),
  };

  const loginSecurity = {
    deviceAutologinDisabled: await DeviceInfo.isPinOrFingerprintSet(),
    deviceLoginMethod: '',
    idleDeviceScreenLockTime: '',
    remoteLoginEnabled: false,
  };

  const netInfo = await NetInfo.fetch();

  const networkInfo = {
    hostname: await DeviceInfo.getHost(),
    interfaceName: netInfo.type || '',
    ipAddress: netInfo.details.ipAddress || '',
    macAddress: await DeviceInfo.getMacAddress(),
    wirelessNetwork: netInfo.type === NetInfoStateType.wifi,
    networkName: netInfo.details.ssid || netInfo.details.carrier || '',
    networkSecurity: '',
  };

  const deviceHygine = {
    deviceInfo,
    deviceOS,
    loginSecurity,
    networkInfo,
    endpointSecurity: {},
  };

  // // //
  // console.log(AppList)

  // AppList.getAll(apps=>{
  //     alert(apps)
  //     info.installedApps=apps
  // })
  // console.log(deviceHygine)

  return deviceHygine;
};

export default getDeviceInfo;
