import React, { Component } from 'react';
import PropTypes from 'prop-types';
import List from '@material-ui/core/List';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

class PrettyPrint extends Component {

  render() {
    let {options} = this.props
    let optskeys = Object.keys(options);
    optskeys = optskeys ? optskeys : []; 
    
    return (
      <List>
        {optskeys.map(k => {
          return (
            <ListItem button key={k}>
              <ListItemText primary={k}></ListItemText>
              <span className="spacer"></span>
              <ListItemText className="right" primary={
                typeof options[k] === 'object' ? JSON.stringify(options[k]) : options[k]
                }></ListItemText>
            </ListItem>
          )
        })}
      </List>
    )
  }

}

export default PrettyPrint