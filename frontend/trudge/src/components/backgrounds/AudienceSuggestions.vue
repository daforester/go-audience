<template>
  <div class="background vw-100 vh-100 position-absolute overflow-hidden">
    <div class="text-center">
      <div class="mt-5">
        <DisplayMessage v-for="(suggestion, i) in suggestions"
                        :key="`${i}-${suggestion}`"
                        effect="suggestion"
                        :message="suggestion"
                        :whole-word="true"
                        class="suggestion"
        />
      </div>
    </div>
  </div>
</template>

<script>

import DisplayMessage from '../display/DisplayMessage.vue';

export default {
  name: 'AudienceSuggestions',
  components: {
    DisplayMessage,
  },
  data() {
    return {
      slotTimers: [],
      usedSlots: [],
      suggestions: ['', '', '', '', '', '', '', '', '', ''],
    };
  },
  created() {
    this.emitter.on('newSuggestion', this.addSuggestion);
  },
  unmounted() {
    this.emitter.off('newSuggestion', this.addSuggestion);
  },
  methods: {
    addSuggestion(v) {
      let p;
      if (this.usedSlots.length >= this.suggestions.length) {
        clearTimeout(this.slotTimers.shift());
        p = this.usedSlots.shift();
      } else {
        do {
          p = Math.floor(Math.random() * this.suggestions.length);
        } while (this.usedSlots.includes(p));
      }
      this.usedSlots.push(p);
      this.suggestions.splice(p, 1, v);
      const t = setTimeout(() => {
        this.slotTimers.shift();
        this.usedSlots.shift();
      }, 4000);
      this.slotTimers.push(t);
    },
  },
};

</script>

<style scoped lang="scss">
.background {
  top: 0;
  left: 0;
  z-index: -98;
}

.bg-icon {
  font-size: 12vw;
  margin: 3rem;
}

.suggestion {
  position: absolute;
}

@for $i from 1 through 4 {
  .suggestion:nth-child(#{$i}) {
    left: 4rem;
    top: $i * 20%;
  }
}

@for $i from 1 through 4 {
  .suggestion:nth-child(#{$i + 4}) {
    right: 4rem;
    top: $i * 20%;
  }
}

.suggestion:nth-child(9) {
  left: 0;
  text-align: center;
  top: 30%;
  width: 100%
}

.suggestion:nth-child(10) {
  left: 0;
  text-align: center;
  top: 70%;
  width: 100%
}
</style>
