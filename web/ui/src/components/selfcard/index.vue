<template lang="pug">
.main
  el-image(
    style="width: 100%; height: 200px",
    :key="profile.header",
    :src="profile.header",
    :preview-src-list="[profile.header]",
    :initial-index="1",
    fit="cover"
  )

  el-row.user-container
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
    el-col(:xs="5", :sm="5", :md="5", :lg="5")
      el-image.avatar(
        :key="profile.avatar"
        :src="profile.avatar",
        :preview-src-list="[profile.avatar]",
        :initial-index="1",
        fit="cover"
      )
    el-col.user-operation(:xs="18", :sm="18", :md="18", :lg="18")
      .user-operation-item
        el-button(
          round,
          style="font-weight: bold",
          @click="isSettingProfile = !isSettingProfile"
        )
          | Set up profile
  el-row.user-info-container
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
    el-col(:xs="22", :sm="22", :md="22", :lg="22")
      .user-name(to="/")
        | {{ profile.name }}
      br
      .user-account(style="padding: 0 0")
        | @{{ profile.user_name }}
      br
      .user-intro(v-if="isMounted" v-html="formatAbout(profile.about)")
      .user-time
        i.el-icon-date
        span.user-date
          | Joined {{ joinAt }}
      .user-follow
        router-link(:to="'/'+profile.user_name+'/following'")
          b
            | {{ profile.following }}
          span(style="margin: 0px 15px 0px 5px")
            | Following
        router-link(:to="'/'+profile.user_name+'/followers'")
          b
            | {{ profile.followers }}
          span(style="margin: 0px 15px 0px 5px")
            | Followers
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
    el-dialog(
      v-model="isSettingProfile",
      width="40%",
      :close-on-click-modal="'false'",
      :close-on-press-escape="'false'",
      :show-close="'false'",
      style="bordor-radius: 30px",
      center
    )
      .setting-profile
        img.logo(src="../../assets/logo.svg")
      transition(name="fade", mode="out-in")
        .step-1(v-if="settingProfile == 1")
          h3.setting-profile-title
            | Pick a profile picture
          p.setting-profile-text
            | Have a favorite selfie? Upload it now.
          .setting-profile
            el-upload.selfcard-uploader(
              :action="'/api/v1/users/' + user_id + '/avatar'",
              :headers="headers",
              :on-success="handleAvatarSuccess",
              :on-error="handleAvatarError",
              :show-file-list="false"
            )
              el-image.selfcard-uploader-avatar(
                v-if="profile.avatar",
                :key="profile.avatar",
                :src="profile.avatar",
                fit="cover"
              )
              img.selfcard-uploader-avatar(v-else, src="../../assets/user.png")
          .setting-profile
            el-button(
            @click="settingProfile = 2",
            round,
            style="font-weight: bold; color: #000"
            )
              | Skip for Now
        .step-2(v-else-if="settingProfile == 2")
          h3.setting-profile-title
            | Pick a header
          p.setting-profile-text
            | People who visit your profile will see it. Show your style.
          .setting-profile
            el-upload.selfcard-uploader(
              :action="'/api/v1/users/' + user_id + '/header'",
              :headers="headers",
              :on-success="handleHeaderSuccess",
              :on-error="handleHeaderError",
              :show-file-list="false"
            )
              el-image.selfcard-uploader-header(
                v-if="profile.header",
                :key="profile.header",
                :src="profile.header",
                fit="cover"
              )
              el-image.selfcard-uploader-header(v-else src="https://www.walldevil.co/wallpapers/w10/grey-filter-gray.png" fit="cover")
          .setting-profile
            el-button(@click="settingProfile=1;isSettingProfile=false", round, style="font-weight: bold; color: #000")
              | Skip for Now
