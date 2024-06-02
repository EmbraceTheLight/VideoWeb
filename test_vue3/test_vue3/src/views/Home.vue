<template>
  <div>
    <div class="header">
      <home_navbar @loginEvent="changeVisible"></home_navbar>
    </div>
    <div class="content">
      <div v-if="isVisible" class="overlay"></div>
      <home_channel></home_channel>
      <recPlate></recPlate>
    </div>
    <loginWidget v-if="isVisible" @closeEvent="changeVisible"></loginWidget>
  </div>
</template>

<script>
import home_channel from '@/components/homeChannel.vue'
import home_navbar from '@/components/homeNavbar.vue'
import recPlate from '@/components/recPlate.vue'
import loginWidget from '@/components/loginWidget.vue'
export default {
  data() {
    return {
      isVisible: false
    }
  },
  methods: {
    changeVisible() {
      this.isVisible = !this.isVisible
    },
  },
  components: {
    home_channel,
    home_navbar,
    recPlate,
    loginWidget,
  },
  mounted() {
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* 黑色透明背景色 */
  z-index: 999; /* 确保覆盖其他内容 */
}
.header {
  height: 200px;
}

.content {
  margin-left: 140px;
}

.el-link {
  display: flex;
  align-items: center;
  white-space: nowrap;
  font-size: 13px;
  text-overflow: ellipsis;
}
</style>
