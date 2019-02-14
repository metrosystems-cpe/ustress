
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
      .get(`http://${window.location.host}:8080/ustress/api/v1/file_reports?file=` + report_id)
      .then(response => {
        console.log(response.data)
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

  getReports = () => {
    Axios
    .get(`http://${window.location.host}/ustress/api/v1/reports`)
    .then(response => {
      console.log(response)
      this.setState({...this.state, data: this.parseReports(response.data.entries)})
    })
    .catch(error => {
      console.log(error)
      console.log(error.response)
      if (error.response && error.response.status === 400) {
        Axios.get(`http://${window.location.host}/ustress/api/v1/file_reports`).then(res => {
          this.setState({data: res.data.length > 0 ? res.data : []})
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