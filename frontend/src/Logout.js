import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { isLoggedIn } from "./UserLibrary";

const Logout = () => {
  useEffect(() => {
    if (isLoggedIn()) {
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("loginUserName");
    }
  }, []);

  function onClickRefreshPage() {
    setTimeout(() => {
      window.location.reload(false);
    }, 300);
  }

  return (
    <div className="text-center">
      <h2>Logout Page</h2>
      <div>
        <Link to="/login" onClick={onClickRefreshPage}>
          Login
        </Link>
      </div>
    </div>
  );
};

export default Logout;
