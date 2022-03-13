import React, { useState, forwardRef } from "react";
import { Link, useHistory } from "react-router-dom";
import { Alert, Button } from "react-bootstrap";
import { IconButton } from "@material-ui/core";
import FavoriteIcon from "@material-ui/icons/Favorite";
import FavoriteBorderIcon from "@material-ui/icons/FavoriteBorder";
import axios from "axios";

import "./Post.css";
import "./common.css";

const Post = forwardRef(
  (
    {
      postingID,
      userName,
      title,
      imageURL,
      uploadedAt,
      likedCount,
      liked,
      loginUserName,
    },
    ref
  ) => {
    const [count, setCount] = useState(likedCount);
    const [isLiked, toggleLiked] = useState(liked);
    const [successMessage, setSuccessMessage] = useState("");
    const [errMessage, setErrMessage] = useState("");
    const history = useHistory();

    const onClickDeletePost = async () => {
      await axios
        .delete(`/postings/${postingID}`)
        .then(() => {
          setSuccessMessage("success delete post");
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
    };

    const onClickChangeLiked = () => {
      if (isLiked === false) {
        // POST
        axios
          .post(`/likes/${postingID}`)
          .then(() => {
            setCount(count + 1);
            toggleLiked(!isLiked);
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
      } else {
        // DELETE
        axios
          .delete(`/likes/${postingID}`)
          .then(() => {
            setCount(count - 1);
            toggleLiked(!isLiked);
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
    };

    return (
      <div className="post" ref={ref}>
        {successMessage && <Alert variant="success">{successMessage}</Alert>}
        {errMessage && <Alert variant="danger">{errMessage}</Alert>}
        <div className="post__body">
          <div className="post__headerText">
            <span className="mini-character">
              <Link to={"/userpage/" + userName}>
                {userName} {uploadedAt.split("T")[0]}
              </Link>
            </span>

            {userName === loginUserName && (
              <Button
                variant="outline-danger"
                size="sm"
                className="float-right"
                onClick={() => onClickDeletePost()}
              >
                Delete
              </Button>
            )}
          </div>

          <div className="post__headerTitle">
            <p>{title}</p>
          </div>

          <img className="post__image" src={imageURL} alt="" />

          <div className="post__footer">
            <IconButton onClick={() => onClickChangeLiked()}>
              {isLiked ? (
                <FavoriteIcon fontSize="small" color={"secondary"} />
              ) : (
                <FavoriteBorderIcon fontSize="small" />
              )}
              {count}
            </IconButton>
          </div>
        </div>
      </div>
    );
  }
);

export default Post;
