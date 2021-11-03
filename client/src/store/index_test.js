import Vue from 'vue';
import Vuex from 'vuex';

import { auth } from './auth.module';
import { core } from './core.module';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    core,
    auth
  }
});