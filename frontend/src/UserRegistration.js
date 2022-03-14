import React, { useState } from "react";
import { registerUser } from "./UserLibrary";
import { Alert } from "react-bootstrap";
import "./UserRegistration.css";

const UserRegistration = (props) => {
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  const [state, setState] = useState({
    userName: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setState((prevState) => ({
      ...prevState,
      [id]: value,
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
        }, 2000);
      })
      .catch((error) => {
        if (error.response) {
          setErrMessage(error.response.data.message);
        } else if (error.request) {
          console.log(error.request);
        } else {
          console.log(error);
        }
      });
  }

  const redirectToLogin = () => {
    props.history.push("/login");
  };

  const onClickHandleSubmitClick = (e) => {
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

      <div className="center">
        <h2>User Registration</h2>
      </div>

      <form>
        <div className="form-group text-left">
          <label htmlFor="exampleInputEmail1">User name</label>
          <input
            type="text"
            className="form-control"
            id="userName"
            placeholder="User name"
            value={state.userName}
            onChange={handleChange}
          />
        </div>

        <div className="form-group text-left">
          <label htmlFor="exampleInputEmail1">Email</label>
          <input
            type="email"
            className="form-control"
            id="email"
            aria-describedby="emailHelp"
            placeholder="Email"
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
          onClick={() => onClickHandleSubmitClick()}
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
