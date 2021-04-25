import React, { useState, forwardRef } from "react";
import { Alert, Button } from "react-bootstrap";
import { Avatar, IconButton } from "@material-ui/core";
import FavoriteIcon from "@material-ui/icons/Favorite";
import FavoriteBorderIcon from "@material-ui/icons/FavoriteBorder";
import axios from "axios";

import "./Post.css";
import "./common.css";

// TODO 優先度低（ユーザー名クリックでユーザー情報詳細画面）
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
      loginUserName
    },
    ref
  ) => {
    const [count, setCount] = useState(likedCount);
    const [isLiked, toggleLiked] = useState(liked);
    const [successMessage, setSuccessMessage] = useState("");
    const [errMessage, setErrMessage] = useState("");

    const deletePost = async () => {
      await axios
        .delete(`/posting/${postingID}`)
        .then(() => {
          setSuccessMessage("success delete post");
        })
        .catch(error => {
          setErrMessage(error.response.data.message);
        });
    };

    const changeLiked = () => {
      if (isLiked === false) {
        // POST
        const reqBody = { posting_id: postingID };
        axios
          .post("/like", reqBody)
          .then(() => {
            setCount(count + 1);
            toggleLiked(!isLiked);
          })
          .catch(error => {
            if (error.response.data.status === 401) {
              localStorage.removeItem("isLoggedIn");
            } else {
              setErrMessage(error.response.data.message);
            }
          });
      } else {
        // DELETE
        axios
          .delete(`/like/${postingID}`)
          .then(() => {
            setCount(count - 1);
            toggleLiked(!isLiked);
          })
          .catch(error => {
            if (error.response.data.status === 401) {
              localStorage.removeItem("isLoggedIn");
            } else {
              setErrMessage(error.response.data.message);
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
              {userName} {uploadedAt.split("T")[0]}
            </span>

            {userName === loginUserName && (
              <Button
                variant="outline-danger"
                size="sm"
                className="float-right"
                onClick={deletePost}
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
            <IconButton onClick={changeLiked}>
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
