import axios from 'axios';
import StorageService from "./storage.service";
import ApiService from "./api.service";


const API_URL = process.env.VUE_APP_API_URL + '/auth/';

class AuthService {
  login(user) {
    return axios
      .post(API_URL + 'login', {
        email: user.email,
        password: user.password
      })
      .then(response => {
        StorageService.setToken(response.data)     
        return response.data;
      });
  }

  logout() {
    return ApiService.get("/auth/logout").then((response) => {
      return Promise.resolve(response.data);
 
    });       
  }

  register(user) {
    return axios.post(API_URL + 'signup', {
      email: user.email,
      password: user.password
    });
  }
}

export default new AuthService();