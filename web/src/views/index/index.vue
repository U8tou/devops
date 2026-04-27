<!-- 布局容器 -->
<template>
  <div class="app-layout">
    <aside id="app-sidebar">
      <ArtSidebarMenu />
    </aside>

    <main id="app-main">
      <div id="app-header">
        <ArtHeaderBar />
      </div>
      <div id="app-content">
        <ArtPageContent />
      </div>
    </main>

    <div id="app-global">
      <ArtGlobalComponent />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import { connectSse, disconnectSse } from '@/utils/sse'
  import { watch } from 'vue'

  defineOptions({ name: 'AppLayout' })

  const userStore = useUserStore()

  // 登录后注册 SSE，监听 accessToken 变化以建立连接（有 token 为登录用户定向推送）
  watch(
    () => userStore.accessToken,
    (token) => {
      if (token) {
        connectSse(token)
      } else {
        disconnectSse()
      }
    },
    { immediate: true }
  )

  onBeforeUnmount(() => {
    disconnectSse()
  })
</script>

<style lang="scss" scoped>
  @use './style';
</style>
