<template lang="pug">
el-row(:gutter="15")
  el-col(:xs="24", :sm="24", :md="14", :lg="14")
    .tweet-list(
      v-if="isMounted"
      v-infinite-scroll="getTweetsByScroll",
      infinite-scroll-distance="200"
      v-loading="!data.tweets[0]"
    )
      el-row.text-area
        el-col(:xs="3", :sm="3", :md="3", :lg="3")
          .text-icon
            el-image.text-avatar(:src="self.avatar", fit="cover")
        el-col(:xs="20", :sm="20", :md="20", :lg="20")
          el-input.text-block(
            v-model="content",
            type="textarea",
            :autosize="{ minRows: 1, maxRows: 6 }",
            maxlength="260",
            show-word-limit,
            placeholder="What's happening?"
          )
          .preview-image(v-if="image", style="position: relative")
            el-button(
              @click="image = ''",
              circle,
              size="small",
              style="background-color: rgba(15, 20, 25, 0.75); color: #fff; border: none; top: 4px; left: 4px; position: absolute; z-index: 1001"
            )
              i.el-icon-close(style="font-size: 15px; font-weight: bold")
            transition(name="fade" mode="out-in")
              el-image(:src="image", fit="cover", style="border-radius: 10px") 
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
                @click="sendTweet(content)",
                :disabled="content == ''",
                size="small",
                round
              )
                | Tweet
        el-col(:xs="1", :sm="1", :md="1", :lg="1")
          | &nbsp;
      tweetlist(:tweets="data.tweets" :type="'tweets'")
  el-col(:xs="24", :sm="24", :md="10", :lg="10")
    | &nbsp;
</template>
<script>
import { avatarProps } from "element-plus";
import rwd from "../../components/rwd/index.vue";
import tweetlist from "../../components/tweetlist/index.vue";
import { getTweets, getTweetByID, createTweet } from "../../api/tweet";
import { getUserByID } from "../../api/user";
export default {
  components: {
    rwd,
    tweetlist,
  },
  data() {
    return {
      failAvatar: "el-icon-user-solid",
      data: {
        tweets: [],
      },
      users: {},
      self: {},
      content: "",
      image: "",
      page: 0,
      isMounted: false
    };
  },
  methods: {
    preview(event) {
      this.image = URL.createObjectURL(event.target.files[0]);
    },
    async sendTweet(text) {
      try {
        let form = new FormData()
        form.append('text', text)
        if(this.$refs.image.files[0]){
          form.append('image', this.$refs.image.files[0]) 
        }
        let res = await createTweet(form);
        this.content = "";
        this.image = "";
        let newTweet = res.data;
        newTweet.user_name = this.self.user_name;
        newTweet.name = this.self.name;
        newTweet.avatar = this.self.avatar;
        newTweet.about = this.self.about;
        newTweet.user_id = this.self.user_id;
        this.data.tweets.unshift(newTweet);
        this.$refs.image.value = null;
      } catch (err) {
        console.log(err);
      }
    },
    async getTweetsByScroll() {
      await this.getAllTweets()
    },
    async getAllTweets() {
      if (this.page >= this.data.total_pages) {
        return;
      } else {
        this.page += 1;
        try {
          let res = await getTweets(this.page.toString());
          this.data.total_pages = res.data.total_pages;
          let tweets = res.data.tweets;
          for (let t of tweets) {
            this.data.tweets.push(t);
          }
        } catch (err) {
          console.log(err);
        }
      }
    }
  },
  async mounted() {
    this.self = this.$store.getters.user;
    await this.getAllTweets()
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