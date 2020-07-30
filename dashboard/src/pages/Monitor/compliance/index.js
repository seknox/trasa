import Paper from '@material-ui/core/Paper';
import { withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';
import Layout from '../../../components/Layout/DashboardBase';
import Headers from '../../../components/Layout/HeaderAndContent';


const styles = theme => ({
  paper: {
      backgroundColor:  '#fdfdfd', 
      marginTop: '10%',
      marginLeft: '5%',
     padding: theme.spacing(3),
    color: theme.palette.text.secondary,
 },
headingText: {

  color: 'black',
  fontSize: '20px',
  fontFamily: 'Open Sans, Rajdhani',
},
      
});

class ManageView extends React.Component {
    constructor(props) {
      super(props) ;
      this.state = {

        open: false 
    }
  
    
    }

  
    handleClickOpen = () => {
      this.setState({ open: true });
    };
  
    handleClose = () => {
      this.setState({ open: false });
    };
  
    render() {  
      return (
        <div >
         <Layout > 
         <Headers pageName={'Compliance'} tabHeaders={['Overview']} Components={[ <ManageHOC /> ]}  />
          </Layout>
        </div>
      );
    }
      
  }


const AccountsMain = (props) => {
    return (
   
              <Switch>         
         
             
               <Route exact path='/monitor/compliance' component={ManageView} />  
              </Switch>
      
    )
  
  }
  
  
  export default withRouter(AccountsMain);


class Manage extends Component {
  render() {
    const { classes } = this.props
    return (
      <div>
         <div className={classes.headingText}>
        Reporting based on PCIDSS 3.2 requirements
        
        
        </div>
            <Paper className={classes.paper}>
       
        </Paper>
      </div>
    )
  }
}

Manage.propTypes = {
  classes: PropTypes.object.isRequired,
  theme: PropTypes.object.isRequired,
};




const ManageHOC = withStyles(styles)(Manage)