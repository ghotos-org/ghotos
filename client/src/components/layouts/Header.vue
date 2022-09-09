
<template>
  <div>
    <v-app-bar app fixed>
      <v-btn  icon v-if="backLink" :disabled="isRequesting"  :to="backLink" >
        <v-icon> mdi-arrow-left</v-icon>
      </v-btn>
      <v-app-bar-nav-icon v-else @click.stop="drawer = !drawer" ></v-app-bar-nav-icon>
        <v-toolbar-title>
            {{ title }} 
        </v-toolbar-title>
        <v-spacer></v-spacer>

        <v-btn
          color="primary"
          class="text-none"
          depressed
          :loading="isSelecting"
          @click="onButtonClick"
        >
    

       upload
        </v-btn>

     <!-- <v-progress-linear v-if="isLoading" fixed top indeterminate ></v-progress-linear> -->
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" absolute temporary>
        <div class="ci-dawer-logo"></div>
        <v-divider></v-divider>
        <v-list flat>
            <v-list-item v-if="loggedIn && 0 == 1" to="/kennwort">
              <v-list-item-icon>
                <v-icon>mdi-edit</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Kennwort ändern</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item  v-if="loggedIn" @click="logout()">
              <v-list-item-icon>
                <v-icon>mdi-logout</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Abmelden</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item  v-if="!loggedIn && 0 == 1"  to="/kennwort-vergessen">
            <v-list-item-icon>
              <v-icon>mdi-pencil</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Kennwort vergessen</v-list-item-title>
            </v-list-item-content>
          </v-list-item>            
        </v-list>
      </v-navigation-drawer>

      <v-dialog v-model="dialog" :persistent="uploading" max-width="600">
        <v-card>
          <v-card-title class="text-h5"> Upload </v-card-title>
          <v-card-text>
            <div v-if="uploading">
              <div v-if="currentFile">
                <div class="ci-image-preview" v-if="imageData.length > 0">
                  <div class="wrapper">
                    <img :src="imageData" />
                  </div>
                  <div>
                    <div class="ci-inner-title">{{ currentFile.name }}</div>
                    <div class="ci-inner-data">
                      <div>
                        {{
                          Math.round(
                            (currentFile.size / (1000 * 1000) +
                              Number.EPSILON) *
                              100
                          ) / 100
                        }}
                        MB
                      </div>
                      <div>{{ img_count }} of {{ img_max }}</div>
                    </div>
                  </div>
                  <v-progress-linear
                    v-model="progress"
                    color="light-blue"
                    height="25"
                  >
                    <strong>{{ Math.ceil(progress) }}%</strong>
                  </v-progress-linear>
                </div>
              </div>
            </div>
            <div v-else>
              <v-alert border="top" colored-border type="info" elevation="2">
                Upload finished!
              </v-alert>
            </div>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="green darken-1"
              text
              @click="close"
            >
              Close
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <input
        ref="uploader"
        class="d-none"
        type="file"
        accept="image/*"
        multiple
        @change="onChangeSelectFiles"
      >
  </div>
</template>
<script>
import { mapGetters } from "vuex";


import FilesService from "../../services/file.service";

export default {
  name: "Header",
  props: ["title", "backLink"],
  data() {
    return {
      isAuthenticated: false,
      uploading: false,
      dialog: false,
      img_max: 0,
      img_count: 0,
      url: process.env.VUE_APP_API_URL,
      currentFile: undefined,
      progress: 0,
      message: "",
      imageData: "",
      drawer: null,
      isSelecting: false,
    }},   
  computed: {
    /*
    ...mapGetters("auth",["loggedIn"]),
    ...mapGetters("fileCount"),
    */
    //"default",["galleryDays"]
    ...mapGetters('auth', ['loggedIn']),

  },      
  mounted(){
  },
  methods: {
    onButtonClick() {
      this.isSelecting = true
      window.addEventListener('focus', () => {
        this.isSelecting = false
      }, { once: true })

      this.$refs.uploader.click()
    },

    loadFiles(){
      this.$store.dispatch("GET_GALLERY");
    },
    close(){
      this.dialog = false
    },
     
    logout() {
      this.$store.dispatch("auth/logout").then(()=>{
        this.$router.push('/login');
      },() => {
        this.$router.push('/login');
      });
    },   
   async onChangeSelectFiles(e) {
      let files = e.target.files
      this.uploading = true;
      this.dialog = true;
      this.img_max = files.length;
      for (let i = 0; i < files.length; ++i) {
        if (!this.dialog) {
          break;
        }
        this.progress = 0;
        this.img_count = i + 1;
        this.currentFile = files[i];

        let reader = new FileReader();
        reader.readAsDataURL(this.currentFile);
        reader.onload = (e) => {
          this.imageData = e.target.result;
        };

        await FilesService.upload(this.currentFile, (event) => {
          this.progress = Math.round((100 * event.loaded) / event.total);
        });
      }
      
      this.uploading = false;
      this.loadFiles()
    },     
  }
};
</script>
