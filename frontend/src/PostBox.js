import React, { useState } from "react";
import { Alert, Form } from "react-bootstrap";
import { Button } from "@material-ui/core";
import axios from "axios";

import { isCorrectSize, getBase64 } from "./libs/fileLibs";

import "./PostBox.css";

// TODO 優先度低（プレビュー、fileボタン）
function PostBox(onFileSelect) {
  const [title, setTitle] = useState("");
  const [image, setImage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");

  function sendPost(e) {
    e.preventDefault();

    if (isCorrectSize(image) === false) {
      setErrMessage(
        "The selected file is too large, maximum file size is 10MB."
      );
    } else {
      getBase64(image, result => {
        const reqBody = {
          title: title,
          image: result
        };
        axios
          .post("/posting", reqBody)
          .then(function () {
            setSuccessMessage("success post");
          })
          .catch(function (error) {
            setErrMessage(error.response.data.message);
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
            placeholder="Title here"
            value={title}
            onChange={e => setTitle(e.target.value)}
            className="mt-2"
          />
        </Form.Group>
        <Form.Group>
          <Form.File
            id="file"
            label={image.name}
            accept="image/png,image/jpeg,image/gif"
            onChange={e => setImage(e.target.files[0])}
          />
        </Form.Group>
        <Button
          onClick={sendPost}
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
