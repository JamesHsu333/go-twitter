<template lang="pug">
.main
  el-image(
    style="width: 100%; height: 200px",
    :src="user.header",
    :preview-src-list="[user.header]",
    :initial-index="1",
    fit="cover"
  )

  el-row.user-container
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
    el-col(:xs="5", :sm="5", :md="5", :lg="5")
      el-image.avatar(
        :src="user.avatar",
        :preview-src-list="[user.avatar]",
        :initial-index="1",
        fit="cover"
      )
    el-col.user-operation(:xs="18", :sm="18", :md="18", :lg="18")
      .user-operation-item(v-if="isExist")
        el-tooltip.box-item(
          effect="dark",
          :content="isNotifiedTip",
          placement="bottom"
        )
          el-button(circle, @click="isNotified = !isNotified")
            transition(name="bounce", mode="out-in")
              i.el-icon-bell(v-if="!isNotified")
              i.el-icon-message-solid(v-else)
      .user-operation-item(v-if="isExist")
        transition(name="fade", mode="out-in")
          el-button(
            v-if="!user.is_following",
            round,
            type="primary",
            style="font-weight: bold",
            @click="followUser(user.user_id)"
          )
            | Follow
          el-button.following-btn(
            v-else,
            round,
            @click="unfollowUser(user.user_id)"
          )
  el-row.user-info-container
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
    el-col(:xs="22", :sm="22", :md="22", :lg="22")
      .user-name(v-if="isExist")
        | {{ user.name }}
      .user-is-not-exist(v-else)
        | This account doesn’t exist
      br
      .user-account(style="padding: 0 0")
        | @{{ user.user_name }}
      br
      .user-intro(v-if="isExist" v-html="formatAbout(user.about)")
      .user-is-not-exist-text(v-else)
        | Try searching for another.
      .user-time(v-if="isExist")
        i.el-icon-date
        span.user-date
          | Joined {{ joinAt }}
      .user-follow(v-if="isExist")
        router-link(:to="'/'+user.user_name+'/following'")
          b
            | {{ user.following }}
          span(style="margin: 0px 15px 0px 5px")
            | Following
        router-link(:to="'/'+user.user_name+'/followers'")
          b
            | {{ user.followers }}
          span(style="margin: 0px 15px 0px 5px")
            | Followers
      .follow-suggestion(v-if="isExist")
        | Not followed by anyone you’re following
    el-col(:xs="1", :sm="1", :md="1", :lg="1")
</template>
<script>
import { follow, deleteFollowing } from "../../api/user";
const URLMatcher = /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/igm
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
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
    isExist: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      csrf_token: "",
      headers: {
        "X-CSRF-Token": "",
      },
      joinAt: "",
      isNotified: false,
      isFollowing: false,
    };
  },
  computed: {
    isNotifiedTip: function () {
      return this.isNotified ? "Turn off notifications" : "Notifiy";
    },
  },
  methods: {
    formatAbout(text) {
      if (text==="") {
        return ""
      }
      text = text.replace(/[\r\n\x0B\x0C\u0085\u2028\u2029]+/g, match => "</br>")
      text = text.replace(URLMatcher, match => "<a href=" + match+">"+match+"</a>")
      return text
    },
    async followUser(user_id) {
      try{
        await follow(this.$store.getters.user.user_id, user_id)
        this.user.is_following = true
        this.user.followers+=1
        this.$store.dispatch("updateUserFollowing", 1)
      }catch(err){
        console.log(err)
      }
    },
    async unfollowUser(user_id) {
      try{
        await deleteFollowing(this.$store.getters.user.user_id, user_id)
        this.user.is_following = false
        this.user.followers-=1
        this.$store.dispatch("updateUserFollowing", -1)
      }catch(err){
        console.log(err)
      }
    },
  },
  async mounted() {
    if (this.isExist) {
      try {
        let tmp = new Date(this.user.created_at);
        this.joinAt = monthNames[tmp.getMonth()] + " " + tmp.getFullYear();
        this.csrf_token = this.$store.getters.csrf_token;
        this.headers["X-CSRF-Token"] = this.csrf_token;
      } catch (err) {
        console.log(err)
      }
    }else{
      this.user.user_name = this.$route.params.user;
    }
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

.user-is-not-exist {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
  line-height: 2rem;
  padding-top: 30px;
  font-size: 20px;
}

.user-is-not-exist-text {
  color: #536471;
  display: inline;
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
  display: inline;
  font-size: 15px;
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
.follow-suggestion {
  color: #536471;
  font-size: 13px;
  padding: 10px 0;
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
</style>