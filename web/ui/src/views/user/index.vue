<template lang="pug">
el-row(:gutter="15")
  el-col(:xs="24", :sm="24", :md="14", :lg="14")
    .tweet-list(v-if="isMounted" v-infinite-scroll="getTweetsByScroll" infinite-scroll-distance="200" v-loading="!isMounted")
      usercard(v-if="isMounted" :user="user" :isExist="isExist")
      el-tabs.tab(v-model="activeTab" stretch="true" @tab-click="tabClick")
        el-tab-pane(label="Tweets" name="tweets")
          tweetlist(v-if="activeTab==='tweets'" :tweets="tweets.tweets" :type="'tweets'")
        el-tab-pane(label="Media" name="media")
          tweetlist(v-if="activeTab==='media'" :tweets="media.tweets" :type="'media'")
        el-tab-pane(label="Likes" name="likes")
          tweetlist(v-if="activeTab==='likes'" :tweets="liked.tweets" :type="'likes'")
  el-col(:xs="24", :sm="24", :md="10", :lg="10")
    | &nbsp;
</template>
<script>
import { avatarProps } from "element-plus";
import rwd from "../../components/rwd/index.vue";
import tweetlist from "../../components/tweetlist/index.vue";
import usercard from "../../components/usercard/index.vue";
import { getUserByName, getLikedTweets, getTweetsByUserID} from "../../api/user";
import {} from "../../api/tweet";
export default {
  components: {
    rwd,
    tweetlist,
    usercard,
  },
  data() {
    return {
      count: 0,
      failAvatar: "el-icon-user-solid",
      activeTab: "tweets",
      user: {},
      tweets: {
        tweets: []
      },
      media: {
        tweets: []
      },
      liked: {
        tweets: []
      },
      page: [0,0],
      isMounted: false,
      isExist: false,
    };
  },
  methods: {
    async getTweetsByScroll() {
      if (!this.isExist) {
        return
      }
      let user_id = this.user.user_id
      if (this.activeTab==="tweets") {
        this.getUserTweets(user_id)
      }else{
        this.getUserLikedTweets(user_id)
      }
    },
    async getUserTweets(user_id) {
      if (this.page[0] >= this.tweets.total_pages) {
        return;
      } else {
        this.page[0] += 1;
        try {
          let res = await getTweetsByUserID(user_id, this.page[0].toString());
          this.tweets.total_pages = res.data.total_pages;
          let tweets = res.data.tweets;
          for (let t of tweets) {
            this.tweets.tweets.push(t);
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
    tabClick(tab, event) {
      if (tab.props.name==="likes") {
        this.getUserLikedTweets(this.user.user_id)
      }
    }
  },
  async mounted() {
    try {
      let res = await getUserByName(this.$route.params.user);
      this.user = res.data;
      this.isExist = true
    }catch(err) {
      console.log(err)
    }
    await this.getUserTweets(this.user.user_id)
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