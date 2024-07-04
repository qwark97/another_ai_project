<template>
    <div class="chat-container">
        <div class="chat-messages">
            <p v-for="msg in history" :class="msg.type">{{ msg.message }}</p>
        </div>
        <div v-if="this.inputEnabled" @keyup.enter="sendMessage" class="chat-input">
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
            inputEnabled: true,
            currentConversationID: null,
            reconecting: false
        }
    },
    methods: {
        sendMessage() {
            if (!this.messageToSent) {
                return;
            }
            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
                let msg = {
                    instruction: this.messageToSent,
                    conversation_id: this.currentConversationID,
                };
                this.socket.send(JSON.stringify(msg));
                let historyMsg = {
                    type: 'user',
                    message: this.messageToSent
                }
                this.history.push(historyMsg);
                this.messageToSent = '';
            }
        },
        receiveMessage(message) {
            console.log('Received message:', message);
            let msg = JSON.parse(message);
            console.log('Parsed message:', msg);
            let historyMsg = {
                type: 'assistant',
                message: msg["answer"]
            }
            this.currentConversationID = msg["conversation_id"];
            this.history.push(historyMsg);
        },
        connect() {
            this.socket = new WebSocket('ws://192.168.0.243:8080/api/v1/chat');
            this.registerEvents();
        },
        registerEvents() {
            this.socket.onopen = (event) => {
                this.inputEnabled = true;
                console.log('WebSocket is open now:', event);
            };
            this.socket.onclose = (event) => {
                console.log('WebSocket is closed now:', event);
                this.inputEnabled = false;
                this.connect();
            };
            this.socket.onerror = (event) => {
                console.error('WebSocket error observed:', event);
            };
            this.socket.onmessage = (event) => {
                console.log('WebSocket message received:', event);
                this.receiveMessage(event.data);
            };
        },
    },
    created() {
        this.connect();
    },

    beforeDestroy() {
        console.log('Chat component is destroyed.');
        if (this.socket) {
            this.socket.close();
        }
    }
}
</script>