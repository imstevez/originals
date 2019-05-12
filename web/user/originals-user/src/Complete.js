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
import Avatar from '@material-ui/core/Avatar';

const styles = theme => ({
    root: {width: 'auto', display: 'block',padding: '50px 0'},
    paper: {width: 600, margin: 'auto', padding: '40px 30px',},
    input: {width: 400},
    button: {margin: 20, width: 150,},
    rightIcon: {marginLeft: theme.spacing.unit,},
    divider: {margin: '0 0 20px 0'},
    bigAvatar: {margin: 10, width: 100, height: 100, cursor: 'pointer'},
    fileInput: {display: 'none',},
});

class Complete extends React.Component {
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <CssBaseline/>
                <Paper className={classes.paper} elevation={1}>
                    <Typography component="h1" variant="h5" color="textPrimary" gutterBottom>
                        Complete
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <Grid container justify="center">
                        <input
                            accept="image/*"
                            className={classes.fileInput}
                            id="contained-button-file"
                            multiple
                            type="file"
                        />
                        <label htmlFor="contained-button-file">
                            <Avatar alt="Remy Sharp" src="go.jpg" className={classes.bigAvatar} />
                        </label>
                    </Grid>
                    <Grid container justify="center">
                        <TextField id="standard-search" disabled label="Email" type="text" value="stevzhang01@gmail.com" className={classes.input} margin="normal"/>
                        <TextField id="standard-search" label="Nickname" type="text" className={classes.input} margin="normal"/>
                        <TextField id="standard-search" label="Password" type="password" className={classes.input} margin="normal"/>
                        <Button variant="contained" color="primary" className={classes.button} size="large">
                            COMPLETE
                        </Button>
                    </Grid>
                </Paper>
            </div>
        );
    }
}

Complete.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Complete);
