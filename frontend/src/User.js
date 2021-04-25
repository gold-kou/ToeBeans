import axios from "axios";

export function isLoggedIn() {
  return localStorage.getItem("isLoggedIn") !== null;
}

export async function login(email, password) {
  const reqBody = { email: email, password: password };
  try {
    return await axios.post("/login", reqBody);
  } catch (error) {
    throw error.response;
  }
}

export async function getMyProfile() {
  try {
    return await axios.get("/user");
  } catch (error) {
    throw error.response;
  }
}

export async function registerUser(userName, email, password) {
  const reqBody = { user_name: userName, email: email, password: password };
  try {
    return await axios.post("/user", reqBody);
  } catch (error) {
    throw error.response;
  }
}

export async function getUser(userName) {
  try {
    return await axios.get(`/user?{userName}`);
  } catch (error) {
    throw error.response;
  }
}

export async function updateUser(newPassword, avator, selfIntroduction) {
  const reqBody = {
    password: newPassword,
    icon: avator,
    self_introduction: selfIntroduction
  };
  try {
    return await axios.put("/user", reqBody);
  } catch (error) {
    throw error.response;
  }
}

export async function deleteUser() {
  try {
    return await axios.delete("/user");
  } catch (error) {
    throw error.response;
  }
}

export async function changePassword(oldPassword, newPassword) {
  const reqBody = {
    old_password: oldPassword,
    new_password: newPassword
  };
  try {
    return await axios.put("/password", reqBody);
  } catch (error) {
    throw error.response;
  }
}
