import React, { useState, useEffect } from "react";
import { Container, Row, Alert } from "react-bootstrap";
import { Link } from "react-router-dom";
import { isLoggedIn, logout } from "./User";

const Logout = () => {
  useEffect(() => {
    if (isLoggedIn()) {
      localStorage.removeItem("isLoggedIn");
    }
  }, []);

  return (
    <div className="text-center">
      <h2>Logout Page</h2>
      <div>
        <Link to="/login">Login</Link>
      </div>
    </div>
  );
};

export default Logout;
