<template>
  <SpotifyAuth />
  <Chat placeholder="How can I help you?" />
</template>

<script>
import Chat from './components/Chat.vue'
import SpotifyAuth from './components/SpotifyAuth.vue'

export default {
  name: 'App',
  components: {
    Chat,
    SpotifyAuth
  }
}
</script>
