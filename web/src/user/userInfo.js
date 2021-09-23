import React from "react";
import {TestData} from "./data"
import Tree from '@naisutech/react-tree'

export class User extends React.Component {
    constructor(props){
      super(props);
      this.handleChange = this.handleChange.bind(this);
      this.handleSubmit = this.handleSubmit.bind(this);
  
      this.state = {
        userID: "",
        show: false
      };
    }

    handleChange = (e) => {
        this.setState({
          [e.target.name]: e.target.value
        })
      }

    handleSubmit(e) {
        e.preventDefault();

        this.setState({show: true })
    }
    

    render() {
      let {show} = this.state;

      let userInfo = null;
      if (show) {
        userInfo = treeComponent(transformDataToTree(TestData))
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