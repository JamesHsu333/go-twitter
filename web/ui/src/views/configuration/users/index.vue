<template lang="pug">
div
  el-table(:data="users", v-loading="loading")
    el-table-column(label="Name", prop="name")
    el-table-column(label="User name", prop="user_name")
    el-table-column(label="Email", prop="email")
    el-table-column(v-if="hasPermission" label="Role", prop="role")
      template(#default="scope")
        el-select(
          v-model="scope.row.role",
          size="small",
          @change="updateUserRole(scope.row)"
        )
          el-option(
            v-for="r in roles",
            :key="r.value",
            :label="r.label",
            :value="r.value"
          )
    el-table-column(v-if="hasPermission" label="Operation")
      template(#default="scope")
        el-button(size="mini", type="primary", @click="")
          | Edit
        el-button(size="mini", type="danger", @click="deleteUser(scope.row)")
          | Delete
  el-pagination(
    v-model:currentPage="currentPage",
    layout="prev, pager, next, jumper",
    :total="totalCount",
    :page-size="size",
    @current-change="getAllUser(currentPage)"
  )
</template>
<script>
import { getAllUser, updateUserRole, deleteUser } from "../../../api/user";
import { getPermission } from "../../../utils/utils"
import { ElMessage, ElMessageBox } from "element-plus";
export default {
  data() {
    return {
      users: [],
      roles: [
        { value: "user", label: "user" },
        { value: "admin", label: "admin" },
      ],
      currentPage: 0,
      totalPages: 0,
      totalCount: 0,
      size: 0,
      loading: true,
    };
  },
  computed: {
    hasPermission() {
      return getPermission()
    }
  },
  methods: {
    async updateUserRole(user) {
      let res = await updateUserRole(user.user_id, user.role);
      if(this.$store.getters.user.user_id === res.data.user_id){
        this.$store.dispatch("updateUserInfo", res.data)
      }
      ElMessage({
        showClose: true,
        message: "Update success",
        type: "success",
      });
    },
    async getAllUser(page) {
      let res = await getAllUser(page);
      this.currentPage = res.data.page;
      this.totalPages = res.data.total_pages;
      this.totalCount = res.data.total_count;
      this.size = res.data.size;
      this.users = res.data.users;
      this.loading = false;
    },
    deleteUser(user) {
      ElMessageBox.confirm(
        "Server will permanently delete the user. Continue?",
        "Warning",
        {
          confirmButtonText: "Delete",
          cancelButtonText: "Cancel",
          type: "warning",
        }
      )
        .then(() => {
          deleteUser(user.user_id)
            .then((res) => {
              console.log(res);
              ElMessage({
                type: "success",
                message: "Delete completed",
              });
              this.getAllUser(this.currentPage)
            })
            .catch((err) => {
              console.log(err);
            });
        })
        .catch(() => {});
    },
  },
  mounted() {
    this.getAllUser('1')
  },
};
</script>