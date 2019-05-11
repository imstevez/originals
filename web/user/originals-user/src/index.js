import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route } from "react-router-dom";
import Invite from "./Invite"
import Login from "./Login"
import Profile from "./Profile"

function App() {
    return (
        <Router>
            <Route path="/invite" component={Invite}/>
            <Route path="/profile" component={Profile}/>
            <Route path="/login" component={Login}/>
        </Router>
    );
}

ReactDOM.render(<App />, document.querySelector('#root'));