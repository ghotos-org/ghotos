import Vue from 'vue'
import Vuex from 'vuex'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import {router} from "./router";
import ApiService from "./services/api.service";

import store from "./store";
import moment from "moment";
import VueMoment from "vue-moment";
import './styles/main.scss';
import VueObserveVisibility from 'vue-observe-visibility'
import VuetifyDialog from 'vuetify-dialog'
import 'vuetify-dialog/dist/vuetify-dialog.css'



//moment.locale("de");
Vue.use(Vuex);
Vue.use(VueMoment, {
  moment
});
Vue.use(VueObserveVisibility)
Vue.use(VuetifyDialog, {
  context: {
    vuetify
  }
})


//this.$moment.locale('de')

ApiService.init();

Vue.config.productionTip = false

new Vue({
  router,
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
