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
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';




const styles = (theme) => ({
    container: {
        marginTop: 30

    },
    paper: {
        maxWidth: 500,
        margin: 'auto',
        padding: 30,
    },
    divider: {
        margin: 0
    },
    grid: {
        display: 'flex',
        margin: '20px 0',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center'
    },
    formControl: {
        width: '70%',
        [theme.breakpoints.down('sm')]: {
            width: '85%',
        },
    },
    button: {
        height: 40,
        marginTop: 20
    }
});

const emailTips = "E.g. example@example.com";
const registerApi = "http://www.koogo.net:8080/user/auth/register";

class Register extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            disabled: false,
            email: "",
            emailErr: false,
            emailValid: false,
            emailTips: emailTips
        };
    }
    handleEmail(e) {
        let reg = new RegExp("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$");
        if(e.target.value !== "" && !reg.test(this.state.email)){
            this.setState({
                email: e.target.value,
                emailErr: true,
                emailValid: false,
            });
        } else if (e.target.value === "") {
            this.setState({
                email: e.target.value,
                emailErr: false,
                emailValid: false,
            })
        } else {
            this.setState({
                email: e.target.value,
                emailErr: false,
                emailValid: true,
            })
        }
    }
    handleSubmit() {
        if(this.state.emailValid){
            fetch(registerApi, {
                method: "POST",
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: this.state.email
                })
            }).then(rsp => {
                if(rsp.status !== 200) {
                    return {
                        code: 500,
                        message: "internal error"
                    };
                } else {
                    return rsp.json();
                }
            }).then(data => {
                if (data.code === 200) {
                    alert("Success, complete email has been send to your mailbox")
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
            this.setState({
                disabled: true
            });
        } else {
            this.setState({
                emailErr: true
            })
        }
    }
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.container}>
                <CssBaseline/>
                <Paper className={classes.paper}>
                    <Typography component="h3" variant="h6" color="textPrimary" align="right" gutterBottom>
                        Register |  <a href="/login">Login</a>
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <Grid container justify="center" className={classes.grid} aria-disabled={this.state.disabled}>
                        <FormControl className={classes.formControl} error={this.state.emailErr}>
                            <InputLabel htmlFor="register-email">Email</InputLabel>
                            <Input
                                id="register-email"
                                value={this.state.email}
                                onChange={e => this.handleEmail(e)}
                                aria-describedby="register-email-error"
                                disabled={this.state.disabled}
                            />
                            <FormHelperText id="register-email-error">{this.state.emailTips}</FormHelperText>
                        </FormControl>
                        <Button onClick={() => this.handleSubmit()} disabled={this.state.disabled} variant="outlined" color="primary" className={classes.button} size="large">
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
