import axios from "axios";

export function isLoggedIn() {
  return localStorage.getItem("isLoggedIn") !== null;
}

export async function login(email, password) {
  const reqBody = { email: email, password: password };
  return await axios.post("/login", reqBody);
}

export async function getMyUserInfo() {
  return await axios.get("/users");
}

export async function registerUser(userName, email, password) {
  const reqBody = { email: email, password: password };
  return await axios.post(`/users/${userName}`, reqBody);
}

export async function getUserInfo(userName) {
  return await axios.get(`/users?user_name=${userName.userName}`);
}

export async function updateUser(
  newPassword,
  avator,
  selfIntroduction,
  userName
) {
  const reqBody = {
    password: newPassword,
    icon: avator,
    self_introduction: selfIntroduction,
  };
  return await axios.put(`/users/${userName}`, reqBody);
}

export async function deleteUser(userName) {
  return await axios.delete(`/users/${userName}`);
}

export async function changePassword(oldPassword, newPassword) {
  const reqBody = {
    old_password: oldPassword,
    new_password: newPassword,
  };
  return await axios.put("/password", reqBody);
}
