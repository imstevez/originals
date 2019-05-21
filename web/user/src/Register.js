import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import Cloud from '@material-ui/icons/Cloud';
import TextField from '@material-ui/core/TextField';
import MySnackbarContent from './MySnackbarContent';
import Snackbar from '@material-ui/core/Snackbar';
import LinearProgress from '@material-ui/core/LinearProgress';
import Link from '@material-ui/core/Link';



const styles = theme => ({
    container: {
        height: '100vh',
        background: 'linear-gradient(to right, #348AC7, #7474BF);',
        overflow: 'hidden'
    },
    progress: {
        position: 'fixed',
        width: '100%',
        background: 'transparent'
    },
    paper: {
        maxWidth: 400,
        margin: '50px auto',
        padding: '20px 0 40px 0',
    },
    typography: {
        textAlign: 'center',
        verticalAlign: 'middle',
        fontSize: '18px',
        lineHeight: '40px'
    },
    cloud: {
        verticalAlign: 'top',
        fontSize: '35px',
        marginRight: 10
    },
    divider: {
        margin: '0 0 40px 0'
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
        marginBottom: 5,
    }
});

const registerApi = "http://www.koogo.net:8080/user/auth/register";
const verifyEmail = function(email){
    let reg = /^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$/;
    return reg.test(email);
};

class Register extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            info: {
                open: false,
                variant: "success",
                message: ""
            },
            infoOpen: true,
            email: "",
            emailErr: false,
            disabled: false,
        }
    }
    handleEmail(e) {
        let emailValue = e.target.value;
        if (emailValue === "") {
            this.setState({
                email: "",
                emailErr: false
            });
            return;
        }
        if (emailValue !== "") {
            if (!verifyEmail(emailValue)) {
                this.setState({
                    email: emailValue,
                    emailErr: true
                });
            } else {
                this.setState({
                    email: emailValue,
                    emailErr: false
                });
            }
        }
    }
    handleInfoClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        let orInfo = this.state.info;
        orInfo.open = false;
        this.setState({
            info: orInfo
        });
    };
    handleKeyDown = e => {
        if (e.keyCode === 13) {
            this.handleRegister();
            e.preventDefault();
        }
    };
    handleRegister() {
        if (!this.state.email) {
            this.setState({
                emailErr: true,
                info: {
                    open: true,
                    variant: "error",
                    message: "邮箱不能为空!"
                }
            });
            return;
        }
        if (this.state.emailErr) {
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "邮箱格式错误"
                }
            });
            return;
        }
        fetch(registerApi, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: this.state.email,
            })
        }).then(rsp => {
            if(rsp.status === 200) {
                return rsp.json();
            } else {
                return {
                    code: rsp.status,
                    message: "internal error"
                }
            }
        }).then(data => {
            if(data.code === 200) {
                this.setState({
                    disabled: false,
                    info: {
                        open: true,
                        variant: "success",
                        message: "注册成功, 验证链接已发送至您的邮箱, 请在30分钟以内前往完成注册!"
                    }
                });
                this.setState({
                    email: "",
                });
                return;
            }
            if(data.code === 300) {
                this.setState({
                    emailErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "邮箱格式不正确!"
                    }
                });
                return;
            }
            if(data.code === 301) {
                this.setState({
                    emailErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "该邮箱已被注册!"
                    }
                });
                return;
            }
            if(data.code === 302) {
                this.setState({
                    passwordErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "验证邮件发送失败, 请稍后再试"
                    }
                });
                return;
            }

            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "error",
                    message: "系统错误: " + data.message
                }
            });
        }).catch(err => {
            alert(err);
            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "error",
                    message: "网络错误: " + err
                }
            });
        });
        this.setState({
            disabled: true
        });
    }
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.container}>
                <CssBaseline/>
                <LinearProgress hidden={!this.state.disabled} className={classes.progress} />
                <Snackbar
                    anchorOrigin={{
                        vertical: 'top',
                        horizontal: 'center',
                    }}
                    open={this.state.info.open}
                    autoHideDuration={2000}
                    onClose={this.handleInfoClose}
                >
                    <MySnackbarContent
                        onClose={this.handleInfoClose}
                        variant={this.state.info.variant}
                        message={this.state.info.message}
                    />
                </Snackbar>
                <Paper className={classes.paper}>
                    <Typography
                        component="h1"
                        variant="h6"
                        color="primary"
                        className={classes.typography}
                    >
                        <Cloud className={classes.cloud}/>云记
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <form className={classes.form} onKeyDown={this.handleKeyDown}>
                        <TextField
                            className={classes.textField}
                            id="login-email"
                            variant="filled"
                            label="邮箱"
                            autoComplete="off"
                            value={this.state.email}
                            onChange={e => this.handleEmail(e)}
                            onBlur={e => this.handleEmail(e)}
                            error={this.state.emailErr}
                            disabled={this.state.disabled}
                            required
                        />
                        <Button
                            variant="outlined" color="primary"
                            className={classes.button}
                            onClick={() => this.handleRegister()}
                            disabled={this.state.disabled}
                        >
                            注册
                        </Button>
                        <Link
                            component="button"
                            variant="body2"
                            onClick={() => {
                                this.props.history.push("/login");
                            }}
                        >
                            已注册, 登陆
                        </Link>
                    </form>
                </Paper>
            </div>
        );
    }
}

Register.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Register);
