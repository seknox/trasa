import React from 'react';
import AccessStats from '../../Overview/AccessStats';
import { LogTableV2 } from '../../Monitor/Session/LogTableV2';

export default function ServiceAccessStats(props: any) {
  return (
    <div>
      <AccessStats
        entityType={props.entityType}
        entityID={props.entityID}
        statusFilter="All"
        timeFilter="Today"
      />
      <br />
      <LogTableV2 entityType={props.entityType} entityID={props.entityID} />
    </div>
  );
}
