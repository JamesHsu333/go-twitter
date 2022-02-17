<template lang="pug">
section.app-main
  router-view(v-slot="{ Component }")
    transition(name="fade" mode="out-in")
      component(:is="Component")
</template>

<script>
export default {
  name: "AppMain"
};
</script>
<style>
.app-main {
  background-color: #fff;
}
</style>