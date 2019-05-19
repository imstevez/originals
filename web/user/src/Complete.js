import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import Cloud from '@material-ui/icons/Cloud';
import TextField from '@material-ui/core/TextField'



const styles = (theme) => ({
    container: {
        height: '100vh',
        padding: 50,
        background: 'linear-gradient(to right, #348AC7, #7474BF);',
    },
    paper: {
        maxWidth: 400,
        margin: 'auto',
        padding: '20px 0 40px 0',
    },
    typography: {
        textAlign: 'center',
        fontSize: '18px',
        lineHeight: '18px'
    },
    cloud: {
        fontSize: 70,
        lineHeight: 70,
    },
    divider: {
        margin: '15px 0 40px 0'
    },
    form: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center'
    },
    textField: {
        width: '70%'
    },
    button: {
        fontSize: '18px',
        height: 45,
        width: '40%',
        marginTop: 20,
        marginBottom: 10,
    }
});

const emailTips = "";
const registerApi = "http://www.koogo.net:8080/user/auth/register";

class Complete extends React.Component {
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
                    <Typography
                        component="h1"
                        variant="h6"
                        color="primary"
                        className={classes.typography}
                    >
                        <Cloud className={classes.cloud}/><br/>我记 ● 云账本
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <form className={classes.form}>
                        <TextField
                            className={classes.textField}
                            id="login-email"
                            variant="filled"
                            label=""
                            autoComplete="off"
                            value="stevzhang01@gmail.com"
                            disabled
                        />
                        <TextField
                            className={classes.textField}
                            id="login-email"
                            variant="filled"
                            label="昵称"
                            autoComplete="off"
                            disabled={this.state.disabled}
                        />
                        <TextField
                            className={classes.textField}
                            id="login-passwd"
                            variant="filled"
                            label="密码"
                            type="password"
                            disabled={this.state.disabled}
                        />
                        <Button
                            onClick={() => this.handleSubmit()}
                            variant="outlined" color="primary"
                            disabled={this.state.disabled}
                            className={classes.button}>
                            登陆
                        </Button>
                        <Typography
                            component="h3"
                            variant="button"
                            color="textPrimary"
                        >
                            <a href="/register">还未注册, 注册</a>
                        </Typography>
                    </form>
                </Paper>
            </div>
        );
    }
}

Complete.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Complete);
