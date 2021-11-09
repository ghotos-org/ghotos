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

  registerRequest(user) {
    return ApiService.post("/auth/signup", user).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }
  registerCheck(link) {
    return ApiService.get("/auth/signup/" + link).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }  

  register(password,link) {
    return ApiService.post("/auth/signup/" + link, {password}).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }

  newPasswordRequest(user) {
    return ApiService.post("/auth/password", user).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }
  newPasswordCheck(link) {
    return ApiService.get("/auth/password/" + link).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }  

  newPassword(password,link) {
    return ApiService.post("/auth/password/" + link, {password}).then((response) => {
      return Promise.resolve(response.data);
    }); 
  }  
}

export default new AuthService();