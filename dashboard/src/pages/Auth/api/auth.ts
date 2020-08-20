import axios from 'axios';
import Constants from '../../../Constants';

export function FetchExternalIdps(setIdps: (v: any) => void) {
  axios
    .get(`${Constants.TRASA_HOSTNAME}/api/woa/providers/uidp/all`)
    .then((response) => {
      if (response.data.status === 'success') {
        // console.log(response.data);
        const data = response.data.data[0];
        setIdps(data);
      }
    })
    .catch((error) => {
      console.log(error);
    });
}

export function VerifyToken(token: any, setShowSetPasswordCard: (v: boolean) => void) {
  axios
    .get(`${Constants.TRASA_HOSTNAME}/api/v1/verify/${token}`)
    .then((response) => {
      // setUserData(response.data);

      if (response.data.status === 'success') {
        setShowSetPasswordCard(true);
      }
    })
    .catch((error: any) => {
      console.error(error);
    });
}

type PassChangeProps = {
  password: string;
  cpassword: string;
  token: string;
};

export function SetupOrChangePass(
  update: boolean,
  passChangeData: PassChangeProps,
  closeUpdatePassDlg?: () => void,
) {
  let url = `${Constants.TRASA_HOSTNAME}/api/v1/setup/password/${passChangeData.token}`;
  if (update === true) {
    url = `${Constants.TRASA_HOSTNAME}/api/v1/my/changepass`;
  }

  axios
    .post(url, passChangeData)
    .then((response) => {
      // if response status is success, close loader and reqirect user to login page.
      if (response.data.status === 'success') {
        if (update && closeUpdatePassDlg) {
          closeUpdatePassDlg();
        }
        window.location.href = '/login';
      }
      // console.log(response.data);
    })
    .catch((e) => {
      console.error(e);
    });
}
