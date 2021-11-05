<template>
  <v-app id="ci-login"   >
    <v-main class="ci-app-content">
      <v-container>
     
        <v-layout align-center justify-center>
            <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
            <v-form @submit.prevent="login">

      
                <v-card-text>
                    <v-text-field
                    prepend-icon="mdi-account"
                    v-model="email"
                    label="Login"
                    type="email"
                    ></v-text-field>
                    <v-text-field
                    id="password"
                    prepend-icon="mdi-lock"
                    label="Password"
                    type="password"
                    v-model="password" 
                    ></v-text-field>
                    <v-text-field
                    id="password_check"
                    prepend-icon="mdi-lock"
                    label="Password Check"
                    type="password_check"
                    v-model="password" 
                    ></v-text-field>
                
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="primary" type="submit">Login</v-btn>
                </v-card-actions>
                </v-form>
            </v-card>
            </v-flex>
        </v-layout>

      </v-container>

    </v-main>
  </v-app>
</template>
<script>

//import { mapGetters } from "vuex";
export default {
  data(){
    return {
      email : "",
      password : "",
      password_check : "",
      loading: false,
      message: ""      
    }    
  },

  mounted(){
   // console.log(this.$store.state.auth.user)
  },
  methods: {
    login: function () {
      let email = this.email
      let password = this.password
      this.$store.dispatch('auth/register', { email, password })
      .then(
        () => this.$router.push('/'),
        error => {
            this.loading = false;
            this.message =
            (error.response && error.response.data && error.response.data.message) ||
            error.message ||
            error.toString();
        }
      )
    }
  },  
};
</script>