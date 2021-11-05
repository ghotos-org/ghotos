
<template>
<v-app id="ci-photo"    >
 <v-app-bar app fixed>
      <v-btn :to="'/'">
        <v-icon> mdi-arrow-left</v-icon>
      </v-btn>
        <v-toolbar-title>
      
        </v-toolbar-title>
        <v-spacer />
          <v-toolbar-items>
            <v-btn  text @click="onClickDelete()"><v-icon dark>mdi-delete</v-icon></v-btn>
            <v-menu bottom left>
            <template v-slot:activator="{ on, attrs }">
              <v-btn  icon v-bind="attrs" v-on="on" >
                <v-icon>mdi-dots-vertical</v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-item>
              <v-list-item-icon>
                  <v-icon >mdi-delete</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title>delete</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
            </v-menu>            
          </v-toolbar-items>

     <!-- <v-progress-linear v-if="isLoading" fixed top indeterminate ></v-progress-linear> -->
    </v-app-bar>
    <v-main class="ci-app-content">
           <v-card>
        <div class="ci-file" v-if="file && file.src">
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
    </v-main>
  </v-app>

</template>
<script>
//import { mapState } from "vuex";
import { mapGetters } from "vuex";



export default {

  name: "Photo",
   data() {
    return {
      url: process.env.VUE_APP_API_URL,      
      preview: null,
      fileID: null,
      file: {},
      showNext: true,
      showPrevious: true
    }
  },
  computed: {
    ...mapGetters(["gallery","galleryDays","files","selectedFile"]),
  },  
  mounted(){
    if (this.$route.params.file) {
       
        this.loadFile(this.$route.params.file)
    }

  },
  methods: {
      async loadFile(id){
          this.fileID = id
          function loadGallery(){
          if (!This.gallery){
              return This.$store.dispatch("GET_GALLERY").then(()=>{
                  return 
              }).catch(()=>{
                return
              })
          }
          return
        }
        var This = this;
        await loadGallery()
        if (!this.files[this.fileID]){
            return this.$store.dispatch("GET_FILE_SRC", this.fileID).then(()=>{
              this.file = this.files[this.fileID]
              this.$store.commit("SET_SELECTED_FILE", this.fileID);
              return 
              
            }).catch(()=>{
              return
            })
        }
        this.file = this.files[this.fileID]
        this.$store.commit("SET_SELECTED_FILE", this.fileID);
  
      },
      async onClickDelete(){
          let res = await this.$dialog.confirm({
              text: 'Delete?',
              persistent: true
            })
            if (res) {        
              this.getNextFile(this.fileID,true).then((id)=>{
              this.$dialog.message.info('deleting',{timeout: 1000})
              this.$store.dispatch("DELETE_FILE", this.fileID ).then(() => {
              this.$store.dispatch("GET_GALLERY_DAY", {day : this.file.day}  ).then(() => {
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
        let imageUrl = this.url + '/f/' + this.file.src
        return imageUrl + '=w' + (this.$vuetify.breakpoint.width ) + '-h' +  (this.$vuetify.breakpoint.height);
      },
      showPreview(){
        let imageUrl = this.url + '/f/' + this.file.src
        return  imageUrl + '=w50-h50';
      },
      previous(){
        this.getNextFile(this.fileID,true).then((id)=>{
          this.showNext = true
          if (id !== 0){
            this.loadFile(id)  
          }else {
            this.showPrevious = false
          }
        })
      },      
      next(){
         this.getNextFile(this.fileID,false).then((id)=>{
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
