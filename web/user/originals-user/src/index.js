import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route } from "react-router-dom";
import Register from "./Register"
import Login from "./Login"
import Profile from "./Profile"
import Complete from "./Complete"

class App extends React.Component {
    render() {
        return (
            <Router>
                <Route path="/" exact component={Profile}/>
                <Route path="/Register" component={Register}/>
                <Route path="/Complete/:token" component={Complete}/>
                <Route path="/login" component={Login}/>
                <Route path="/profile" component={Profile}/>
            </Router>
        );
    }
}

ReactDOM.render(<App />, document.querySelector('#root'));