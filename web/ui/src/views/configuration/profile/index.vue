<template lang="pug">
el-row(:gutter="30")
  el-col(:xs="24", :sm="24", :md="1", :lg="1")
    | &nbsp;
  el-col(:xs="24", :sm="24", :md="22", :lg="22")
    el-row
      el-col(:xs="12", :sm="12", :md="12", :lg="12")
        .title
          | Profile
      el-col(:xs="12", :sm="12", :md="12", :lg="12")
        .switch
          el-switch(
            v-model="isEditting",
            active-text="Edit",
            inactive-text="View"
          )
    el-row.profile-container
      el-col(:xs="24", :sm="24", :md="24", :lg="24")
        el-form(
          :model="profile",
          ref="profile",
          :rules="profileRules",
          status-icon,
          label-width="180px",
          :label-position="'top'",
          :disabled="!isEditting"
        )
          el-row(:gutter="30")
            el-col(:xs="24", :sm="24", :md="12", :lg="12")
              el-form-item(label="User name", prop="user_name")
                el-input(
                  v-model="profile.user_name",
                  placeholder="User name",
                  suffix-icon="el-icon-avatar"
                )
              el-form-item(label="Email", prop="email")
                el-input(
                  v-model="profile.email",
                  placeholder="Email",
                  suffix-icon="el-icon-message"
                )
            el-col(:xs="24", :sm="24", :md="12", :lg="12")
              el-form-item(label="Name", prop="name")
                el-input(
                  v-model="profile.name",
                  placeholder="Name",
                  suffix-icon="el-icon-avatar"
                )
              el-form-item(label="Phone number", prop="phone_number")
                el-input(
                  v-model="profile.phone_number",
                  placeholder="Phone number",
                  suffix-icon="el-icon-phone"
                )
          el-row(:gutter="15")
            el-col(:xs="24", :sm="24", :md="12", :lg="12")
              el-form-item(label="Birthday", prop="birthday")
                el-date-picker(
                  v-model="profile.birthday",
                  type="date",
                  placeholder="Select your birthday",
                  format="YYYY/MM/DD",
                  value-format="YYYY-MM-DDTHH:mm:ssZ"
                )
            el-col(:xs="24", :sm="24", :md="12", :lg="12")
              el-form-item(label="Gender")
                el-radio-group(v-model="profile.gender")
                  el-radio(label="male")
                  el-radio(label="female")
                  el-radio(label="other")
          el-form-item(label="About", prop="about")
            el-input(
              v-model="profile.about",
              placeholder="About",
              maxlength="160",
              show-word-limit,
              type="textarea"
            )
        el-divider
        .button-container
          el-button(
            round,
            type="primary",
            style="width: 100%",
            @click="updateProfile",
            :disabled="!isEditting"
          )
            | Update Profile
  el-col(:xs="24", :sm="24", :md="1", :lg="1")
    | &nbsp;
</template>
<script>
import rwd from "../../../components/rwd/index.vue";
import { updateUser } from "../../../api/user";
import { ElMessage } from "element-plus";
export default {
  components: {
    rwd,
  },
  data() {
    let validateName = (rule, value, callback) => {
      if (value.match(/\d/g)) {
        callback(new Error("Must not contain numbers"));
      } else {
        callback();
      }
    };
    let validatePhone = (rule, value, callback) => {
      if (!value) {
        callback();
      } else {
        if (!value.match(/\d/g)) {
          callback(new Error("Phone number must be a number"));
        } else {
          if (value.length !== 10) {
            callback(new Error("Phone number must be 10 digits "));
          } else {
            callback();
          }
        }
      }
    };
    return {
      isEditting: false,
      user_id: "",
      csrf_token: "",
      headers: {
        "X-CSRF-Token": "",
      },
      profile: {
        user_name: "",
        name: "",
        gender: "male",
        email: "",
        about: "",
        phone_number: "",
        country: "",
        birthday: "",
        avatar: "",
      },
      profileRules: {
        user_name: [
          {
            required: true,
            message: "Please type your user name",
            trigger: "blur",
          }
        ],
        name: [
          {
            required: true,
            message: "Please type your name",
            trigger: "blur",
          },
          { validator: validateName, trigger: "blur" },
        ],
        gender: [
          {
            required: true,
            message: "Please choose your gender",
            trigger: "blur",
          },
        ],
        email: [
          {
            required: true,
            message: "Please type your email address",
            trigger: "blur",
          },
          { type: "email", message: "Type must be email address" },
        ],
        phone_number: [{ validator: validatePhone, trigger: "blur" }],
      },
    };
  },
  methods: {
    async updateProfile() {
      let res = await updateUser(this.user_id, this.profile);
      ElMessage({
        showClose: true,
        message: "Update success",
        type: "success",
      });
    }
  },
  async mounted() {
    this.profile = this.$store.getters.user;
    this.user_id = this.profile.user_id;
    this.csrf_token = this.$store.getters.csrf_token;
    this.headers["X-CSRF-Token"] = this.csrf_token;
  },
};
</script>
<style>
.title {
  line-height: 1.25;
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 32px 0;
}
.switch {
  line-height: 1.25;
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 32px 0;
  right: 0;
  position: absolute;
}
.profile-container {
  position: relative;
  padding: 20px;
  margin-bottom: 30px;
  background-color: #f0f6f9;
  box-shadow: 0px -1px 3px rgb(0 0 0 / 8%), 0px 3px 6px rgb(0 0 0 / 12%);
  border-radius: 4px;
}
.button-container {
  margin: 5px 0;
}
.instruction-button {
  display: inline-block;
  padding-left: 10px;
  line-height: 19px;
  font-size: 14px;
}
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 20px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}
.avatar-uploader-icon {
  font-size: 40px;
  color: #8c939d;
  width: 100%;
  height: 100%;
  line-height: 100%;
  text-align: center;
}
.avatar-uploader-image {
  width: 200px;
  height: 200px;
  display: block;
}
</style>