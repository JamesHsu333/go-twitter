<template lang="pug">
.app-wrapper
  .logo-container
    .logo-wrapper
      router-link(to="/")
        img.logo(src="./../assets/logo.svg")
  sidebar.sidebar-container
  .main-container
    .fixed-header
      navbar
    app-main
</template>
<script>
import AppMain from "./AppMain.vue";
import Sidebar from "./sidebar/index.vue";
import Navbar from "./Navbar.vue";
export default {
  components: {
    AppMain,
    Sidebar,
    Navbar,
  },
};
</script>