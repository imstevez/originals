import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

const styles = theme => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
        textAlign: 'center'
    },
    textField: {
        width: '50%',
        marginRight: '10px',
    },
    button: {
        marginTop: '16px',
        marginBottom: '8px',
    }
});

class Invite extends React.Component {
    render() {
        const { classes } = this.props;
        return (
            <form className={classes.container} noValidate autoComplete="off">
                <TextField
                    id="outlined-email-input"
                    label="Email"
                    className={classes.textField}
                    type="email"
                    name="email"
                    autoComplete="email"
                    margin="normal"
                    variant="outlined"
                />
                <Button variant="contained" color="primary" className={classes.button}>
                    GET STARTING
                </Button>
            </form>
        );
    }
}

Invite.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Invite);