<template>
  <v-app id="ci-check-account">
    <Header title="Ghotos: Register" backLink="/" />
    <v-main class="ci-app-content">
      <v-container>
        <v-layout align-center justify-center>
          <v-flex xs12 sm8 md4>
            <v-card class="elevation-12">
              <div v-if="!check_link">
                <v-card-text>
                  Link from E-Mail is not valid or expired!
                 <br> <v-btn :to="'/register'" depressed> Try again </v-btn>
                </v-card-text>
              </div>
              <div v-else>
                <v-form
                  ref="form"
                  @submit.prevent="register"
                  v-model="valid"
                  v-if="!register_success"
                >
                  <v-card-text>
                    <v-text-field
                      :name="Math.random()"
                      prepend-icon="mdi-account"
                      v-model="password"
                      label="Password"
                      type="password"
                      :rules="passwordRules"
                      :error-messages="serverError.password"
                      @blur="serverError.password = null"
                    ></v-text-field>

                    <v-text-field
                      :name="Math.random()"
                      prepend-icon="mdi-account"
                      v-model="password_check"
                      label="Password Check"
                      type="password"
                      :rules="passwordRulesCheck"
                    ></v-text-field>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn :disabled="!valid" color="primary" type="submit"
                      >Login</v-btn
                    >
                  </v-card-actions>
                </v-form>
                <v-card-text v-else>
                  <p>Register Successful, go to Login Page</p>
                  <v-btn :to="'/'" depressed> Login </v-btn>
                </v-card-text>
              </div>
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
  data() {
    return {
      register_success: false,
      valid: true,
      check_link: false,
      link: null,
      password: "",
      password_check: "",
      loading: false,
      message: "",
      passwordRules: [
        (v) => !!v || "Password is required",
        (v) => v.length >= 5 || "Password must be more than 5 characters",
        (v) => v.length <= 255 || "Password must be less than 20 characters",
      ],
      passwordRulesCheck: [(v) => !!v || "Password is required"],
      serverError: {
        password: null,
      },
    };
  },
  components: {
    Header,
  },
  mounted() {
    if (this.$route.params.link) {
      this.link = this.$route.params.link;
    }
    this.$store.dispatch("auth/checkLink", { link: this.link }).then(
      () => {
        this.check_link = true;
      },
      () => {}
    );
  },
  methods: {
    register: function () {
      if (!this.$refs.form.validate()) {
        return;
      }

      let password = this.password;
      let link = this.link;
      this.$store.dispatch("auth/create", { password, link }).then(
        () => {
          this.register_success = true;
        },
        (error) => {
          if (error.response.data.error) {
            let errResponse = error.response.data.error;
            if (errResponse.fields) {
              this.serverError = errResponse.fields;
            }
            if (error.response.data.error.message) {
              this.$dialog.info({
                title: "Error",
                text: error.response.data.error.message,
              });
            }
            return;
          }

          if (error.message) {
            this.$dialog.info({ title: "Error", text: error.message });
            return;
          }
        }
      );
    },
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