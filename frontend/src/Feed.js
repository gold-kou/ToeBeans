import React, { useState, useEffect } from "react";
import { useHistory } from "react-router-dom";
import axios from "axios";
import InfiniteScroll from "react-infinite-scroller";
import { Container, Row, Col, Alert } from "react-bootstrap";

import Sidebar from "./Sidebar";
import PostBox from "./PostBox";
import Post from "./Post";
import { getMyUserInfo } from "./UserLibrary";

import "./Feed.css";
import "./common.css";

function Feed() {
  const [posts, setPosts] = useState([]);
  const [sinceAt, setSinceAt] = useState("2100-01-01T00:00:00+09:00");
  const [hasMore, setHasMore] = useState(true); //再読み込み判定
  const [errMessage, setErrMessage] = useState("");
  const history = useHistory();

  useEffect(() => {
    getMyUserInfo()
      .then((response) => {
        localStorage.setItem("loginUserName", response.data.user_name);
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
  }, [history]);

  const getPosts = async () => {
    await axios
      .get(`/postings?since_at=${sinceAt}&limit=10`)
      .then((response) => {
        if (response.data.postings.length < 10) {
          setHasMore(false);
        }
        if (response.data.postings.length !== 0) {
          setPosts([...posts, ...response.data.postings]);
          // 取得データのうち一番古い uploaded_at を次のリクエスト用に保持しておく
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
        } else {
          console.log(error);
        }
      });
  };

  return (
    <div className="main">
      <Container className="background">
        <Row>
          <Col xs={4} sm={4} md={3} lg={3} xl={3}>
            <Sidebar />
          </Col>

          <Col xs={8} sm={8} md={6} lg={6} xl={6}>
            <div className="feed">
              <div className="content_header">
                <h2>Home</h2>
              </div>
              {errMessage && <Alert variant="danger">{errMessage}</Alert>}

              <PostBox />

              <InfiniteScroll
                loadMore={getPosts} //項目を読み込む際に処理するコールバック関数
                hasMore={hasMore} //読み込みを行うかどうかの判定
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
}

export default Feed;
