<template>
  <v-app id="ci-login">
    <Header title="Ghotos: Login" />
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
                    label="E-Mail"
                    type="email"
                    :rules="emailRules"
                    :error-messages="serverError.email"
                    @blur="serverError.email = null"                         
                  ></v-text-field>
                  <v-text-field
                    id="password"
                    prepend-icon="mdi-lock"
                    label="Password"
                    type="password"
                    v-model="password"
                    :rules="passwordRules"
                    :error-messages="serverError.password"
                    @blur="serverError.password = null"                    
                  ></v-text-field>
                </v-card-text>
                <v-card-actions>
                  <router-link :to="'/register'"
                    ><span small>Create Account</span></router-link
                  >
                  | <router-link :to="'/password'">New Password</router-link>
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
import Header from "@/components/layouts/HeaderAuth.vue";
export default {
  data() {
    return {
      email: "",
      password: "",
      loading: false,
      message: "",
      passwordRules: [
        (v) => !!v || "Password is required",
      ],
      emailRules: [
        v => !!v || 'E-mail is required',
      ],    
      serverError: {
        password: null,
        email: null,
      },                
    };
  },
  components: {
    Header,
  },
  mounted() {
    // console.log(this.$store.state.auth.user)
  },
  methods: {
    login: function () {
      let email = this.email;
      let password = this.password;
      this.$store.dispatch("auth/login", { email, password }).then(
        () => this.$router.push("/"),
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
  },
};
</script>