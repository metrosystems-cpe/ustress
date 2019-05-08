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
import ReportsDashboard from './utils/reportsDashboard';
import { Line, Scatter } from 'react-chartjs-2';
import * as zoom from 'chartjs-plugin-zoom'
import * as downsample from 'chartjs-plugin-responsive-downsample';
import Switch from '@material-ui/core/Switch';
import SimpleSlider from './utils/slider';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import classNames from 'classnames';
import PropTypes from 'prop-types';

import LinearProgress from '@material-ui/core/LinearProgress';
import Chip from '@material-ui/core/Chip';



import {CurrentDomain} from '../index';
import { withSnackbar } from 'notistack';




const uuidv4 = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c === 'x' ? r : (r && 0x3 | 0x8);
    return v.toString(16);
  });
}

const httpMethods = ["GET", "POST", "PUT", "DELETE"];

class Stress extends Component {
  constructor(props) {
    super(props);
    this.headerContainer = React.createRef();
    this.expansionPanel = React.createRef();
    WS.enqueueSnackbar = this.props.enqueueSnackbar
  }

  defaultconfig = {
    url: `${CurrentDomain}/api/v1/test`,
    method: "GET",
    headers: {},
    threads: 4,
    requests: 4,
    resolve: "",
    payload: "",
    withResponse: false,
    duration: 0,
    frequency: 0,
    durationBased: false,
  }

  defaultDurationConfig = {
    url: `${CurrentDomain}/api/v1/test`,
    method: "GET",
    headers: {},
    threads: 4,
    requests: 0,
    resolve: "",
    payload: "",
    withResponse: false,
    duration: 60,
    frequency: 1000,
    durationBased: true,
  }

  state = {
    config: {
      url: `${CurrentDomain}/api/v1/test`,
      method: "GET",
      headers: {},
      threads: 4,
      requests: 4,
      resolve: "",
      payload: "",
      withResponse: false,
      duration: 0,
      frequency: 0,
      durationBased: false,
    },
    data: [],
    stats: {},
    details: {},
    headerElems: {},
    expansionPanelExpanded: true
  };

  wsMessages = WS.feed.subscribe( message => {
    console.log(message)
    if (message.data) {
      let sortedData = message.data.sort((a, b) => {
        return a - b
      })
      if (!sortedData) {
        return
      }
      if (message.completed) {
          message.config.requests = message.data.length;

      }
      this.setState({
        data: sortedData,
        stats: message.stats,
        details: {uuid: message.uuid, timestamp: message.timestamp, duration: message.durationTotal},
        config: {...message.config}
      })

    }
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
    if (this.state.config.durationBased) {
      this.state.config.requests = 0;
    } else {
      this.state.config.duration = 0;
      this.state.config.frequency = 0;
    }
    this.toggleExpansionPanel();
    WS.feed.next(this.state.config);
  }
  clearConfig = () => {
    this.setState({
      config: this.defaultconfig,
      headerElems: {},
      data: []
      
    })


  }

  toggleExpansionPanel = () => {
      this.setState({expansionPanelExpanded: !this.state.expansionPanelExpanded})

  }
  handleHeaderChange = (name, id) => event =>  {
    this.setState({
      ...this.state, 
      headerElems: { ...this.state.headerElems, [id]: {...this.state.headerElems[id], [name]: event.target.value}}
    })
  }
  
  
  handleChange = name => (event, val) =>  {
    if (name != "duration" && name != "frequency") {
      val = event.target.value
    }
    let ints = ["requests", "threads", "duration", "frequency"]
    if (ints.indexOf(name) != -1) {val = parseInt(val)}
    setTimeout(() => {
      this.setState({
        config: {...this.state.config, [name]: val}
      });
    })
  }

  toggleCheckbox = name => event =>  {
    var req = this.state.config.requests;

    this.setState({
      config: {...this.state.config, [name]: !this.state.config[name]}
    });

  }

  removeHeader = (h) => () => {


    let temp = this.state.headerElems;
    delete temp[h]
    this.setState({headerElems: temp})

  }

  getFormFields = () => {

    if (!this.state.config.durationBased) {
      return (
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
      )
    } else {
      this.state.config.requests = 0;
      return (
        <div>
          <FormControl  
          className="text-field"
          >
            <SimpleSlider
                sliderLabel="Duration"
                sliderDescription="Seconds"
                val={this.state.config.duration}
                min={1}
                max={3600}
                handleChange={this.handleChange("duration")}>

            </SimpleSlider>
          </FormControl>
          <FormControl  
          className="text-field"
          >
            <SimpleSlider
                sliderLabel="Frequency"
                sliderDescription="Milliseconds"
                val={this.state.config.frequency}
                min={1}
                max={3000}
                handleChange={this.handleChange("frequency")}>

            </SimpleSlider>
          </FormControl>

        </div>
      )
    }

  }

  getForm = () => {
    var durationBased = this.state.config.durationBased
    return ( 
      <form className="text-field">
        {/*<Card className="paper" elevation={1}>*/}
        {/*  <CardContent>*/}
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

            {this.getFormFields()}

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
                  value={'false'}
                  color="primary"
                />
              }
              label="Return response"
          />

          <FormControlLabel control={
              <Switch color="primary" onChange={this.toggleCheckbox('durationBased')} checked={this.state.config.durationBased}/>
              }
              label="Duration based"
          />
          {/*</CardContent>*/}
          {/*<CardActions>*/}
          {/*  <Button onClick={this.submitConfig} color="primary">Submit</Button>*/}
          {/*  <Button color="primary" onClick={this.clearConfig}>Clear</Button>*/}
          {/*</CardActions>*/}

        {/*</Card>*/}

      </form>
    )

  }

