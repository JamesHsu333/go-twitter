<template lang="pug">
el-row(:gutter="15")
  el-col(:xs="24", :sm="24", :md="14", :lg="14")
    el-row.tweet-container(v-loading="!isMounted" :key="tweet")
      el-col(:xs="3", :sm="3", :md="3", :lg="3")
        userpop(v-if="isMounted", :user="tweet")
      el-col(:xs="18", :sm="18", :md="18", :lg="18")
        router-link.tweet-user(to="/")
          | {{ tweet.name }}
        br
        .tweet-account
          | @{{ tweet.user_name }}
      el-col(:xs="3", :sm="3", :md="3", :lg="3")
        el-tooltip.tweet-delete(
              v-if="me.user_id === tweet.user_id",
              effect="dark",
              content="Delete tweet",
              placement="top",
            )
          span(@click="confirmDelete(tweet.id)")
            i.el-icon-close(style="font-size: 1.1rem")
      el-col.tweet-main(:xs="24", :sm="24", :md="24", :lg="24")
        .tweet(v-if="isMounted" v-html="formatText(tweet.text)")
        br
        el-image.tweet-image(
          v-if="tweet.image",
          :src="'/' + tweet.image",
          :preview-src-list="['/' + tweet.image]",
          :initial-index="1",
          fit="cover"
        )
        br
        .tweet-time
          | {{ fullDate(tweet.created_at)}}
        span(style="padding: 0 10px")
          | ·
        .tweet-time(@click="isCheckLiked = true")
          b(style="color: black")
            | {{tweet.likes + " "}}
          | Likes
      el-col.tweet-operation(:xs="8", :sm="8", :md="8", :lg="8")
        .reply(@click="isActive = true")
          i.el-icon-chat-round
          span(style="padding: 0 10px")
            | {{ tweet.replys }}
      el-col.tweet-operation(:xs="8", :sm="8", :md="8", :lg="8")
        .like(@click="like(tweet.already_liked, tweet.id)")
          transition(name="bounce", mode="out-in")
            i.el-icon-star-on(
              v-if="tweet.already_liked",
              style="color: rgb(249, 24, 128)"
            )
            i.el-icon-star-off(v-else)
          transition(name="bounce", mode="out-in")
            span(
              v-if="tweet.already_liked",
              style="padding: 0 10px; color: rgb(249, 24, 128)"
            )
              | {{ tweet.likes }}
            span(v-else, style="padding: 0 10px")
              | {{ tweet.likes }}
      el-col.tweet-operation(:xs="8", :sm="8", :md="8", :lg="8")
        i.el-icon-upload2(@click="copyURL(tweet.user_name, tweet.id)")
      el-col(:xs="24", :sm="24", :md="24", :lg="24")
        el-row(style="padding: 10px 0")
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
          el-col(:xs="20", :sm="20", :md="20", :lg="20")
            span.reply-list-instruction
              | Replying to
            router-link.reply-list-user-link(:to="'/' + tweet.user_name")
              | @{{ tweet.user_name }}
        el-row.text-area
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
            .text-icon
              el-image.text-avatar(
                v-if="isMounted",
                :src="'/' + me.avatar",
                fit="cover"
              )
          el-col(:xs="20", :sm="20", :md="20", :lg="20")
            el-input.text-block(
              v-model="content",
              type="textarea",
              :autosize="{ minRows: 1, maxRows: 6 }",
              maxlength="260",
              show-word-limit,
              placeholder="Tweet your reply"
            )
            .preview-image(v-if="image", style="position: relative")
              el-button(
                @click="image = ''",
                circle,
                size="small",
                style="background-color: rgba(15, 20, 25, 0.75); color: #fff; border: none; top: 4px; left: 4px; position: absolute; z-index: 1001"
              )
                i.el-icon-close(style="font-size: 15px; font-weight: bold")
              transition(name="fade", mode="out-in")
                el-image(
                  :src="image",
                  fit="cover",
                  style="border-radius: 10px"
                ) 
            .divider(style="border-bottom: solid 1px; padding: 5px 0")
            .divider(style="padding: 5px 0")
            el-row
              el-col.text-attachment(:xs="3", :sm="3", :md="3", :lg="3")
                el-button(type="text", @click="$refs.image.click()")
                  i.el-icon-picture-outline(style="font-size: 1.3rem")
                input(
                  ref="image",
                  type="file",
                  accept="image/*",
                  @change="preview",
                  style="display: none"
                )
              el-col.text-send(:xs="21", :sm="21", :md="21", :lg="21")
                el-button.send-button(
                  @click="replyTweet(tweet.id, content)",
                  :disabled="content == ''",
                  size="small",
                  round
                )
                  | Reply
        replylist(:replys="replys.tweets", :tweet="tweet" :key="replys.tweets" :nocontent="'This tweet does not have any replies.'")
        replycard(
          v-if="isMounted",
          :tweet="tweet",
          :isActive="isActive",
          @close="updateIsActive",
          @update="updateNewTweet"
        )
        likedcard(
          v-if="isCheckLiked",
          :tweetID="tweet.id"
          :isActive="isCheckLiked",
          @close="updateIsCheckLiked",
        )
  el-col(:xs="24", :sm="24", :md="10", :lg="10")
    | &nbsp;
