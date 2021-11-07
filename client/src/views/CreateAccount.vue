<template>
  <v-app id="ci-check-account"   >
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
                    label="Password"
                    type="password"
                    :rules="passwordRules"
                    :error-messages="serverError.password"
                    @blur="serverError.password = null"                    
                    ></v-text-field>

                     <v-text-field
                    :name="Math.random()"
                    prepend-icon="mdi-account"
                    v-model="email_check"
                    label="Password Check"
                    type="password"
                    :rules="passwordRulesCheck"               
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
      link: null,
      email : "",
      email_check: "",
      loading: false,
      message: ""  ,
      passwordRules: [
        v => !!v || 'Password is required',
        v => v.length >= 5 || 'Password must be more than 5 characters',
        v => v.length <= 255 || 'Password must be less than 20 characters',
      ],      
      passwordRulesCheck: [
        v => !!v || 'Password is required',
      ],      
      serverError: {
        "password": null,
      }     
    }    
  },
  components: {
    Header
  },
  mounted(){
    if (this.$route.params.link) {
       this.link = this.$route.params.link;
    }
  },
  methods: {
    register: function () {
      if (!this.$refs.form.validate()){
        return
      }  

      let email = this.email
      let link = this.link
      this.$store.dispatch('auth/create', { email, link })
      .then(
        () => {
          this.register_success = true
        },
        error => {
          if (error.response.data.error) {
              let errResponse = error.response.data.error
              if (errResponse.fields) {
                    this.serverError = errResponse.fields
              }
              if (error.response.data.error.message){
                this.$dialog.info({title: "Error", text: error.response.data.error.message})        
              }
            return 
          }    

          if (error.message){
            this.$dialog.info({title: "Error", text: error.message})        
            return 
          }
        }
      )  
          
    }
    /*
    register: function () {
      if (!this.$refs.form.validate()){
        return
      }
      let email = this.email
      this.$store.dispatch('auth/register', { email })
      .then(
        () => {
          this.register_success = true
        },
        error => {
          if (error.response.data.error) {
              let errResponse = error.response.data.error
              if (errResponse.fields) {
                    this.serverError = errResponse.fields
              }
              if (error.response.data.error.message){
                this.$dialog.info({title: "Error", text: error.response.data.error.message})        
              }
            return 
          }    

          if (error.message){
            this.$dialog.info({title: "Error", text: error.message})        
            return 
          }
        }
      )
    }
      */
  },  
};
</script>