<template lang="pug">
el-row(:gutter="15")
  el-col(:xs="24", :sm="24", :md="14", :lg="14")
    .follow-list(v-if="isMounted" v-infinite-scroll="getUsersByScroll" infinite-scroll-distance="200" v-loading="!isMounted")
      el-tabs.tab(v-model="activeTab" stretch="true" @tab-click="tabClick")
        el-tab-pane(label="Followers" name="followers")
          userlist(v-if="activeTab==='followers'" :users="followers.users" :nocontent="'@'+user.user_name+ ' don’t have any followers yet'")
        el-tab-pane(label="Following" name="following")
          userlist(v-if="activeTab==='following'" :users="following.users" :nocontent="'@'+user.user_name+ ' don’t have any following yet'")
  el-col(:xs="24", :sm="24", :md="10", :lg="10")
    | &nbsp;
</template>
<script>
import { avatarProps } from "element-plus";
import rwd from "../../components/rwd/index.vue";
import userlist from "../../components/userlist/index.vue";
import { getUserByName, getFollowers, getFollowing} from "../../api/user";
export default {
  components: {
    rwd,
    userlist,
  },
  data() {
    return {
      activeTab: "",
      user: {},
      followers: {
        users: []
      },
      following: {
        users: []
      },
      page: [0,0],
      isMounted: false,
      isExist: false,
    };
  },
  methods: {
    async getUsersByScroll() {
      if (!this.isExist) {
        return
      }
      let user_id = this.user.user_id
      if (this.activeTab==="followers") {
        this.getUserFollowers(user_id)
      }else{
        this.getUserFollowing(user_id)
      }
    },
    async getUserFollowers(user_id) {
      if (this.page[0] >= this.followers.total_pages) {
        return;
      } else {
        this.page[0] += 1;
        try {
          let res = await getFollowers(user_id, this.page[0].toString());
          this.followers.total_pages = res.data.total_pages;
          let users = res.data.users;
          for (let u of users) {
            this.followers.users.push(u);
          }
        } catch (err) {
          console.log(err);
        }
      }
    },
    async getUserFollowing(user_id) {
     if (this.page[1] >= this.following.total_pages) {
        return;
      } else {
        this.page[1] += 1;
        try {
          let res = await getFollowing(user_id, this.page[1].toString());
          this.following.total_pages = res.data.total_pages;
          let users = res.data.users;
          for (let u of users) {
            this.following.users.push(u);
          }
        } catch (err) {
          console.log(err);
        }
      } 
    },
    tabClick(tab, event) {
      if (tab.props.name==="following") {
        this.$router.replace("/"+this.$route.params.user+"/following")
        this.getUserFollowing(this.user.user_id)
      }else{
        this.$router.replace("/"+this.$route.params.user+"/followers")
        this.getUserFollowers(this.user.user_id)
      }
    }
  },
  async mounted() {
    let route = this.$route
    try {
      let res = await getUserByName(route.params.user);
      this.user = res.data;
      this.isExist = true
    }catch(err) {
      console.log(err)
    }
    this.activeTab = route.path.includes("followers") ? "followers": "following"
    await this.getUserFollowers(this.user.user_id)
    this.isMounted = true
  },
};
</script>
<style>
.follow-list {
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