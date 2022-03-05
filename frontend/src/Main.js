import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Feed from "./Feed";
import UserPage from "./UserPage";
import Settings from "./Settings";
import ChangePassword from "./ChangePassword";
import "./Main.css";

function Main() {
  const loginUserName = localStorage.getItem("loginUserName");
  return (
    <div>
      <Router>
        <Switch>
          <Route exact path="/home" component={Feed}></Route>
          <Route
            exact
            path="/mypage"
            render={() => <UserPage userName={loginUserName} />}
          ></Route>
          <Route
            path="/userpage/:userName"
            render={(props) => (
              <UserPage userName={props.match.params.userName} />
            )}
          ></Route>
          <Route exact path="/settings" component={Settings}></Route>
          <Route
            exact
            path="/change_password"
            component={ChangePassword}
          ></Route>
        </Switch>
      </Router>
    </div>
  );
}

export default Main;
