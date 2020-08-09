export async function checkresp(resp) {
    if(!(resp.status() === 200 ){
        return false
    }

    let data=await resp.json()

    if (data.status!=="success"){
    return false
    }


}

