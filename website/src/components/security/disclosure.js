import { Typography } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import ThemeBase from '../muiTheme';

const useStyles = makeStyles(() => ({
  ctaPad: {
    marginTop: 50,
    textAlign: 'center',
  },
}));

export default function Disclosure() {
  const classes = useStyles();
  return (
    <ThemeBase>
      <Grid container spacing={2} direction="column" justify="center" alignItems="center">
        <Grid item xs={12} sm={12}>
          <div className={classes.ctaPad}>
            <Typography variant="h1">Responsible disclosure</Typography>
          </div>
        </Grid>
        <Grid item xs={12} sm={7}>
          <div className={classes.ctaPad}>
            <Typography variant="body1">
              We take security very seriously at Seknox. We apply threat driven security design and
              risk modeled feature implementation to develop and secure our products. We also
              acknowledge that security is a continuous process, and as such, we believe in an
              offensive security approach to find and patch security flaws sooner than later.
            </Typography>
            <Typography variant="body1">
              If you would like to report a security vulnerability, send us an email to address
              security at seknox dot com. If you want to secure the message, please use the
              following PGP key:{' '}
            </Typography>
            <div style={{ color: 'silver', backgroundColor: 'black' }}>
              <code style={{ color: 'silver', backgroundColor: 'black' }}>
                -----BEGIN PGP PUBLIC KEY BLOCK-----
                mQGNBF+LQU4BDADK1syqmLe+0e8OmJiCDfHsp5k2Tt6urhuUo+DySYIR0DvtUKH/
                xW6tlOH548Zk3AOlbEqCdXJYjcnSXsGkFbUFoPy9K5Z2e7USjtaqvmiVEGAT/y+7
                r3KKVggVxIFf+04MxmHGdaWPks706B90p2LABtPdMLKP9pSPyCYS1iGfKzpBdUfG
                BRK/vCoyBrJZKHyREIYLuBaSpClbNLcadDic2kzuPokkRcSqILPot57XC8+vIdYA
                zyOJiaFfkYJWHhallw/cXiLtBIyLqcgke2PB7PdFcJH5kumYluJyAwAiiQWE5s5i
                IHH0hkdWh+hzqbVKkCWK9mvdU/pD7NJeuinvSy/06Si6ub3Cm/E6tsaDacGzN+IX
                2VmZg9i+HU8uogm9a2TsWb4ou8mpt95fNw3d1OiwwBxrfpBP3N8BCRogD6QRVJvx
                Qqvzaptb0Q2Bvyx8ufyHGVc6GSXG3uTN90c9mvD7FRonicCrSoLIm0RE7yrnipss
                WeKjcsSlhPXbZNEAEQEAAbQqU2Vrbm94IHNlY3VyaXR5IHRlYW0gPHNlY3VyaXR5
                QHNla25veC5jb20+iQHUBBMBCgA+FiEEDugsgGM+/oH7841q3xsa357zMF0FAl+L
                QU4CGwMFCQPCZwAFCwkIBwIGFQoJCAsCBBYCAwECHgECF4AACgkQ3xsa357zMF2w
                bQv5AccmqEd7vBhlthW4GYnHcqESxHXCcauje8582olt0mXZLm8Z+8avpBOSkSlq
                kSzPVPS6SylRcnVVTKVEQB0r0o/GTXgxGiqycWtzjXDkLtCjnhHj5aKmyKmL+MqH
                5o639ihnvlpuLEpnCOhj8GZGfkEhtF2lh9gPU61OlE1XdKnriCrmEgqiOoQ43N+R
                Zb0Vr3pj7GKzW1YXdL/zIMIL9mxcqJJgsMA4tb76huVj8tZIuhGEYSQJsFQtUdCt
                G7Vl8vDJYdjPD0PIvWSjXeez0HHy52/D25AJND3TgW3VLoNdwFrmzapNhuVuNCTj
                KI2gkGN1ApZqGH3rwgFtGgYJZ70egQYNofB5HbJ9+W1+AOJlWojrzgrlDBgO7DGY
                3WJwtezoCO/TQNTzqAMQMKN/c1VjyBfEGcU/3CxIbPJfc/1LpW5StYX2EPNsTD3y
                IzxojYaSpl9iIvB4g3LVrvy5J3cbvsBi+zgcDiNHusN3GXk6C7Yb7Kkik+k+NS3d
                qW/tuQGNBF+LQU4BDACxONoR+Z0s4Uiqze0BIGP9LX7+ad57jweMFBZr7MVFulzk
                xbwQv0/JG1B+YrjQljLuM6gRoneh3rvy6WDrnrwF9zQqaxbW9l6+LiQr2PZCLBhJ
                EhzSsDMYg79YxdQ5iGCiNymnPv3n4pzy5Lvbk0vFD2FmsjwfkHxoBCb2mGVG2Zvh
                BIdw3fcgIUUW30mFq0OD3d6BFxSm+eE9TJ/PCi1ZrVGHoSU46IEULHWsbWLK2LhP
                gDADlhgp0aQpaklGGyNN9JY1Swn5X9AF9KEVLrImuRmX5GLkH4jowLNONNAZfatz
                DjAUgRRc5CJvj+Pi1IF95J7eHn7d8YNXxKrXjhE3GUVN8Sv98jjiXVjdBfbUam51
                GmQmKk+ldyGleyCjCKRVMnOGUUvtVC0COX1lNCAEPdW9s5AilAR/m7rnSr7E1iQO
                CpffMJKPWC4rhPThMaUaxtLVORdekwB4Y8g+r8cqC361e4I6jAB8B4+w4lV6aJwh
                zVi8EZGsLf9sKhFLKx0AEQEAAYkBvAQYAQoAJhYhBA7oLIBjPv6B+/ONat8bGt+e
                8zBdBQJfi0FOAhsMBQkDwmcAAAoJEN8bGt+e8zBd3sIL+gNNKSNNxTVLDrb973G/
                2Fa4WfaUhMIamMadvAaZ6LLorWAReof6kiD0506mq5ZD//FWurXUZIJIK+JxDdab
                OPHZBEJKoPqevAFK/OhHtyaYhZMOzhFPPbDiXwQrLew6qy09AoXcNRlwRUibz7zd
                9Lv48ZkE578UG27i6BLUie6OelgBbdMs9POOpcUxUuwMmP8mwlZGIPQX64UE2iJh
                xyVydhG8c/2GLoQRisJqM0TkfhsUaAyK+jxa07XBmRbpQfJae+gmIH1oGjdV3nDh
                rEeIVtaRc9IVk3VDP9UWMAZEiT+pPDzJWxuijD3XyqFLPnI2C4xaCUjZgGyFVk8j
                6v6QYWF9bAhj+EXiBvrVzxTJXqduy6J8kBB2rOP1gO7NpmJqZAkxdf6a9pn07hN0
                jhU67vViOyE2gr9ZNzC90yiTNdw2JrVCUXb29DHFCUeHEefMto2Vdt+nJR7zxDM6
                cOChHxEct0QsxqM5JOAr9ST4hU6LMbYstSJV0Aquiy8LuA== <br />
                =Mxvb <br />
                -----END PGP PUBLIC KEY BLOCK-----
              </code>
            </div>
          </div>
        </Grid>
      </Grid>
    </ThemeBase>
  );
}
