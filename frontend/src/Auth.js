import React from "react";
import { Redirect } from "react-router-dom";
import { isLoggedIn } from "./User";

const Auth = props =>
  isLoggedIn() ? props.children : <Redirect to={"/login"} />;

export default Auth;
