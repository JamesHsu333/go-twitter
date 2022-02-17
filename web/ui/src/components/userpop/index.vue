<template lang="pug">
el-popover(placement="bottom", width="300px", trigger="hover" @show="getUserDetail(user.user_id)")
  el-row
    el-col(:xs="16", :sm="16", :md="16", :lg="16")
      el-image.pop-avatar(
        :src="'/'+user.avatar",
        fit="cover"
        )
      br
      router-link.pop-user(:to="'/'+user.user_name")
        | {{ user.name }}
      br
      .pop-account(style="padding: 0 0")
        | @{{ user.user_name }}
    el-col(:xs="8", :sm="8", :md="8", :lg="8")
        div(v-if="me.user_id != user.user_id")
          el-button(
              v-if="!isLoading && !isFollowing",
              round,
              size="small",
              type="primary",
              style="font-weight: bold",
              @click="followUser(user.user_id)"
            )
              | Follow
          el-button.following-btn(
              v-if="!isLoading && isFollowing",
              round,
              size="small",
              @click="unfollowUser(user.user_id)"
            )
  el-row
    el-col(:xs="24", :sm="24", :md="24", :lg="24")
      .pop-tweet(v-html="formatAbout(user.about)")
  el-row
    .pop-follow
      router-link(:to="'/'+user.user_name+'/following'")
        b
          | {{following}}
        span(style="margin: 0px 15px 0px 5px")
          | Following
      router-link(:to="'/'+user.user_name+'/followers'")
        b
          | {{followers}}
        span(style="margin: 0px 15px 0px 5px")
          | Followers
  template(#reference)
    router-link(:to="'/' + user.user_name")
      .pop-icon
        el-image.tweet-avatar(
        :src="'/'+user.avatar",
        fit="cover"
        )
</template>
<script>
import { getUserByID, follow, deleteFollowing} from "../../api/user";
const URLMatcher = /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/igm
export default {
  props: {
    user: {
      type: Object,
      required: true, 
    },
  },
  data() {
    return {
      followers: 0,
      following: 0,
      isFollowing: "",
      isLoading: true,
      me: {},
    }
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
    async getUserDetail(user_id) {
      let res = await getUserByID(user_id)
      let user = res.data
      this.following = user.following
      this.followers = user.followers
      this.isFollowing = user.is_following
      this.isLoading = false
    },
    async followUser(user_id) {
      try{
        await follow(this.$store.getters.user.user_id, user_id)
        this.isFollowing = true
        this.followers+=1
        this.$store.dispatch("updateUserFollowing", 1)
      }catch(err){
        console.log(err)
      }
    },
    async unfollowUser(user_id) {
      try{
        await deleteFollowing(this.$store.getters.user.user_id, user_id)
        this.isFollowing = false
        this.followers-=1
        this.$store.dispatch("updateUserFollowing", -1)
      }catch(err){
        console.log(err)
      }
    },
  },
};
</script>
<style>
.pop-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.pop-icon .tweet-avatar {
  border-radius: 50%;
  width: 50px;
  height: 50px;
}

.pop-icon:hover {
  transform: scale(1.2);
}

.tweet-list-dropdown {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.tweet-list-dropdown:hover {
  color: #1ea2f1;
  transform: scale(1.3);
}

.pop-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
}
.pop-user:hover {
  border-bottom: solid 1px;
}

.pop-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.pop-tweet {
  color: black;
  display: inline-block;
  overflow: hidden;
  font-size: 15px;
}

.pop-follow {
  display: inline-block;
  font-weight: 500;
  overflow: hidden;
}

.pop-avatar {
  border-radius: 50%;
  width: 70px;
  height: 70px;
}

.pop-follow a {
  color: #536471;
  transition:  0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.pop-follow a:hover {
  font-size: 1.05rem;
}

.pop-follow b {
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
</style>