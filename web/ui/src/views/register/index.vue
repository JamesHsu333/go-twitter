<template lang="pug">
.register-main
  .register-container
    .title
      | Sign up
    el-form(
      :model="registerForm",
      ref="registerForm",
      :rules="registerFormRules",
      status-icon,
      label-width="180px",
      :label-position="'top'"
    )
      el-row(:gutter="15")
        el-col(:xs="24", :sm="24", :md="12", :lg="12")
          el-form-item(label="User name", prop="user_name")
            el-input(
              v-model="registerForm.user_name",
              placeholder="User name",
              suffix-icon="el-icon-avatar"
            )
        el-col(:xs="24", :sm="24", :md="12", :lg="12")
          el-form-item(label="Name", prop="name")
            el-input(
              v-model="registerForm.name",
              placeholder="Name",
              suffix-icon="el-icon-avatar"
            )
      el-form-item(label="Email", prop="email")
        el-input(
          v-model="registerForm.email",
          placeholder="Email",
          suffix-icon="el-icon-message"
        )
      el-form-item(label="Password", prop="password")
        el-input(
          v-model="registerForm.password",
          placeholder="Password",
          show-password
        )
      el-form-item(label="Confirm Password", prop="repassword")
        el-input(
          v-model="registerForm.repassword",
          placeholder="Confirm password",
          show-password
        )
    el-divider
    .button-container
      el-button(
        round,
        style="width: 100%",
        @click="submitRegisterForm('registerForm')"
      )
        b
          | Sign up
    el-row(:gutter="15")
      el-col.button-container(:xs="24", :sm="24", :md="18", :lg="18")
        .instruction-button
          | Already have a account?
      el-col.button-container(:xs="24", :sm="24", :md="6", :lg="6")
        router-link(to="/login")
          .instruction-button(style="float: right")
            | Log in
</template>
<script>
import rwd from "../../components/rwd/index.vue";
export default {
  components: {
    rwd,
  },
  data() {
    let validatePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error("Please type your password"));
      } else {
        if (this.registerForm.checkPass !== "") {
          this.$refs.registerForm.validateField("repassword");
        }
        callback();
      }
    };
    let validateRePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error("Please retype your password to confirm"));
      } else {
        if (value !== this.registerForm.password) {
          callback(new Error("Different Password!!"));
        }
        callback();
      }
    };
    return {
      registerForm: {
        user_name: "",
        name: "",
        email: "",
        password: "",
        repassword: "",
      },
      registerFormRules: {
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
          }
        ],
        email: [
          {
            required: true,
            message: "Please type your email address",
            trigger: "blur",
          },
          { type: "email", message: "Type must be email address" },
        ],
        password: [
          {
            required: true,
            message: "Please type your password",
            trigger: "blur",
          },
          { validator: validatePassword, trigger: "blur" },
        ],
        repassword: [
          {
            required: true,
            message: "Please retype your password to confirm",
            trigger: "blur",
          },
          { validator: validateRePassword, trigger: "blur" },
        ],
      },
    };
  },
  methods: {
    submitRegisterForm(form) {
      this.$refs[form].validate((valid) => {
        if (valid) {
          console.log(this.$refs[form].model);
          this.register(this.$refs[form].model);
        } else {
          return false;
        }
      });
    },
    async register(user) {
      await this.$store.dispatch("register", user);
      this.$router.push("/");
    },
  },
};
</script>
<style>
.register-main {
  display: flex;
  align-items: center;
  justify-content: center;
}
.title {
  line-height: 1.5;
  font-size: 28px;
  font-weight: 700;
  color: #333;
}
.register-container {
  padding: 20px 30px;
  background-color: #fff;
  box-shadow: 0px -1px 3px rgb(0 0 0 / 8%), 0px 3px 6px rgb(0 0 0 / 12%);
  border-radius: 20px;
  width: 40%;
  margin: 50px 0;
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
</style>