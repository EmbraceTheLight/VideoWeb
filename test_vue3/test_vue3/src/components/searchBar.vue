<template>
<div :style="{position:'relative'}">
    <el-autocomplete
      v-model="state"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      @select="handleSelect"
      :style="{width: '500px',height: '40px',marginTop: '10px'}"
      clearable="true"
      trigger-on-focus="false"
      placement="bottom-end"
      hide-loading="true"
      ref="autocomplete"
      @keydown.native.enter="handleEnter"
    >
     <template #suffix>
              <el-icon class="el-input__icon" >
        <search />
      </el-icon>
     </template>
         <template #default>

    </template>
    </el-autocomplete>
    <div class="search-history"></div>
</div>
</template>
<script lang="ts" setup >
import {  Search  } from '@element-plus/icons-vue'
import { onMounted, ref  } from 'vue'

const state = ref('')
interface LinkItem {
  value: string
  link: string
}

const links = ref<LinkItem[]>([])
const autocomplete = ref(null)
const loadAll = () => {
  return [
    { value: 'vue', link: 'https://github.com/vuejs/vue' },
  ]
}

let timeout: ReturnType<typeof setTimeout>
const querySearchAsync = (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
    ? links.value.filter(createFilter(queryString))
    : null

  clearTimeout(timeout)
  timeout = setTimeout(() => {
  const value: any = autocomplete.value;
  value.close();
    cb(results)
  }, 0)
}
const createFilter = (queryString: string) => {
  return (restaurant: LinkItem) => {
    return (
      restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

const handleSelect = (item: Record<string, any>) => {
  console.log(item)
}
const handleEnter = () => {
  const inputValue: any = state;
  window.location.href = `#/search/keyword=${inputValue.value}`;
  window.location.reload();
}
onMounted(() => {
  links.value = loadAll();
})
 
</script>

<style  scoped>
  *{
    margin: 0;
    padding: 0;
  }
  
  .search-history{
    display: flex;
    width:500px;
    height: 500px;
    position: absolute;
    margin-top:55px;
    z-index:1;
    pointer-events: none;
  }
</style>

