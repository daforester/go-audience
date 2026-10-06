<template>
  <div />
</template>
<script lang="ts">
import { defineComponent, type PropType } from 'vue';

export default defineComponent({
  name: 'WebSocket',
  props: {
    // A function is called on every (re)connect, so a reconnect uses a fresh token
    wsuri: {
      type: [String, Function] as PropType<string | (() => string)>,
      required: true,
    },
  },
  emits: ['sockopen', 'sockclose', 'sockerror', 'sockmsg'],
  data() {
    return {
      iAmConnecting: false,
      iAmClosing: false,
      socket: null as WebSocket | null,
      timer: {
        ping: null as ReturnType<typeof setInterval> | null,
        reconnect: null as ReturnType<typeof setTimeout> | null,
      },
    };
  },
  mounted() {
    this.setupSocket();
  },
  beforeUnmount() {
    this.closeSocket();
    this.clearTimers();
  },
  methods: {
    clearTimers() {
      if (this.timer.ping) {
        clearInterval(this.timer.ping);
        this.timer.ping = null;
      }
      if (this.timer.reconnect) {
        clearTimeout(this.timer.reconnect);
        this.timer.reconnect = null;
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
    reconnect() {
      this.clearTimers();
      this.timer.reconnect = setTimeout(this.setupSocket, 1000);
    },
    setupSocket() {
      if (this.iAmConnecting) {
        return;
      }
      this.iAmConnecting = true;
      this.iAmClosing = false;
      const uri = typeof this.wsuri === 'function' ? this.wsuri() : this.wsuri;
      const socket = new WebSocket(uri);
      this.socket = socket;

      socket.onopen = () => {
        this.iAmConnecting = false;
        this.$emit('sockopen');
        this.timer.ping = setInterval(this.pingSocket, 30000);
      };

      socket.onclose = () => {
        this.iAmConnecting = false;
        if (this.socket === socket) {
          this.socket = null;
        }
        this.$emit('sockclose');
        if (!this.iAmClosing) {
          this.reconnect();
        }
      };

      socket.onerror = () => {
        this.iAmConnecting = false;
        this.$emit('sockerror');
      };

      socket.onmessage = (e) => {
        const msg = e.data as string;
        if (msg.substring(0, 4) === 'PING') {
          socket.send(msg.replace('PING', 'PONG'));
          return;
        }
        if (msg.substring(0, 4) === 'PONG') {
          return;
        }

        this.$emit('sockmsg', msg);
      };
    },
  },
});
</script>
