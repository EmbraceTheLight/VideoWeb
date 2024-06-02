<template>
  <div>
    <div class="box">
      <el-input placeholder="请输入账号或者用户名" v-model="account">
        <template slot="prepend">账号/用户名</template>
      </el-input>
    </div>
    <div class="box">
      <el-input
        placeholder="请输入密码"
        :show-password="true"
        v-model="password"
      >
        <template slot="prepend">密码</template>
      </el-input>
    </div>
    <div class="button-box">
      <el-button round @click="handleReg">注册</el-button>
      <el-button
        type="primary"
        round
        :style="{ 'margin-left': '40px' }"
        @click="handleLogin"
        >登录</el-button
      >
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      account: '',
      password: '',
    }
  },
  methods: {
    handleReg() {
      this.$emit('toggle')
    },
    handleLogin() {
      // console.log('test')
      this.$emit('loginEvent', {
        account: this.account,
        password: this.password,
      })
    },
  },
}
</script>

<style scoped>
.box {
  width: 500px;
  position: relative;
  transform: translate(15%, 15%);
  height: 60px;
}

.button-box {
  position: relative;
  margin-left: 160px;
  margin-top: 80px;
}

.el-button {
  width: 150px;
  height: 60px;
}
</style>