  getCurrentProgress = () => {
    var pgrs = (this.state.data.length/ this.state.config.requests) * 100
    return pgrs
  }

  getExpansionPanelSummary = () => {
    const currentProgress = this.getCurrentProgress();
    const status = currentProgress < 100 ? "Attack in progress (" + currentProgress+"%)"
        : "Attack completed (" + this.state.data.length + " requests made in " + this.state.details.duration + "s)" ;

    let rps = this.state.data.length > 0 ? Math.floor(this.state.data.length / this.state.details.duration) : 0;

    if (this.state.expansionPanelExpanded) {
      return (
          <div>
            <Typography variant="caption">Configure Stress Test</Typography>
          </div>
      )
    }

    return (
        <div className="text-field">
          <div className="cstm-flex">
            <div className="cstm-flex-item">

              <Typography variant="subtitle2">
                Avg. Response time
              </Typography>
              <Typography variant="title" className="center">
                { this.state.stats.median }
              </Typography>
            </div>
            <div className="cstm-flex-item">
              <Typography variant="subtitle2">
                RPS
              </Typography>
              <Typography variant="title" className="center">
                { rps }
              </Typography>
            </div>
            <div className="cstm-flex-item">
              <Typography variant="subtitle2">
                Error rate
              </Typography>
              <Typography variant="title" className="center">
                { this.state.stats.error_percentage }%
              </Typography>
            </div>
          </div>
          <Typography variant="caption">{status}</Typography>
          <br/>
          <LinearProgress variant="determinate" value={this.getCurrentProgress()} />

        </div>
    )


  }


  render() {

    const {classes} = this.props;
    
    return (
      <div>

          <ExpansionPanel ref={this.expansionPanel} expanded={this.state.expansionPanelExpanded} onChange={this.toggleExpansionPanel}>
            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
              {this.getExpansionPanelSummary()}
            </ExpansionPanelSummary>
            <ExpansionPanelDetails >
              {this.getForm()}
            </ExpansionPanelDetails>
            <Divider />
            <ExpansionPanelActions>
              <Button size="small" onClick={this.clearConfig}>Clear</Button>
              <Button size="small" color="primary" onClick={this.submitConfig}>
                Submit
              </Button>
            </ExpansionPanelActions>
          </ExpansionPanel>

        {/*{this.getForm()}*/}

        <ReportsDashboard
          details={this.state.details}
          data={this.state.data}
          stats={this.state.stats}
          config={this.state.config}/>
      </div>
    )
  }
}


export default withSnackbar(Stress);