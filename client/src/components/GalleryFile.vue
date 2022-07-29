
<template>
  <div  v-observe-visibility="visibilityChanged" >

    <div v-if="image"  v-bind:id="file"  ref="imagebox" v-bind:class="{'ci-selected-file': selectedFile && selectedFile == file}"  >
      <v-responsive :aspect-ratio="4 / 3" >    
      <router-link :to="{ name: 'photo', params: { file }}"> 
        <v-img
                :src="image"
                :lazy-src="preview"
                :aspect-ratio="4 / 3"
                class="grey lighten-2"
              > 
        </v-img>
        </router-link>
      </v-responsive>
   
    </div>
    <div v-else  v-bind:id="file" ref="imagebox">    
        <v-responsive :aspect-ratio="4 / 3">
          <div style="background: #cccccc; height: 100%" width="400px">&nbsp;</div>
        </v-responsive>
    </div>    
  </div>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: ["file"],
  name: "GalleryFile",
   data() {
    return {
      url: process.env.VUE_APP_API_URL,      
      preview: null,
      image: null
    }
  },
  mounted() {
   // console.log("test",this.file)
    this.loadFile()

  },
  computed: {
    ...mapGetters(["files", "selectedFile"]),
  },  
  methods: {
    openMedia(file){  
      //  @click="openMedia(file)"
        this.$store.dispatch("SHOW_FILE", file)
    },    
    loadFile(){
      if (this.file && this.files[this.file]) {
        let imageUrl = this.url + '/f/' + this.files[this.file].src
        this.preview = imageUrl + '=w50-h50';
        this.image = imageUrl + '=w' +  this.$refs.imagebox.clientWidth + '-h' +  this.$refs.imagebox.clientHeight;
      }
    },
    visibilityChanged(visible){
      if (!this.image) {
        if (visible) {        
          this.$store.dispatch("GET_FILE_SRC", this.file).then((/*resp*/) => {
              this.loadFile()
              /*
              let imageUrl = this.url + '/f/' + resp.data.src
              this.preview = imageUrl + '=w50-h50';
              this.image = imageUrl + '=w' +  this.$refs.imagebox.clientWidth + '-h' +  this.$refs.imagebox.clientHeight;
              */
          });

          
        }
      }
    }
  },
};
</script>
