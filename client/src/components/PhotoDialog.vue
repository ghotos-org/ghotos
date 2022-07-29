<template>
    <v-dialog
      v-model="dialog"
      fullscreen
      hide-overlay
      transition="scale-transition"
    >
      <template v-slot:activator="{ on, attrs }">
            <div  v-if="fileID" v-observe-visibility="visibilityChanged" >     
                <div v-if="imageThumb"  v-bind:id="fileID"  ref="imagebox" v-bind:class="{'ci-selected-file': selectedFile && selectedFile == fileID}"  >
                <v-responsive :aspect-ratio="4 / 3" >    
                    <v-img  
                            @click="selectFile()"
                            v-bind="attrs"
                            v-on="on"
                            :src="imageThumb"
                            :lazy-src="previewThumb"
                            :aspect-ratio="4 / 3"
                            class="grey lighten-2"
                        > 
                    </v-img>
                </v-responsive>            
                </div>
                <div v-else  v-bind:id="fileID" ref="imagebox">    
                    <v-responsive :aspect-ratio="4 / 3">
                    
                    <div style="background: #cccccc; height: 100%" width="400px">&nbsp;</div>

                    </v-responsive>
                </div>    
            </div>      

      </template>
      <v-card v-if="dialog">
        <v-toolbar
          dark
          color="primary"
        >
          <v-btn
            icon
            dark
            @click="dialog = false"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>Settings</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-toolbar-items>
            <!--
            <v-btn
              dark
              text
              @click="dialog = false"
            >
              Save
            </v-btn>
            -->
            <v-btn
              dark
              text
              @click="onClickDelete()"
            >
                <v-icon>mdi-delete</v-icon>
            </v-btn>            
          </v-toolbar-items>
        </v-toolbar>
        
         <div class="ci-file" v-if="bigFileID">
          <v-img :src="showImage()"  
          :max-height="(this.$vuetify.breakpoint.height - 130)"
          :max-width="(this.$vuetify.breakpoint.width)"
          contain
          class="grey lighten-2"
          >
          <div v-if="showPrevious" class="ci-file-navi ci-file-p" @click="previous()"><div  class="btn-wrapper"><v-btn fab dark secondary><v-icon>mdi-chevron-left-circle</v-icon></v-btn></div></div>
          <div v-if="showNext" class="ci-file-navi ci-file-n" @click="next()"><div  class="btn-wrapper" ><v-btn fab dark secondary><v-icon>mdi-chevron-right-circle</v-icon></v-btn></div></div>
          </v-img>

        </div>

      </v-card>
    </v-dialog>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: ["fileID"],
  name: "GalleryFile",
   data() {
    return {
      dialog: false,
      url: process.env.VUE_APP_API_URL,      
      previewThumb: null,
      imageThumb: null,
      preview: null,
      bigFileID: null,
      bigFile: {},
      showNext: true,
      showPrevious: true,
    }
  },
  mounted() {
    this.bigFileID = this.fileID
  },
  computed: {
    ...mapGetters(["gallery","galleryDays","files","selectedFile"]),
  },  
  methods: {  
    selectFile(){
       this.$store.commit("SET_SELECTED_FILE", this.fileID);

    },
    loadFileThumb(){
      if (this.fileID && this.files[this.fileID]) {
        let imageUrl = this.url + '/f/' + this.files[this.fileID].src
        this.previewThumb = imageUrl + '=w50-h50';
        this.imageThumb = imageUrl + '=w' +  this.$refs.imagebox.clientWidth + '-h' +  this.$refs.imagebox.clientHeight;
      }
    },
    visibilityChanged(visible){
      if (!this.imageThumb) {
        if (visible) {        
          this.$store.dispatch("GET_FILE_SRC", this.fileID).then((/*resp*/) => {
              this.loadFileThumb()
          });
        }
      }
    },
    async loadFile(id){
          this.bigFileID = id
          function loadGallery(){
          if (!This.gallery){
              return This.$store.dispatch("GET_GALLERY").then(()=>{
                    This.$forceUpdate();
                  return 
              }).catch(()=>{
                return
              })
          }
          return
        }
        var This = this;
        await loadGallery()
        if (!this.files[this.bigFileID]){
            return this.$store.dispatch("GET_FILE_SRC", this.bigFileID).then(()=>{
                this.bigFile = this.files[this.bigFileID]
                this.$store.commit("SET_SELECTED_FILE", this.bigFileID);
                This.$forceUpdate();
              return 
              
            }).catch(()=>{
              return
            })
        }
        this.bigFile = this.files[this.bigFileID]
        this.$store.commit("SET_SELECTED_FILE", this.bigFileID);

  
      },
      async onClickDelete(){
          let res = await this.$dialog.confirm({
              text: 'Delete?',
              persistent: true
            })
            if (res) {        
              this.getNextFile(this.bigFileID,true).then((id)=>{
              this.$dialog.message.info('deleting',{timeout: 1000})
              this.$store.dispatch("DELETE_FILE", this.bigFileID ).then(() => {
              this.$store.dispatch("GET_GALLERY_DAY", {day : this.bigFile.day}  ).then(() => {
                  if (id === 0){
                      this.$router.push('/')
                  }else {
                    this.loadFile(id)
                  }
                }).catch( () => {
                  this.$dialog.error({
                    title: 'Error',
                    text: 'Error on Delete',
                    persistent: true
                })});    
              });                        
              });
            }
      },
      showImage(){
        if (this.bigFileID && this.files[this.bigFileID]) {
            let imageUrl = this.url + '/f/' + this.files[this.bigFileID].src
            return imageUrl + '=w' + (this.$vuetify.breakpoint.width ) + '-h' +  (this.$vuetify.breakpoint.height);
        }
        return
      },
      showPreview(){
        let imageUrl = this.url + '/f/' + this.bigFile.src
        return  imageUrl + '=w50-h50';
      },
      previous(){
        this.getNextFile(this.bigFileID,true).then((id)=>{
          this.showNext = true
          if (id !== 0){
            this.loadFile(id)  
          }else {
            this.showPrevious = false
          }
        })
      },      
      next(){

         this.getNextFile(this.bigFileID,false).then((id)=>{
           this.showPrevious = true
           if (id !== -1){
              this.loadFile(id)  
           }else {
             this.showNext = false
           }
        })
      },
      getNextFile(fileID,previous){
          return new Promise((resolve, reject) =>{
            let file = this.files[fileID]

            function pageing(This) {

              function newIndex(idx){
                if (!previous){
                    return idx + 1;
                }
                return idx -1;
              }
              function firstLastFileFromDay(galleryDayFiles) {
                let id                
                if (!previous) {
                  id = galleryDayFiles[0];
                }else {
                  id = galleryDayFiles[galleryDayFiles.length - 1];     
                }
                return id
              }

              let day = This.galleryDays[file.day]      
              let dayFiles = day.files;
              let fileIdx = dayFiles.indexOf(fileID)
              if (fileIdx == -1) {
                  return reject("no Index FOUND")
              }

              // Ende oder Anfang des Tags erreicht
              if ((!previous && fileIdx == dayFiles.length -1) || (previous && fileIdx == 0)){

                // naechster Tag finden
                let curDayIndex = -1; // aktueller indexDay
                curDayIndex = This.gallery.findIndex(x => x.day == file.day);
                if (curDayIndex == -1) {
                  return reject("Index Day nicht vorhanden!")
                }
                // naechster Tag setzen
                if (!This.gallery[newIndex(curDayIndex)]){
                    if (previous && curDayIndex == 0){
                      // im ersten eintrag
                      return resolve(0)
                    }

                    if (!previous && (This.gallery.length - 1) == curDayIndex){
                      // im letzen eintrag eintrag
                      return resolve(-1)
                    }
                    return reject("Naechster Tag nicht gefnunden")
                }
                let nextDay = This.gallery[newIndex(curDayIndex)].day;   
                // nachter Tag ueberprüfen ob vorhandne ist und ggf. laden   
                let gallerDay = This.galleryDays[nextDay];  
                                                 
                if (!gallerDay || !gallerDay.files ){
                     return This.$store.dispatch("GET_GALLERY_DAY", {day : nextDay}).then(() => {
                        gallerDay = This.galleryDays[nextDay];
                        if (!gallerDay || !gallerDay.files ){
                            return reject("Naechster Tag enhaelt keine Dateien")
                        }
                        return resolve(firstLastFileFromDay(gallerDay.files))
                    });
                }
                return resolve(firstLastFileFromDay(gallerDay.files))

              }

              return resolve(dayFiles[newIndex(fileIdx)])
            }

            if (!this.galleryDays || !this.galleryDays[file.day]){
              return this.$store.dispatch("GET_GALLERY_DAY", {day : file.day}).then(() => {
                  return pageing(this)
              });
            }
            return pageing(this)

          })


      },    
  },
};
</script>