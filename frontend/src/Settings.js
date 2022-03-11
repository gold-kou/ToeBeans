import React, { useEffect, useState } from "react";
import { useHistory, Link } from "react-router-dom";
import { Button } from "@material-ui/core";
import { Alert, Container, Row, Col } from "react-bootstrap";

import Sidebar from "./Sidebar";
import { deleteUser } from "./UserLibrary";

import "./Settings.css";
import "./common.css";

const Settings = (props) => {
  useEffect(() => {}, []);

  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");
  const history = useHistory();

  function onClickRefreshPage() {
    setTimeout(() => {
      window.location.reload(false);
    }, 300);
  }

  async function onClickDeleteAccount() {
    deleteUser(localStorage.getItem("loginUserName"))
      .then(() => {
        setSuccessMessage("delete success");
        setTimeout(() => {
          props.history.push("/login");
        }, 1500);
      })
      .catch((error) => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          } else {
            setErrMessage(error.response.data.message);
          }
        } else if (error.request) {
          setErrMessage(error.request.data.message);
        } else {
          console.log(error);
        }
      });
  }

  return (
    <div className="main">
      <Container className="background">
        <Row>
          <Col xs={4} sm={4} md={3} lg={3}>
            <Sidebar />
          </Col>

          <Col xs={8} sm={8} md={6} lg={6}>
            {successMessage && (
              <Alert variant="success">{successMessage}</Alert>
            )}
            {errMessage && <Alert variant="danger">{errMessage}</Alert>}
            <div className="settings">
              <div className="content_header">
                <h2>Settings</h2>
              </div>

              <div className="mt-3 text-center">
                <Button
                  variant="contained"
                  color="primary"
                  component={Link}
                  to="/change_password"
                  style={{ textTransform: "none" }}
                >
                  Change password
                </Button>
              </div>

              <div className="mt-3 mb-2 text-center">
                <Button
                  variant="contained"
                  color="primary"
                  component={Link}
                  to="/logout"
                  onClick={() => onClickRefreshPage()}
                  style={{ textTransform: "none" }}
                >
                  Logout
                </Button>
              </div>

              <div className="mt-3 mb-2 text-center">
                <Button
                  variant="contained"
                  color="primary"
                  onClick={() => onClickDeleteAccount()}
                  style={{ textTransform: "none" }}
                >
                  Delete an account
                </Button>
              </div>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default Settings;
