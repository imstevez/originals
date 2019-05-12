import React from 'react';
import PropTypes from 'prop-types';
import {withStyles} from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Grid from "@material-ui/core/Grid";
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';

const styles = theme => ({
    root: {width: 'auto', display: 'block', padding: '50px 0'},
    paper: {width: 600, margin: 'auto', padding: '40px 30px',},
    input: {width: 400},
    button: {margin: 20, width: 150,},
    rightIcon: {marginLeft: theme.spacing.unit,},
    divider: {margin: '0 0 10px 0'},
    bigAvatar: {margin: 10, width: 100, height: 100, cursor: 'pointer'},
    fileInput: {display: 'none',},
});

const obj2params = obj =>  {
    let result = '';
    let item;
    for(item in obj){
        result += '&' + item + '=' +encodeURIComponent(obj[item]);
    }
    if(result) {
        result = result.slice(1);
    }
    return result;
};

const storage = window.localStorage;

class Complete extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            emailValidate: false,
            passValidate: false,
            email: "",
            password: "",
            emailErr: false,
            passwordErr: false
        };
    }
    componentDidMount() {
        if(storage.token !== "") {
            this.props.history.push("/profile")
        }
    }

    handleEmail(e){
        let reg = new RegExp("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$");
        this.setState({email: e.target.value})
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
    handlePass(e){
        this.setState({password: e.target.value})
        if(this.state.password !== ""){
            this.setState({
                passwordErr: false,
                passValidate: true
            });
        } else {
            this.setState({
                passwordErr: false,
                passValidate: false
            })
        }
    }
    handleLogin(){
        if(this.state.emailValidate && this.state.passValidate){
            fetch("http://localhost:8080/user/login", {
                method: "POST",
                mode: "cors",
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: obj2params({
                    email: this.state.email,
                    password: this.state.password
                }),
            }).then(rsp => {
                return rsp.json()
            }).then(data => {
                if(data.code === 200) {
                    storage.token = data.result.token;
                    this.props.history.push("/profile");
                } else if(data.code === 402) {
                    this.setState({passwordErr: true});
                    alert(data.message);
                } else if(data.code === 401) {
                    this.setState({emailErr: true});
                    alert(data.message);
                } else {
                    alert(data.message);
                }
            }).catch( err => {
                alert(err);
            });
        } else {
            this.setState({
                emailErr: !this.state.emailValidate,
                passwordErr: !this.state.passValidate
            })
        }

    }
    render() {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <CssBaseline/>
                <Paper className={classes.paper}>
                    <Typography component="h1" variant="h5" color="textPrimary" gutterBottom align="right">
                        Login | <a href="/register">Register</a>
                    </Typography>
                    <Divider variant="middle" className={classes.divider} />
                    <Grid container justify="center">
                        <TextField error={this.state.emailErr} id="login-email" label="Email" value={this.state.email} onChange={e => this.handleEmail(e)} onBlur={e => this.handleEmail(e)} type="text" className={classes.input} margin="normal"/>
                        <TextField error={this.state.passwordErr} id="login-pwd" label="Password" value={this.state.password} onChange={e => this.handlePass(e)} onBlur={e => this.handlePass(e)} type="password" className={classes.input} margin="normal"/>
                        <Button variant="contained" color="primary" className={classes.button} size="large" onClick={e => this.handleLogin()}>
                            LOGIN
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
