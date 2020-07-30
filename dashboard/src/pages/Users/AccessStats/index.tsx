import React from 'react';
import AccessStats from '../../Overview/AccessStats';
// import { LogTableV2 } from '../../../monitor/session/LogTableV2';

export default function UserAccessStats(props: any) {
  return (
    <div>
      <AccessStats
        entityType="user"
        entityID={props.userID}
        statusFilter="All"
        timeFilter="Today"
      />
      <br />
      {/* <LogTableV2 entityType="user" entityID={props.userID} /> */}
    </div>
  );
}
