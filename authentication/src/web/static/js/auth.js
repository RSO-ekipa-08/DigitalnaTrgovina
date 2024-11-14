const auth = {
  setTokens(data) {
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("id_token", data.id_token);
    localStorage.setItem("user_profile", JSON.stringify(data.profile));
  },

  getToken() {
    return localStorage.getItem("access_token");
  },

  getProfile() {
    const profile = localStorage.getItem("user_profile");
    return profile ? JSON.parse(profile) : null;
  },

  logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("id_token");
    localStorage.removeItem("user_profile");
  },
};
