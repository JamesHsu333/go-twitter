<template lang="pug">
.menu-wrapper(v-if="!item.hidden")
  router-link(:to="item.path")
    el-menu-item(:index="item.path")
      template(#title)
        i(:class="'el-icon-' + item.children[0].meta.icon")
        span
          | {{ item.children[0].meta.title }}
</template>
<script>
export default {
  name: "SideBarItem",
  props: {
    // route object
    item: {
      type: Object,
      required: true,
    },
    basePath: {
      type: String,
      default: "",
    },
  },
};
</script>