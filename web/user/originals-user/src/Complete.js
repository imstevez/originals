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

class Complete extends React.Component {
    constructor(props) {
        super(props);
        let err = {
            nickname: false,
            password: false
        };
        let data = {
            token: "",
            email: "",
            nickname: "",
            password: "",
            avatar: "./go.jpg"
        };
        let validate = {
            nickname: false,
            password: false
        };
        this.state = {
            err: err,
            validate: validate,
            data: data
        }
    }

    componentDidMount() {
        let payload = this.props.match.params.token.split(".")[1];
        let tokenObj = JSON.parse(atob(payload));
        let now = (new Date()).getTime();
        if (!tokenObj.email) {
            alert("token invalid");
            this.props.history.push("/register");
            return;
        }
        // if (now > tokenObj.exp * 1000) {
        //     alert("token expired, please register again");
        //     this.props.history.push("/register");
        //     return;
        // }
        let data = this.state.data;
        data.token = this.props.match.params.token;
        data.email = tokenObj.email;
        this.setState({
            data: data
        });
    }

    handleName(e) {
        let data = this.state.data;
        data.nickname = e.target.value;
        this.setState({
            data: data,
        });
        if(this.state.data.nickname === "") {
            let err = this.state.err;
            let validate = this.state.validate;
            err.nickname = false;
            validate.nickname = false;
            this.setState({
                err: err,
                validate: validate
            });
        } else if (this.state.data.nickname.length < 4) {
            let err = this.state.err;
            let validate = this.state.validate;
            err.nickname = true;
            validate.nickname = false;
            this.setState({
                err: err,
                validate: validate
            });
        } else {
            let err = this.state.err;
            let validate = this.state.validate;
            err.nickname = false;
            validate.nickname = true;
            this.setState({
                err: err,
                validate: validate
            });
        }
    }

    handlePassword(e) {
        let data = this.state.data;
        data.password = e.target.value;
        this.setState({
            data: data,
        });
        if(this.state.data.password === "") {
            let err = this.state.err;
            let validate = this.state.validate;
            err.password = false;
            validate.password = false;
            this.setState({
                err: err,
                validate: validate
            });
        } else if (this.state.data.password.length < 4) {
            let err = this.state.err;
            let validate = this.state.validate;
            err.password = true;
            validate.password = false;
            this.setState({
                err: err,
                validate: validate
            });
        } else {
            let err = this.state.err;
            let validate = this.state.validate;
            err.password = false;
            validate.password = true;
            this.setState({
                err: err,
                validate: validate
            });
        }
    }
    handleSubmit() {
        if(this.state.validate.nickname && this.state.validate.password) {
            fetch("http://localhost:8080/user/register", {
                method: "POST",
                mode: "cors",
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
                body: obj2params({
                    token: this.state.data.token,
                    nickname: this.state.data.nickname,
                    password: this.state.data.password
                }),
            }).then(rsp => {
                return rsp.json();
            }).then(data => {
                if(data.code === 200) {
                    alert("complete success");
                    this.props.history.push("/login");
                } else {
                    alert(data.message);
                }
            }).catch(err => {
                alert(err);
            });
        }
    }
    render() {
        const { classes } = this.props;
        const data = this.state.data;
        const err = this.state.err;
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
                            <Avatar alt="" src={data.avatar} className={classes.bigAvatar} />
                        </label>
                    </Grid>
                    <Grid container justify="center">
                        <TextField id="standard-search" disabled label="Email" type="text" value={data.email} className={classes.input} margin="normal"/>
                        <TextField id="standard-search" error={err.nickname} value={data.nickname} onChange={e => this.handleName(e)}  label="Nickname" type="text" className={classes.input} margin="normal"/>
                        <TextField id="standard-search" error={err.password} value={data.password} onChange={e => this.handlePassword(e)}  label="Password" type="password" className={classes.input} margin="normal"/>
                        <Button variant="contained" onClick={() => this.handleSubmit()} color="primary" className={classes.button} size="large">
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
