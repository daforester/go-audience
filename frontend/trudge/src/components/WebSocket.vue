<template>
  <div />
</template>
<script>
export default {
  name: 'WebSocket',
  props: {
    wsuri: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      iAmConnecting: false,
      iAmClosing: false,
      socket: null,
      timer: {
        ping: null,
        reconnect: null,
      },
    };
  },
  mounted() {
    this.setupSocket();
  },
  beforeDestroy() {
    this.closeSocket();
    this.clearTimers();
  },
  methods: {
    clearTimers() {
      if (this.timer.ping) {
        clearInterval(this.timer.ping);
      }
      if (this.timer.reconnect) {
        clearInterval(this.timer.reconnect);
      }
    },
    closeSocket() {
      this.iAmClosing = true;
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
    },
    pingSocket() {
      if (!this.socket) {
        if (this.timer.ping) {
          clearInterval(this.timer.ping);
        }
        return;
      }

      this.socket.send(`PING ${new Date().getTime()}`);
    },
    setupSocket() {
      if (this.iAmConnecting) {
        return;
      }
      this.iAmConnecting = true;
      this.socket = new WebSocket(this.wsuri);

      this.socket.onopen = () => {
        this.iAmConnecting = false;
        this.$emit('sockopen');
        this.timer.ping = setInterval(this.pingSocket, 30000);
      };

      this.socket.onclose = () => {
        this.iAmConnecting = false;
        this.$emit('sockclose');
        this.socket = null;
        if (!this.iAmClosing) {
          this.clearTimers();
          this.timer.reconnect = setTimeout(this.setupSocket, 1000);
        }
      };

      this.socket.onerror = () => {
        this.iAmConnecting = false;
        this.$emit('sockerror');
        this.socket = null;
        this.clearTimers();
        this.timer.reconnect = setTimeout(this.setupSocket, 1000);
      };

      this.socket.onmessage = (e) => {
        const msg = e.data;
        if (msg.substr(0, 4) === 'PING') {
          this.socket.send(msg.replace('PING', 'PONG'));
          return;
        }
        if (msg.substr(0, 4) === 'PONG') {
          return;
        }

        this.$emit('sockmsg', msg);
      };
    },
  },
};
</script>
