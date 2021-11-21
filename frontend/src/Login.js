import React, { useState } from "react";
import { Form, Button, Alert } from "react-bootstrap";
import { withRouter } from "react-router-dom";
import { login } from "./User";

import "./Login.css";

const Login = props => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errMessage, setErrMessage] = useState("");

  async function loginClick() {
    login(email, password)
      .then(response => {
        localStorage.setItem("isLoggedIn", true);
        props.history.push({ pathname: "home" });
      })
      .catch(error => {
        localStorage.removeItem("isLoggedIn");
        setErrMessage(error.data.message);
      });
  }

  async function guestUserLoginClick() {
    login("guestUser@example.com", "Guest1234")
      .then(response => {
        localStorage.setItem("isLoggedIn", true);
        props.history.push({ pathname: "home" });
      })
      .catch(error => {
        localStorage.removeItem("isLoggedIn");
        setErrMessage(error.data.message);
      });
  }

  return (
    <Form className="registrationLoginForm">
      {errMessage && <Alert variant="danger">{errMessage}</Alert>}
      <p>
        <b>Login</b>
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
      <Button
        variant="primary"
        type="button"
        onClick={loginClick}
        className="loginButton"
      >
        Login
      </Button>
      <br />
      <br />
      <br />
      <br />
      <Button
        variant="info"
        type="button"
        onClick={guestUserLoginClick}
        className="loginButton"
      >
        Guest Login
      </Button>
      <div className="text-muted">You can also login as a guest user.</div>
    </Form>
  );
};

export default withRouter(Login);