</template>
<script>
import likedcard from "../../components/likedcard/index.vue"
import rwd from "../../components/rwd/index.vue";
import replycard from "../../components/replycard/index.vue";
import replylist from "../../components/replylist/index.vue";
import userpop from "../../components/userpop/index.vue";
import {
  createReplyTweet,
  getTweetByID,
  getReplyTweets,
  deleteTweet
} from "../../api/tweet";
import { likeTweet, deleteLike } from "../../api/user";
import { ElMessage, ElMessageBox } from "element-plus";
const monthNames = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];
const URLMatcher = /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/igm
export default {
  components: {
    likedcard,
    rwd,
    replycard,
    replylist,
    userpop,
  },
  data() {
    return {
      tweet: {},
      isMounted: false,
      isActive: false,
      isCheckLiked: false,
      me: {},
      replys: {
        tweets: [],
      },
      content: "",
      image: "",
      activeTweet: "",
    };
  },
  methods: {
    formatDate(v) {
      let tmp = new Date(v);
      let now = new Date();
      let diff = now - tmp;
      if ((diff / (1000 * 3600)) >> 0 < 24) {
        if (diff / ((1000 * 60) >> 0) < 60) {
          return ~~(diff / (1000 * 60)) + "m";
        }
        return ~~(diff / (1000 * 3600)) + "h";
      }
      return (
        monthNames[tmp.getMonth()] +
        " " +
        tmp.getDate() +
        " " +
        tmp.getFullYear()
      );
    },
    fullDate(v) {
      let tmp = new Date(v);
      return (
        tmp.getHours() +
        ":" +
        tmp.getMinutes() +
        " · " +
        monthNames[tmp.getMonth()] +
        " " +
        tmp.getDate() +
        ", " +
        tmp.getFullYear()
      );
    },
    formatText(text) {
      text = text.replace(/[\r\n\x0B\x0C\u0085\u2028\u2029]+/g, match => "</br>")
      text = text.replace(URLMatcher, match => "<a href=" + match+">"+match+"</a>")
      return text
    },
    async delete(id) {
      try {
        await deleteTweet(id);
        this.$router.push("/")
      } catch (err) {
        console.log(err);
      }
    },
    async like(already_liked, tweet_id) {
      try {
        if (already_liked) {
          await deleteLike(this.me.user_id, tweet_id);
          this.tweet.already_liked = !already_liked;
          this.tweet.likes -= 1;
        } else {
          await likeTweet(this.me.user_id, tweet_id);
          this.tweet.already_liked = !already_liked;
          this.tweet.likes += 1;
        }
      } catch (err) {
        console.log(err);
      }
    },
    preview(event) {
      this.image = URL.createObjectURL(event.target.files[0]);
    },
    async replyTweet(id, text) {
      try {
        let form = new FormData();
        form.append("text", text);
        if (this.$refs.image.files[0]) {
          form.append("image", this.$refs.image.files[0]);
        }
        let res = await createReplyTweet(id, form);
        this.content = "";
        this.image = "";
        let newTweet = res.data;
        newTweet.user_name = this.me.user_name;
        newTweet.name = this.me.name;
        newTweet.avatar = this.me.avatar;
        newTweet.about = this.me.about;
        newTweet.user_id = this.me.user_id;
        this.replys.tweets.unshift(newTweet);
        this.$refs.image.value = null;
      } catch (err) {
        console.log(err);
      }
    },
    async copyURL(user_name, tweet_id) {
      let url =
        window.location.protocol +
        "//" +
        window.location.host +
        "/" +
        user_name +
        "/status/" +
        tweet_id;
      try {
        await navigator.clipboard.writeText(url);
        ElMessage.success("Copied to clipboard");
      } catch ($e) {
        ElMessage.error("Cannot copied to clipboard");
      }
    },
    updateIsActive(v) {
      this.isActive = v;
    },
    updateIsCheckLiked(v) {
      this.isCheckLiked = v;
    }, 
    updateNewTweet(v) {
      this.replys.tweets.unshift(v);
      this.isActive = false;
      this.tweet.replys += 1;
    },
    confirmDelete(tweet_id) {
      ElMessageBox.confirm(
        "This tweet will be permanently deleted. Continue?",
        "Warning",
        {
          confirmButtonText: "Delete",
          cancelButtonText: "Cancel",
          type: "error",
        }
      )
        .then(() => {
          this.delete(tweet_id)
        })
        .catch(() => {
        });
    },
  },
  async mounted() {
    this.me = this.$store.getters.user;
    try {
      let res = await getTweetByID(this.$route.params.tweet);
      if (this.$route.params.user != res.data.user_name.toLowerCase()) {
        this.$router.replace(
          "/" +
            res.data.user_name +
            "/status/" +
            this.$route.params.tweet
        );
      }
      this.tweet = res.data;
      res = await getReplyTweets(this.tweet.id);
      this.replys = res.data;
    } catch (err) {
      console.log(err);
    }
    this.isMounted = true;
  },
  watch: {
    async $route(to, from) {
      this.isMounted=false
      this.tweet={}
      this.replys={}
      if(!to.params.user || !to.params.tweet){
        return
      }
      try {
        let res = await getTweetByID(to.params.tweet);
        if (to.params.user != res.data.user_name) {
          this.$router.replace(
            "/" +
              res.data.user_name +
              "/status/" +
              to.params.tweet
          );
        }
        this.tweet = res.data;
        res = await getReplyTweets(this.tweet.id);
        this.replys = res.data;
      } catch (err) {
        console.log(err);
      }
      this.isMounted=true
    },
  },
};
</script>
<style>
.tweet-container {
  color: #e6e6e6;
  border-right: solid 1px;
  padding-top: 15px;
}

