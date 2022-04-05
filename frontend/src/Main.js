import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Feed from "./Feed";
import UserPage from "./UserPage";
import MyPage from "./MyPage";
import UserReport from "./UserReport";
import PostingReport from "./PostingReport";
import Settings from "./Settings";
import ChangePassword from "./ChangePassword";
import UserDelete from "./UserDelete";
import Logout from "./Logout";
import "./Main.css";

function Main() {
  return (
    <div>
      <Router>
        <Switch>
          <Route exact path="/home" component={Feed}></Route>
          <Route exact path="/mypage" render={() => <MyPage />}></Route>
          <Route
            path="/userpage/:userName"
            render={(props) => (
              <UserPage userName={props.match.params.userName} />
            )}
          ></Route>
          <Route exact path="/reports/users/:userName" component={UserReport} />
          <Route
            exact
            path="/reports/postings/:postingID"
            component={PostingReport}
          />
          <Route exact path="/settings" component={Settings}></Route>
          <Route
            exact
            path="/change_password"
            component={ChangePassword}
          ></Route>
          <Route exact path="/delete_user" component={UserDelete} />
          <Route exact path="/logout" component={Logout} />
        </Switch>
      </Router>
    </div>
  );
}

export default Main;
