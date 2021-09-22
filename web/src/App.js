import React from "react";
import {Login} from './login/login'

class App extends React.Component {
  constructor() {
    super();
    this.onLogin = this.onLogin.bind(this);


    // Define the initial state:
    this.state = {
      token: "",
    };
  }

  onLogin = (token) => {
    this.setState({token:token})
  }

  render(){
    console.log("state", this.state)
    return (
      <div>
        <Login handleUserLogin={this.onLogin}/>
        <hr/>
      </div>
    );
  }
}

export default App;
