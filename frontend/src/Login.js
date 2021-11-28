import React, { useState } from "react";
import { Form, Button, Alert } from "react-bootstrap";
import { useHistory, withRouter } from "react-router-dom";
import { login } from "./User";

import "./Login.css";

const Login = props => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errMessage, setErrMessage] = useState("");
  const history = useHistory();

  async function loginClick() {
    login(email, password)
      .then(() => {
        localStorage.setItem("isLoggedIn", true);
        props.history.push({ pathname: "home" });
      })
      .catch(error => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          }
          else {
            setErrMessage(error.response.data.message);
          }
        }
        else if (error.request) {
          setErrMessage(error.request.data.message);
        }
        else {
          console.log(error);
        }
      });
  }

  async function guestUserLoginClick() {
    login("guestUser@example.com", "Guest1234")
      .then(() => {
        localStorage.setItem("isLoggedIn", true);
        props.history.push({ pathname: "home" });
      })
      .catch(error => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          }
          else {
            setErrMessage(error.response.data.message);
          }
        }
        else if (error.request) {
          setErrMessage(error.request.data.message);
        }
        else {
          console.log(error);
        }
      });
  }

  return (
    <Form className="registrationLoginForm">
      {errMessage && <Alert variant="danger">{errMessage}</Alert>}
      <p>
        <div className="center">
          <h2>Login</h2>
        </div>
      </p>

      <Form.Group controlId="email">
        <Form.Label>Email</Form.Label>
        <Form.Control
          type="email"
          placeholder="Email"
          onChange={e => {
            setEmail(e.target.value);
          }}
          value={email}
        />
        <Form.Text className="text-muted">
          We'll never share your email with anyone else.
        </Form.Text>
      </Form.Group>
      <Form.Group controlId="password">
        <Form.Label>Password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Password"
          onChange={e => {
            setPassword(e.target.value);
          }}
          value={password}
        />
      </Form.Group>

      <div className="center">
        <Button
          variant="primary"
          type="button"
          onClick={loginClick}
          className="loginButton"
        >
          Login
        </Button>
      </div>

      <br />
      <br />
      <br />

      <div className="center">
        <Button
          variant="info"
          type="button"
          onClick={guestUserLoginClick}
          className="loginButton"
        >
          Guest Login
        </Button>
      </div>
      <div className="text-muted center">You can also login as a guest user.</div>

      <br />
      <br />
      <br />

      <div className="center">
        Sign up for <a href="/user-registration">Toe Beans</a>
      </div>
    </Form>
  );
};

export default withRouter(Login);
