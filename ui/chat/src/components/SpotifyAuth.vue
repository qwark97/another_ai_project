<template>
    <div class="spotify-auth-container">
        <button @click="auth()">SpotifyAuth</button>
    </div>
</template>

<script>
export default {
    name: 'SpotifyAuth',
    props: {
        placeholder: String
    },
    data() {
        return {
            
        }
    },
    methods: {
        auth() {
            window.location.href = process.env.SPOTIFY_REDIRECT ;
        }
    },
    created() {
    },

    beforeDestroy() {
        console.log('Chat component is destroyed.');
    }
}
</script>