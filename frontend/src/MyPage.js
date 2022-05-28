import React, { useState, useEffect } from "react";
import { Button } from "@material-ui/core";
import { Container, Row, Col, Alert, Form } from "react-bootstrap";
import axios from "axios";
import { useHistory } from "react-router-dom";
import InfiniteScroll from "react-infinite-scroller";

import Sidebar from "./Sidebar";
import { getUserInfo, updateUser } from "./UserLibrary";
import Post from "./Post";

import "./UserPage.css";
import "./common.css";

/*
異なるユーザのユーザページからマイページへ遷移する場合に、useStateが実行されないため、ユーザページとマイページで別々にコンポーネントを用意する。
*/

const MyPage = () => {
  const userName = localStorage.getItem("loginUserName");
  const loginUserName = localStorage.getItem("loginUserName");
  const [selfIntroduction, setSelfIntroduction] = useState("");
  const history = useHistory();
  const [postingCount, setPostingCount] = useState(0);
  const [likeCount, setLikeCount] = useState(0);
  const [likedCount, setLikedCount] = useState(0);
  const [followCount, setFollowCount] = useState(0);
  const [followedCount, setFollowedCount] = useState(0);
  const [createdAt, setCreatedAt] = useState("");
  const [posts, setPosts] = useState([]);
  const [sinceAt, setSinceAt] = useState("2100-01-01T00:00:00+09:00");
  const [hasMore, setHasMore] = useState(true);
  const [errMessage, setErrMessage] = useState("");

  useEffect(() => {
    getUserInfo({ userName })
      .then((response) => {
        setSelfIntroduction(response.data.self_introduction);
        setPostingCount(response.data.posting_count);
        setLikeCount(response.data.like_count);
        setLikedCount(response.data.liked_count);
        setFollowCount(response.data.follow_count);
        setFollowedCount(response.data.followed_count);
        setCreatedAt(response.data.created_at);
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
      });
  }, [userName, history]);

  const getUserPosts = async () => {
    await axios
      .get(`/postings?since_at=${sinceAt}&limit=10&user_name=${userName}`)
      .then((response) => {
        if (response.data.postings.length < 10) {
          setHasMore(false);
        }
        if (response.data.postings.length !== 0) {
          setPosts([...posts, ...response.data.postings]);
          setSinceAt(
            response.data.postings[response.data.postings.length - 1]
              .uploaded_at
          );
        }
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
      });
  };

  async function onClickUpdateSelfIntroduction() {
    updateUser("", "", selfIntroduction, userName)
      .then(() => {
        setSelfIntroduction(selfIntroduction);
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
            {errMessage && <Alert variant="danger">{errMessage}</Alert>}
            <div className="userpage">
              <div className="content_header">
                <h2>User Page</h2>
              </div>

              <Container className="mt-3 mb-5 ml-2">
                <Row>
                  <Col className="ml-1">
                    {userName}
                    <div className="mini-character">
                      since {createdAt.split("T")[0]}
                    </div>
                  </Col>
                </Row>

                <Row className="mt-5 userpage-data">
                  <Col sm-3="true" md-5="true">
                    like count
                  </Col>
                  <Col sm-3="true" md-5="true">
                    liked count
                  </Col>
                </Row>
                <Row className="mt-1 userpage-data">
                  <Col sm-3="true" md-5="true">
                    {likeCount}
                  </Col>
                  <Col sm-3="true" md-5="true">
                    {likedCount}
                  </Col>
                </Row>
                <Row className="mt-3 userpage-data">
                  <Col sm-3="true" md-5="true">
                    follow count
                  </Col>
                  <Col sm-3="true" md-5="true">
                    followed count
                  </Col>
                </Row>
                <Row className="mt-1 userpage-data">
                  <Col sm-3="true" md-5="true">
                    {followCount}
                  </Col>
                  <Col sm-3="true" md-5="true">
                    {followedCount}
                  </Col>
                </Row>
                <Row className="mt-3 userpage-data">
                  <Col sm-3="true" md-5="true">
                    post count
                  </Col>
                  <Col sm-3="true" md-5="true"></Col>
                </Row>
                <Row className="mt-1 userpage-data">
                  <Col sm-3="true" md-5="true">
                    {postingCount}
                  </Col>
                  <Col sm-3="true" md-5="true"></Col>
                </Row>

                <Form className="mt-5">
                  <Form.Group>
                    <Form.Label>Self introduction</Form.Label>
                    <Form.Control
                      as="textarea"
                      rows={3}
                      placeholder={selfIntroduction}
                      value={selfIntroduction}
                      onChange={(e) => setSelfIntroduction(e.target.value)}
                    />
                  </Form.Group>
                  <Button
                    onClick={() => onClickUpdateSelfIntroduction()}
                    variant="contained"
                    color="primary"
                    size="small"
                    className="float-right"
                  >
                    Update
                  </Button>
                </Form>
              </Container>
              <br></br>
              <br></br>

              <InfiniteScroll
                loadMore={getUserPosts} // 項目を読み込む際に処理するコールバック関数
                hasMore={hasMore} // 読み込みを行うかどうかの判定
              >
                {posts.map((post) => (
                  <Post
                    key={post.posting_id}
                    postingID={post.posting_id}
                    userName={post.user_name}
                    title={post.title}
                    imageURL={post.image_url}
                    uploadedAt={post.uploaded_at}
                    likedCount={post.liked_count}
                    liked={post.liked}
                    loginUserName={loginUserName}
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
