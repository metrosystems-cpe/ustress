
import React, { Component } from 'react';
import { Typography, Card, CardContent, Select, InputLabel, FormControl,  ListItem, ListItemText } from '@material-ui/core';
import Axios from 'axios';
import PrettyPrint from './utils/prettyprint';
import CustomTable from './utils/table';

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
      .get('http://localhost:8080/ustress/api/v1/reports?file=' + report_id)
      .then(response => {
        this.setState({...this.state, report: response.entries})
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

  getReports = () => {
    Axios
    .get('http://localhost:8080/ustress/api/v1/reports')
    .then(response => {
      this.setState({...this.state, data: this.parseReports(response.data.entries)})
    })
    .catch(error => {
      console.error(error)
    })
    .finally(() => this.loading = false)
  }

  handleChange = event =>  {
    console.log(event.target.value)
    this.setState({
      selectedReport: event.target.value,
      report: event.target.value
    })
    // FILE BACKUP
    // this.getReport(event.target.value)
  }

  render() {
    return (
      <div>
        <Card>
          <CardContent>
            <FormControl className="text-field">
              <InputLabel>Select a report</InputLabel>
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
        <Card className="paper">
          <CardContent>
            <Typography variant="title"> Config </Typography>
            <PrettyPrint options={this.state.report.config}></PrettyPrint>
          </CardContent>
        </Card>
        <Card className="paper">
          <CardContent>
            <Typography variant="title"> Stats </Typography>
            <PrettyPrint options={this.state.report.stats}></PrettyPrint>

          </CardContent>
        </Card>
        <CustomTable data={this.state.report.data}></CustomTable>

      </div>
    )
  }
}
export default Reports