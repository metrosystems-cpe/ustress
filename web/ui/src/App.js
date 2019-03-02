import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Link } from "react-router-dom";
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Reports from './components/Reports';
import Stress from './components/Stress';
import Home from './components/Home';
import Icon from '@material-ui/core/Icon';
import Toolbar from '@material-ui/core/Toolbar';
import AppBar from '@material-ui/core/AppBar';
import IconButton from '@material-ui/core/IconButton';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import { SnackbarProvider, withSnackbar } from 'notistack';
import CssBaseline from '@material-ui/core/CssBaseline';
import {WS} from './index'

import './App.scss';



const styles = theme => ({
  root: {
    display: 'flex',
  },
  appBar: {
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: drawerWidth,
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginLeft: 12,
    marginRight: 20,
  },
  hide: {
    display: 'none',
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: '0 8px',
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing.unit * 3,
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: -drawerWidth,
  },
  contentShift: {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: 0,
  },
});



const drawerWidth = 240;


class App extends Component {

  state = {
    open: true,
    title: "Home",
  };

  routesList  = [
    {
      name: "home",
      path: "/",
      icon: "home",
      text: "Home",
      component: Home,
    },
    {
      name: "stress",
      path: "/stress",
      icon: "broken_image",
      text: "Stress",
      component: Stress,
    },
    {
      name: "reports",
      path: "/reports",
      icon: "show_chart",
      text: "Reports",
      component: Reports,
    }
  ]

  handleDrawerOpen = () => {
    this.setState({ open: true });
  };

  handleDrawerClose = () => {
    this.setState({ open: false });
  };


  updatePageTitle =  () => {
    let filtered = this.routesList.filter(e => {
      return window.location.pathname.indexOf(e.path) !== -1 ? e : "";
    })
    
    let page = filtered.length > 0 ? filtered[filtered.length-1].text : null;  
    this.setState({ title: page })
  }
  render() {
    
    const { classes, theme } = this.props;
    const { open } = this.state;
    return (
      <SnackbarProvider maxSnack="3" anchorOrigin={{vertical:"bottom", horizontal:"right"}}>

      <Router path="/" basename="/ustress/ui/public">
        <div className={classes.root}>
          <CssBaseline />
          <AppBar
            position="fixed"
            className={classNames(classes.appBar, {
              [classes.appBarShift]: open,
            })}
          >
            <Toolbar>

              <IconButton
                color="inherit"
                aria-label="Open drawer"
                onClick={this.handleDrawerOpen}
                className={classNames(classes.menuButton, open && classes.hide)}
              >
                <MenuIcon />
              </IconButton>
              <Typography variant="h6" color="inherit" noWrap>
                uStress Tool - {this.state.title}
              </Typography>
            </Toolbar>
          </AppBar>
          <Drawer 
            className={classes.drawer}
            variant="persistent"
            anchor="left"
            open={open}
            classes={{
              paper: classes.drawerPaper,
            }}
          >
            <div className={classes.drawerHeader}>
              <IconButton onClick={this.handleDrawerClose}>
                {theme.direction === 'ltr' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
              </IconButton>
            </div>
            <Divider />
            <div
              tabIndex={0}
              role="button"
            >
              <List onClick={this.updatePageTitle}>
                {this.routesList.map(elem => {
                  return (
                    <Link to={elem.path} key={elem.name}>
                      <ListItem button key={elem.name}>
                        <ListItemIcon><Icon>{elem.icon}</Icon></ListItemIcon>

                        <ListItemText primary={elem.text} />
                      </ListItem>
                    </Link>
                  )
                })}
              </List>
            </div>
          </Drawer>
          <main
            className={classNames(classes.content, {
              [classes.contentShift]: open,
            })}
          >
          <div className="router-content">
            {this.routesList.map( r => {
              return <Route key={r.name} path={r.path} exact component={r.component} />
            })}
          </div>

          </main>
        </div>
      </Router>
      </SnackbarProvider>
    );
  }
}
App.propTypes = {
  classes: PropTypes.object.isRequired,
  theme: PropTypes.object.isRequired,
};

export default withStyles(styles, { withTheme: true })(App);






