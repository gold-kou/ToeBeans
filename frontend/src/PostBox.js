import React, { useState } from "react";
import { Alert, Form } from "react-bootstrap";
import { useHistory } from "react-router-dom";
import { Button } from "@material-ui/core";
import axios from "axios";

import { isCorrectSize, getBase64 } from "./libs/fileLibs";

import "./PostBox.css";

function PostBox() {
  const [title, setTitle] = useState("");
  const [image, setImage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");
  const history = useHistory();

  function onClickSendPost(e) {
    e.preventDefault();

    if (isCorrectSize(image) === false) {
      setErrMessage(
        "The selected file is too large, maximum file size is 10MB."
      );
    } else {
      getBase64(image, (result) => {
        const reqBody = {
          title: title,
          image: result,
        };
        axios
          .post("/postings", reqBody)
          .then(function () {
            setSuccessMessage("success post");
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
          })
          .finally(function () {
            setTitle("");
            setImage("");
          });
      });
    }
  }

  function isValid() {
    if (title !== "" && image !== "") {
      return true;
    } else {
      return false;
    }
  }

  return (
    <div className="post-box">
      {successMessage && <Alert variant="success">{successMessage}</Alert>}
      {errMessage && <Alert variant="danger">{errMessage}</Alert>}
      <Form>
        <Form.Group controlId="title">
          <Form.Control
            type="text"
            placeholder="short message here"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="mt-2"
          />
        </Form.Group>
        <Form.Group>
          <Form.File
            id="file"
            label={image.name}
            accept="image/png,image/jpeg,image/gif"
            onChange={(e) => setImage(e.target.files[0])}
          />
        </Form.Group>
        <Button
          onClick={() => onClickSendPost()}
          type="submit"
          variant="contained"
          color="primary"
          size="small"
          className="mt-5 mr-2 float-right"
          disabled={!isValid()}
        >
          Post
        </Button>
      </Form>
    </div>
  );
}

export default PostBox;
