import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Grid from "@material-ui/core/Grid";
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';




const styles = theme => ({
    root: {width: 'auto', display: 'block', padding: '50px 0'},
    paper: {width: 600, margin: 'auto', padding: 30,},
    input: {width: 300},
    button: {margin: '20px 0 0 10px', height: 42},
    divider: {margin: '0 0 40px 0'}
});

class Register extends React.Component {
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <CssBaseline/>
                <Paper className={classes.paper}>
                    <Typography component="h1" variant="h5" color="textPrimary" align="right" gutterBottom>
                        Register |  <a href="/login">Login</a>
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <Grid container justify="center">
                        <TextField id="standard-search" label="Email" type="email" className={classes.input} margin="normal"/>
                        <Button variant="contained" color="primary" className={classes.button} size="large">
                            REGISTER
                        </Button>
                    </Grid>
                </Paper>
            </div>
        );
    }
}

Register.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Register);
