import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, Divider } from '@material-ui/core';
import PrettyPrint from './prettyprint';

import { Line, Scatter, Bar } from 'react-chartjs-2';

import { Doughnut } from 'react-chartjs-2';
import CustomTable from './table';
import {color_pool} from '../../index';

class ReportsDashboard extends Component {
  constructor(props) {
    super(props)
  }
  state = {
    config: {
      url: '',
      method: '',
      headers: {},
      threads: 0,
      requests: 0,
      resolve: '',
      payload: '',
      withResponse: false,
    },
    data: [],
    stats: {},
    // details: {},
  };

  getChartOptions = () => {
    return {
      tooltips: {
        mode: 'index',
        intersect: false,
      },
      scales: {
        xAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Request number'
          },
        }],
        yAxes: [{
          display: true,
          scaleLabel: {
            display: true, 
            labelString: "Latency",
          },

        }]

      },

      responsiveDownsample: {
        enabled: true,
        /**
         * Choose aggregation algorithm 'AVG'(Average values) or
         * 'LTTB' (Largest-Triangle-Three-Buckets). Default: 'LTTB'
         */
        aggregationAlgorithm: 'LTTB',
        /**
         * The desired minimal distance between data points in pixels.
         * The plugin downsamples the data and tries to match this threshold.
         * Default: 1 pixel
         */
        desiredDataPointDistance: 1,
        /**
         * The minimal number of data points. The chart data is not downsampled further than
         * this threshold. Default: 100
         */
        minNumPoints: 100,
         /**
          * Cull data to displayed range of x scale. Default: true
          */
        cullData: false
      },
      pan: {
        enabled: true,
        mode: 'y',
        drag: true, 
      },
      zoom: {
        enabled: false,
        mode: 'xy',
        drag: false,
      }
    }
  }


  getChartLabels = () => {
    return this.state.data.map(e => {
      return e.request
    })
  }

  getChartData = () => {
    return this.state.data.map( e => {
      return e.duration
    })
  }

  getChartDatasets = (data) => {
    let workers = {}
    data.map( e => {
      if (workers[e.thread] === undefined)  {
        workers[e.thread] = [{x:e.request, y: e.duration}]
      } else {
        workers[e.thread].push({x: e.request, y: e.duration})
      }
    })
    return workers
  }

  getChartData = (data) => {
    let workersData = this.getChartDatasets(data)
    let workers = Object.keys(workersData);
    let dataSetsObj = []

    for (let index = 0; index < workers.length; index++) {
      let worker = workers[index];
      dataSetsObj.push({
        label: `Worker ${worker}`,
        data: workersData[worker], 
        backgroundColor: color_pool[worker],
      })
    }

    return {
      labels: this.getChartLabels(),
      datasets: dataSetsObj
    }


  }

  getBarChartData = (stats) => {

    let keys = Object.keys(stats);
    let ep = keys.indexOf("error_percentage");
    let cc = keys.indexOf("codes_count");
    if (cc != -1) { keys.pop(cc)}
    if (ep != -1) { keys.pop(ep)};

    let data = [];
    keys.map( k => {
      data.push(stats[k])
    })

    return {
      labels:keys, 
      datasets:[{
        data:data,
        backgroundColor: [
          "rgb(171, 228, 247)", 
          "rgb(130, 217, 247)", 
          "rgb(59, 199, 246)",
          "rgb(6, 164, 217)",
          "rgb(0, 149, 198)",
          // "rgb(255, 112, 112)",
        ]

      }]
    }

  }

  getBarChartOptions = () => {
    return {legend:{display:false}}
  }

  getDoughnutData = (data) => {
    if (!data.codes_count) {return}
    let labels = Object.keys(data.codes_count);
    return {
      labels: labels,
      datasets:[{
        data:Object.values(data.codes_count),
        backgroundColor: [
          "rgb(171, 228, 247)",
          "rgb(130, 217, 247)",
          "rgb(59, 199, 246)",
        ]

      }]
    }
  }

  render() {
    let { details, config, stats, data  } = this.props
    this.state.data = data

    const getDetails = (details) => {
      if (details != undefined) {
        return (
        <Card elevation={1} className="paper">
          <Typography variant="title">
            Report details
          </Typography>
          <Divider />

          <PrettyPrint options={details}>
          </PrettyPrint>
        </Card>
        )
      } 
      return (<div></div>)
    }

    const getRPS = (data) => {

      return 1

    }

    return (
      <div>
        {getDetails(details)}

        <Card elevation={1} className="paper">
            
            <Typography variant="title">
              Statistics
            </Typography>
            <Divider className="divider" />
            <PrettyPrint options={stats}>
            </PrettyPrint>
            
        </Card>


        <div className="cstm-flex">

          <Card elevation={1} className="paper text-field-s">

            <Typography variant="title">
              Code Count
            </Typography>
            <Typography variant="caption">
                Number of requests with status code
            </Typography>
            <Divider className="divider" />
            <Doughnut data={this.getDoughnutData(stats)}></Doughnut>

          </Card>
          <Card elevation={1} className="paper text-field-s">

            <Typography variant="title">
              Hist 99th Chart
            </Typography>
            <Typography variant="caption">
              Displays the latency statistics
            </Typography>
            <Divider className="divider" />
            <Bar data={this.getBarChartData(stats)} options={this.getBarChartOptions(data)}>
            </Bar>

          </Card>

        </div>
        <Card elevation={1} className="paper">
            
            <Typography variant="title">
              Latency Chart
            </Typography>
            <Typography variant="caption">
              Displays all the requests made and their durations
            </Typography>
            <Divider className="divider" />
            <Line data={this.getChartData(data)} options={this.getChartOptions(data)}>
            </Line>
            
        </Card>
        <CustomTable data={this.state.data}></CustomTable>

        <Card elevation={1} className="paper">

          <Typography variant="title">
            uStress Config
          </Typography>
          <Divider />
          <PrettyPrint options={config}>
          </PrettyPrint>

        </Card>

      </div>
    )

  }

}

export default ReportsDashboard