import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Icon from '@material-ui/core/Icon';
import Grid from "@material-ui/core/Grid";
import Typography from '@material-ui/core/Typography';




const styles = theme => ({
    root: {
        width: 'auto',
        display: 'block',
        padding: '50px 0',

    },
    paper: {
        width: 600,
        margin: 'auto',
        padding: '50px 20px',
    },
    input: {
        width: 300
    },
    button: {
        margin: theme.spacing.unit,
    },
    rightIcon: {
        marginLeft: theme.spacing.unit,
    }

});

class Invite extends React.Component {
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <CssBaseline/>
                <Typography component="h1" variant="h4" align="center" color="textPrimary" gutterBottom>
                    Originals Beta v1.0
                </Typography>
                <Paper className={classes.paper} elevation={1}>
                    <Grid container justify="center">
                        <TextField id="standard-search" label="Email" type="email" className={classes.input} margin="normal"/>
                        <Button variant="contained" color="primary" className={classes.button}>
                            GET STATING
                            <Icon className={classes.rightIcon}>send</Icon>
                        </Button>
                    </Grid>
                </Paper>
            </div>
        );
    }
}

Invite.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Invite);
