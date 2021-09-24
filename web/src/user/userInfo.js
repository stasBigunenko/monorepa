import React from "react";
import {TestData} from "./data"
import Tree from '@naisutech/react-tree'

const urlGetUser = "http://localhost:8081/users"
export class User extends React.Component {
    constructor(props){
      super(props);
      this.handleChange = this.handleChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
  
      this.state = {
        userID: "",
        userData: []
      };
    }

    handleChange = (e) => {
        this.setState({
          [e.target.name]: e.target.value
        })
      }

    handleSubmit(e) {
        e.preventDefault();

        let {userID} = this.state;
        if (userID == "" || userData.length == 0){
          console.log("wrong data")
          console.log("userData: ", userID)
          console.log("userData: ", userData)
          return
        } 

        axios({
          method: 'GET',
          url: `${urlGetUser}/${userID}`,
          headers: {
            "Authorization": `bearer ${this.props.token}`
          },
          data: user,
        }).then((response) => {
          console.log("status", response.status)
          console.log("headers", response.headers)
          console.log("body", response.data)

          this.setState({userData: response.data})
  
        }).catch((error) => {
          console.log(error);
        });

        this.setState({show: true })
    }
    

    render() {
      let {userData} = this.state;

      let userInfo = null;
      if (userData.length !== 0) {
        userInfo = treeComponent(transformDataToTree(userData))
      }
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

            {userInfo}
          </div>        
          )
    }
}


function treeComponent(data) {
  return <Tree nodes={data} />
}

function transformDataToTree(data) {
  let {user, accounts} = data;

  let items = accounts.map( (el, id)  => {
    return {
        "id": id,
        "label": `Account id: ${el.userID}; Balance: ${el.balance}`,
        "parentId": 0
      }
    })

  return [
    {
      "id": 0,
      "parentId": null,
      "label": `Username: ${user.name}`,
      "items": items,
    }
  ]
}