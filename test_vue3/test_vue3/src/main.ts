import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router' // 假设这是你的路由文件路径
import * as Icons from "@element-plus/icons-vue";
import SvgIcon from './components/iconfont/SvgIcon.vue'
import './assets/iconfont/iconfont.js'
import Axios from 'axios'
import store from './store'
const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(Axios)
app.use(store)
// 注册全局组件
Object.keys(Icons).forEach((key) => {
  app.component(key, Icons[key as keyof typeof Icons]);
});
app.component('SvgIcon',SvgIcon)
app.mount('#app')
