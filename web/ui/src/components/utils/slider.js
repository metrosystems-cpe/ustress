import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Slider from '@material-ui/lab/Slider';

const styles = {
    root: {
        width: '100%'
    },
    slider: {
        padding: '22px 0px',
    },
};

class SimpleSlider extends React.Component {
    state = {
        value: 0,
        min: 0,
        max: 100
    };

    render() {
        const { classes } = this.props;
        const { sliderLabel, sliderDescription, val, min, max } = this.props;

        return (
            <div className={classes.root}>
                <Typography id="label">{ sliderLabel } ({ val } {sliderDescription})</Typography>
                <Slider
                    classes={{ container: classes.slider }}
                    value={val}
                    aria-labelledby="label"
                    onChange={this.props.handleChange}
                    min={min}
                    max={max}
                />
            </div>
        );
    }
}

SimpleSlider.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(SimpleSlider);