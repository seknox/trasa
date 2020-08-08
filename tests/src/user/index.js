import create from "./create";
import Constants from "../../Constants";


export default  ()=>{
    it('Should load service page on url address: '+Constants.TRASA_DASHBOARD+'/users', async () => {
        await page.goto(Constants.TRASA_DASHBOARD+'/services');
        await expect(page).toMatch('HTTPs applications')
    })



    it('does user operations correctly',async ()=>{
        await create()

    })
}
