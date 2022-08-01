
<template>
    <v-row v-observe-visibility="visibilityChanged"  :dense="true" v-bind:id="item.day" ref="days">
      <v-col cols="12">{{ item.day | moment("dddd, Do MMMM YYYY") }}</v-col>
     
      <template  v-if="files" >
        <v-col 
          v-for="file in files"
          v-bind:key="file"
          :ref="file"
          :id="'th-' + file"
          class="d-flex child-flex"
          cols="6 col-sm-3 col-md-3 col-lg-2  col-xl-1"          
        >     
          <PhotoDialog :fileID="file" @closePhotoDialog="closePhotoDialog" />
        </v-col>
      </template>        
      <template v-else>
        <v-col 
          v-for="index in item.count"
          v-bind:key="index"
          class="d-flex child-flex"
          cols="6 col-sm-3 col-md-3 col-lg-2  col-xl-1"
          
        >       
          <div ref="infoBox">    
            <v-responsive :aspect-ratio="4 / 3">
            
              <div style="background: #cccccc; height: 100%" width="400px">&nbsp;</div>

            </v-responsive>
          </div>
        </v-col>
      </template>
    
    </v-row>
</template>
<script>
import { mapGetters } from "vuex";
//import GalleryFile from "@/components/GalleryFile.vue";
import PhotoDialog from "@/components/PhotoDialog.vue";


export default {
  props: ["item"],
  name: "GalleryFiles",
   data() {
    return {
      url: process.env.VUE_APP_API_URL,
      visible: false,
      files: null,
    }
  },
  components: {
   // GalleryFile,
    PhotoDialog

  },  
  

  watch: {
      update: function (newValue, oldValue) {
        if (newValue != oldValue){
          console.log("watche" + this.item.day, newValue, oldValue)
              return  this.$store.dispatch("GET_GALLERY_DAY", {day : this.item.day}).then(() => {
            this.files = this.galleryDays[this.item.day].files

          });
        }
    }

  },



  computed: {
    ...mapGetters(["galleryDays","pageIsScrolling", "selectedFile"]),
    update: {
        get: function (){
          if (this.item.day && this.galleryDays && this.galleryDays[this.item.day] && this.galleryDays[this.item.day].update) {
              return this.galleryDays[this.item.day].update.last
          }
          return 1
        },
        set: function() {},        
    }


  },  
  mounted(){
   // console.log("mounted galleryfiles",this.item, this.item.day)
    //console.log("mounted galleryfiles")
    //console.log(this.update, this.lastUpdate)
    //this.update = this.lastUpdate
    /*
    console.log(this.visible)
    if (this.visible && !this.galleryDays[this.item.day] || (this.galleryDays[this.item.day] && this.galleryDays[this.item.day].update.last !=  this.galleryDays[this.item.day].update.new)) {
      console.log("innen")
      this.visible = false
        return  this.$store.dispatch("GET_GALLERY_DAY", {day : this.item.day}).then(() => {
          this.galleryDays[this.item.day].update.last = this.galleryDays[this.item.day].update.new
          this.files = this.galleryDays[this.item.day].files
          this.visible = true
        return;
      });
    }
    */
  }, 
  methods: {
    closePhotoDialog(){
     // console.log("test", $e.day)

      /*
      console.log("TEST Close")
                          this.$forceUpdate();

      let a = document.getElementById('th-' + this.selectedFile);

        console.log(a)

         a.scrollIntoView({behavior: 'smooth'}, true);
*/
      this.$forceUpdate()
      this.$emit('closePhotoDialog')

    },
    openFile(file){  
        this.$emit('openFile', file)
    },  
    openMedia(file){  
      this.fileDialog = true
      this.fileName = file
    },    
    visibilityChanged(visible){
      if (visible) {
          if (!this.galleryDays[this.item.day] || (this.galleryDays[this.item.day] && this.galleryDays[this.item.day].update.last !=  this.galleryDays[this.item.day].update.new)) {
            return  this.$store.dispatch("GET_GALLERY_DAY", {day : this.item.day}).then(() => {
              this.galleryDays[this.item.day].update.last = this.galleryDays[this.item.day].update.new
              this.files = this.galleryDays[this.item.day].files
              this.visible = true

            return;
          });
          }
          
          if (this.galleryDays[this.item.day] && this.galleryDays[this.item.day].files  ) {
            this.files = this.galleryDays[this.item.day].files
            this.visible = true

            return;
          } 
        /*
          if (this.galleryDays[this.item.day] && this.galleryDays[this.item.day].files  ) {
            this.files = this.galleryDays[this.item.day].files
            this.visible = true
            return;
          } 
          return  this.$store.dispatch("GET_GALLERY_DAY", {day : this.item.day}).then(() => {
            this.files = this.galleryDays[this.item.day].files
            this.visible = true
            return;
          });

          */

       }
    }
  },
};
</script>
