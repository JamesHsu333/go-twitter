<template lang="pug">
.navbar
  breadcrumb.breadcrumb-container
  .right-menu
    .right-menu-item
      el-input.w-50.m-2(
        v-model="input1",
        size="large",
        placeholder="Search Twitter",
        prefix-icon="el-icon-search",
        style="border-radius: 40px"
      )
    el-dropdown.right-menu-item(trigger="click")
      .avatar-wrapper
        el-image.nav-avatar(
        v-if="avatar",
        :src="'/'+avatar",
        fit="cover"
      )
        el-avatar(v-else, icon="el-icon-user-solid")
        p.user-info
          | {{ name }}
          br
          .nav-user-name
            | @{{ user_name }}
      template(#dropdown)
        el-dropdown-menu
          el-dropdown-item(command="logout")
            span(style="display: block", @click="logout()")
              | Log Out @{{ user_name }}
</template>

<script>
import Breadcrumb from "../components/breadcrumb/index.vue";
import { logout } from "../api/user";
export default {
  components: {
    Breadcrumb,
  },
  data() {
    return {};
  },
  computed: {
    name() {
      let user = this.$store.getters.user;
      return user.name;
    },
    user_name() {
      let user = this.$store.getters.user;
      return user.user_name;
    },
    avatar() {
      let user = this.$store.getters.user;
      return user.avatar;
    },
  },
  methods: {
    async logout() {
      await this.$store.dispatch("removeUserInfo");
      await logout();
      this.$router.go("/login");
    },
  },
};
</script>

<style>
.navbar {
  height: 56px;
  overflow: hidden;
  position: relative;
  z-index: 1001;
  background-color: #fff;
}
.breadcrumb-container {
  float: left;
  padding-left: 10px;
}
.right-menu {
  float: right;
  height: 100%;
  display: flex;
  align-items: center;
}
.right-menu :focus {
  outline: none;
}
.right-menu-item {
  display: inline-block;
  text-align: center;
  height: max-content;
  color: #5a5e66;
  padding: 0 10px;
}
.avatar-wrapper {
  display: flex;
  border-radius: 30px;
  align-items: center;
  justify-content: center;
  margin: auto;
}
.avatar-wrapper:hover {
  background: rgba(0, 0, 0, 0.075);
  border-radius: 30px;
}
.nav-avatar {
  border-radius: 50%;
  width: 50px;
  height: 50px;
}
.user-info {
  display: inline-block;
  font-weight: bold;
  font-size: 15px;
  line-height: 1.2rem;
  padding: 0px 10px;
}
.nav-user-name {
  font-weight: 300;
  font-size: 13px;
  color:rgb(83, 100, 113);
}
</style>