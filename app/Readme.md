# Trasa mobile app


## Initial setup

1. Setup react native development environment. You can follow [this](https://reactnative.dev/docs/environment-setup) link.
2. Run `npm install`
3. Sign up/in to Firebase using your Google account
4. Create a firebase project


#### Android
1. Create firebase app in firebase project you created
2. Edit `android/app/src/main/AndroidManifest.xml` and edit the package name to match package name of the firebase app you just created.
3. Download `firebase-services.json` file  
4. Add the Firebase Android config file (google-services.json) to `android/app/` 

#### IOS  
1. Install [cocoapods](https://cocoapods.org/)  if not installed
2. Go to `ios` directory and run `pod install`
3. Create firebase app in firebase project you created
4. Download `GoogleService-Info.plist` file 
5. Edit bundle ID  to match "iOS bundle ID" of firebase app you just created. The bundle ID can be found within the "General" tab when opening the project with Xcode.
5. Import GoogleService-Info.plist into Xcode project
6. Link Apple's APN with FCM. Follow [this](https://rnfirebase.io/messaging/usage/ios-setup)


## Run
* Android: `npx react-native run-android`  
* IOS: `npx react-native run-ios`

