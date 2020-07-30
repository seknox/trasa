export type UserData = {
  idpName: string;
  userName: string;
  email: string;
  password: string;
  orgID: string;
  intent: string;
};

export type LoginProps = {
  userData: UserData;
  autofillEmail: boolean;
  intent: 'AUTH_REQ_DASH_LOGIN' | 'AUTH_REQ_CHANGE_PASS' | 'AUTH_REQ_ENROL_DEVICE';
  title: string;
  showForgetPass: boolean;
  tfaRequired?: boolean;
  // loader: boolean;
  setData: (value: any) => void;
  changeHasAuthenticated: () => void | null;
};

export type TfaProps = {
  tfaMethod: string;
  totpCode: string;
  token: string;
  intent: string;
};

export type LoginIntent =
  | 'AUTH_REQ_DASH_LOGIN'
  | 'AUTH_REQ_CHANGE_PASS'
  | 'AUTH_REQ_ENROL_DEVICE'
  | 'AUTH_REQ_FORGOT_PASS';
