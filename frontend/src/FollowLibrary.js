import axios from "axios";

export async function follow(userName) {
  return await axios.post(`/follows/${userName.userName}`);
}

export async function getFollowState(userName) {
  return await axios.get(`/follows/${userName.userName}`);
}

export async function unfollow(userName) {
  return await axios.delete(`/follows/${userName.userName}`);
}
