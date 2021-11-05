import Vue from "vue";
import Router from "vue-router";
import Home from '@/views/Home.vue';
import Login from '@/views/Login.vue';
import store from "../store";

Vue.use(Router);

export const router = new Router({
    mode: 'history',
    routes: [        
        {
            path: "/",
            name: "gallery",
            meta: {
                requiresAuth: true
            },
            component: Home
        },
        {
            path: '/login',
            component: Login
        },
        {
            path: '/register',
            component: () => import("@/views/Register"),
        },
        {
            path: "/photo/:file",
            name: "photo",
            meta: {
                requiresAuth: true
            },
            component: () => import("@/views/Photo"),
        },        
    ]
});

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
      if(store.state.auth.loggedIn){
        next()
        return
      }

      next('/login')
      return
    }
    next()
  });
  
// router.beforeEach((to, from, next) => {
//   const publicPages = ['/login', '/register', '/home'];
//   const authRequired = !publicPages.includes(to.path);
//   const loggedIn = localStorage.getItem('user');

//   // trying to access a restricted page + not logged in
//   // redirect to login page
//   if (authRequired && !loggedIn) {
//     next('/login');
//   } else {
//     next();
//   }
// });