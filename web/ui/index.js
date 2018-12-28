// var dataFiles = {
//     data: function () {
//         return {
//             selected: '',
//             reports: [],
//         }
//     },
//     mounted() {
//         axios
//             // .get('/api/v1/reports?file=4a0e87b9-1940-4486-af2c-2de55dded5f0.json')
//             .get('/api/v1/reports')
//             .then(response => {
//                 this.reports = response.data
//                 // this.reports = response.data.map(element => element.file);
//                 // this.files = response.data.map(element => element.file);
//             })
//             .catch(error => {
//                 console.log(error)
//                 this.errored = true
//             })
//             .finally(() => this.loading = false)
//     },
//     template: `
//             <ul>
//                 <li v-for="(report, key) in reports">
//                     {{ key }} {{ report.file }} <br> {{ report.time }}
//                 </li>
//               </ul>`
// }


// Vue.component('line-chart', {
//     extends: VueChartJs.Line,
//     mounted() {
//         this.renderChart({
//             labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
//             datasets: [{
//                 label: 'Data One',
//                 backgroundColor: '#f87979',
//                 data: []
//             }]
//         }, {
//             responsive: true,
//             maintainAspectRatio: false
//         })
//     }
// })


var mapQueryParams = function (url) {
    if (url == "") return {};
    var queryMap = {};

    for (var i = 0; i < url.length; ++i) {
        var p = url[i].split('=', 2);
        if (p.length == 1)
            queryMap[p[0]] = "";
        else
            queryMap[p[0]] = decodeURIComponent(p[1].replace(/\+/g, " "));
    }
    return queryMap;
}

// bootstrap the demo
var workers = new Vue({
    el: '#monkey-data',
    components: {
        // 'data-files': dataFiles,
    },
    props: {
        source: String,
    },
    data() {
        return {
            loading: true,
            errored: false,
            searchQuery: '',
            drawer: "",
            configInfo: false,
            reportID: null,
            monkeyWorkerDataTableHeader: [{
                    text: 'request',
                    align: 'right',
                    value: 'request',
                    width: '90'
                },
                {
                    text: 'status',
                    align: 'right',
                    value: 'status',
                    width: '90'
                },
                {
                    text: 'thread',
                    align: 'right',
                    value: 'thread',
                    width: '90'
                },
                {
                    text: 'duration',
                    align: 'right',
                    value: 'duration',
                    width: '150'
                },
                {
                    text: 'error',
                    align: 'left',
                    value: 'error',
                }
            ],
            monkeyData: {},
            monkeyWorkerDataTableData: [],
            reports: [],
            chartData: [40, 39, 10, 40, 39, 80, 40]
        }
    },
    methods: {
        getReport: function (value) {
            axios
                .get('/restmonkey/api/v1/reports?file=' + value)
                .then(response => {
                    this.monkeyData = response.data
                    this.monkeyWorkerDataTableData = response.data.data
                })
                .catch(error => {
                    console.log(error)
                    this.errored = true
                })
                .finally(() => this.loading = false)
        },
        formatTimeStamp: function (value) {
            var date = new Date(value);
            var formattedDate = ('0' + date.getDate()).slice(-2) + '.' + ('0' + (date.getMonth() + 1)).slice(-2) + '.' + date.getFullYear() + ', ' + ('0' + date.getHours()).slice(-2) + ':' + ('0' + date.getMinutes()).slice(-2) + ':' + ('0' + date.getSeconds()).slice(-2);
            //  + ':' + ('0' + date.getMilliseconds());
            return formattedDate
        },
        removeChip: function (value) {
            this.reportID = null;
        },
    },
    mounted() {
        var params = mapQueryParams(window.location.search.substr(1).split('&'));
        if (params.report_id) {
            this.reportID = params.report_id + ".json";
            this.getReport(this.reportID);
        }
        axios
            .get('/restmonkey/api/v1/reports')
            .then(response => {
                this.reports = response.data
            })
            .catch(error => {
                console.log(error)
                this.errored = true
            })
            .finally(() => this.loading = false)
    }
})