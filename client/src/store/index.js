

import Vue from 'vue';
import Vuex from 'vuex';

import { auth } from './auth.module';
import FilesService from "../services/file.service";

Vue.use(Vuex);

//Vue.use(Vuex);


const state = {
  loggedIn:"test",
    selectedFile: null,
    localUpdate: {last: 0, new: 1},
    newLocalUpdate: null,    
    lastUpdate: 0,    
    gallery: null,
    galleryDays: {},
    files: {},
    loaderVerzoegert: false,
    isRequesting: false,
    isLoading: false,
    pageIsScrolling: false,
    doGalleryUpdate: false,
};
  
const getters = {
  loggedIn(state) {
    return state.loggedIn;
  },     
  localUpdate(state) {
    return state.localUpdate;
  },      
  selectedFile(state) {
    return state.selectedFile;
  },    
  pageIsScrolling(state) {
    return state.pageIsScrolling;
  },    
  gallery(state) {
    return state.gallery;
  },
  files(state) {
    return state.files;
  },  
  galleryDays(state) {
    return state.galleryDays;
  },  
  isLoading(state) {
    return state.isLoading;
  },
  isRequesting(state) {
    return state.isRequesting;
  },
  doGalleryUpdate(state) {
    return state.doGalleryUpdate;
  } 

}

/*
function getIndexDayFromGalleryArray(array, day){
  for (const [index, el] of array) {
    if (el.day === day) {
      return index
    }
  }
  return -1
}
*/

const actions = {
  ["GET_GALLERY"]({ commit }) {
    commit("FETCH_START");
    return FilesService.getGallery().then((response) => {
      commit("FETCH_END");   
      commit("SET_GALLERY_LIST", response.data);
    }).catch(({ response }) => {
      commit("FETCH_END");
      console.log("ERROR:",response)
    });
  },   
  ["GET_GALLERY_DAY"]({ commit }, data) {
    commit("FETCH_START");
    //console.log("GET_GALLERY_DAY", data)
    return FilesService.getGalleryDay(data.day).then((response) => {
      commit("FETCH_END");      
      commit("SET_GALLERY_DAY", {response, day: data.day});
      return response
    }).catch(({ response  }) => {
      commit("FETCH_END");
      return response
    });
  },
  ["GET_FILE_SRC"]({ commit }, id) {
    commit("FETCH_START");
    return FilesService.getFileSrc(id).then((response) => {
      commit("FETCH_END");      

      if (response.data.src){
        commit("SET_FILE", {id, file: response.data});  
      }
      return response
    }).catch(({ response  }) => {
      commit("FETCH_END");
      return response
    });
  },
  ["DELETE_FILE"]({ commit }, id) {

    commit("FETCH_START");
   return FilesService.deleteFile(id).then(() => {
      commit("FETCH_END");      
      commit("DO_GALLERY_UPDATE");
      
     // commit("SET_LOCAL_UPDATE");
     let pfile = this.state.files[id];
     //console.log(file.day)
     if (pfile.day){

      
     // commit("SET_LOCAL_DAY_UPDATE",file.day);
    //  delete this.state.files[id];
    //  console.log(this.state.galleryDays[file.day])
    
      let files = {}
      Object.assign(files, this.state.files);
      delete files[id];

      Object.assign(this.state.files, files);

  
      let  galleryFiles = []

       this.state.galleryDays[pfile.day].files.forEach((file) => {
          if (file != id) {
            galleryFiles.push(file);
          }
       });

       this.state.galleryDays[pfile.day].files = galleryFiles;

     }else{
      commit("SET_LOCAL_UPDATE");

     }

    }).catch(() => {
      commit("FETCH_END");
    });
  },        
}


const mutations = { 
  ["FETCH_START"](state) {
    state.isRequesting = true;
    if (state.loaderVerzoegert) {
      clearTimeout(state.loaderVerzoegert);
    }

    state.loaderVerzoegert = setTimeout(() => {
      if (state.isRequesting) {
        state.isLoading = true;
      }
    }, 200);
  },
  ["FETCH_END"](state) {
    clearTimeout(state.loaderVerzoegert);
    state.isRequesting = false;
    state.isLoading = false;
  },

  ["SET_GALLERY_LIST"](state, data) {
      state.lastUpdate = data.updated
      state.lastLocalUpdate = new Date().getTime()
      state.doGalleryUpdate = false   
      state.gallery = data.days;
      if (state.lastUpdate != 0) {
        data.days.forEach(d => {
            if (state.galleryDays[d.day]) {
                if (state.galleryDays[d.day].files.length != d.count) {
                    state.galleryDays[d.day].update.last = new Date().getTime()
                }
            }
        });
      }

  },
  ["DO_GALLERY_UPDATE"](state) {
    state.doGalleryUpdate = true
  },    
  ["SET_FILE"](state, {id,file}) {
    state.files[id] = file
  },  
  ["REMOVE_GALLERY_DAY"](state, day) {
   

    console.log(state.gallery)
    state.gallery.forEach((item,index) => {
            if (item.day == day) {
               delete state.gallery[index];
               return
            }
    });
  },    
  ["SET_SELECTED_FILE"](state, file) {
    state.selectedFile = file
  },   
  ["SET_GALLERY_DAY"](state, {day, response}) {
    let files = [];
    response.data.forEach(file => {
      files.push(file.id)
      if (file.src){
        state.files[file.id] = {src: file.src, day: day} 
      }
    });
    let update = {
      new:  new Date().getTime(),
      last: 0
    }
    state.galleryDays[day] = {files, update}
  },    
  ["SET_LOCAL_DAY_UPDATE"](state,day) {
    if (state.galleryDays[day]){
      state.galleryDays[day].update.new = new Date().getTime()
    }
    
  },
  ["SET_LOCAL_UPDATE"](state) {
    state.localUpdate.new = new Date().getTime()
    /*
    let file = state.files[id];
    if (state.galleryDays[file.day]) {

      if (state.galleryDays[file.day].files && state.galleryDays[file.day].files){
        let index = state.galleryDays[file.day].files.indexOf(id)
        state.galleryDays[file.day].files.splice(index, 1)



        let dayIndex =  state.gallery.findIndex(x => x.day == file.day);

        if ()

        if (state.galleryDays[file.day].files.length == 0){

          
          delete state.galleryDays[file.day];
          

        }else {
          state.galleryDays[file.day].updated = new Date().getTime()
          state.gallery
        }


      }

      delete state.files[id]

        }
      */



/*

    let file = state.files[id];

    //this.state.gallery

    if (state.galleryDays[file.day]) {
        console.log( state.galleryDays[file.day], state.galleryDays[file.day].updated)

        if (state.galleryDays[file.day].files.length) {
          state.galleryDays[file.day].updated = new Date().getTime()
      }   
          console.log( state.galleryDays[file.day].updated)
        //  delete state.files[id]

    }

    /*
    data.days.forEach(d => {
      if (state.galleryDays[file.day]) {
          if (state.galleryDays[file.day].files.length != d.count) {
              state.galleryDays[file.day].updated = new Date().getTime()

              delete state.files[id]
              return;
          }
      }
    });
    */

  },   
  ["SET_PAGESCROLL"](state, scrolling) {    
    state.pageIsScrolling = scrolling
  }
  
}

export default new Vuex.Store({
    state,
    getters,
    actions,
    mutations,
    modules: {
      auth
    }
});
