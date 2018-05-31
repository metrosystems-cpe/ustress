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
            reportID: null,
            monkeyWorkerDataTableHeader: [{
                        text: 'request',
                        align: 'left',
                        value: 'request'
                    },
                    {
                        text: 'status',
                        value: 'status'
                    },
                    {
                        text: 'thread',
                        value: 'thread'
                    },
                    {
                        text: 'duration',
                        value: 'duration'
                    },
                    {
                        text: 'error',
                        value: 'error'
                    }
                ],
                monkeyData: {},
                monkeyWorkerDataTableData: [],
                reports: []
            }
    },
    methods: {
        getReport: function(value) {
            // console.log(value);
            axios
             .get('/api/v1/reports?file=' + value)
                 .then(response => {
                     this.monkeyData = response.data
                     this.monkeyWorkerDataTableData = response.data.data
                 })
                 .catch(error => {
                     console.log(error)
                     this.errored = true
                 })
                 .finally(() => this.loading = false)
        }
    },
    mounted() {
        axios
            // .get('/api/v1/reports?file=4a0e87b9-1940-4486-af2c-2de55dded5f0.json')
            .get('/api/v1/reports')
            .then(response => {

                // try { // it is, so now let's see if its valid JSON
                //     var myJson = resp = JSON.parse(response.data);
                //     // yep, we're working with valid JSON
                // } catch (e) {
                //     // nope, we got what we thought was JSON, it isn't; let's handle it.
                //     console.log(e)
                // }

                this.reports = response.data
                // this.monkeyData = response.data
                // this.monkeyWorkerDataTableData = response.data.data



            })
            .catch(error => {
                console.log(error)
                this.errored = true
            })
            .finally(() => this.loading = false)
    }
})