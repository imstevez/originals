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
import Avatar from "@material-ui/core/Avatar";
import FormControl from "@material-ui/core/FormControl";
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
        margin: '0 0 20px 0'
    },
    form: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center'
    },
    textField: {
        width: '70%',
        background: 'transparent'
    },
    avatarLabel: {
        textAlign: 'center',
        cursor: 'pointer',
        marginBottom: 10
    },
    button: {
        fontSize: '18px',
        height: 45,
        width: '40%',
        marginTop: 30,
        marginBottom: 5,
    },
    bigAvatar: {width: 100, height: 100, cursor: 'pointer'},
    fileInput: {display: 'none',},
});

const completeApi = "http://www.koogo.net:8080/user/auth/complete";

class Complete extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            info: {
                open: false,
                variant: "success",
                message: ""
            },
            avatar: null,
            avatarSrc: "http://www.koogo.net:8080/user/statics/avatar/avatar_default.png",
            email: "",
            nickname: "",
            password: "",
            passwordConfirm: "",
            nicknameErr: false,
            passwordErr: false,
            passwordConfirmErr: false,
            disabled: false,
        }
    }
    componentDidMount() {
        let payload = this.props.match.params.token.split(".")[1];
        let tokenObj = null;
        try {
            tokenObj = JSON.parse(atob(payload));
        } catch (e) {
            let that = this;
            setTimeout(function () {
                that.props.history.push("/register");
            }, 1500);
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "注册验证token非法!"
                }
            });
            return;
        }
        if (!tokenObj.email) {
            let that = this;
            setTimeout(function () {
                that.props.history.push("/register");
            }, 1500);
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "注册验证token非法!"
                }
            });
            return;
        }
        let now = (new Date()).getTime();
        if (now > tokenObj.exp * 1000) {
            let that = this;
            setTimeout(function () {
                that.props.history.push("/register");
            }, 1500);
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "注册验证token已过期, 请重新注册!"
                }
            });
            return;
        }
        console.log(tokenObj.email);
        this.setState({
            email: tokenObj.email,
            token: this.props.match.params.token
        });
    }
    handleAvatar(e){
        let src = URL.createObjectURL(e.target.files[0]);
        this.setState({
            avatar: e.target.files[0],
            avatarSrc: src
        });
    }
    handleNickname(e) {
        let err = false;
        let nicknameValue = e.target.value;
        if (nicknameValue && nicknameValue < 4) {
            err = true
        }
        this.setState({
            nickname: nicknameValue,
            nicknameErr: err
        })
    }
    handlePassword(e) {
        let passwordValue = e.target.value;
        if (passwordValue === "") {
            this.setState({
                password: "",
                passwordErr: false
            });
            return;
        }
        if (passwordValue.length < 6) {
            this.setState({
                password: passwordValue,
                passwordErr: true,
            });
            return;
        }
        this.setState({
            password: passwordValue,
            passwordErr: false
        })
    }
    handlePasswordConfirm(e) {
        let passwordConfirmValue = e.target.value;
        if (passwordConfirmValue !== this.state.password) {
            this.setState({
                passwordConfirm: passwordConfirmValue,
                passwordConfirmErr: true
            });
            return;
        }
        this.setState({
            passwordConfirm: passwordConfirmValue,
            passwordConfirmErr: false
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
    handleComplete() {
        if (!this.state.nickname || this.state.nicknameErr) {
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "昵称不能为空且长度不能小于4位!"
                }
            });
            return;
        }
        if (!this.state.password || this.state.passwordErr) {
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "密码不能为空且长度不能小于6位!"
                }
            });
            return;
        }
        if (this.state.passwordConfirmErr) {
            this.setState({
                info: {
                    open: true,
                    variant: "error",
                    message: "两次输入密码不-致!"
                }
            });
            return;
        }

        let formData = new FormData();
        formData.append('nickname', this.state.nickname);
        formData.append('password', this.state.password);
        if(this.state.avatar) {
            formData.append('avatar', this.state.avatar);
        }

        fetch(completeApi, {
            method: "POST",
            headers: {
                "x-register-token": this.state.token
            },
            body: formData
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
                let that = this;
                setTimeout(function () {
                    that.props.history.push("/profile");
                }, 1500);
                this.setState({
                    disabled: false,
                    info: {
                        open: true,
                        variant: "success",
                        message: "注册完成, 请登陆!"
                    }
                });
                return;
            }
            if(data.code === 301) {
                let that = this;
                setTimeout(function () {
                    that.props.history.push("/profile");
                }, 1500);
                this.setState({
                    emailErr: true,
                    passwordErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "该邮箱注册已经完成, 请登陆!"
                    }
                });
                return;
            }
            if(data.code === 401) {
                let that = this;
                setTimeout(function () {
                    that.props.history.push("/register");
                }, 1500);
                this.setState({
                    emailErr: true,
                    disabled: false,
                    info: {
                        open: true,
                        variant: "error",
                        message: "注册验证token非法或已过期, 请重新注册!"
                    }
                });
                return;
            }
            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "error",
                    message: "未知错误: " + data.message
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
                        <FormControl>
                            <input
                                accept="image/*"
                                className={classes.fileInput}
                                id="avatar-file"
                                multiple
                                type="file"
                                onChange={e => this.handleAvatar(e)}
                            />
                            <label htmlFor="avatar-file" className={classes.avatarLabel}>
                                <Avatar alt="" src={this.state.avatarSrc} className={classes.bigAvatar} />
                                <Link variant="body2">
                                    点击上传
                                </Link>
                            </label>
                        </FormControl>
                        <TextField
                            className={classes.textField}
                            id="complete-nickname"
                            variant="filled"
                            label="昵称"
                            autoComplete="off"
                            defaultValue={this.state.nickname}
                            error={this.state.nicknameErr}
                            disabled={this.state.disabled}
                            onChange={e => this.handleNickname(e)}
                            onBlur={e => this.handleNickname(e)}
                            required
                        />
                        <TextField
                            className={classes.textField}
                            variant="filled"
                            label="邮箱"
                            autoComplete="off"
                            value={this.state.email}
                            disabled={true}
                        />
                        <TextField
                            className={classes.textField}
                            id="complete-password"
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
                        <TextField
                            className={classes.textField}
                            id="complete-password-confirm"
                            variant="filled"
                            label="密码确认"
                            autoComplete="new-password"
                            type="password"
                            value={this.state.passwordConfirm}
                            onChange={e => this.handlePasswordConfirm(e)}
                            onBlur={e => this.handlePasswordConfirm(e)}
                            error={this.state.passwordConfirmErr}
                            disabled={this.state.disabled}
                            required
                        />
                        <Button
                            variant="outlined" color="primary"
                            className={classes.button}
                            onClick={() => this.handleComplete()}
                            disabled={this.state.disabled}
                        >
                            完成注册
                        </Button>
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
