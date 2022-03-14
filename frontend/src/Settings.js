import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@material-ui/core";
import { Container, Row, Col } from "react-bootstrap";
import Sidebar from "./Sidebar";

import "./Settings.css";
import "./common.css";

const Settings = () => {
  return (
    <div className="main">
      <Container className="background">
        <Row>
          <Col xs={4} sm={4} md={3} lg={3}>
            <Sidebar />
          </Col>

          <Col xs={8} sm={8} md={6} lg={6}>
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
                  to="/delete_user"
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
