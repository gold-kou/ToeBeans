import React, { useState, useEffect } from "react";
import { useHistory } from "react-router-dom";
import axios from "axios";
import FlipMove from "react-flip-move";
import { Container, Row, Col, Alert } from "react-bootstrap";

import Sidebar from "./Sidebar";
import PostBox from "./PostBox";
import Post from "./Post";
import { getMyProfile } from "./User";

import "./Feed.css";
import "./common.css";

function Feed() {
  const [userName, setUserName] = useState("");
  const [avator, setAvator] = useState("");
  const [posts, setPosts] = useState([]);
  const [errMessage, setErrMessage] = useState("");

  const history = useHistory();

  useEffect(() => {
    getMyProfile()
      .then(response => {
        setUserName(response.data.user_name);
        setAvator(response.data.icon);
      })
      .catch(error => {
        if (error.response.data.status === 401) {
          localStorage.removeItem("isLoggedIn");
          history.push({ pathname: "login" });
        } else {
          setErrMessage(error.response.data.message);
        }
      });

    const getPosts = async () => {
      await axios
        .get("/postings?since_at=2100-01-01T00:00:00Z&limit=50")
        .then(response => {
          setPosts(response.data.postings);
        })
        .catch(error => {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            history.push({ pathname: "login" });
          } else {
            setErrMessage(error.response.data.message);
          }
        });
    };
    getPosts();
  }, []);

  return (
    <div className="main">
      <Container className="background">
        <Row>
          <Col xs={4} sm={4} md={3} lg={3}>
            <Sidebar />
          </Col>

          <Col xs={8} sm={8} md={6} lg={6}>
            <div className="feed">
              <div className="content_header">
                <h2>Home</h2>
              </div>
              {errMessage && <Alert variant="danger">{errMessage}</Alert>}

              <PostBox />

              <FlipMove>
                {posts.map(post => (
                  <Post
                    key={post.posting_id}
                    postingID={post.posting_id}
                    userName={post.user_name}
                    title={post.title}
                    imageURL={post.image_url}
                    uploadedAt={post.uploaded_at}
                    likedCount={post.liked_count}
                    liked={post.liked}
                    loginUserName={userName}
                  />
                ))}
              </FlipMove>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default Feed;
