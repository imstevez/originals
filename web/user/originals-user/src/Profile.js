import React from 'react';
import Typography from "@material-ui/core/Typography";

const storage = window.localStorage;
storage.token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InN0ZXZ6aGFuZzAxQGdtYWlsLmNvbSIsIm1vYmlsZSI6IjEzOTE5MTU0OTI1Iiwibmlja19uYW1lIjoiU3RldkUuWiIsImltYWdlX3VybCI6IiIsImV4cCI6MTU1NzU2NDkxNCwiaWF0IjoxNTU3NTYzNzE0fQ.YJrkmx6SPOxYTfCMmoHdRGOui1_DsQSjDHDy4GwGsdk";

class Profile extends React.Component{
    constructor(props) {
        super(props);
        this.state = {
            user: {}
        };
    }
    componentDidMount() {
        if (storage.token === "") {
            this.props.history.push("/login")
        }
        fetch('http://localhost:8080/user/profile/', {
            method: "GET",
            mode: "cors",
            cache: 'default',
            headers: {
                "x-originals-token": storage.token
            }
        }).then(rsp => {
            return rsp.json()
        }).then(data => {
            if (data.code === 200) {
                console.log("abc")
                this.setState({
                    user: data.result
                })
            } else {
                storage.token = "";
                this.props.history.push("/login")
            }
        }).catch(err => {
            console.log(err);
            this.props.history.push("/login")
        })
    }

    render() {
        return (
            <div>
                <Typography component="h1" variant="h5">
                    {this.state.user.email}
                    {this.state.user.mobile}
                    {this.state.user.nickname}
                    {this.state.user.image_url}
                </Typography>
            </div>
        );
    }
}

export default Profile;
