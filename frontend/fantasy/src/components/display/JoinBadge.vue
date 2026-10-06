<template>
  <div class="join-badge">
    <JoinQr :room-key="roomKey" :size="72" />
    <div class="text">
      <div class="small">Join at {{ host }}</div>
      <div class="code">{{ roomKey }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import JoinQr from './JoinQr.vue';
import { joinHost } from '@/libs/join';

// JoinBadge keeps the way in on screen for latecomers once the round screens take over
export default defineComponent({
  name: 'JoinBadge',
  components: {
    JoinQr,
  },
  props: {
    roomKey: {
      type: String,
      required: true,
    },
  },
  computed: {
    host(): string {
      return joinHost();
    },
  },
});
</script>

<style scoped lang="scss">
.join-badge {
  align-items: center;
  background: rgba(255, 255, 255, 0.92);
  border: 2px solid #2b2233;
  border-radius: 0.8rem;
  bottom: 1rem;
  display: flex;
  gap: 0.6rem;
  padding: 0.4rem 0.75rem 0.4rem 0.4rem;
  position: fixed;
  right: 1rem;
  z-index: 50;
}
.code {
  font-size: 1.3rem;
  font-weight: 800;
  letter-spacing: 0.12em;
  line-height: 1;
}
</style>
