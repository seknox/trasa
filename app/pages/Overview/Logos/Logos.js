



export default (name)=>{
    let facebook=require("./Facebook.png")
    let google=require("./Google.png")
    let github=require("./Github.png")
    let flickr=require("./Flickr.png")
    let insta=require("./Instagram.png")
    let yahoo=require("./Yahoo.png")
    let def=require("./Default/default.png")

    switch (name){
        case 'Facebook':return facebook
        case 'Google' : return google
        case 'Github' : return github

        default: return def
    }

}