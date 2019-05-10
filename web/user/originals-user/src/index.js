import React from 'react';
import ReactDOM from 'react-dom';
import ButtonAppBar from './button_app_bar'
import Paper from '@material-ui/core/Paper';
import Invite from './invite.js'
import { styled } from '@material-ui/styles';

const MyPaper = styled(Paper)({
    with: '40%',
    height: '400px',
    margin: '20px',
    padding: '30px'
});


function App() {
    return (
        <div>
            <ButtonAppBar/>
            <MyPaper>
                <Invite/>
            </MyPaper>

        </div>
    );
}

ReactDOM.render(<App />, document.querySelector('#root'));
