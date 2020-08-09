import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import PolicyIcon from '@material-ui/icons/Assignment';
import axios from 'axios';
// import mixpanel from 'mixpanel-browser';
import React, { useEffect, useState } from 'react';
import DialogueWrapper from '../../../utils/Components/DialogueWrapComponent';
import Constants from '../../../Constants';
import CreatePolicy, { ReviewAccess } from './newCreatePolicy';
import PolicyTable from './newPolicyTable';

type singlePolicytype = {
  policyID: string;
  policyName: string;
  dayAndTime: any;
  tfaRequired: boolean;
  recordSession: boolean;
  fileTransfer: boolean;
  ipSource: string;
  expiry: string;
  devicePolicy: DevicePolicyProps;
};

export type DevicePolicyProps = {
  blockUntrustedDevices: boolean;
  blockAutologinEnabled: boolean;
  blockTfaNotConfigured: boolean;
  blockIdleScreenLockDisabled: boolean;
  blockRemoteLoginEnabled: boolean;
  blockJailBroken: boolean;
  blockDebuggingEnabled: boolean;
  blockEmulated: boolean;
  blockEncryptionNotSet: boolean;
  blockOpenWifiConn: boolean;
  blockFirewallDisabled: boolean;
  blockAntivirusDisabled: boolean;
};

const policyObj = {
  policyID: '',
  policyName: '',
  dayAndTime: {},
  tfaRequired: false,
  recordSession: false,
  fileTransfer: false,
  ipSource: '',
  expiry: '',
  devicePolicy: {
    blockAutologinEnabled: false,
    blockUntrustedDevices: false,
    blockTfaNotConfigured: false,
    blockIdleScreenLockDisabled: false,
    blockRemoteLoginEnabled: false,
    blockJailBroken: false,
    blockDebuggingEnabled: false,
    blockEmulated: false,
    blockEncryptionNotSet: false,
    blockOpenWifiConn: false,
    blockFirewallDisabled: false,
    blockAntivirusDisabled: false,
  },
};

export default function Policies() {
  const [policyArr, setPolicyArr] = useState([]); // useState<Array<object[]>>([]);
  const [policyObject, setpolicyObject] = useState([policyObj]); // useState<object[]>([{}]);
  const [viewDlg, openViewDlg] = useState(false);
  const [createPolicyDlgState, changeCreatePolicyDlgState] = useState(false);
  const [update, setUpdate] = useState(false);
  const [singlePolicy, setSinglePolicy] = useState(policyObj);

  const changeViewDlgState = (rowIndex: number) => {
    openViewDlg(!viewDlg);
    setSinglePolicy(policyObject[rowIndex]);
  };
  const closeViewDlg = () => {
    openViewDlg(false);
  };
  const openCreatePolicyDlg = () => {
    setSinglePolicy(policyObj);
    changeCreatePolicyDlgState(true);
    setUpdate(false);
  };

  const handleClose = () => {
    changeCreatePolicyDlgState(false);
    setUpdate(false);
  };

  // };
  const changeUpdatePolicyState = (rowIndex: number) => {
    setUpdate(true);
    changeCreatePolicyDlgState(true);
    setSinglePolicy(policyObject[rowIndex]);
  };

  const handleDeletePermission = (rowsDeleted: any) => {
    const policy = policyObject[rowsDeleted.data[0].index];

    const req = { policyID: [policy.policyID] };
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/delete`, req)
      .then(() => {
        window.location.reload();
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const getAllPolicies = () => {
    axios
      .get(`${Constants.TRASA_HOSTNAME}/api/v1/groups/policy/all`)
      .then((response) => {
        // this.setState({allUsers: response.data})
        const resp = response.data.data[0];
        setpolicyObject(resp);
        let dataArr = [];
        dataArr = resp.map(function (n: any) {
          // let edate = new Date(n.expiry*1000)
          const udate = new Date(n.updatedAt * 1000);
          return [
            n.policyName,
            n.expiry,
            udate.toDateString(),
            n.usedBy,
            n.isExpired,
            n.policyID,
            n.policyID,
          ];
        });
        setPolicyArr(dataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    // mixpanel.track('control-policies');
    getAllPolicies();
  }, []);

  return (
    <div>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <CreatePolicy
            open={createPolicyDlgState}
            update={update}
            handleClose={handleClose}
            updateData={singlePolicy as singlePolicytype}
            // updatePoliciesArr={updatePoliciesArr}
          />
          <br />
          <Button variant="contained" size="small" onClick={openCreatePolicyDlg}>
            <PolicyIcon />
            Create New Policy
          </Button>
        </Grid>
        <Grid item xs={12}>
          <PolicyTable
            policies={policyArr}
            changeUpdatePolicyState={changeUpdatePolicyState}
            handleDeletePermission={handleDeletePermission}
            changeViewDlgState={changeViewDlgState}
          />
        </Grid>
        <DialogueWrapper
          open={viewDlg}
          handleClose={closeViewDlg}
          title="Policy"
          maxWidth="lg"
          fullScreen={false}
        >
          <ReviewAccess
            policyName={singlePolicy.policyName}
            dayAndTime={singlePolicy.dayAndTime}
            tfaRequired={singlePolicy.tfaRequired}
            recordSession={singlePolicy.recordSession}
            fileTransfer={singlePolicy.fileTransfer}
            ipSource={singlePolicy.ipSource}
            expiry={singlePolicy.expiry}
            devicePolicy={singlePolicy.devicePolicy}
          />
        </DialogueWrapper>
      </Grid>
    </div>
  );
}
