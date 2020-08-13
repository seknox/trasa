export async function checkresp(resp) {
    if(!(resp.status() === 200 )){
        return false
    }

    let data=await resp.json()

    if (data.status!=="success"){
    return false
    }


}
//
// export const clickWithText=(text )=>{
//         let btns = document.querySelectorAll("button");
//         for (let i=0;i<btns.length;i++){
//             if (btns[i].toString().includes(text) ){
//                 btns[i].click();
//                 return
//             }
//
//         }
// }


const escapeXpathString = str => {
    const splitedQuotes = str.replace(/'/g, `', "'", '`);
    return `concat('${splitedQuotes}', '')`;
};


export const clickWithText = async function(page, text, element) {
    element = element || 'a';
    const escapedText = escapeXpathString(text);
    let xpath = `//*[text()[contains(., ${escapedText})]]`;
    const elements = await page.$x(xpath);
    if(elements.length > 0) {
        for(let i in elements) {
            let e = elements[i];

            // await e.click();
            // return;
            if(await e.isIntersectingViewport()) {
                await e.click();
                return;
            }
        }
    }
    else {
        console.log(xpath);
    }
    throw new Error(`Link not found: ${text}`);
};
