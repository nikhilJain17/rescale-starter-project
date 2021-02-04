import React from "react";
import logo from "./logo.png";
import "./App.css";

function App() {
    return (
        <div className="App">
            <header className="App-header">
                <img src={logo} className="App-logo" alt="logo" />
                <p>
                    <code>Codebase Transfer Manager</code>
                </p>
                <button onClick={() => {
                    window.api.send('notify', 'Hello there');
                }}>
                    Notify
                </button>
                <button onClick={() => {
                    window.api.send('upload', '');
                }}>
                    Upload File
                </button>
            </header>
        </div>
    );
}

export default App;
