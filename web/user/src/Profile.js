import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import CssBaseline from "@material-ui/core/CssBaseline";
import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import Avatar from "@material-ui/core/Avatar";
import Divider from "@material-ui/core/Divider";

const styles = theme => ({
    root: {flexGrow: 1,},
    grow: {flexGrow: 1,},
    appBar: {height: 60},
    menuButton: {marginLeft: -12, marginRight: 20},
    paper: {width: 600, margin: 'auto', padding: 30, marginTop: 20},
    input: {width: 400},
    button: {margin: 20, height: 50, width: 200,},
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
            return rsp.json()
        }).then(data => {
            alert(data.message);
            localStorage.token = "";
            this.props.history.push("/login");
        }).catch(err => {
            localStorage.token = "";
            alert(err);
            this.props.history.push("/login");
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
                <AppBar position="static" className={classes.appBar}>
                    <Toolbar>
                        <Typography variant="h6" color="inherit" className={classes.grow}>
                            Profile
                        </Typography>
                        <Button color="inherit" onClick={() => {this.handleLogout()}}>LOGOUT</Button>
                    </Toolbar>
                </AppBar>
                <Paper className={classes.paper} elevation={1}>
                    <Grid container justify="center">
                        <Avatar src={this.state.avatar} className={classes.bigAvatar} />
                    </Grid>
                    <Divider variant="middle" className={classes.divider} />
                    <Grid container justify="center" direction="column">
                        <Typography gutterBottom variant="h6" component="h3">{this.state.id}</Typography>
                        <Divider variant="middle" className={classes.divider} />
                        <Typography gutterBottom variant="h6" component="h3">{this.state.nickname}</Typography>
                        <Divider variant="middle" className={classes.divider} />
                        <Typography gutterBottom variant="h6" component="h3">{this.state.email}</Typography>
                    </Grid>
                </Paper>
            </div>
        );
    }

}

ButtonAppBar.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ButtonAppBar);
