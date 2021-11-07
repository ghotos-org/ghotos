import StorageService from "./storage.service";
import ApiService from "./api.service";


class AuthService {
  login(user) {
    return ApiService.post("/auth/login", user).then((response) => {
      StorageService.setToken(response.data)     
      return Promise.resolve();
    }); 
  }

  logout() {
    return ApiService.get("/auth/logout").then((response) => {
      return Promise.resolve(response.data);
 
    });       
  }

  register(user) {
    return ApiService.post("/auth/signup", user).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }

  create(password,link) {
    return ApiService.put("/auth/signup/" + link, {password}).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }
}

export default new AuthService();