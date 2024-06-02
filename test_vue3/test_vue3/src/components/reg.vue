<template>
  <div>
    <div class="box">
      <el-input placeholder="请输入用户名" v-model="username">
        <template slot="prepend">用户名</template>
      </el-input>
    </div>
    <div class="box">
      <el-input placeholder="请输入密码" :show-password="true" v-model="password">
        <template slot="prepend">密码</template>
      </el-input>
    </div>
    <div class="box">
      <el-input placeholder="请再次确认密码" :show-password="true" v-model="repassword">
        <template slot="prepend" >确认密码</template>
      </el-input>
    </div>

    <div class="box">
      <el-input placeholder="请输入邮箱" v-model="Email">
        <template slot="prepend">邮箱</template>
      </el-input>
    </div>

    <div class="box">
      <el-input placeholder="请输入验证码" v-model="code">
        <template #append>
          <el-button type="primary" @click="sendCode"
            :style="{ 'background-color': 'rgb(0, 174, 236)', color: 'white' }">获取验证码</el-button>
        </template>
      </el-input>
    </div>
    <div class="box">
      <el-input placeholder="请输入个性签名" v-model="Signature">
        <template slot="prepend">个性签名</template>
      </el-input>
    </div>
    <div class="sbox">
      <el-input placeholder="输入图形验证码" v-model="pCode" :style="{ width: '300px' }">
      </el-input>
       <el-image :style="{ width: '150px', height: '75px', marginLeft:'50px' }" :src="picturebase64" @click.native="refreshCode"></el-image>
    </div>
    <div class="reg">
      <el-button class="regButton" round @click="handleReg">注册</el-button>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import axios from 'axios'
import {ElMessage} from "element-plus";

export default {
  props: ['answer', 'b64lob'],
  setup(props, { emit }) {
    const picturebase64 = ref(props.b64lob);

    watch(() => props.b64lob, (newValue) => {
      picturebase64.value = newValue;
    });

    const password = ref('');
    const username = ref('');
    const repassword = ref('');
    const Email = ref('');
    const code = ref('');
    const Signature = ref('');
    const pCode = ref('');

    const sendCode = () => {
      emit('sendCode', { Email: Email.value });
      console.log('发送验证码');
    };

    const handleReg = () => {
      emit('regEvent', {
        username: username.value,
        password: password.value,
        repassword: repassword.value,
        email: Email.value,
        code: code.value,
        signature: Signature.value,
        pCode: pCode.value,
      });
    };

    const refreshCode = () => {
      axios.get('/yanxi/Captcha/GenerateGraphicCaptcha').then(res => {
        const captcha = res.data.captcha;
        picturebase64.value = captcha.b64lob;
        console.log(captcha);
      });
    };

    return {
      picturebase64,
      password,
      username,
      repassword,
      Email,
      code,
      Signature,
      sendCode,
      handleReg,
      refreshCode,
      pCode
    };
  }
}
</script>

<style scoped>
.box {
  width: 500px;
  position: relative;
  margin-left: 100px;
  height: 40px;
}
.sbox {
  width: 500px;
  position: relative;
  margin-left: 100px;
  display: flex;
  align-items: center;
  margin-top: 20px;
  height: 40px;
}
.sgbutton {
  height: 30px;
  position: relative;
  transform: translate(0%, 20%);
}
.sgbutton:hover {
  background-color: initial !important;
}
.el-divider {
  height: 30px;
  position: relative;
  bottom: 16px;
}

.text {
  display: inline-block;
  position: relative;
  bottom: 13px;
  left: 5px;
}

.reg {
  position: relative;
  margin-left: 260px;
  margin-top:20px;
}

.regButton {
  position: relative;
  margin-top: 20px;
  width: 150px;
  height: 40px;
}
</style>
