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
      return Promise.resolve(response);
    }); 
  }
}

export default new AuthService();