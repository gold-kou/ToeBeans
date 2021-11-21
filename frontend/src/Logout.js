import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { isLoggedIn } from "./User";

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
