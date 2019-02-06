import React, { Component } from 'react';
import TextField from '@material-ui/core/TextField';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import CustomTable from './utils/table';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl'
import IconButton from '@material-ui/core/IconButton';
import Icon from '@material-ui/core/Icon';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import PrettyPrint from './utils/prettyprint';
import { Button, Divider } from '@material-ui/core';
import {WS} from '../index';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';




const uuidv4 = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c === 'x' ? r : (r && 0x3 | 0x8);
    return v.toString(16);
  });
}


class Stress extends Component {
  constructor(props) {
    super(props);
    this.headerContainer = React.createRef();
  }

  defaultconfig = {
      url: "http://localhost:8080/api/v1/test",
      method: "GET",
      headers: {},
      threads: 4,
      requests: 4,
      resolve: "",
  }

  state = {
    config: {
      url: "http://localhost:8080/api/v1/test",
      method: "GET",
      headers: {},
      threads: 4,
      requests: 4,
      resolve: "",
      payload: "",
      withResponse: false,
    },
    data: [],
    stats: {},
    details: {},
    headerElems: {},
  };
  
  wsMessages = WS.feed.subscribe( message => {
    this.setState({
      data: message.data,
      stats: message.stats,
      details: {uuid: message.uuid, timestamp: message.timestamp, duration: message.durationTotal}
    })
  })

  componentWillUnmount() {
    this.wsMessages.unsubscribe()
  }

  addHeaderBox = () => {
    let uuid = uuidv4();
    let obj = {id: uuid, name:"key", value: "val"}
    this.setState({headerElems: {...this.state.headerElems, [uuid]: {...obj}} })
  }



  submitConfig = () => {
    let headers = Object.keys(this.state.headerElems);
    let headersObj = {}
    for (let index = 0; index < headers.length; index++) {
      let key = headers[index];
      let comp = this.state.headerElems[key];
      headersObj[comp.name] = comp.value 
      
    }
    this.state.config.headers = headersObj
    WS.feed.next(this.state.config);
  }
  clearConfig = () => {
    this.setState({
      config: this.defaultconfig,
      headerElems: {}
    })


  }

  handleHeaderChange = (name, id) => event =>  {
    this.setState({
      ...this.state, 
      headerElems: { ...this.state.headerElems, [id]: {...this.state.headerElems[id], [name]: event.target.value}}
    })
  }
  
  
  handleChange = name => event =>  {
    let val = event.target.value
    if (name === "requests" || name === "threads") {val = parseInt(val)}
    this.setState({
      config: {...this.state.config, [name]: val}
    });
  }

  toggleCheckbox = name => event =>  {
    this.setState({
      config: {...this.state.config, [name]: !this.state.config[name]}
    });
  }

  removeHeader = (h) => () => {


    let temp = this.state.headerElems;
    delete temp[h]
    this.setState({headerElems: temp})

  }
  
  render() {

    const httpMethods = ["GET", "POST", "PUT", "DELETE"];
    
    return (
      <div>
        <form>
          <Card className="paper" elevation={1}>
            <CardContent>
              <FormControl  
              className="text-field"
              >
                <TextField
                id="url"
                label="URL"
                className="text-field"
                value={this.state.config.url}
                onChange={this.handleChange('url')}
                margin="normal"
                />
              </FormControl>
              <FormControl  
              className="text-field"
              >
                <InputLabel htmlFor="">Method</InputLabel>
                <Select
                value={this.state.config.method}
                onChange={this.handleChange('method')}
                >
                {httpMethods.map( e => {
                  return (
                    <MenuItem key={e} value={e}>{e}</MenuItem>
                  )
                  
                })}
                </Select>
              </FormControl>

              <IconButton color="primary" onClick={this.addHeaderBox}><Icon>add_circle</Icon></IconButton>
              <label htmlFor=""> Add Header</label>
              <div ref={this.headerContainer} className="headers-container"> 
                {Object.keys(this.state.headerElems).map( k => {
                  return (
                    <div key={k} className="text-field">
                      <TextField
                      label="Name"
                      value={this.state.headerElems[k].name}
                      className="text-field-s"
                      onChange={this.handleHeaderChange('name', k)}
                      margin="normal"
                      />

                      <TextField
                      label="Value"
                      className="text-field-s"
                      value={this.state.headerElems[k].value}
                      onChange={this.handleHeaderChange('value', k)}
                      margin="normal"
                      />
                      <IconButton color="primary"  onClick={this.removeHeader(k)}><Icon>remove</Icon></IconButton>
                    </div>
                  )
                })}
              </div>

              <FormControl  
              className="text-field"
              >
                <TextField
                id="method"
                label="Threads"
                className="text-field"
                value={this.state.config.threads}
                onChange={this.handleChange('threads')}
                margin="normal"
                type="number"
                />
              </FormControl>

              <FormControl  
              className="text-field"
              >
                <TextField
                id="requests"
                label="Requests"
                className="text-field"
                value={this.state.config.requests}
                onChange={this.handleChange('requests')}
                margin="normal"
                type="number"
                />
              </FormControl>

              <FormControl  
              className="text-field"
              >
                <TextField
                id="resolve"
                label="Resolve"
                className="text-field"
                value={this.state.config.resolve}
                onChange={this.handleChange('resolve')}
                margin="normal"
                />
              </FormControl>
              { 
                this.state.config.method !== 'GET' && 
                this.state.config.method !== 'DELETE' ?
                <TextField
                  className="text-field"
                  placeholder="Payload"
                  value={this.state.config.payload}
                  onChange={this.handleChange('payload')}
                  multiline={true}
                  rows={8}
                  rowsMax={16}
                /> : null
              }
            <FormControlLabel control={
                <Checkbox
                    checked={this.state.config.withResponse}
                    onChange={this.toggleCheckbox('withResponse')}
                    value={false}
                    color="primary"
                  />
                }
                label="Return response"
            />
            </CardContent>
            <CardActions>
              <Button onClick={this.submitConfig} color="primary">Submit</Button>
              <Button color="primary" onClick={this.clearConfig}>Clear</Button>
            </CardActions>

          </Card>
          <Card elevation={1} className="paper">
            <Typography variant="title">
              Report details
            </Typography>
            <Divider />

            <PrettyPrint options={this.state.details}>
            </PrettyPrint>
          </Card>

          <Card elevation={1} className="paper">
              
              <Typography variant="title">
                uStress Config
              </Typography>
              <Divider />
              <PrettyPrint options={this.state.config}>
              </PrettyPrint>
              
          </Card>

          <Card elevation={1} className="paper">
              
              <Typography variant="title">
                Statistics
              </Typography>
              <Divider className="divider" />
              <PrettyPrint options={this.state.stats}>
              </PrettyPrint>
              
          </Card>
          <CustomTable data={this.state.data}></CustomTable>
        </form>
      </div>
    )
  }
}

export default Stress;