<template lang="pug">
.login-main
  .login-container
    .title
      | Log in
      el-form(
        :model="loginForm",
        ref="loginForm",
        :rules="loginFormRules",
        status-icon,
        label-width="180px",
        :label-position="'top'"
      )
        el-form-item(label="Email", prop="email")
          el-input(
            v-model="loginForm.email",
            placeholder="Email",
            suffix-icon="el-icon-message"
          )
        el-form-item(label="Password", prop="password")
          el-input(
            v-model="loginForm.password",
            placeholder="Password",
            show-password
          )
      el-row
        el-col(:xs="24", :sm="24", :md="12", :lg="12")
          el-checkbox(v-model="checked")
            |
            | Remember me
        el-col(:xs="24", :sm="24", :md="12", :lg="12")
          router-link(to="/home")
            .instruction-button(style="float: right; padding: 10px 0")
              |
              | Forgot your password?
      el-divider
      .button-container
        el-button(
          round,
          style="width: 100%",
          @click="submitLoginForm('loginForm')"
        )
          b
            | Log in
      el-row(:gutter="15")
        el-col.button-container(:xs="24", :sm="24", :md="18", :lg="18")
          .instruction-button
            | No account yet? Sign up now!
        el-col.button-container(:xs="24", :sm="24", :md="6", :lg="6")
          router-link(to="/register")
            .instruction-button(style="float: right")
              | Sign up
</template>
<script>
import rwd from "../../components/rwd/index.vue";
export default {
  components: {
    rwd,
  },
  data() {
    return {
      checked: true,
      loginForm: {
        email: "",
        password: "",
      },
      loginFormRules: {
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
        ],
      },
    };
  },
  methods: {
    submitLoginForm(form) {
      this.$refs[form].validate((valid) => {
        if (valid) {
          this.login(this.$refs[form].model);
        } else {
          return false;
        }
      });
    },
    async login(form) {
      await this.$store.dispatch("login", form);
      this.$router.push("/");
    },
  },
};
</script>
<style>
.login-main {
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
.login-container {
  padding: 20px 30px;
  background-color: #fff;
  box-shadow: 0px -1px 3px rgb(0 0 0 / 8%), 0px 3px 6px rgb(0 0 0 / 12%);
  border-radius: 20px;
  width: 40%;
  margin-top: 50px;
}
.button-container {
  margin: 5px 0;
}
.instruction-button {
  display: inline-block;
  padding-left: 10px;
  font-size: 14px;
}
</style>