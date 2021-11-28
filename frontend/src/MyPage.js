import React, { useState, useEffect } from "react";
import { Button } from "@material-ui/core";
import { Container, Row, Col, Alert, Form } from "react-bootstrap";
import axios from "axios";
import { useHistory } from "react-router-dom";
import InfiniteScroll from "react-infinite-scroller";

import Sidebar from "./Sidebar";
import { getMyUserInfo, updateUser } from "./User";
import Post from "./Post";

import "./MyPage.css";
import "./common.css";

const MyPage = () => {
  // const [avator, setAvator] = useState("");
  const [selfIntroduction, setSelfIntroduction] = useState("");
  const [postingCount, setPostingCount] = useState(0);
  const [likeCount, setLikeCount] = useState(0);
  const [likedCount, setLikedCount] = useState(0);
  const [followCount, setFollowCount] = useState(0);
  const [followedCount, setFollowedCount] = useState(0);
  const [createdAt, setCreatedAt] = useState("");
  const [posts, setPosts] = useState([]);
  const [sinceAt, setSinceAt] = useState("2100-01-01T00:00:00+09:00");
  const [hasMore, setHasMore] = useState(true);
  const [successMessage, setSuccessMessage] = useState("");
  const [errMessage, setErrMessage] = useState("");
  const history = useHistory();

  useEffect(() => {
    getMyUserInfo()
      .then(response => {
        // setAvator(response.data.icon);
        setSelfIntroduction(response.data.self_introduction);
        setPostingCount(response.data.posting_count);
        setLikeCount(response.data.like_count);
        setLikedCount(response.data.liked_count);
        setFollowCount(response.data.follow_count);
        setFollowedCount(response.data.followed_count);
        setCreatedAt(response.data.created_at);
      })
      .catch(error => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          }
          else {
            setErrMessage(error.response.data.message);
          }
        }
        else if (error.request) {
          setErrMessage(error.request.data.message);
        }
        else {
          console.log(error);
        }
      });
  }, [history]);

  const getUserPosts = async () => {
    await axios
      .get(`/postings?since_at=${sinceAt}&limit=10&user_name=${localStorage.getItem("loginUserName")}`)
      .then(response => {
        console.log("hatano2")
        console.log(response.data.postings)
        if (response.data.postings == null) {
          setHasMore(false);
          return;
        }

        setPosts([...posts, ...response.data.postings]);

        setSinceAt(
          response.data.postings[response.data.postings.length - 1].uploaded_at
        );
      })
      .catch(error => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          }
          else {
            setErrMessage(error.response.data.message);
          }
        }
        else if (error.request) {
          setErrMessage(error.request.data.message);
        }
        else {
          console.log(error);
        }
      });
  };

  const loader = (
    <div className="loader" key={0}>
      Loading ...
    </div>
  );

  async function updateSelfIntroduction() {
    updateUser("", "", selfIntroduction)
      .then(() => {
        setSuccessMessage("update success");
      })
      .catch(error => {
        if (error.response) {
          if (error.response.data.status === 401) {
            localStorage.removeItem("isLoggedIn");
            localStorage.removeItem("loginUserName");
            history.push({ pathname: "login" });
          }
          else {
            setErrMessage(error.response.data.message);
          }
        }
        else if (error.request) {
          setErrMessage(error.request.data.message);
        }
        else {
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
            <div className="mypage">
              <div className="content_header">
                <h2>My Page</h2>
              </div>

              <Container className="mt-5 mb-5 ml-2">
                {/* <img src={avator} alt="" /> */}
                {localStorage.getItem("loginUserName")}
                <div className="mini-character">
                  since {createdAt.split("T")[0]}
                </div>

                <Row className="mt-5 mypage-data">
                  <Col sm-3 md-5>
                    like count
                  </Col>
                  <Col sm-3 md-5>
                    liked count
                  </Col>
                </Row>
                <Row className="mt-1 mypage-data">
                  <Col sm-3 md-5>
                    {likeCount}
                  </Col>
                  <Col sm-3 md-5>
                    {likedCount}
                  </Col>
                </Row>
                <Row className="mt-3 mypage-data">
                  <Col sm-3 md-5>
                    follow count
                  </Col>
                  <Col sm-3 md-5>
                    followed count
                  </Col>
                </Row>
                <Row className="mt-1 mypage-data">
                  <Col sm-3 md-5>
                    {followCount}
                  </Col>
                  <Col sm-3 md-5>
                    {followedCount}
                  </Col>
                </Row>
                <Row className="mt-3 mypage-data">
                  <Col sm-3 md-5>
                    post count
                  </Col>
                  <Col sm-3 md-5></Col>
                </Row>
                <Row className="mt-1 mypage-data">
                  <Col sm-3 md-5>
                    {postingCount}
                  </Col>
                  <Col sm-3 md-5></Col>
                </Row>

                <Form className="mt-5">
                  <Form.Group>
                    <Form.Label>Self introduction</Form.Label>
                    <Form.Control
                      as="textarea"
                      rows={3}
                      placeholder={selfIntroduction}
                      value={selfIntroduction}
                      onChange={e => setSelfIntroduction(e.target.value)}
                    />
                  </Form.Group>
                  <Button
                    onClick={updateSelfIntroduction}
                    variant="contained"
                    color="primary"
                    size="sm"
                    className="float-right"
                  >
                    Update
                  </Button>
                </Form>
              </Container>
              <br></br>
              <br></br>

              <InfiniteScroll
                loadMore={getUserPosts} //項目を読み込む際に処理するコールバック関数
                hasMore={hasMore} //読み込みを行うかどうかの判定
                loader={loader}
              >
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
                    loginUserName={localStorage.getItem("loginUserName")}
                  />
                ))}
              </InfiniteScroll>
            </div>
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default MyPage;
