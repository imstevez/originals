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
import Link from "@material-ui/core/Link";



const styles = (theme) => ({
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

const loginApi = "http://www.koogo.net:8080/user/auth/login";
const storage = window.localStorage;

class Login extends React.Component {
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
            password: "",
            emailErr: false,
            passwordErr: false,
            disabled: false,
        }
    }
    componentDidMount() {
        if (storage.token) {
            this.props.history.push("/profile")
        }
    }

    handleEmail(e) {
        this.setState({
            email: e.target.value,
            emailErr: false
        })
    }
    handlePassword(e) {
        this.setState({
            password: e.target.value,
            passwordErr: false
        })
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
            this.handleLogin();
            e.preventDefault();
        }
    };
    handleLogin() {
        if(!this.state.email) {
            this.setState({
                emailErr: true,
                info: {
                    open: true,
                    variant: "error",
                    message: "邮箱不能为空!"
                }
            });
            return
        }
        if(!this.state.password) {
            this.setState({
                passwordErr: true,
                info: {
                    open: true,
                    variant: "error",
                    message: "密码不能为空!"
                }
            });
            return;
        }
        fetch(loginApi, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: this.state.email,
                password: this.state.password
            })
        }).then(rsp => {
            if(rsp.status === 200) {
                let token = rsp.headers.get("X-Login-token");
                storage.token = token;
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
                        message: "登陆成功!"
                    }
                });

                let that = this;
                setTimeout(function () {
                    that.props.history.push("/profile");
                }, 1000);

                return;
            }
            if(data.code === 300) {
                this.setState({
                    emailErr: true,
                    passwordErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "邮箱或密码格式不正确!"
                    }
                });
                return;
            }
            if(data.code === 303) {
                this.setState({
                    emailErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "用户不存在!"
                    }
                });
                return;
            }
            if(data.code === 304) {
                this.setState({
                    passwordErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "密码错误!"
                    }
                });
                return;
            }

            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "error",
                    message: "未知错误!"
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
            <div className={classes.container} onSubmit={e => {
                alert(e)
            }}>
                <CssBaseline/>
                <LinearProgress hidden={!this.state.disabled} className={classes.progress} />
                <Snackbar
                    anchorOrigin={{
                        vertical: 'top',
                        horizontal: 'center',
                    }}
                    open={this.state.info.open}
                    autoHideDuration={1000}
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
                    <Divider variant="middle" className={classes.divider}/>
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
                        <TextField
                            className={classes.textField}
                            id="login-password"
                            variant="filled"
                            label="密码"
                            autoComplete="new-password"
                            type="password"
                            value={this.state.password}
                            onChange={e => this.handlePassword(e)}
                            onBlur={e => this.handlePassword(e)}
                            error={this.state.passwordErr}
                            disabled={this.state.disabled}
                            required
                        />
                        <Button
                            variant="outlined" color="primary"
                            className={classes.button}
                            onClick={() => this.handleLogin()}
                            disabled={this.state.disabled}
                        >
                            登陆
                        </Button>
                        <Link
                            component="button"
                            variant="body2"
                            onClick={() => {
                                this.props.history.push("/register");
                            }}
                        >
                            没有账号, 注册
                        </Link>
                    </form>
                </Paper>
            </div>
        );
    }
}

Login.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Login);
