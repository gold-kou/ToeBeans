import React, { useEffect } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import axios from "axios";

import Auth from "./Auth";
import Landing from "./Landing";
import UserRegistration from "./UserRegistration";
import UserPage from "./UserPage";
import MyPage from "./MyPage";
import Login from "./Login";
import Logout from "./Logout";
import Main from "./Main";
import UserReport from "./UserReport";
import PostingReport from "./PostingReport";
import Settings from "./Settings";
import ChangePassword from "./ChangePassword";
import UserDelete from "./UserDelete";

axios.defaults.baseURL = process.env.REACT_APP_BACK_BASE_URL;
// 異なるオリジンへのリクエストでもCookieを許可する
axios.defaults.withCredentials = true;

function App() {
  useEffect(() => {
    const getCsrfToken = async () => {
      const { data } = await axios.get("/csrf-token");
      axios.defaults.headers.common["X-CSRF-Token"] = data.csrf_token;
    };
    getCsrfToken();
  }, []);

  return (
    <Router>
      <Switch>
        <Route exact path="/landing" component={Landing} />
        <Route exact path="/user-registration" component={UserRegistration} />
        <Route exact path="/login" component={Login} />
        <Route exact path="/logout" component={Logout} />
        <Auth>
          <Route exact path="/home" component={Main} />
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
          <Route exact path="/settings" component={Settings} />
          <Route exact path="/change_password" component={ChangePassword} />
          <Route exact path="/delete_user" component={UserDelete} />
        </Auth>
      </Switch>
    </Router>
  );
}

export default App;
