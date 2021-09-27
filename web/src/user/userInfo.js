import React from "react";
// import {TestData} from "./data"
import Tree from '@naisutech/react-tree'
import axios from 'axios';

const urlGetUser = "http://localhost:8081/accounts_and_user"
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
        console.log("userID", userID)
        if (userID === ""){
          console.log("wrong data")
          console.log("userData: ", userID)
          return
        } 

        axios({
          method: 'GET',
          url: `${urlGetUser}/${userID}`,
          headers: {
            'Access-Control-Allow-Origin': '*',
            "Authorization": `bearer ${this.props.token}`
          },
          data: {},
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
        "label": `Account id: ${el.id}; Balance: ${el.balance}`,
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