import AsyncStorage from "@react-native-community/async-storage";


/*
*
* NotificationStorage is a circular list data structure with fixed size
* once size is full new item replaces oldest item in the list
*
*
*
*            notifStackInfo contains variables needed to implement circular list
*            stackStart : index of first element
*            stackEnd  : index of last element
*            stackSize : total number of items in the list
*
*
*
* items are stored in AsyncStorage with keys "notifStack{index}"
* eg notifStack1,notifStack2,notifStack3 etc
*
*
* */






//getStoredNotifications retrives all notifications from notification stack
    export let getStoredNotifications=(callback)=>{
        AsyncStorage.getItem("notifStackInfo").then(info => {

            info=JSON.parse(info)
            let keys=[]
            let stackStart=info.stackStart
            let stackEnd=info.stackEnd

            if(stackEnd>=stackStart){
                for (let i=stackStart;i<=stackEnd;i++){
                    keys.push("notifStack"+i)
                }
            }else {
                for (let i=stackStart;i<=15;i++){
                    keys.push("notifStack"+i)
                }

                for (let i=1;i<=stackEnd;i++){
                    keys.push("notifStack"+i)
                }
            }

            keys.reverse()


            AsyncStorage.multiGet(keys).then(values=>{
                let val=values.map((v,i)=>{
                    let temp=JSON.parse(v[1])
                    temp["ASkey"]=keys[i]
                    return temp

                })
                callback(val)
            })

        })
    }
    //storeNotification pushes a new notification into notification stack
   export let storeNotification= async (notification)=>{
        //Get notification Info from storage
        let notifInfo= await AsyncStorage.getItem("notifStackInfo")

        //If notifInfo is not initialized initialize it else parse it
        if(!notifInfo){
            notifInfo={stackStart:1,stackEnd:0,stackSize:0}
            //console.log("notifInfo initialized")
        }else {
            notifInfo=JSON.parse(notifInfo)
            //console.log("notifInfo parsed")
        }
        let tempNotif=JSON.stringify(notification.data)
        //console.log("temp notif data", tempNotif)

        let stackStart=notifInfo.stackStart
       let stackEnd=notifInfo.stackEnd
       let stackSize=notifInfo.stackSize
        if(stackEnd<15){
            stackEnd=stackEnd+1
        }else {
            stackEnd=1
        }

        if(stackSize>=15){
            stackStart=stackStart+1
        }else {
            stackSize=stackSize+1
        }

        notifInfo.stackStart=stackStart
        notifInfo.stackEnd=stackEnd
        notifInfo.stackSize=stackSize


        await AsyncStorage.setItem("notifStack"+stackEnd,tempNotif)
        await AsyncStorage.mergeItem("notifStackInfo",JSON.stringify(notifInfo))

    }
