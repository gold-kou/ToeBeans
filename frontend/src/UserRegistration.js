import React, { useState } from "react";
import { registerUser } from "./User";
import { Alert } from "react-bootstrap";
import "./UserRegistration.css";

const UserRegistration = props => {
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  const [state, setState] = useState({
    userName: "",
    email: "",
    password: "",
    confirmPassword: ""
  });

  const handleChange = e => {
    const { id, value } = e.target;
    setState(prevState => ({
      ...prevState,
      [id]: value
    }));
  };

  async function registerUserServer() {
    registerUser(state.userName, state.email, state.password)
      .then(() => {
        setSuccessMessage(
          "Registration successful. Redirecting to login page.."
        );
        setTimeout(() => {
          redirectToLogin();
        }, 1500);
      })
      .catch(error => {
        setErrMessage(error.data.message);
      });
  }

  const redirectToLogin = () => {
    props.history.push("/login");
  };

  const handleSubmitClick = e => {
    e.preventDefault();
    if (state.password === state.confirmPassword) {
      registerUserServer();
    } else {
      setErrMessage("Passwords do not match");
    }
  };
  return (
    <div className="registrationLoginForm card col-12 col-lg-4 login-card mt-2 hv-center">
      {successMessage && <Alert variant="success">{successMessage}</Alert>}
      {errMessage && <Alert variant="danger">{errMessage}</Alert>}

      <p>
        <b>User registration</b>
      </p>

      <form>
        <div className="form-group text-left">
          <label htmlFor="exampleInputEmail1">User name</label>
          <input
            type="text"
            className="form-control"
            id="userName"
            aria-describedby="userNameHelp"
            placeholder="Enter user name"
            value={state.userName}
            onChange={handleChange}
          />
          <small id="userNameHelp" className="form-text text-muted">
            This is displayed name.
          </small>
        </div>

        <div className="form-group text-left">
          <label htmlFor="exampleInputEmail1">Email address</label>
          <input
            type="email"
            className="form-control"
            id="email"
            aria-describedby="emailHelp"
            placeholder="Enter email"
            value={state.email}
            onChange={handleChange}
          />
          <small id="emailHelp" className="form-text text-muted">
            We'll never share your email with anyone else.
          </small>
        </div>

        <div className="form-group text-left">
          <label htmlFor="exampleInputPassword1">Password</label>
          <input
            type="password"
            className="form-control"
            id="password"
            placeholder="Password"
            value={state.password}
            onChange={handleChange}
          />
        </div>

        <div className="form-group text-left">
          <label htmlFor="exampleInputPassword1">Confirm Password</label>
          <input
            type="password"
            className="form-control"
            id="confirmPassword"
            placeholder="Confirm Password"
            value={state.confirmPassword}
            onChange={handleChange}
          />
        </div>

        <button
          type="submit"
          className="btn btn-primary"
          onClick={handleSubmitClick}
        >
          Register
        </button>
      </form>

      <br />
      <div className="mt-2">
        <span>Already have an account? </span>
        <br />
        <span className="loginText" onClick={() => redirectToLogin()}>
          Login here
        </span>
      </div>
    </div>
  );
};

export default UserRegistration;
