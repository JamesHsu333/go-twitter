<template lang="pug">
el-dialog(
  v-model="isActive",
  width="40%",
  :close-on-click-modal="'false'",
  :close-on-press-escape="'false'",
  :show-close="'false'",
  style="bordor-radius: 30px",
  center
)
  el-row(v-if="tweet")
    el-col(:xs="4", :sm="4", :md="4", :lg="4")
      .text-icon
        el-image.text-avatar(
          v-if="isActive",
          :src="'/' + tweet.avatar",
          fit="cover"
        )
    el-col(:xs="18", :sm="18", :md="18", :lg="18")
      .reply-card-user
        | {{ tweet.name }}
      span.reply-card-account
        | @{{ tweet.user_name }}
      span.reply-card-account
        | ·
      .reply-card-time
        el-tooltip(
          effect="dark",
          :content="fullDate(tweet.created_at)",
          placement="top"
        )
          span
            | {{ formatDate(tweet.created_at) }}
      br
      .reply-card-intro
        | {{ tweet.text }}
      .reply-card-intro(v-if="tweet.image" v-html="formatUrl(tweet.image)")
      br
      span.reply-card-instruction(style="padding-top: 25px")
        | Replying to
      .reply-card-user-link
        | @{{ tweet.user_name }}
  el-row(v-if="tweet", style="padding-top: 25px")
    el-col(:xs="4", :sm="4", :md="4", :lg="4")
      .text-icon
        el-image.text-avatar(
          v-if="isActive",
          :src="'/' + tweet.avatar",
          fit="cover"
        )
    el-col(:xs="18", :sm="18", :md="18", :lg="18")
      el-input.text-block(
        v-model="content",
        type="textarea",
        :autosize="{ minRows: 3, maxRows: 6 }",
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
          el-image(:src="image", fit="cover", style="border-radius: 10px") 
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
</template>
<script>
import {createReplyTweet} from "../../api/tweet"
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
export default {
  props: {
    tweet: {
      type: Object,
      required: true,
    },
    isActive: {
      type: Boolean,
      required: true,
    }
  },
  data() {
    return {
      content: "",
      image: "",
      me: {},
    };
  },
  methods: {
    close() {
      this.$emit("close", this.isActive);
    },
    update() {
      this.$emit("update", this.newTweet)
    },
    preview(event) {
      this.image = URL.createObjectURL(event.target.files[0]);
    },
    async replyTweet(id, text) {
      try {
        let form = new FormData()
        form.append('text', text)
        if(this.$refs.image.files[0]){
          form.append('image', this.$refs.image.files[0]) 
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
        this.newTweet = {...newTweet}
        this.update()
        this.$refs.image.value = null;
        this.$emit("close", false)
      } catch (err) {
        console.log(err);
      }
    },
    formatUrl(v) {
        return "<a href=" + window.location.protocol + "//" + window.location.host + "/" + v + ">" + window.location.host + "/" + v + "</a>";
    },
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
  },
  mounted() {
    this.me = this.$store.getters.user;
  },
};
</script>
<style>
.el-dialog {
  border-radius: 40px;
}
.reply-card-item {
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
.reply-card-item:hover {
  background-color: rgba(0, 0, 0, 0.02);
}

.reply-card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}

.reply-card-icon:hover {
  transform: scale(1.2);
}

.reply-card-dropdown {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
}
.reply-card-dropdown:hover {
  color: #1ea2f1;
  transform: scale(1.3);
}

.reply-card-user {
  color: rgba(0, 0, 0, 1);
  display: inline-block;
  font-weight: bold;
  overflow: hidden;
}
.reply-card-user:hover {
  border-bottom: solid 1px;
}

.reply-card-account {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.reply-card-time {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-size: 0.9rem;
  padding: 0 5px;
}

.reply-card-time:hover {
  border-bottom: solid 1px;
}

.reply-card-instruction {
  color: #536471;
  display: inline-block;
  overflow: hidden;
  font-weight: 600;
  font-size: 0.9rem;
}

.reply-card-user-link {
  color: #1d9bf0;
  display: inline-block;
  overflow: hidden;
  font-weight: 600;
  font-size: 0.9rem;
}
.reply-card-user-link:hover {
  border-bottom: solid 1px;
}

.reply-card-intro {
  color: black;
  display: inline-block;
  overflow: hidden;
}

.reply-card-operation {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #5a5e66;
  transition: transform 0.15s cubic-bezier(0.11, 0.65, 1, 1.69);
  padding: 15px 5px;
  border-bottom: solid 1px;
  border-bottom-color: #e6e6e6;
}

.reply-card-image {
  border-radius: 15px;
  max-height: 300px;
}

.reply-card-operation i {
  font-size: 1.5rem;
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