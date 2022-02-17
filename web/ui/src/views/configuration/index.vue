<template lang="pug">
rwd
  el-row(:gutter="30")
    el-col(:xs="24", :sm="24", :md="24", :lg="24")
      .title
        i.el-icon-setting
        |  Configuration
      el-menu(:default-active="activeMenu",
      class="el-menu-demo"
      mode="horizontal",
      )
        menu-item(
        v-for="route in routes",
        :key="route.path",
        :item="route",
        :base-path="route.path"
        )
  el-row(:gutter="30")
    el-col(:xs="24", :sm="24", :md="24", :lg="24")
      router-view(v-slot="{ Component }")
        transition(name="fade" mode="out-in")
          component(:is="Component")
</template>
<script>
import rwd from "../../components/rwd/index.vue";
import MenuItem from "./MenuItem.vue";
export default {
  components: {
    rwd,
    MenuItem,
  },
  computed: {
    routes() {
      let config = {}
      for(const route of this.$router.options.routes) {
          if(route.path === '/configuration'){
              config = route.children
          }
      }

      return config;
    },
    activeMenu() {
      const route = this.$route;
      const { meta, name } = route;
      if (meta.activeMenu) {
        return meta.activeMenu;
      }
      return name;
    },
  },
};
</script>
<style>
.title {
  line-height: 1.25;
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 32px 0;
}
</style>