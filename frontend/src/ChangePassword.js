import React, { useState } from "react";
import {
  FormGroup,
  FormControl,
  FormLabel,
  Alert,
  Container,
  Row,
  Col,
} from "react-bootstrap";
import { useHistory } from "react-router-dom";
import Sidebar from "./Sidebar";
import { changePassword } from "./UserLibrary";
import LoaderButton from "./LoaderButton";
import { useFormFields } from "./libs/hooksLib";

import "./ChangePassword.css";

function ChangePassword() {
  const [fields, handleFieldChange] = useFormFields({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  });
  const [isChanging, setIsChanging] = useState(false);
  const history = useHistory();
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  function validateForm() {
    return (
      fields.oldPassword.length > 0 &&
      fields.newPassword.length > 0 &&
      fields.newPassword === fields.confirmPassword
    );
  }

  async function updatePassword(event) {
    event.preventDefault();

    setIsChanging(true);

    changePassword(fields.oldPassword, fields.newPassword)
      .then(() => {
        setSuccessMessage("update success");
      })
      .catch((error) => {
        setIsChanging(false);
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
          setErrMessage("failed");
        } else {
          console.log(error);
          setErrMessage("failed");
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
          <div className="ChangePassword">
            <div className="content_header">
              <h2>Change Password</h2>
            </div>

            <div className="change-password">
              <form onSubmit={updatePassword}>
                <FormGroup bssize="large" controlId="oldPassword">
                  <FormLabel>Old Password</FormLabel>
                  <FormControl
                    type="password"
                    onChange={handleFieldChange}
                    value={fields.oldPassword}
                  />
                </FormGroup>
                <hr />
                <FormGroup bssize="large" controlId="newPassword">
                  <FormLabel>New Password</FormLabel>
                  <FormControl
                    type="password"
                    onChange={handleFieldChange}
                    value={fields.newPassword}
                  />
                </FormGroup>
                <FormGroup bssize="large" controlId="confirmPassword">
                  <FormLabel>Confirm Password</FormLabel>
                  <FormControl
                    type="password"
                    onChange={handleFieldChange}
                    value={fields.confirmPassword}
                  />
                </FormGroup>
                <LoaderButton
                  block
                  type="submit"
                  bssize="large"
                  disabled={!validateForm()}
                  isLoading={isChanging}
                >
                  Change Password
                </LoaderButton>
              </form>
            </div>
          </div>
        </Col>
      </Row>
    </Container>
  );
}

export default ChangePassword;
