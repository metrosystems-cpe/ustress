import React, { Component } from 'react';
import { Card, CardContent, Typography } from '@material-ui/core';
import Snackbar from '@material-ui/core/Snackbar';
import { withSnackbar } from 'notistack';


class Home extends Component {

    render() {
        return (
            <Card>
                <CardContent>
                    <Typography variant="title">Documentation</Typography>
                </CardContent>
            </Card>
        )
    }
}
export default withSnackbar(Home);

// export default Home