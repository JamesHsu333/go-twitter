<template lang="pug">
el-row(:gutter="15")
  el-col(:xs="24", :sm="24", :md="14", :lg="14")
    .tweet-list(v-if="isMounted" v-infinite-scroll="getTweetsByScroll" infinite-scroll-distance="200" v-loading="!isMounted")
      selfcard
      el-tabs.tab(v-model="activeTab" stretch="true" @tab-click="tabClick")
        el-tab-pane(label="Tweets" name="tweets")
          tweetlist(v-if="activeTab==='tweets'" :tweets="user.tweets")
        el-tab-pane(label="Media" name="media")
          tweetlist(v-if="activeTab==='media'" :tweets="media.tweets")
        el-tab-pane(label="Likes" name="likes")
          tweetlist(v-if="activeTab==='likes'" :tweets="liked.tweets")
  el-col(:xs="24", :sm="24", :md="10", :lg="10")
    | &nbsp;
</template>
<script>
import { avatarProps } from "element-plus";
import rwd from "../../components/rwd/index.vue";
import tweetlist from "../../components/tweetlist/index.vue";
import selfcard from "../../components/selfcard/index.vue";
import { getTweetsByUserID, getLikedTweets } from "../../api/user";
export default {
  components: {
    rwd,
    tweetlist,
    selfcard,
  },
  data() {
    return {
      count: 0,
      failAvatar: "el-icon-user-solid",
      activeTab: "tweets",
      user_id: "",
      user: {
        tweets: []
      },
      media: {
        tweets: []
      },
      liked: {
        tweets: []
      },
      page: [0,0],
      isMounted: false
    };
  },
  methods: {
    async getTweetsByScroll() {
      let user_id = this.user_id
      if (this.activeTab==="tweets") {
        this.getUserTweets(user_id)
      }else{
        console.log("switch")
        this.getUserLikedTweets(user_id)
      }
    },
    async getUserTweets(user_id) {
      if (this.page[0] >= this.user.total_pages) {
        return;
      } else {
        this.page[0] += 1;
        try {
          let res = await getTweetsByUserID(user_id, this.page[0].toString());
          this.user.total_pages = res.data.total_pages;
          let tweets = res.data.tweets;
          for (let t of tweets) {
            this.user.tweets.push(t);
            if (t.image) {
              this.media.tweets.push(t);
            }
          }
        } catch (err) {
          console.log(err);
        }
      }
    },
    async getUserLikedTweets(user_id) {
     if (this.page[1] >= this.liked.total_pages) {
        return;
      } else {
        this.page[1] += 1;
        try {
          let res = await getLikedTweets(user_id, this.page[1].toString());
          this.liked.total_pages = res.data.total_pages;
          let tweets = res.data.tweets;
          for (let t of tweets) {
            this.liked.tweets.push(t);
          }
        } catch (err) {
          console.log(err);
        }
      } 
    },
    async tabClick(tab, event) {
      if (tab.props.name==="likes") {
        this.getUserLikedTweets(this.user_id)
      }
    },
  },
  async mounted() {
    this.user_id = this.$store.getters.user.user_id
    this.getUserTweets(this.user_id)
    this.isMounted = true
  },
};
</script>
<style>
.tweet-list {
  overflow: auto;
  position: relative;
  width: 100%;
  max-height: calc(100vh - 56px);
  min-height: 500px;
  list-style: none;
}
.tab {
  align-content: center;
  justify-content: center;
  border-right: solid 1px;
  color: #e6e6e6;
}
</style>