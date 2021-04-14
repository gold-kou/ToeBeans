import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Feed from "./Feed";
import MyPage from "./MyPage";
import Settings from "./Settings";
import ChangePassword from "./ChangePassword";

import "./Main.css";

// TODO TypeScript化
function Main() {
  return (
    <div>
      <Router>
        <Switch>
          <Route exact path="/home" component={Feed}></Route>
          <Route exact path="/mypage" component={MyPage}></Route>
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
