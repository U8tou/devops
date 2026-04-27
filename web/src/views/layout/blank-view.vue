<!-- 空白布局：仅用于嵌套路由的 RouterView 占位，不包含 sidebar/header -->
<template>
  <RouterView />
</template>

<script setup lang="ts">
  defineOptions({ name: 'LayoutBlankView' })
</script>
