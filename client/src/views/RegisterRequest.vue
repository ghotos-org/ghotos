<template>
  <v-app id="ci-login"   >
     <Header v-if="page_register" title="Ghotos: Register" backLink="/" />
     <Header v-else-if="page_newpassword" title="Ghotos: New Password" backLink="/" />
     <Header v-else title="Ghotos" backLink="/" />
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
                    label="Email"
                    type="email"
                    :rules="emailRules"
                    :error-messages="serverError.email"
                    @blur="serverError.email = null"                    
                    ></v-text-field>
                
                    <v-text-field
                    :name="Math.random()"
                    prepend-icon="mdi-account"
                    v-model="email_check"
                    label="Email Repeat"
                    type="email"
                    :rules="emailCheckRules"             
                    ></v-text-field>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn v-if="page_register" :disabled="!valid" color="primary" type="submit">Request Login</v-btn>
                    <v-btn v-if="page_newpassword" :disabled="!valid" color="primary" type="submit">Request New Password</v-btn>
                </v-card-actions>
                </v-form>
                 <v-card-text v-else>
                    <div v-if="page_register">
                        <p>Register Successful, check your e-mail account please</p>
                        <v-btn :to="'/'" depressed>Login</v-btn>
                    </div>
                    <div v-if="page_newpassword">
                      <p>Send New Password Request, check your e-mail account please</p>
                      <v-btn :to="'/'" depressed>Login</v-btn>
                    </div>
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
  props: ['page_register', 'page_newpassword'],
  data(){
    return {
      register_success: false,
      valid: true,
      email : "",
      email_check : "",
      loading: false,
      message: ""  ,
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
        v => v.length <= 255 || 'E-mail must be less than 255 characters',
      ],      
      emailCheckRules: [
        v => !!v || 'E-mail is required',
      ],      
      serverError: {
        "email": null,
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
      let api = ""
      if (this.page_register){
        api = 'auth/registerRequest'
      }
      if (this.page_newpassword){
        api = 'auth/newPasswordRequest'
      }      
      this.$store.dispatch(api, { email })
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
  },  
};
</script>