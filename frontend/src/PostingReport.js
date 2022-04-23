import React, { useState } from "react";
import { Button } from "@material-ui/core";
import { Alert, Container, Form, Row, Col } from "react-bootstrap";
import { useHistory } from "react-router-dom";
import axios from "axios";
import Sidebar from "./Sidebar";

function PostingReport(props) {
  const postingID = props.match.params.postingID;
  const [detail, setDetail] = useState("");
  const history = useHistory();
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  const onClickSendPostingReport = async (e) => {
    e.preventDefault();

    const reqBody = {
      detail: detail,
    };
    await axios
      .post(`/reports/postings/${postingID}`, reqBody)
      .then(function () {
        setSuccessMessage("success");
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
          setErrMessage("failed");
        } else {
          console.log(error);
          setErrMessage("failed");
        }
      })
      .finally(function () {
        setDetail("");
      });
  };

  return (
    <div className="">
      <Container className="background">
        {successMessage && <Alert variant="success">{successMessage}</Alert>}
        {errMessage && <Alert variant="danger">{errMessage}</Alert>}
        <Row>
          <Col xs={4} sm={4} md={3} lg={3}>
            <Sidebar />
          </Col>

          <Col xs={8} sm={8} md={6} lg={6}>
            <div className="">
              <div className="content_header">
                <h2>Posting Report</h2>
              </div>
              <Form>
                <Form.Group>
                  <Form.Control
                    as="textarea"
                    placeholder="Any problem? Please describe the detail."
                    value={detail}
                    onChange={(e) => setDetail(e.target.value)}
                    className="mt-2"
                    style={{ height: 200 }}
                  />
                </Form.Group>
                <Button
                  onClick={(e) => onClickSendPostingReport(e)}
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="small"
                  className="mt-5 mr-2 float-right"
                >
                  Submit
                </Button>
              </Form>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default PostingReport;
