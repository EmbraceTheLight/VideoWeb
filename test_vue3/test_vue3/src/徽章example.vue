<template>
  <div >

    <div>
    <el-badge is-dot :offset="[-10,5]">      
    <el-button  :icon="Message" text class="icons"></el-button>
     <el-button   text class="text">消息</el-button>
      </el-badge>
    </div>
    
  </div>





</template>

<script lang="ts" setup>
import {
  Check,
  Delete,
  Edit,
  Message,
  Position,
  Search,
  Star,
} from '@element-plus/icons-vue'
</script>

<style>

  .text{
    position: absolute;
    left:-15px;
    top:20px;
    font-size: 15px;
  }

  .el-button .el-icon{
  font-size: 20px;
  }

  .el-button:hover{
  background-color: transparent !important;
  color: none !important;
  box-shadow: none !important;
  }
</style>