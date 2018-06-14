// when an update is received via ws connection, we update the model
var socket;
// var socketResponse = {}

socket = new WebSocket("ws://" + location.host + "/restmonkey/api/v1/ws");
socket.onopen = function (event) {
    store.wsConn(true);
}
socket.onmessage = function (event) {
    var monkeyData = JSON.parse(event.data);
    store.setMessageAction(event);
}
socket.onerror = function (e) {
    console.log(e); //TODO: Remove in production
}
socket.onclose = function () {
    store.wsConn(false);
}


function validateNum(input, min, max) {
    var num = +input;
    return num >= min && num <= max && input === num.toString();
}



var store = {
    debug: false,
    state: {
        wsConnection: false,
        monkeyData: {},
        monkeyWorkerDataTableData: [],
    },
    wsConn(status) {
        if (this.debug) console.log('setMessageAction triggered with', newValue)
        this.state.wsConnection = status;
    },

    setMessageAction(newValue) {
        if (this.debug) console.log('setMessageAction triggered with', newValue)
        // console.log(newValue);
        foo = JSON.parse(newValue.data);
        socketResponse = JSON.parse(foo);
        // console.log(socketResponse.timestamp);

        // console.log(socketResponse);
        // console.log(typeof socketResponse);
        this.state.monkeyData = socketResponse;
        this.state.monkeyWorkerDataTableData = socketResponse.data;
        // console.log(this.state.monkeyData);
        // console.log(this.state.monkeyWorkerDataTableData);
    },
    clearMessageAction() {
        if (this.debug) console.log('clearMessageAction triggered')
        this.state.monkeyData = {}
        this.state.monkeyWorkerDataTableData = []
    }
}


var worker = new Vue({
    debug: true,
    el: '#monkey-data',
    components: {
        // 'data-files': dataFiles,
    },
    data() {
        return {
            store: store.state,
            searchQuery: '',
            drawer: "",
            monkeyconfig :{
                url: 'http://' + location.host + '/restmonkey/api/v1/test',
                requests: 16,
                threads: 4,
                insecure: false,
                resolve: ''
            },
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
            IPaddressRule: [
                v => !!v || 'IP:PORT is required',
                v => /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):[0-9]{1,5}$/.test(v) || 'IP:PORT is required'
            ],
        }
    },
    methods: {
        formatTimeStamp: function (value) {
            var date = new Date(value);
            var formattedDate = ('0' + date.getDate()).slice(-2) + '/' + ('0' + (date.getMonth() + 1)).slice(-2) + '/' + date.getFullYear() + ' ' + ('0' + date.getHours()).slice(-2) + ':' + ('0' + date.getMinutes()).slice(-2) + ':' + ('0' + date.getSeconds()).slice(-2) + ':' + ('0' + date.getMilliseconds());
            return formattedDate
        },
        submitNewVictim: function () {
            // TODO SOME VALIDATION
            console.log(this.monkeyconfig)
            console.log(JSON.stringify(this.monkeyconfig))
            socket.send(JSON.stringify(this.monkeyconfig))
        },
        clearSubmitForm: function() {
            if (this.debug) console.log('clearMessageAction triggered')
            this.monkeyconfig =  {
                url: '',
                requests: 4,
                threads: 4,
                insecure: false,
                resolve: ''
            }
        }
    }
    // ,
    // mounted() {
    //     initWS();
    // }
})