<template lang="pug">
rwd
  .err-container
    el-button.pan-back-btn(icon="el-icon-arrow-left", @click="back")
      | Back
    el-row(:gutter="30")
      el-col(:span="12")
        h1.text-jumbo.text-ginormous.text
          | Oops!
        h2.text
          | You don't have permission to access this resource
        h5.text
          | Please contact admin to solve your problem
      el-col(:span="12")
        img(:src="errGif" width="313" height="428" alt="Girl has dropped her ice cream.")
</template>
<script>
import rwd from "../../components/rwd/index.vue";
import errGif from "../../assets/forbidden.gif"
export default {
  name: 'forbidden',
  components: {
    rwd,
  },
  data() {
    return {
      errGif: errGif + '?' + +new Date()
    }
  },
  methods: {
    back() {
      if (this.$route.query.noGoBack) {
        this.$router.push({ path: '/home' })
      } else {
        this.$router.go(-1)
      }
    }
  }
}
</script>
<style>
.err-container {
  width: 800px;
  max-width: 100%;
  margin: 100px auto;
}
.pan-back-btn {
  background: #008489;
  color: #fff;
  border: none !important;
}
.text-jumbo {
  font-size: 60px;
  font-weight: 1000;
  color: #484848;
}
.text {
  padding: 30px 0;
}
</style>