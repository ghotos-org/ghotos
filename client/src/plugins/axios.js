import axios from 'axios';
import StorageService from "../services/storage.service";
import {router} from "../router/index.js";
// doing something with the request
axios.interceptors.request.use(
  (config) => {

      const token = StorageService.getAccessToken();
    //  console.log(config.headers)
      if (token) {
       
          config.headers['Authorization'] = 'Bearer ' + token;
          //config.headers['x-access-token'] = token;
      } 

      return config;
  },
  error => {
      Promise.reject(error)
  }
);

// doing something with the response
axios.interceptors.response.use(
  (response) => {

     return response;
  },
  (error) => {
      const originalRequest = error.config;



     // all 4xx/5xx responses will end here
      if (error.response.status === 401 && originalRequest.url.includes("/auth/login")) {
          StorageService.clearToken()
          router.push('/login');
          return Promise.reject(error);
      }


      if (error.response.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;
          const refreshToken = StorageService.getRefreshToken();

          return axios.post(process.env.VUE_APP_API_URL + "/auth/refresh",
          {
              "refresh_token": refreshToken
          })
          .then(res => {
              console.log("res",res)

              if (res.status === 200) {
                  StorageService.setToken(res.data);
                  axios.defaults.headers.common['Authorization'] = 'Bearer ' + StorageService.getAccessToken();
                  return axios(originalRequest);
              }
          }).catch((e) => {
              console.log("e",e)
          })
  
      }
     
         
     return Promise.reject(error);
  }
);


export default axios;