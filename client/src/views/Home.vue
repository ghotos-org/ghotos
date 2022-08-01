<template>
  <v-app id="ci-home"  v-scroll="pageScroll"  >
    <Header />
    <v-main class="ci-app-content">
      <v-container>
                    <PhotoDialog v:on="closeDialog" />

        <div v-for="item in gallery" v-bind:key="item.day">        
            <GalleryFiles :id="'g-' + item.day"  :item="item"  @closePhotoDialog="closePhotoDialog" />   
                     
        </div>
      </v-container>

      <v-btn  v-scroll="onScroll"  v-show="fab" fab dark fixed bottom right color="primary" @click="toTop">
        <v-icon>mdi-chevron-up</v-icon>
      </v-btn>
    </v-main>
  </v-app>
</template>
<script>
import { mapGetters } from "vuex";
import Header from "@/components/layouts/Header.vue";
import GalleryFiles from "@/components/GalleryFiles.vue";
import PhotoDialog from "@/components/PhotoDialog.vue";

export default {
  name: "home",
  template: "#home-template",
  data() {
    return {
      fileDialog: false,
      fileName: null,
      url: process.env.VUE_APP_API_URL,
      fab: false,
      scrollTimeout: null,

 

    };
  },

  components: {
    Header,
    GalleryFiles,
    PhotoDialog
   // 'ci-view':  checkView,
  },
  computed: {
    ...mapGetters(["gallery","pageIsScrolling", "selectedFile", "files", "localUpdate", "doGalleryUpdate"]),
  },
  mounted() {
    //console.log("home mounted", this.gallery)
    if (this.localUpdate.last != this.localUpdate.new){
        return this.$store.dispatch("GET_GALLERY").then(() => {
          this.localUpdate.last = this.localUpdate.new
        });    

    }



    var This = this
    setTimeout(() => {
      let id = 0
      if (document.getElementById(This.selectedFile)) {
        id = This.selectedFile
      }else if (This.files && This.files[This.selectedFile] && This.files[This.selectedFile].day && document.getElementById(This.files[This.selectedFile].day)) {
        id = This.files[This.selectedFile].day
      }
      if (id) {
        
        document.getElementById(id).scrollIntoView({
            block: "center"
        });
      }
    },400)



  },
  methods: {
    closePhotoDialog(){

      if  (this.doGalleryUpdate) {
               return this.$store.dispatch("GET_GALLERY").then(() => {
          this.localUpdate.last = this.localUpdate.new
            this.$forceUpdate()

        });
      }
  

    },
    pageScroll(){
      const vm = this
      vm.$store.commit("SET_PAGESCROLL", true); 
      clearTimeout(vm.scrollTimeout);
      vm.scrollTimeout = setTimeout(function(){
          vm.$store.commit("SET_PAGESCROLL", false); 
      },300) 
   },

    openMedia(file){  
      this.fileDialog = true
      this.fileName = file
    },
    onScroll(e) {
      if (typeof window === "undefined") return;
      const top = window.pageYOffset || e.target.scrollTop || 0;
      this.fab = top > 20;
    },

    toTop() {
      let options = {
        duration: 0,
      };
      this.$vuetify.goTo(0, options);
    },
    viewHandler(e){
          console.log(e.type) // 'enter', 'exit', 'progress'
      console.log(e.percentInView) // 0..1 how much element overlap the viewport
      console.log(e.percentTop) // 0..1 position of element at viewport 0 - above , 1 - below
      console.log(e.percentCenter) // 0..1 position the center of element at viewport 0 - center at viewport top, 1 - center at viewport bottom
      console.log(e.scrollPercent) // 0..1 current scroll position of page
      console.log(e.scrollValue) // 0..1 last scroll value (change of page scroll offset)
      console.log(e.target.rect) // element.getBoundingClientRect() result
    },
    handleScroll: function(evt, el) {
      console.log(evt)
      console.log(el)
    },    
  },
};
</script>