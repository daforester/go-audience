<template>
  <img v-if="src"
       class="join-qr"
       :src="src"
       :width="size"
       :height="size"
       :alt="`QR code to join room ${roomKey}`"
  >
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import QRCode from 'qrcode';
import { joinUrl } from '@/libs/join';

// JoinQr draws the join link as a QR code in the browser; nothing is sent elsewhere
export default defineComponent({
  name: 'JoinQr',
  props: {
    roomKey: {
      type: String,
      required: true,
    },
    size: {
      type: Number,
      default: 220,
    },
  },
  data() {
    return {
      src: '',
    };
  },
  watch: {
    roomKey: {
      immediate: true,
      handler() {
        this.draw();
      },
    },
  },
  methods: {
    draw() {
      if (!this.roomKey) {
        this.src = '';
        return;
      }
      QRCode.toDataURL(joinUrl(this.roomKey), {
        errorCorrectionLevel: 'M',
        margin: 1,
        // Render at twice the display size so it stays sharp on a projector
        width: this.size * 2,
        color: { dark: '#2b2233', light: '#ffffff' },
      })
        .then((src) => {
          this.src = src;
        })
        .catch(() => {
          this.src = '';
        });
    },
  },
});
</script>

<style scoped lang="scss">
.join-qr {
  background: #fff;
  border-radius: 0.6rem;
  image-rendering: pixelated;
}
</style>
