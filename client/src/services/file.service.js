
import ApiService from "./api.service";


class FilesService {
    upload(file, onUploadProgress) {
      let formData = new FormData();
      formData.append("file", file);
      formData.append("time", file.lastModified);
      formData.append("date", file.lastModifiedDate);
      
      return ApiService.post("/file", formData, {
        headers: {
          "Content-Type": "multipart/form-data"
          
        },
        onUploadProgress
      });
    }
  
    getGallery() {
      return ApiService.get("/gallery").then((resp) => {
       // console.log(resp)
       return resp
      }).catch(() => {
        //  alert(e.message)
      });
    }
    getGalleryDay(day) {
      return ApiService.get("/gallery/" + day).then((resp) => {
       // console.log(resp)
       return resp
      }).catch(() => {
        //  alert(e.message)
      });
    }    
    getFileSrc(id) {
      return ApiService.get("/file/src/" + id).then((resp) => {
       // console.log(resp)
       return resp
      }).catch(() => {
        //  alert(e.message)
      });
    }       
    deleteFile(id) {
      return ApiService.delete("/file/" + id);
    }      
    
  }
  
  export default new FilesService();