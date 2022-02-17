<template lang="pug">
div
  el-empty(v-if="users.length===0" :description="nocontent")
  article.user-list-item(v-for="(u, index) in users", :key="u")
    transition(name="bounce")
      el-row
        el-col(:xs="3", :sm="3", :md="3", :lg="3")
          userpop(:user="u")
        el-col(:xs="17", :sm="17", :md="17", :lg="17")
          router-link.user-list-user(:to="'/' + u.user_name")
            | {{ u.name }}
          br
          span.user-list-account
            | @{{ u.user_name }}
          br
          .user-list-about(v-html="formatText(u.about)")
        el-col(:xs="4", :sm="4", :md="4", :lg="4")
          .user-operation-item(v-if="me.user_id != u.user_id")
            transition(name="fade", mode="out-in")
              el-button(
                v-if="!u.is_following",
                round,
                type="primary",
                size="small",
                style="font-weight: bold",
                @click="followUser(index, me.user_id, u.user_id)"
              )
                | Follow
              el-button.following-btn(
                v-else,
                round,
                size="small",
                @click="unfollowUser(index, me.user_id, u.user_id)"
              )
</template>
<script>
import userpop from "../userpop/index.vue";
import { ElMessage } from "element-plus";
import { follow, deleteFollowing } from "../../api/user";
const URLMatcher =
  /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/gim;
export default {
  props: {
    users: {
      type: Object,
      required: true,
    },
    nocontent: {
      type: String,
      required: false
    }
  },
  components: {
    userpop,
  },
  data() {
    return {
      isMounted: false,
      isActive: false,
      me: {}
    };
  },
  methods: {
    formatText(text) {
      if (text == "") {
        return;
      }
      text = text.replace(
        /[\r\n\x0B\x0C\u0085\u2028\u2029]+/g,
        (match) => "</br>"
      );
      text = text.replace(
        URLMatcher,
        (match) => "<a href=" + match + ">" + match + "</a>"
      );
      return text;
    },
    async followUser(index, follower, following) {
      try{
        await follow(follower, following)
        this.users[index].is_following = true
        this.$store.dispatch("updateUserFollowing", 1)
      }catch(err){
        console.log(err)
      }
    },
    async unfollowUser(index, follower, following) {
      try{
        await deleteFollowing(follower, following)
        this.users.splice(index, 1)
        this.$store.dispatch("updateUserFollowing", -1)
      }catch(err){
        console.log(err)
      }
    },
  },
  mounted() {
    this.me = this.$store.getters.user;
  },
};
</script>
<style>
.user-list-item {
  color: #e6e6e6;
  position: relative;
  padding-bottom: 50px;
  padding: 15px 0;
  max-height: calc(100vh - 56px);
  scrollbar-width: none;
  -ms-overflow-style: none;
  overflow-x: hidden;
  overflow-y: scroll;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.user-list-item:hover {
  background-color: rgba(0, 0, 0, 0.02);
}

.user-list-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
}
.user-list-user:hover {
  border-bottom: solid 1px;
}

.user-list-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.user-list-time {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.user-list-time:hover {
  border-bottom: solid 1px;
}

.user-list-about {
  color: black;
  display: inline-block;
  overflow: hidden;
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
