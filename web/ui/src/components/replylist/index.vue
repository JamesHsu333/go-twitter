<template lang="pug">
div
  el-empty(v-if="replys.length===0" :description="nocontent")
  .reply-list-item(v-for="(r, index) in replys", :key="r")
    transition(name="bounce")
      div
        replycard(
          :tweet="activeTweet",
          :isActive="isActive",
          @close="updateIsActive",
          @update="updateNewTweet",
          :key="activeTweet"
        )
        el-row
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
            userpop(:user="r")
          el-col(:xs="18", :sm="18", :md="18", :lg="18")
            router-link.reply-list-user(:to="'/' + r.user_name")
              | {{ r.name }}
            span.reply-list-account
              | @{{ r.user_name }}
            span.reply-list-account
              | ·
            .reply-list-time
              el-tooltip(
                effect="dark",
                :content="fullDate(r.created_at)",
                placement="top"
              )
                span
                  | {{ formatDate(r.created_at) }}
            br
            span.reply-list-instruction
              | Replying to
            router-link.reply-list-user-link(:to="'/' + tweet.user_name")
              | @{{ tweet.user_name }}
            br
            .reply-list-intro(v-html="formatText(r.text)")
            br
            el-image.reply-list-image(
              v-if="r.image",
              :src="'/' + r.image",
              :preview-src-list="['/' + r.image]",
              :initial-index="1",
              fit="cover"
            )
          el-col(:xs="3", :sm="3", :md="3", :lg="3")
            el-tooltip.reply-list-delete(
              v-if="me.user_id === r.user_id",
              effect="dark",
              content="Delete tweet",
              placement="top",
            )
              span(@click="confirmDelete(index,r.id)")
                i.el-icon-close(style="font-size: 1.1rem")
          el-col.reply-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            .reply(@click="activate(r)")
              i.el-icon-chat-round(style="font-size: 1rem")
              span(style="font-size: 1rem; padding: 0 10px")
                | {{ r.replys }}
          el-col.reply-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            .like(@click="like(index, r.id)")
              transition(name="bounce", mode="out-in")
                i.el-icon-star-on(
                  v-if="r.already_liked",
                  style="font-size: 1rem; color: rgb(249, 24, 128)"
                )
                i.el-icon-star-off(v-else, style="font-size: 1rem")
              transition(name="bounce", mode="out-in")
                span(
                  v-if="r.already_liked",
                  style="padding: 0 10px; color: rgb(249, 24, 128)"
                )
                  | {{ r.likes }}
                span(v-else, style="padding: 0 10px")
                  | {{ r.likes }}
          el-col.reply-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            i.el-icon-upload2(
              @click="copyURL(r.user_name, r.id)",
              style="font-size: 1rem"
            )
          el-col.reply-list-operation(:xs="6", :sm="6", :md="6", :lg="6")
            router-link(:to="direct(r.user_name, r.id)" style="color: #5a5e66")
              i.el-icon-view(style="font-size: 1rem")
</template>
<script>
import replycard from "../replycard/index.vue";
import userpop from "../userpop/index.vue";
import { deleteTweet } from "../../api/tweet";
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
  props: {
    replys: {
      type: Object,
      required: true,
    },
    tweet: {
      type: Object,
      required: true,
    },
    nocontent: {
      type: String,
      required: false,
    },
  },
  components: {
    replycard,
    userpop,
  },
  data() {
    return {
      me: {},
      isMounted: false,
      isActive: false,
      activeTweet: {},
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
    async delete(index, id) {
      try {
        await deleteTweet(id);
        this.replys.splice(index, 1);
      } catch (err) {
        console.log(err);
      }
    },
    async like(index, tweet_id) {
      try {
        if (this.replys[index].already_liked) {
          await deleteLike(this.me.user_id, tweet_id);
          this.replys[index].already_liked = !this.replys[index].already_liked;
          this.replys[index].likes -= 1;
        } else {
          await likeTweet(this.me.user_id, tweet_id);
          this.replys[index].already_liked = !this.replys[index].already_liked;
          this.replys[index].likes += 1;
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
      return "/"+name.toLowerCase()+'/status/'+id
    },
    updateIsActive(v) {
      this.isActive = v;
    },
    updateNewTweet(v) {
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
.reply-list-item {
  color: #e6e6e6;
  position: relative;
  padding-bottom: 50px;
  padding: 15px 0 0 0;
  max-height: calc(100vh - 56px);
  scrollbar-width: none;
  -ms-overflow-style: none;
  overflow-x: hidden;
  overflow-y: scroll;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.reply-list-item:hover {
  background-color: rgba(0, 0, 0, 0.02);
}

.reply-list-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.reply-list-icon:hover {
  transform: scale(1.2);
}

.reply-list-delete {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.reply-list-delete:hover {
  color: red;
  transform: scale(1.3);
}

.reply-list-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
}
.reply-list-user:hover {
  border-bottom: solid 1px;
}

.reply-list-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.reply-list-time {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.reply-list-time:hover {
  border-bottom: solid 1px;
}

.reply-list-instruction {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-weight: 600;
  font-size: 0.9rem;
}

.reply-list-user-link {
  color: #1d9bf0;
  display: inline-block;
  overflow: hidden;
  font-weight: 600;
  font-size: 0.9rem;
}
.reply-list-user-link:hover {
  border-bottom: solid 1px;
}

.reply-list-intro {
  color: black;
  display: inline-block;
  overflow: hidden;
}

.reply-list-operation {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
  padding: 15px 5px;
  border-bottom: solid 1px;
  border-bottom-color: #e6e6e6;
}

.reply-list-image {
  border-radius: 15px;
  max-height: 300px;
}

.reply-list-operation i {
  font-size: 1.5rem;
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
</style>