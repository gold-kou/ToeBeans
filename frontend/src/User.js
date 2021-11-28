import axios from "axios";

export function isLoggedIn() {
  return localStorage.getItem("isLoggedIn") !== null;
}

export async function login(email, password) {
  const reqBody = { email: email, password: password };
  return await axios.post("/login", reqBody);
}

export async function getMyUserInfo() {
  return await axios.get("/user");
}

export async function registerUser(userName, email, password) {
  const reqBody = { user_name: userName, email: email, password: password };
  return await axios.post("/user", reqBody);
}

export async function getUser(userName) {
  return await axios.get(`/user?${userName}`);
}

export async function updateUser(newPassword, avator, selfIntroduction) {
  const reqBody = {
    password: newPassword,
    icon: avator,
    self_introduction: selfIntroduction
  };
  return await axios.put("/user", reqBody);
}

export async function deleteUser() {
  return await axios.delete("/user");
}

export async function changePassword(oldPassword, newPassword) {
  const reqBody = {
    old_password: oldPassword,
    new_password: newPassword
  };
  return await axios.put("/password", reqBody);
}
