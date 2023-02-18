import React from 'react';
import './App.css';
import { BrowserRouter, Route } from 'react-router-dom'
import { Login, Users, Register } from './pages';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Route path={'/'} exact component={Users}/>
        <Route path={'/login'} exact component={Login}/>
        <Route path={'/register'} exact component={Register}/>
      </BrowserRouter>
    </div>
  );
}

export default App;
