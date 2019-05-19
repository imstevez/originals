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
class Ca extends React.Component {
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
            avatar: "",
            avatarSrc: "http://www.koogo.net:8080/user/statics/avatar/avatar_default.png"
        };
        let validate = {
            nickname: false,
            password: false
        };
        this.state = {
            disabled: false,
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
        if (now > tokenObj.exp * 1000) {
            alert("token expired, please register again");
            this.props.history.push("/register");
            return;
        }
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
            this.setState({
                disabled: true,
            });
            let formData = new FormData();

            formData.append('nickname', this.state.data.nickname);
            formData.append('password', this.state.data.password);
            if(this.state.data.avatar) {
                formData.append('avatar', this.state.data.avatar);
            }

            fetch('http://www.koogo.net:8080/user/auth/complete', {
                method: 'POST',
                headers: {
                    "x-register-token": this.state.data.token
                },
                body: formData
            }).then(rsp => {
                    return rsp.json();
            }).then(data => {
                if(data.code === 200) {
                    alert("register complete successfully, please login")
                    this.props.history.push("/login")
                } else {
                    alert(data.code + ": " + data.message);
                    if(data.code === 301) {
                        this.props.history.push("/login")
                    }
                    let err = this.state.err;
                    err.nickname = true;
                    err.password = true;
                    this.setState({
                        disabled: false,
                        err: err
                    });
                }
            }).catch(err => {
                this.setState({
                    disabled: false,
                });
                alert(err)
            });
        }
    }

    handleAvatar(e){
        let src = URL.createObjectURL(e.target.files[0]);
        let data = this.state.data;
        data.avatarSrc = src;
        data.avatar = e.target.files[0];
        this.setState({
            data: data
        });
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
                            id="avatar-file"
                            multiple
                            type="file"
                            onChange={e => this.handleAvatar(e)}
                        />
                        <label htmlFor="avatar-file">
                            <Avatar alt="" src={data.avatarSrc} className={classes.bigAvatar} />
                        </label>
                    </Grid>
                    <Grid container justify="center">
                        <TextField id="complete-email" disabled label="Email" type="text" value={data.email} className={classes.input} margin="normal"/>
                        <TextField id="complete-nickname" error={err.nickname} value={data.nickname} onChange={e => this.handleName(e)}  label="Nickname" type="text" className={classes.input} margin="normal"/>
                        <TextField id="complete-password" error={err.password} value={data.password} onChange={e => this.handlePassword(e)}  label="Password" type="password" className={classes.input} margin="normal"/>
                        <Button variant="contained" disabled={this.state.disabled} onClick={() => this.handleSubmit()} color="primary" className={classes.button} size="large">
                            COMPLETE
                        </Button>
                    </Grid>
                </Paper>
            </div>
        );
    }
}

Ca.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Ca);
