<template>
  <div
   
    class="ci-main-frame"
    v-bind:class="{
      'ci-theme--light': !theme.isDark,
      'ci-theme--dark': theme.isDark,
    }"
  >
    <div class="ci-main-container">
      <router-view></router-view>
   
    </div>
  </div>
</template>


<script>


export default {
  name: "App",
  components: {

  },
  data() {
    return {
      timer: '',
      theme: this.$vuetify.theme,
    };
  },
  mounted() {
    this.$moment.locale(process.env.VUE_APP_LOCALE)
    this.fetchEventsList();
    this.timer = setInterval(this.fetchEventsList, 10000);    
  },
  beforeDestroy () {
    this.cancelAutoUpdate();
  },
  methods: {
    fetchEventsList () {
      //console.log("TEST")
     
    },
    cancelAutoUpdate () {
        clearInterval(this.timer);
    },
  },  
};
</script>
