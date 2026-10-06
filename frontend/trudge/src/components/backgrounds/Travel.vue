<template>
  <div>
    <div ref="bglayer" class="background vw-100 vh-100 position-absolute overflow-hidden">
      <img ref="map" class="map" src="@/assets/maps/01.svg" />
    </div>
  </div>
</template>

<script>

import map01 from "@/assets/maps/01.svg";
import map02 from "@/assets/maps/02.svg";
import map03 from "@/assets/maps/03.svg";
import map04 from "@/assets/maps/04.svg";
import map05 from "@/assets/maps/05.svg";
import map06 from "@/assets/maps/06.svg";
import map07 from "@/assets/maps/07.svg";
import map08 from "@/assets/maps/08.svg";
import map09 from "@/assets/maps/09.svg";
import map10 from "@/assets/maps/10.svg";
import map11 from "@/assets/maps/11.svg";

export default {
  name: 'Travel',
  data() {
    return {
      maxDuration: 30,
      resetTime: null,
    };
  },
  mounted() {
    this.resetTimer = setInterval(this.resetBoard, this.maxDuration * 1000);
    setTimeout(this.resetBoard, 1);
  },
  unmounted() {
    if (this.resetTimer) {
      clearInterval(this.resetTimer);
    }
  },
  methods: {
    resetBoard() {
      if (this.$refs.map === null) {
        return;
      }
      return new Promise((r) => {
        const fileNo = Math.floor(Math.random() * 10) + 1;
        switch (fileNo) {
          case 1:
            this.$refs.map.src = map01;
            break;
          case 2:
            this.$refs.map.src = map02;
            break;
          case 3:
            this.$refs.map.src = map03;
            break;
          case 4:
            this.$refs.map.src = map04;
            break;
          case 5:
            this.$refs.map.src = map05;
            break;
          case 6:
            this.$refs.map.src = map06;
            break;
          case 7:
            this.$refs.map.src = map07;
            break;
          case 8:
            this.$refs.map.src = map08;
            break;
          case 9:
            this.$refs.map.src = map09;
            break;
          case 10:
            this.$refs.map.src = map10;
            break;
          case 11:
            this.$refs.map.src = map11;
            break;
        }
        r();
      });
    },
  },
};

</script>

<style scoped lang="scss">
@keyframes moveMap {
  0% {
    transform: scale(4) translate(-15%, 0%);
  }
  25% {
    transform: scale(4) translate(-70%, -15%);
  }
  50% {
    transform: scale(4) translate(-55%, -70%);
  }
  75% {
    transform: scale(4) translate(0%, -55%);
  }
  100% {
    transform: scale(4) translate(-15%, 0%);
  }
}

button {
  z-index: 99;
}

.background {
  top: 0;
  left: 0;
  z-index: -99;
}

.map {
  animation-name: moveMap;
  animation-duration: 45s;
  animation-iteration-count: infinite;
  animation-timing-function: linear;
  transform-origin: 0 0;
  opacity: 0.2;
  position: absolute;
  top: 0;
  left: 0;
}

</style>
