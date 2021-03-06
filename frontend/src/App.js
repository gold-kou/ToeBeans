import React, { useEffect } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import axios from "axios";

import Auth from "./Auth";
import Landing from "./Landing";
import UserRegistration from "./UserRegistration";
import Login from "./Login";
import Logout from "./Logout";
import Main from "./Main";
import MyPage from "./MyPage";
import Settings from "./Settings";
import ChangePassword from "./ChangePassword";

axios.defaults.baseURL = process.env.REACT_APP_BACK_BASE_URL;
axios.defaults.withCredentials = true;

// TODO ユーザ登録、サービストップ
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
          <Route exact path="/mypage" component={MyPage} />
          <Route exact path="/settings" component={Settings} />
          <Route exact path="/change_password" component={ChangePassword} />
        </Auth>
      </Switch>
    </Router>
  );
}

export default App;
