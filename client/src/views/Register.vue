<template>
  <v-app id="ci-login"   >
     <Header title="Ghotos: Register" backLink="/" />
    <v-main class="ci-app-content">
      <v-container>
        <v-layout align-center justify-center>
            <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
            <v-form   ref="form" @submit.prevent="register"  v-model="valid"  v-if="!register_success">

      
                <v-card-text>
                    <v-text-field
                    :name="Math.random()"
                    prepend-icon="mdi-account"
                    v-model="email"
                    label="Login"
                    type="email"
                    :rules="emailRules"
                    :error-messages="serverError.email"
                    @blur="serverError.email = null"                    
                    ></v-text-field>
                    <v-text-field
                    :name="Math.random()"
                    autocomplete="false"
                    id="password"
                    prepend-icon="mdi-lock"
                    label="Password"
                    type="password"
                    v-model="password" 
                    :rules="passwordRules"
                    :error-messages="serverError.password"
                    @change="serverError.password = null"                    
                    ></v-text-field>
                    <v-text-field
                    :name="Math.random()"
                    autocomplete="false"
                    id="password_check"
                    prepend-icon="mdi-lock"
                    label="Password Check"
                    type="password"
                    v-model="password_check" 
                    ></v-text-field>
                
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn :disabled="!valid" color="primary" type="submit">Login</v-btn>
                </v-card-actions>
                </v-form>
                 <v-card-text v-else>
                    
                    <p>
                    Register Successful, check your e-mail account please
                    </p>
                    <v-btn :to="'/'" depressed>
                       Login
                    </v-btn>
                </v-card-text>
            </v-card>
            </v-flex>
        </v-layout>


      </v-container>

    </v-main>
  </v-app>
</template>
<script>

import Header from "@/components/layouts/HeaderAuth.vue";

export default {
  data(){
    return {
      register_success: false,
      valid: true,
      email : "",
      password : "",
      password_check : "",
      loading: false,
      message: ""  ,
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
        v => v.length <= 255 || 'E-mail must be less than 255 characters',
      ],      
      passwordRules: [
        v => !!v || 'Password is required',
        v => v.length > 5 || 'Password must be greater than 5 characters',
        v => v.length <= 20 || 'Password must be less than 20 characters',
      ], 
      serverError: {
        "email": null,
        "password": null
      }     
    }    
  },
  components: {
    Header
  },
  mounted(){
  },
  methods: {
    register: function () {
      if (!this.$refs.form.validate()){
        return
      }
      let email = this.email
      let password = this.password
      this.$store.dispatch('auth/register', { email, password })
      .then(
        (data) => {
          if (data !== true){
            this.$dialog.info({title: "Error", text: "Error..."})        
            return 
          }
          this.register_success = true

        },
        error => {
          if (error.message){
            this.$dialog.info({title: "Error", text: error.message})        
            return 
          }
          if (error.response.data.error) {
              let errResponse = error.response.data.error
              if (errResponse.fields) {
                    this.serverError = errResponse.fields
              }
   
          }
          /*
            this.loading = false;
            this.message =
            (error.response && error.response.data && error.response.data.message) ||
            error.message ||
            error.toString();
            */
        }
      )
    }
  },  
};
</script>