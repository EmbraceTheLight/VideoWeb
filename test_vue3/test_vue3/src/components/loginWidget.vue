<template>
  <div class="widget">
    <el-tabs v-model="activeName" :style="{ width: '630px' }">
      <el-tab-pane label="登录" name="first">
        <login @toggle="toggleReg" @loginEvent="login"></login>
      </el-tab-pane>
      <el-tab-pane label="注册" name="second">
        <reg :answer="answer" :b64lob="b64lob" @sendCode="sendCode" @regEvent="reg"></reg>
      </el-tab-pane>
    </el-tabs>

    <el-button type="text" @click="handleClose">
      <el-icon size="20px" :class="iconClasses">
        <SvgIcon iconName="icon-guanbixiao" fill="black" />
      </el-icon>
    </el-button>


  </div>
</template>

<script>
import { Close as ElIconClose } from '@element-plus/icons-vue'
import login from '@/components/login.vue'
import reg from '@/components/reg.vue'
import axios from 'axios'
import SvgIcon from './iconfont/SvgIcon.vue';
import { ElMessage } from 'element-plus'
export default {
  components: {
    login,
    reg,
    ElIconClose,
  },
  data() {
    return {
      password: '',
      account: '111',
      Email: '',
      activeName: 'first',
      answer: '',
      b64lob: '',
      dataloaded: false,
      userLevel: '',
      followers: '',
      username: '',
      signature: '',
      shells: ''
    }
  },
  methods: {
    handleClose() {
      this.$emit('closeEvent')
    }, // 关闭登录小窗口
    login(data) {
      // 登录请求
       console.log('test2');
      this.account = data.account
      this.password = data.password

      // 将数据拼接为字符串
      const formData = new URLSearchParams()
      formData.append('Username', this.account)
      formData.append('password', this.password)

      axios({
        method: 'post',
        url: '/yanxi/User/Login',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded', // 设置 Content-Type 为 application/x-www-form-urlencoded
        },
        data: formData, // 发送拼接的字符串数据
      }).then((res) => {
        console.log(res.data.code)
        const token = res.data.data.token
        console.log(token)
        if (res.data.code === 200) {
          this.$emit('closeEvent')
          this.login_success()
          this.$store.state.isLogin = true
          window.location.reload()
          localStorage.setItem('token', token)
          const parts = token.split('.');
          const header = JSON.parse(atob(parts[0]));
          const payload = JSON.parse(atob(parts[1]));
          console.log(payload.userId)
          localStorage.setItem('userId', payload.userId)

        }
      })
    },
    sendCode(data) {
      // 发送验证码
      console.log('test1');
      this.Email = data.Email
      // 将数据拼接到URL中
      const urlWithParams = `/yanxi/Captcha/Send-code?email=${encodeURIComponent(
        this.Email
      )}`

      axios({
        method: 'get',
        url: urlWithParams,
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      }).then((res) => {
        ElMessage.success('验证码发送成功，请注意查收')
        console.log(res.data.code)
      })
    },
    reg(data) {
      // 注册
      // console.log('test2');
      const username = data.username
      const password = data.password
      const repassword = data.repassword
      const email = data.email
      const code = data.code
      const Signature = data.signature
      const pCode = data.pCode
      // 将数据拼接为字符串
      const formData = new URLSearchParams()
      formData.append('userName', username)
      formData.append('password', password)
      formData.append('repeatPassword', repassword)
      formData.append('Email', email)
      formData.append('Code', code)
      formData.append('Signature', Signature)


        axios({
        method: 'post',
        url: '/yanxi/User/Register',
        data: formData, // 发送拼接的字符串数据
      }).then((res) => {
        if (res.data.code === 200) {
          ElMessage.success('注册成功');
          this.activeName = 'first'
        }
      })
    },
    login_success() {
      console.log('登录成功')
    },
  },
  created() { 
    axios.get('/yanxi/Captcha/GenerateGraphicCaptcha').then(res => {
      const captcha = res.data.captcha;
      console.log(captcha);
      this.b64lob = captcha.b64lob;
      console.log(this.b64lob);
       this.answer = captcha.answer;
      console.log(this.answer);
       this.dataloaded = true;
    })

  },
  computed: {
    iconClasses() {
      return {
        'icon-close-login': this.activeName === 'first',
        'icon-close-reg': this.activeName ==='second'
      }
    }
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}

.widget {
  /*登录小窗口样式*/
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 630px; /* 设置宽度 */
  height: 440px; /* 设置高度 */
  transform: translate(-50%, -50%); /* 居中显示 */
  position: fixed;
  top: 50%;
  left: 50%;
  z-index: 1000;
}
.el-menu-item {
  padding-left: 90px; /* 根据需要调整 */
  font-size: 16px;
  width: 50%;
  height: 60px;
  padding-right: 100px;
  user-select: none;
}

.el-menu {
  padding-left: 0px;
  text-align: center;
  margin-left: 100px;
}
i {
  position: absolute;
  right: -103px;
  bottom: 49px;
  font-size: 25px;
}

i:hover {
  cursor: pointer;
}

.icon-close-login {
  position: relative;
  z-index:999;
  transform:translate(522px,-280px)
}
.icon-close-reg{
  position: relative;
  z-index:999;
  transform: translate(522px,-400px);
}
</style>
