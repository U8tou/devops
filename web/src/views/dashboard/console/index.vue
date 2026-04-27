<!-- 工作台页面 -->
<template>
  <div>
    <CardList></CardList>
  </div>
</template>

<script setup lang="ts">
  import CardList from './modules/card-list.vue'

  defineOptions({ name: 'Console' })
</script>
