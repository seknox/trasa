export const policies =[
    {
        name:"full-policy",
        tfaEnabled:true,
        filetransfer:true,
        recordSession:true,
        ipPolicy:'0.0.0.0/0',
        dayTime:[{"days": ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"], "fromTime": "01:00", "toTime": "23:59"}],

    },
    {
        name:"wo tfa",
        tfaEnabled:false,
        filetransfer:true,
        recordSession:true,
        ipPolicy:'0.0.0.0/0',
        dayTime:[{"days": ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"], "fromTime": "01:00", "toTime": "23:59"}],

    },


     {
        name:"nil",
        tfaEnabled:true,
        filetransfer:true,
        recordSession:true,
        ipPolicy:'0.0.0.0/0',
        dayTime:[],

    }


]
