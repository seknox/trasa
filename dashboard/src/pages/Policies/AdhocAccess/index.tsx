import axios from 'axios';
import React, { useState, useEffect } from 'react';
import Constants from '../../../Constants';
import AccessRequests from './AdhocRequest';
import AccessRequestHistory from './AdhocRequestHistory';

export default function AdocAccess() {
  const [orgID, setOrgID] = useState('');

  useEffect(() => {
    axios
      .get(Constants.TRASA_HOSTNAME + '/api/v1/my')
      .then((response) => {
        if (response.data.status === 'success') {
          setOrgID(response.data.data[0].Org.ID);
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <div>
      <AccessRequests />
      <AccessRequestHistory orgID={orgID} />
    </div>
  );
}
