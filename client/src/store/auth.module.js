import AuthService from '../services/auth.service';
import VueJwtDecode from 'vue-jwt-decode'
import StorageService from "../services/storage.service";
function decodeTokens(token){
  if (!token || token === null) {
    return null
  }

  try {
    let tokenObj =  VueJwtDecode.decode(token);

    if (tokenObj && tokenObj.id) {     
      return {
        id: tokenObj.id
      }  
    }
  }catch(e){
    //console.log(e)    
  }
  return null
}

//const user = JSON.parse(StorageService.getItem('user'));
const user = decodeTokens(StorageService.getAccessToken());
const initialState = user
  ? { loggedIn: true, user }
  : { loggedIn: false, user: null };

export const auth = {
  namespaced: true,
  state: initialState,
  getters: {
    loggedIn(state) {
      return state.loggedIn;
    },  
  },      
  actions: {
    login({ commit }, user) {
      return AuthService.login(user).then(
        user => {
          commit('loginSuccess', user);
          return Promise.resolve(user);
        },
        error => {
          commit('loginFailure');
          return Promise.reject(error);
        }
      );
    },
    logout({ commit }) {
      return AuthService.logout().then(
        () => {
          commit('logout');
          return Promise.resolve();
        },
        error => {
          commit('logout');
          return Promise.reject(error);
        }
      );
      /*
      return AuthService.logout(() => {
        console.log("fdafaf22222")
     
      }).catch(() => {
        console.log("eeefdafaf22222")

      });
      */
    },
    registerRequest({ commit }, user) {
      return AuthService.registerRequest(user).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    register({ commit }, {password, link}) {
      return AuthService.register(password, link).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    registerCheck({ commit }, {link}) {
      return AuthService.registerCheck( link).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    newPasswordRequest({ commit }, user) {
      return AuthService.newPasswordRequest(user).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    newPassword({ commit }, {password, link}) {
      return AuthService.newPassword(password, link).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    },
    newPasswordCheck({ commit }, {link}) {
      return AuthService.newPasswordCheck( link).then(
        () => {
          return Promise.resolve();
        },
        error => {
          console.log(error)
          commit('registerFailure');
          return Promise.reject(error);
        }
      );
    }  
  },
  mutations: {
    loginSuccess(state, tokens) {
      state.loggedIn = true;
      state.user = decodeTokens(tokens);
    },
    loginFailure(state) {
      state.loggedIn = false;
      state.user = null;
    },
    logout(state) {
      state.loggedIn = false;
      state.user = null;
      
      StorageService.clearToken()
    },
    registerSuccess(state) {
      state.loggedIn = false;
    },
    registerFailure(state) {
      state.loggedIn = false;
    }
  }
};