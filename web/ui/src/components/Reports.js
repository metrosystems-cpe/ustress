
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
        this.setState({report: response.data})
      })
      .catch(error => {
        this.errored = true
      })
      .finally(() => this.loading = false)



  }

  getReports = () => {
    Axios
    .get('http://localhost:8080/ustress/api/v1/reports')
    .then(response => {
      this.setState({data: response.data})
    })
    .catch(error => {
      console.error(error)
    })
    .finally(() => this.loading = false)

  }
  handleChange = event =>  {
    this.setState({
      selectedReport: event.target.value
    })
    this.getReport(event.target.value)


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
                    <ListItem key={r.file}  value={r.file}>
                      <ListItemText primary={r.file} secondary={r.time}>
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