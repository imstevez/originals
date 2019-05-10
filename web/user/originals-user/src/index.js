import React from 'react';
import ReactDOM from 'react-dom';
import Button from '@material-ui/core/Button';
import ButtonAppBar from './button_app_bar'

function App() {
    return (
        <div>
            <ButtonAppBar/>
            <Button variant="contained" color="primary">
                Hello World
            </Button>
        </div>
    );
}

ReactDOM.render(<App />, document.querySelector('#root'));
