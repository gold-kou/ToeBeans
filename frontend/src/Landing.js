import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@material-ui/core";

import "./Landing.css";
import image from "./images/landing_top.png";

const Landing = () => {
  return (
    <div className="landing">
      <div className="image-center">
        <img src={image} alt="top image" className="landing-image"></img>
      </div>

      <br />
      <br />
      <br />
      <Button
        variant="outlined"
        component={Link}
        to="/user-registration"
        size="large"
        style={{ textTransform: "none" }}
      >
        Getting started
      </Button>
    </div>
  );
};

export default Landing;
