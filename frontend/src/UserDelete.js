import React, { useState } from "react";
import { Alert, Button, Container, Form, Row, Col } from "react-bootstrap";
import { useHistory } from "react-router-dom";
import Sidebar from "./Sidebar";
import { deleteUser } from "./UserLibrary";

const UserDelete = (props) => {
  const history = useHistory();
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  async function onClickDeleteMyAccount() {
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
          console.log(error.request);
        } else {
          console.log(error);
        }
      });
  }

  return (
    <Container>
      <Row>
        <Col xs={4} sm={4} md={3} lg={3}>
          <Sidebar />
        </Col>
        <Col xs={8} sm={8} md={6} lg={6}>
          {successMessage && <Alert variant="success">{successMessage}</Alert>}
          {errMessage && <Alert variant="danger">{errMessage}</Alert>}

          <div className="center">
            <h2>Account delete</h2>
          </div>

          <div>
            Do you really delete your account?
            <br></br>
            If you delete your account, all your postings and so on will be
            deleted.
          </div>
          <Form>
            <div className="center mt-5">
              <Button
                variant="danger"
                type="button"
                onClick={() => onClickDeleteMyAccount()}
                className="loginButton"
              >
                Delete
              </Button>
            </div>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default UserDelete;
