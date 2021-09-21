import React from "react";
import axios from 'axios';

const urlLogin = "http://127.0.0.1:8080/login";

export class Login extends React.Component {
  constructor(props){
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);

    this.state = {
      name: "",
      password: ""
    };
  }
  
  handleChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value
    })
  }

  handleSubmit(e) {
    e.preventDefault();
    console.log(this.state)

    let {name, password} = this.state;
    let user = {
      user: name,
      password: password
    }

    console.log("start login")

    console.log("request login")
    console.log(urlLogin)
    console.log(user)

    console.log("resp login")
    axios({
        method: 'POST',
        url: urlLogin,
        data: user,
      }).then((response) => {
        console.log("status", response.status)
        console.log("headers", response.headers)
        console.log("body", response.data)

        this.props.handleUserLogin(response.headers.token)
      }).catch((error) => {
        console.log(error);
    });
  };

  render() {
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <label>
            Name:
            <input type="text" name="name" onChange={this.handleChange} />
          </label>
          
          <br/>
          <label>
          Password:
          <input type="text" name="password" onChange={this.handleChange} />
          </label>
          <br/>
          <button type="submit">Add</button>
        </form>
      </div>
    );
  }
}


// render(<Login />, document.getElementById("root"));
