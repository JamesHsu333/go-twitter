<template lang="pug">
div
  el-scrollbar(wrap-class="scrollbar-wrapper")
    el-menu(
      :default-active="activeMenu",
      text-color="#2e2a2a",
      active-text-color="#409eff",
      :collapse-transition="false",
      mode="vertical"
    )
      side-bar-item(
        v-for="route in routes",
        :key="route.path",
        :item="route",
        :base-path="route.path"
      )
</template>
<script>
import SideBarItem from "./SideBarItem.vue";
export default {
  components: {
    SideBarItem,
  },
  computed: {
    routes() {
      return this.$router.options.routes;
    },
    activeMenu() {
      const route = this.$route;
      const { meta, path } = route;
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu;
      }
      return path;
    },
  },
};
</script>