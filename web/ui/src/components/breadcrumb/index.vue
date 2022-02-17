<template lang="pug">
el-breadcrumb.app-breadcrumb(separator="/")
  transition-group(name="breadcrumb")
    el-breadcrumb-item(v-for="(item, index) in levelList", :key="item.path")
      span.no-redirect(
        v-if="item.redirect === 'noRedirect' || index == levelList.length - 1"
      )
        | {{ item.meta.title }}
      router-link(v-else, :to="handleLink(item)")
        | {{ item.meta.title }}
</template>

<script>
export default {
  data() {
    return {
      levelList: null,
    };
  },
  watch: {
    $route() {
      this.getBreadcrumb();
    },
  },
  created() {
    this.getBreadcrumb();
  },
  methods: {
    getBreadcrumb() {
      let matched = this.$route.matched.filter(
        (item) => item.meta && item.meta.title
      );
      const first = matched[0];
      if (!this.isHome(first)) {
        matched = [
          { path: "/home", meta: { title: "Home" } },
        ].concat(matched);
      }
      this.levelList = matched.filter(
        (item) => item.meta && item.meta.title && item.meta.breadcrumb !== false
      );
    },
    isHome(route) {
      const name = route && route.name;
      if (!name) {
        return false;
      }
      return (
        name.trim().toLocaleLowerCase() === "Home".toLocaleLowerCase()
      );
    },
    pathCompile(path) {
      const { params } = this.$route
      var toPath = pathToRegexp.compile(path)
      return toPath(params)
    },
    handleLink(item) {
      const { redirect, path } = item;
      if (redirect) {
        return redirect;
      }
      return path;
    },
  },
};
</script>

<style>
.app-breadcrumb {
  display: inline-block;
  margin-left: 8px;
}
.el-breadcrumb {
  line-height: 50px !important;
}
.no-redirect {
  color: #303133;
  cursor: text;
}
</style>