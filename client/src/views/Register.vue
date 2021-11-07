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
      loading: false,
      message: ""  ,
      emailRules: [
        v => !!v || 'E-mail is required',
        v => /.+@.+/.test(v) || 'E-mail must be valid',
        v => v.length <= 255 || 'E-mail must be less than 255 characters',
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
      this.$store.dispatch('auth/register', { email })
      .then(
        () => {
          /*
          if (data !== true){
            this.$dialog.info({title: "Error", text: "Error..."})        
            return 
          }
          */
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