import Button from '@material-ui/core/Button';
import { makeStyles } from '@material-ui/core/styles';
import MUIDataTable, { MUIDataTableColumn, MUIDataTableMeta } from 'mui-datatables';
import React from 'react';


const useStyles = makeStyles((theme) => ({
  root: {
    flexgrow: 1
  },
  policyName: {
    marginLeft: theme.spacing(6)
  },
  paper: {
    //  backgroundColor:  '#fdfdfd',
    textAlign: 'center',
    padding: theme.spacing(2)
  },
  paperTrans: {
    backgroundColor: 'transparent',
    textAlign: 'center'
  },
  success: {
    paddingLeft: 5,
    paddingRight: 0,
    maxWidth: 50,
    background: 'green',
    color: 'white'
  },
  failed: {
    paddingLeft: 10,
    // paddingRight: 5,
    maxWidth: 50,
    background: 'maroon',
    color: 'white'
  }
}));




type policyTableProps = {
  changeViewDlgState: (v: number) => void;
  changeUpdatePolicyState: (v: number) => void;
  handleDeletePermission: (rowsDeleted: {
    lookup: { [dataIndex: number]: boolean };
    data: Array<{ index: number; dataIndex: number }>;
}) => void;
  policies: Array<object | number[] | string[]>;
}



export default function PolicyTable(props: policyTableProps) {
  const classes = useStyles();
  function statusDiv(val: boolean) {
    if (!val) {
      return <div className={classes.success}>active</div>;
    } else {
      return <div className={classes.failed}>expired</div>;
    }
  }

  const columns = [
    {
      name: 'Policy Name',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        }
      }
    },
    {
      name: 'Expiry',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        }
      }
    },
    {
      name: 'Last Updated',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        }
      }
    },
    {
      name: 'Assigned Services',
      options: {
        filter: true,
        customBodyRender: (value: any) => {
          return value;
        }
      }
    },
    {
      name: 'Status',
      options: {
        filter: true,
        filterOptions: {
          status: [false, true]
        },
        customBodyRender: (value: any) => {
          return statusDiv(value);
        }
      }
    },

    {
      name: 'View Details',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta, updateValue:(value: string) => void)=> {
          // console.log('value', value, 'tableMeta', tableMeta, 'updateValue', updateValue )
          return (
            <Button
              variant="outlined"
              color="secondary"
              size="small"
              onClick={() => {
                props.changeViewDlgState(tableMeta.rowIndex);
              }}
            >
              View
            </Button>
          );
        }
      }
    },
    {
      name: 'Edit',
      options: {
        filter: true,
        customBodyRender: (value: any, tableMeta: MUIDataTableMeta, updateValue:(value: string) => void) => {
          return (
            <Button
              variant="outlined"
              color="secondary"
              size="small"
              onClick={() => {
                props.changeUpdatePolicyState(tableMeta.rowIndex);
              }}
            >
              {/* <EditIcon fontSize="small" /> */}
              Edit
            </Button>
            // <IconButton color="primary" onClick={() => {this.props.changeUpdatePolicyState(value) }} >
            // <EditIcon />
            // </IconButton>
          );
        }
      }
    }
  ];

  return (
    <MUIDataTable
      title={'Policies'}
      data={props.policies}
      columns={columns  as MUIDataTableColumn[]}
      options={{
        filter: true,
        responsive: 'scrollMaxHeight',
        //resizableColumns: true,
        selectableRows: 'single',
        onRowsDelete: props.handleDeletePermission
      }}
    />
  );
}
