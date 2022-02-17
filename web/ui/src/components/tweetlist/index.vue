<template lang="pug">
div
  replycard(
    :tweet="activeTweet",
    :isActive="isActive",
    @close="updateIsActive",
    @update="updateNewTweet",
    :key="activeTweet"
  )
  el-empty(
    v-if="tweets.length === 0",
    :description="'This user don’t have any ' + type + ' yet.'"
  )
  article.tweet-list-item(v-for="(t, index) in tweets", :key="t")
    transition(name="bounce")
      div
        el-row
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
            userpop(:user="t")
          el-col(:xs="18", :sm="18", :md="18", :lg="18")
            router-link.tweet-list-user(:to="'/' + t.user_name")
              | {{ t.name }}
            span.tweet-list-account
              | @{{ t.user_name }}
            span.tweet-list-account
              | ·
            .tweet-list-time
              el-tooltip(
                effect="dark",
                :content="fullDate(t.created_at)",
                placement="top"
              )
                span
                  | {{ formatDate(t.created_at) }}
            br
            .tweet-list-tweet(v-html="formatText(t.text)")
            br
            el-image.tweet-list-image(
              v-if="t.image",
              :src="t.image",
              :preview-src-list="[t.image]",
              :initial-index="1",
              fit="cover"
            )
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
            el-tooltip.tweet-list-delete(
              v-if="me.user_id === t.user_id",
              effect="dark",
              content="Delete tweet",
              placement="top",
            )
              span(@click="confirmDelete(index,t.id)")
                i.el-icon-close(style="font-size: 1.1rem")
        el-row
          el-col.tweet-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            .reply(@click="activate(t)", :key="isActive")
              i.el-icon-chat-round(style="font-size: 1.1rem")
              span(style="font-size: 0.8rem; padding: 0 10px")
                | {{ t.replys }}
          el-col.tweet-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            .like(@click="like(index, t.id)")
              transition(name="bounce", mode="out-in")
                i.el-icon-star-on(
                  v-if="t.already_liked",
                  style="font-size: 1.1rem; color: rgb(249, 24, 128)"
                )
                i.el-icon-star-off(v-else, style="font-size: 1.1rem")
              transition(name="bounce", mode="out-in")
                span(
                  v-if="t.already_liked",
                  style="font-size: 0.8rem; padding: 0 10px; color: rgb(249, 24, 128)"
                )
                  | {{ t.likes }}
                span(v-else, style="font-size: 0.8rem; padding: 0 10px")
                  | {{ t.likes }}
          el-col.tweet-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            i.el-icon-upload2(
              @click="copyURL(t.user_name, t.id)",
              style="font-size: 1.1rem"
            )
          el-col.tweet-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            router-link(
              :to="direct(t.user_name, t.id)",
              style="color: #5a5e66"
            )
              i.el-icon-view(style="font-size: 1.1rem")
</template>
<script>
import replycard from "../replycard/index.vue";
import userpop from "../userpop/index.vue";
import { deleteTweet } from "../../api/tweet";
import { likeTweet, deleteLike } from "../../api/user";
import { ElMessageBox, ElMessage } from "element-plus";
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
const URLMatcher =
  /(?:(?:https?|ftp|file|http):\/\/|www\.|ftp\.)(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])/gim;
export default {
  props: {
    tweets: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: false,
    },
  },
  components: {
    userpop,
    replycard,
  },
  data() {
    return {
      me: {},
      activeTweet: {},
      isMounted: false,
      isActive: false,
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
    async delete(index, id) {
      try {
        await deleteTweet(id);
        this.tweets.splice(index, 1);
      } catch (err) {
        console.log(err);
      }
    },
    async like(index, tweet_id) {
      try {
        if (this.tweets[index].already_liked) {
          await deleteLike(this.me.user_id, tweet_id);
          this.tweets[index].already_liked = !this.tweets[index].already_liked;
          this.tweets[index].likes -= 1;
        } else {
          await likeTweet(this.me.user_id, tweet_id);
          this.tweets[index].already_liked = !this.tweets[index].already_liked;
          this.tweets[index].likes += 1;
        }
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
        user_name.toLowerCase() +
        "/status/" +
        tweet_id;
      try {
        await navigator.clipboard.writeText(url);
        ElMessage.success("Copied to clipboard");
      } catch ($e) {
        ElMessage.error("Cannot copied to clipboard");
      }
    },
    direct(name, id) {
      return "/" + name.toLowerCase() + "/status/" + id;
    },
    updateIsActive(v) {
      this.isActive = v;
    },
    updateNewTweet(v) {
      this.tweets.unshift(v);
      this.isActive = false;
    },
    activate(tweet) {
      this.activeTweet = tweet;
      this.isActive = true;
    },
    confirmDelete(index, tweet_id) {
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
          this.delete(index,tweet_id)
        })
        .catch(() => {
        });
    },
  },
  mounted() {
    this.me = this.$store.getters.user;
  },
};
</script>
<style>
.tweet-list-item {
  color: #e6e6e6;
  position: relative;
  padding-bottom: 50px;
  padding: 15px 0;
  max-height: calc(100vh - 56px);
  scrollbar-width: none;
  -ms-overflow-style: none;
  overflow-x: hidden;
  overflow-y: scroll;
  border-right: solid 1px;
  border-bottom: solid 1px;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.tweet-list-item:hover {
  background-color: rgba(0, 0, 0, 0.02);
}
.tweet-list-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.tweet-list-icon:hover {
  transform: scale(1.2);
}

.tweet-list-delete {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.tweet-list-delete:hover {
  color: red;
  transform: scale(1.3);
}

.tweet-list-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
}
.tweet-list-user:hover {
  border-bottom: solid 1px;
}

.tweet-list-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.tweet-list-time {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.tweet-list-time:hover {
  border-bottom: solid 1px;
}

.tweet-list-tweet {
  color: black;
  display: inline-block;
  overflow: hidden;
}

.tweet-list-operation {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
  padding: 5px 0;
}

.tweet-list-image {
  border-radius: 15px;
  border: solid 1px;
  max-height: 300px;
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

.el-icon-view:hover {
  color: lightgreen;
  transform: scale(1.2);
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.link-to-tweet {
  background: linear-gradient(to right, transparent 50%, black 0);
  display: inline-block;
  position: absolute;
  left: -100%;
  margin: 0;
  width: 3rem;
  height: 3rem;
  border-radius: 100%;
  transition: transform 800ms;
}
</style>
