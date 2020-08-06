import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import Constants from '../../../../Constants';
import Layout from '../../../../Layout/DashboardBase';
import Headers from '../../../../Layout/HeaderAndContent';
import AssignedUserGroupPage from './AssignedUserGroupsPage';
import GroupPage from './GroupPage';

function GroupPageIndex(props: any) {
  const [addServicesDlg, setaddServicesDlg] = useState(false);
  const [dispayList, setdispayList] = useState(true);
  const [addedservicesArray, setaddedservicesArray] = useState([]);
  const [groupMeta, setgroupMeta] = useState({ groupID: '', groupName: '' });
  // const [addedServices, setaddedServices] = useState([]);
  const [addedServicesIndex, setaddedServicesIndex] = useState([]);
  const [allAddedservicesObj, setallAddedservicesObj] = useState({});
  const [unaddedservices, setunaddedservices] = useState([]);
  // const [serviceName, setserviceName] = useState('');

  useEffect(() => {
    fetchservices(props.Servicegroupid);
  }, [props.Servicegroupid]);

  const removeServices = (rowsDeleted: any) => {
    const Services = rowsDeleted.data.map((v: any) => {
      return addedservicesArray[v.index][3];
    });

    const reqData = {
      groupID: groupMeta.groupID,
      authID: Services,
      updateType: 'remove',
    };
    const url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/service/update`;
    axios
      .post(url, reqData)
      .then(() => {})
      .catch((error) => {
        console.log(error);
      });
  };

  const onServiceselected = (index: any) => {
    // const addedServices = ServiceData
    //   .filter((u: any, i: any) => addedServicesIndex.indexOf(i) > -1)
    //   .map((app: any) => ({ serviceName: app.serviceName, ID: app.ID }));

    setaddedServicesIndex(index);
    // setaddedServices(addedServices);
  };

  const openAddServicesDlg = () => {
    setaddServicesDlg(!addServicesDlg);
  };

  const changeServiceDisplay = () => {
    setdispayList(!dispayList);
  };

  const fetchservices = (isGroupOrAllservices: any) => {
    let url = '';
    if (isGroupOrAllservices === 'allservices') {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/services/all`;
    } else {
      url = `${Constants.TRASA_HOSTNAME}/api/v1/groups/service/${isGroupOrAllservices}`;
    }
    axios
      .get(url)
      .then((response) => {
        const resp = response.data.data[0];
        setallAddedservicesObj(resp.addedServices);
        setunaddedservices(resp.unaddedServices);
        setgroupMeta(resp.groupMeta);

        const ddataArr = resp.addedServices.map(function (n: any) {
          return [n.serviceName, n.hostname, n.serviceType, n.ID];
        });

        setaddedservicesArray(ddataArr);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const DeleteAuthServicegroup = () => {
    axios
      .post(`${Constants.TRASA_HOSTNAME}/api/v1/groups/delete/${props.Servicegroupid}`)
      .then(() => {
        window.location.href = '/services#Service%20Groups%20(clusters)';
      })
      .catch((error) => {
        console.log(error);
      });
  };

  const getRoute = (staticVal: any, dynamicVal: any) => {
    return staticVal + dynamicVal;
  };

  return (
    <Layout>
      <Headers
        pageName={[
          { name: 'Service Group', route: '/services#Service%20Groups%20(clusters)' },
          {
            name: groupMeta.groupName,
            route: getRoute('/services/groups/group/', props.Servicegroupid),
          },
        ]}
        tabHeaders={['Overview', 'Assigned User Groups']}
        Components={[
          <GroupPage
            groupID={props.Servicegroupid}
            allAddedservicesObj={allAddedservicesObj}
            dispayList={dispayList}
            changeServiceDisplay={changeServiceDisplay}
            addedservicesArray={addedservicesArray}
            openAddServicesDlg={openAddServicesDlg}
            open={addServicesDlg}
            onServiceselected={onServiceselected}
            addedServicesIndex={addedServicesIndex}
            unaddedservices={unaddedservices}
            groupMeta={groupMeta}
            DeleteAuthServicegroup={DeleteAuthServicegroup}
            removeServices={removeServices}
          />,

          <AssignedUserGroupPage groupMeta={groupMeta} groupID={props.Servicegroupid} />,
        ]}
      />
    </Layout>
  );
}

type TParams = { ID: string };

const GroupPageHOC = ({ match }: RouteComponentProps<TParams>) => {
  return (
    <div>
      <GroupPageIndex Servicegroupid={match.params.ID} />
    </div>
  );
};

export default GroupPageHOC;
