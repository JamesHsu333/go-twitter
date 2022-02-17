<template lang="pug">
el-dialog(
  v-model="isActive",
  width="45%",
  :close-on-click-modal="'false'",
  :close-on-press-escape="'false'",
  :show-close="'false'",
  style="bordor-radius: 30px",
  @close="close"
  center
)
  .liked-card-title
    | Liked By
  .liked-card-list(v-infinite-scroll="getUserByScroll" infinite-scroll-distance="200")
    userlist(:users="liked_users.users" :nocontent="'ddd'")
</template>
<script>
import userlist from "../userlist/index.vue"
import {getLikedUser} from "../../api/tweet"
export default {
  props: {
    isActive: {
      type: Boolean,
      required: true,
    },
    tweetID: {
      type: Number,
      required: true,
    }
  },
  components: {
    userlist,
  },
  data() {
    return {
      liked_users:{
          users: []
      },
      page: 0
    };
  },
  methods: {
    close() {
      this.$emit("close", false);
    },
    getUserByScroll() {
      this.getUsers(this.tweetID)
    },
    async getUsers(tweet_id) {
     if (this.page >= this.liked_users.total_pages) {
        return;
      } else {
        this.page += 1;
        try {
          let res = await getLikedUser(tweet_id, this.page.toString());
          this.liked_users.total_pages = res.data.total_pages;
          let users = res.data.users;
          for (let u of users) {
            this.liked_users.users.push(u);
          }
        } catch (err) {
          console.log(err);
        }
      } 
    },
  },
};
</script>
<style>
.liked-card-title {
  display: flex;
  align-items: left;
  justify-content: left;
  padding: 0 0 20px 0;
  color: #0F1419;
  font-weight: bold;
  font-size: 20px;
}
.liked-card-list {
  overflow: auto;
  position: relative;
  width: 100%;
  max-height: calc(100vh - 56px);
  min-height: 500px;
  list-style: none;
}
</style>