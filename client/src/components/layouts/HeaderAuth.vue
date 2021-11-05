
<template>
    <v-app-bar app fixed>       
        <v-btn  icon v-if="backLink"   :to="backLink" >
            <v-icon> mdi-arrow-left</v-icon>
        </v-btn>        
        <v-toolbar-title>{{ title }}</v-toolbar-title>
    </v-app-bar> 
</template>
<script>

export default {
  name: "HeaderAuth",
  props: ["title", "backLink"],
  data() {
    return {
    }},   
  computed: {

  },      
  mounted(){
  },
};
</script>
