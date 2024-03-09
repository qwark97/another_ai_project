<template>
    <div class="chat-container">
        <div class="chat-messages">
            <p v-for="msg in history" :class="msg.type">{{ msg.message }}</p>
        </div>
        <div v-if="this.inputEnabled"  @keyup.enter="sendMessage" class="chat-input">
            <input type="text" v-model="messageToSent" id="message-input" :placeholder="placeholder">            
            <button id="send-button" @click="sendMessage">Send</button>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Chat',
    props: {
        placeholder: String
    },
    data() {
        return {
            socket: null,
            messageToSent: '',
            history: [],
            inputEnabled: true
        }
    },
    methods: {
        sendMessage() {
            if (!this.messageToSent) {
                return;
            }
            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
                this.socket.send(this.messageToSent);
                let historyMsg = {
                    type: 'user',
                    message: this.messageToSent
                }
                this.history.push(historyMsg);
                this.messageToSent = '';
            }
        },
        receiveMessage(message) {
            let historyMsg = {
                    type: 'assistant',
                    message: message
                }
            this.history.push(historyMsg);
        }
    },
    created() {
        this.socket = new WebSocket('ws://localhost:8080/api/v1/chat');
        this.socket.onopen = (event) => {
            console.log('WebSocket is open now.');
        };
        this.socket.onclose = (event) => {
            console.log('WebSocket is closed now:', event);
            this.inputEnabled = false;
        };
        this.socket.onerror = (event) => {
            console.error('WebSocket error observed:', event);
        };
        this.socket.onmessage = (event) => {
            console.log('WebSocket message received:', event);
            this.receiveMessage(event.data);
        };
    },
    beforeDestroy() {
        console.log('Chat component is destroyed.');
        if (this.socket) {
            this.socket.close();
        }
    }
}
</script>