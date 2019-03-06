
import React, { Component } from 'react';
import { Typography, Card, CardContent, Select, InputLabel, FormControl,  ListItem, ListItemText } from '@material-ui/core';
import Axios from 'axios';
import PrettyPrint from './utils/prettyprint';
import CustomTable from './utils/table';
import {CurrentDomain} from '../index';
import { withSnackbar } from 'notistack';
import ReportsDashboard from './utils/reportsDashboard';

class Reports extends Component {

  constructor(props) {
    super(props)
    this.getReports()
  }

  state = {
    selectedReport: "",
    data: [],
    report: {
      data: [],
      stats: {},
      config: {},
    }
  }


  
  getReport = report_id => {
    if (!report_id) { return }
    Axios
      .get(`${CurrentDomain}/ustress/api/v1/file_reports?file=` + report_id)
      .then(response => {
        
        this.props.enqueueSnackbar(`Fetched report ${report_id}`, {variant:"success"})
        this.setState({...this.state, report: response.data})
      })
      .catch(error => {
        this.errored = true
      })
      .finally(() => this.loading = false)
  }

  parseReports = (reports) => {
    return reports.map( e => {
      return JSON.parse(e.report)
    })
  }

  /*
  @TODO!
  Error messages and success messages should not be hardcoded, 
  those should come directly from backend 
  */
  getReports = () => {
    Axios
    .get(`${CurrentDomain}/ustress/api/v1/reports`)
    .then(response => {
      this.props.enqueueSnackbar("Reports fetched", {variant:"success"})
      this.setState({...this.state, data: this.parseReports(response.data.entries)})
    })
    .catch(error => {
      this.props.enqueueSnackbar("Retrieving reports from cassandra failed", {variant:"error"})
      if (error.response && error.response.status === 400) {
        Axios.get(`${CurrentDomain}/ustress/api/v1/file_reports`).then(res => {
          this.props.enqueueSnackbar("Retrieved reports from local storage", {variant:"success"})
          this.setState({data: res.data.length > 0 ? res.data : []})
        }).catch( err => {
          this.props.enqueueSnackbar(err.error, {variant: "error"})
        })

      }
    })
    .finally(() => this.loading = false)
  }

  handleChange = event =>  {
    console.log(event.target.value)
    if (event.target.value.uuid.indexOf(".json") != -1) {
      this.getReport(event.target.value.uuid)
      return
    }
    this.setState({
      selectedReport: event.target.value,
      report: event.target.value
    })
  }

  render() {
    return (
      <div>
        <Card>
          <CardContent>
            <FormControl className="text-field">
              <InputLabel>{this.state.report.data.length != 0 ? this.state.report.uuid : "Select a report"}</InputLabel>
              <Select className="text-field" value={this.state.selectedReport} onChange={this.handleChange}>
                {this.state.data.map( r => {
                  return (
                    <ListItem key={r.uuid}  value={r}>
                      <ListItemText primary={r.uuid} secondary={r.timestamp}>
                      </ListItemText>
                    </ListItem>
                  ) 
                })}
              </Select>

            </FormControl>

          </CardContent>
        </Card>

        <ReportsDashboard 
          data={this.state.report.data} 
          config={this.state.report.config}
          stats={this.state.report.stats}
          ></ReportsDashboard>
      </div>
    )
  }
}
export default withSnackbar(Reports)