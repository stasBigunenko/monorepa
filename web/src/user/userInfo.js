import React from "react";


export class User extends React.Component {
    constructor(props){
      super(props);
      this.handleChange = this.handleChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
  
      this.state = {
        userID: "",
      };
    }

    handleChange = (e) => {
        this.setState({
          [e.target.name]: e.target.value
        })
      }

    handleSubmit(e) {
        e.preventDefault();

        console.log('token from props', this.props.token)
        // console.log('token', props.token)
        console.log("sumbit", e)
    }
    

    render() {
        return (
            <div>
            <form onSubmit={this.handleSubmit}>
              <label>
                User ID:
                <input type="text" name="userID" onChange={this.handleChange} />
              </label>

              <br/>

              <button type="submit">Get user</button>
            </form>
          </div>        
          )
    }
}    