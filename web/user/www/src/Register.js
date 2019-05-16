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
    constructor(props) {
        super(props);
        this.state = {
            disabled: false,
            email: "",
            emailErr: false,
            emailValidate: false
        };
    }
    handleEmail(e) {
        this.setState({
            email: e.target.value,
        });
        let reg = new RegExp("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$");
        if(this.state.email !== "" && !reg.test(this.state.email)){
            this.setState({
                emailErr: true,
                emailValidate: false
            });
        } else if (this.state.email === "") {
            this.setState({
                emailErr: false,
                emailValidate: false
            })
        } else {
            this.setState({
                emailErr: false,
                emailValidate: true
            })
        }
    }
    handleSubmit() {
        if(this.state.emailValidate){
            this.setState({
                disabled: true
            });
            fetch("http://www.koogo.net:8080/user/auth/register", {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: this.state.email
                })
            }).then(rsp => {
                return rsp.json();
            }).then(data => {
                if (data.code === 200) {
                    alert("register success, an validation email has been send to your mail box")
                    this.props.history.push("/login");
                } else {
                    this.setState({
                        disabled: false,
                        emailErr: true
                    });
                    alert(data.message);
                }
            }).catch(err => {
                this.setState({
                    disabled: false
                });
                alert(err);
            });
        } else {
            this.setState({
                emailErr: true
            });
        }
    }
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
                        <TextField id="register-email" error={this.state.emailErr} value={this.state.email} onChange={e => this.handleEmail(e)} onBlur={e => this.handleEmail(e)} label="Email" type="email" className={classes.input} margin="normal"/>
                        <Button onClick={() => this.handleSubmit()} disabled={this.state.disabled} variant="contained" color="primary" className={classes.button} size="large">
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