</template>
<script>
import { getFollowing, getFollowers} from "../../api/user";
import { ElMessage } from "element-plus";
const monthNames = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];
const URLMatcher = /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/igm
export default {
  data() {
    return {
      user_id: "",
      csrf_token: "",
      headers: {
        "X-CSRF-Token": "",
      },
      profile: {
        user_name: "",
        name: "",
        gender: "male",
        email: "",
        about: "",
        phone_number: "",
        country: "",
        birthday: "",
        avatar: "",
        header: "",
        created_at: "",
        followers: 0,
        following: 0,
      },
      joinAt: "",
      isSettingProfile: false,
      settingProfile: 1,
      isMounted: false,
    };
  },
  methods: {
    handleAvatarSuccess(res) {
      this.profile.avatar = res.avatar;
      this.$store.dispatch("updateUserInfo", res);
      ElMessage({
        showClose: true,
        message: "Update Avatar success",
        type: "success",
      });
      this.settingProfile=2
    },
    handleAvatarError(err) {
      ElMessage({
        showClose: true,
        message: err,
        type: "error",
      });
    },
    handleHeaderSuccess(res) {
      this.profile.header = res.header;
      this.$store.dispatch("updateUserInfo", res);
      ElMessage({
        showClose: true,
        message: "Update Header success",
        type: "success",
      });
      this.settingProfile=1
      this.isSettingProfile=false
    },
    handleHeaderError(err) {
      ElMessage({
        showClose: true,
        message: err,
        type: "error",
      });
    },
    formatAbout(text) {
      if (text==="") {
        return ""
      }
      text = text.replace(/[\r\n\x0B\x0C\u0085\u2028\u2029]+/g, match => "</br>")
      text = text.replace(URLMatcher, match => "<a href=" + match+">"+match+"</a>")
      return text
    },
  },
  async mounted() {
    this.profile = this.$store.getters.user;
    this.user_id = this.profile.user_id;
    let tmp = new Date(this.profile.created_at);
    this.joinAt = monthNames[tmp.getMonth()] + " " + tmp.getFullYear();
    this.csrf_token = this.$store.getters.csrf_token;
    this.headers["X-CSRF-Token"] = this.csrf_token;
    this.isMounted = true
  },
};
</script>
<style>
.main {
  border-right: solid 1px;
  color: #e6e6e6;
}

.avatar {
  border-radius: 50%;
  border: solid #fff 5px;
  width: 140px;
  height: 140px;
  position: absolute;
  margin-top: -75px;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.avatar:hover {
  border: solid #e6e6e6 5px;
}

.user-operation {
  display: flex;
  align-items: center;
  justify-content: right;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
  padding: 5px;
}

.user-operation-item {
  padding: 0 5px;
}

.user-info-container {
  padding-top: 25px;
}

.user-name {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
  line-height: 1.3rem;
  font-size: 20px;
}

.user-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  line-height: 1.3rem;
  font-size: 15px;
}

.user-intro {
  border: 0 solid black;
  box-sizing: border-box;
  color: rgba(0, 0, 0, 1);
  font-size: 15px;
  display: inline;
}

.user-time {
  color: #536471;
  border: 0 solid black;
  box-sizing: border-box;
  padding: 10px 0;
}

.user-date {
  padding: 0 5px;
  font-weight: 500;
  font-size: 15px;
}

.user-follow {
  display: inline-block;
  font-weight: 500;
  overflow: hidden;
}

.user-follow a {
  color: #536471;
  transition: 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.user-follow a:hover {
  font-size: 1.05rem;
}

.user-follow b {
  color: #000;
}
.following-btn {
  font-weight: bold;
}
.following-btn:hover {
  background-color: rgb(241, 194, 216);
  color: rgb(249, 24, 128);
  border-color: rgb(241, 194, 216);
}
.following-btn:after {
  content: "Following";
}
.following-btn:hover:after {
  content: "Unfollow";
}
.el-dialog {
  border-radius: 40px;
}
.setting-profile {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 0 20px 0;
}
.selfcard-uploader {
  padding: 50px 0;
}
.selfcard-uploader-avatar {
  border-radius: 50%;
  width: 180px;
  height: 180px;
}
.selfcard-uploader-header-avatar {
  border-radius: 50%;
  width: 130px;
  height: 130px;
}
.selfcard-uploader-header {
  height: 200px;
}
.setting-profile-title {
  font-size: 23px;
  font-weight: bold;
  color: #0f1419;
  padding: 0 0 10px 0;
}
.setting-profile-text {
  font-size: 15px;
  font-weight: 600;
  color: #606266;
  padding: 0 0 10px 0;
}
</style>