.tweet-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.tweet-icon:hover {
  transform: scale(1.2);
}

.tweet-delete {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.tweet-delete:hover {
  color: red;
  transform: scale(1.3);
}

.tweet-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
  line-height: 1.3rem;
}
.tweet-user:hover {
  border-bottom: solid 1px;
}

.tweet-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  line-height: 1.3rem;
}

.tweet-time {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 1rem;
  padding-top: 15px;
}

.tweet-time:hover {
  border-bottom: solid 1px;
}

.tweet {
  color: black;
  display: inline-block;
  overflow: hidden;
  font-size: 23px;
}

.tweet-intro {
  color: black;
  display: inline-block;
  overflow: hidden;
}

.tweet-main {
  padding: 3% 3%;
  color: #e6e6e6;
  border-bottom: solid 1px;
}

.tweet-operation {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
  padding: 15px 5px;
  border-bottom: solid 1px;
  border-bottom-color: #e6e6e6;
}

.tweet-image {
  border-radius: 15px;
  max-height: 500px;
}

.tweet-operation i {
  font-size: 1.3rem;
}

.reply:hover {
  color: #1ea2f1;
  transform: scale(1.2);
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.like:hover {
  color: rgb(249, 24, 128);
  transform: scale(1.2);
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.el-icon-upload2:hover {
  color: #1ea2f1;
  transform: scale(1.2);
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.text-avatar {
  border-radius: 50%;
  width: 50px;
  height: 50px;
}
.text-area {
  display: flex;
  color: #e6e6e6;
  position: relative;
  padding: 5px 0;
  max-height: calc(100vh - 56px);
  scrollbar-width: none;
  -ms-overflow-style: none;
  overflow-x: hidden;
  overflow-y: scroll;
  border-right: solid 1px;
  border-bottom: solid 1px;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.text-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.text-block .el-textarea__inner {
  border: none;
}

.text-block textarea {
  resize: none;
  font-size: 20px;
  color: black;
}

.text-block textarea::placeholder {
  color: #536471;
  font-weight: 500;
}

.text-send {
  display: flex;
  align-items: center;
  justify-content: right;
  color: #5a5e66;
}

.text-attachment {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(29, 155, 240);
}

.send-button {
  background-color: rgb(29, 155, 240);
  color: #fff;
  font-size: 15px;
  font-weight: bold;
  border: none;
  transition: 0.3s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.send-button button:disabled {
  filter: contrast(1.1);
}

.preview-image {
  width: 100%;
}

.preview-image button {
  transition: 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.preview-image button:hover {
  transform: scale(1.1);
}
</style>