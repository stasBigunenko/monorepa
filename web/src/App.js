import React from "react";
import {Login} from './login/login'
import {User} from './user/userInfo'

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
    let {token} = this.state;
    
    return (
      <div>
        <Login handleUserLogin={this.onLogin}/>
        <hr/>
        <User token={token}/>
      </div>
    );
  }
}

export default App;
