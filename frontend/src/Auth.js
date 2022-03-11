import React from "react";
import { Redirect } from "react-router-dom";
import { isLoggedIn } from "./UserLibrary";

const Auth = (props) =>
  isLoggedIn() ? props.children : <Redirect to={"/login"} />;

export default Auth;
