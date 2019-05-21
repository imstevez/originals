import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import CssBaseline from "@material-ui/core/CssBaseline";
import Grid from "@material-ui/core/Grid";
import Avatar from "@material-ui/core/Avatar";
import Divider from "@material-ui/core/Divider";
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import ExitToApp from '@material-ui/icons/ExitToApp';
import MySnackbarContent from "./MySnackbarContent";
import Snackbar from "@material-ui/core/Snackbar";


const styles = theme => ({
    root: {flexGrow: 1,},
    toolbar: {display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center'},
    wrapMenu: {display: 'flex', flexDirection: 'row', justifyContent: 'center', alignItems: 'center'},
    menuButton: {marginLeft: -12, marginRight: 20},
    paper: {width: 450, margin: 'auto', padding: 30, marginTop: 20, textAlign: 'center'},
    rightIcon: {marginLeft: theme.spacing.unit,},
    divider: {margin: '0 0 20px 0'},
    bigAvatar: {margin: 10, width: 100, height: 100, cursor: 'pointer'},
    fileInput: {display: 'none',},
    b: {display: 'flex', width: 300}
});

const storage = window.localStorage;
class ButtonAppBar extends React.Component {
    image_url;
    user_id;
    constructor(props) {
        super(props);
        this.state = {
            avatar: "",
            id: "",
            email: "",
            nickname: "",
            info: {
                open: false,
                variant: "success",
                message: ""
            },
            disabled: false
        }
    }
    redirectLogin(){
        storage.token = "";
        this.props.history.push("/login");
    }
    handleLogout() {
        fetch("http://www.koogo.net:8080/user/auth/logout", {
            method: "POST",
            mode: "cors",
            headers: {
                "x-login-token": storage.token,
            }
        }).then(rsp => {
            let freshToken = rsp.headers.get("x-login-token")
            if(freshToken) {
                storage.token = freshToken;
            }
            return rsp.json()
        }).then(data => {
            storage.token = "";
            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "success",
                    message: "注销登陆成功!"
                }
            });

            let that = this;
            setTimeout(function () {
                that.props.history.push("/login");
            }, 1000);
        }).catch(err => {
            storage.token = "";
            this.setState({
                disabled: false,
                info: {
                    open: true,
                    variant: "error",
                    message: "注销登陆失败!"
                }
            });
        })
    }
    componentDidMount() {
        if (storage.token === "") {
            this.props.history.push("/login")
        } else {
            fetch("http://www.koogo.net:8080/user/info/", {
                headers: {
                    "x-login-token": storage.token
                },
                mode: "cors"
            }).then(rsp => {
                return rsp.json()
            }).then(data => {
                if(data.code === 200) {
                    this.setState({
                        email: data.data.email,
                        id: data.data.user_id,
                        nickname: data.data.nickname,
                        avatar: data.data.avatar,
                    })
                } else {
                    alert(data.message);
                    this.redirectLogin()
                }
            }).catch(err => {
                alert(err);
                this.redirectLogin();
            });
        }
    }

    render(){
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <CssBaseline/>
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
                <AppBar position="static">
                    <Toolbar variant="dense" className={classes.toolbar}>
                        <Grid className={classes.wrapMenu}>
                            <IconButton className={classes.menuButton} color="inherit" aria-label="Menu">
                                <MenuIcon />
                            </IconButton>
                            <Typography variant="h6" color="inherit">
                                用户信息
                            </Typography>
                        </Grid>
                        <Button color="inherit" onClick={() => {this.handleLogout()}}><ExitToApp/>退出登陆</Button>
                    </Toolbar>
                </AppBar>
                <Paper className={classes.paper} elevation={1}>
                    <Grid container justify="center">
                        <Avatar src={this.state.avatar} className={classes.bigAvatar} />
                    </Grid>
                    <Divider variant="middle" className={classes.divider} />
                    <Table className={classes.table}>
                        <TableBody>
                            <TableRow>
                                <TableCell component="th" scope="row">
                                    用户ID
                                </TableCell>
                                <TableCell align="right">{this.state.id}</TableCell>
                            </TableRow>
                            <TableRow>
                                <TableCell component="th" scope="row">
                                    昵称
                                </TableCell>
                                <TableCell align="right">{this.state.nickname}</TableCell>
                            </TableRow>
                            <TableRow>
                                <TableCell component="th" scope="row">
                                    邮箱
                                </TableCell>
                                <TableCell align="right">{this.state.email}</TableCell>
                            </TableRow>
                        </TableBody>
                    </Table>
                </Paper>
            </div>
        );
    }

}

ButtonAppBar.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ButtonAppBar);
