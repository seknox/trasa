import React, { useEffect, useState } from 'react';
import MUIDataTable, {MUIDataTableColumn, MUIDataTableMeta, MUIDataTableOptions} from "mui-datatables";
import Button from "@material-ui/core/Button";
import {Link} from "react-router-dom";
import {createMuiTheme, MuiThemeProvider} from "@material-ui/core";

const theme = createMuiTheme({
    typography: { fontFamily: 'Open Sans, Rajdhani' },
    palette: {
        type: 'light',
        primary: { 500: '#000080' },
        secondary: { A400: '#000080' }, // '#000080' },
    },
});

export default function (props:any) {

    const columns = [
        {
            name: 'Group name',
            options: {
                filter: true,
                customBodyRender: (value: any) => {
                    return value;
                },
            },
        },
        {
            name: 'Policy',
            options: {
                filter: true,
                customBodyRender: (value: any) => {
                    return value;
                },
            },
        },
        {
            name: 'Delete',
            options: {
                filter: false,
                customBodyRender: (
                    value: any,
                    tableMeta: MUIDataTableMeta,
                    updateValue: (value: string) => void,
                ) => {
                    return (
                        <Button
                            onClick={()=>{props.deleteRule(value)}}
                            variant="outlined"
                            color="secondary"
                        >
                            Delete
                        </Button>
                    );
                },
            },
        },
    ];

    const options = {
        filter: true,
        responsive: 'scrollMaxHeight',
        onRowsDelete: props.removeServices,
    };

    return (
        <div>
            <MuiThemeProvider theme={theme}>
                <MUIDataTable
                    title="Dynamic Access Rules"
                    data={props.data}
                    columns={columns as MUIDataTableColumn[]}
                    options={options as MUIDataTableOptions}
                />
            </MuiThemeProvider>
        </div>
    );
}
