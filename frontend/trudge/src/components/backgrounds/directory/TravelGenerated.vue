<template>
  <div>
    <div ref="bglayer" class="background vw-100 vh-100 position-absolute overflow-hidden">
      <svg ref="travelMap"
           class="travelMap"
           preserveAspectRatio="none"
      ></svg>
    </div>
  </div>
</template>

<script>

import * as d3 from 'd3';
import { doMap, defaultParams } from '../../../libs/terrain-master';

export default {
  name: 'TravelGenerated',
  created() {
    window.addEventListener("resize", this.adjustViewBox);
  },
  unmounted() {
    window.removeEventListener("resize", this.adjustViewBox);
  },
  mounted() {
    setTimeout(this.resetBoard, 1);
    this.adjustViewBox();
  },
  methods: {
    adjustViewBox() {
      this.$refs.travelMap.setAttribute('viewBox', `-450 -450 ${window.innerWidth} ${window.innerHeight}`);
    },
    resetBoard() {
      return new Promise((r) => {
        const svg = d3.select(this.$refs.travelMap);
        doMap(svg, defaultParams);
        r();
      });
    },
  },
};

</script>

<style scoped lang="scss">
@keyframes moveMap {
  0% {
    transform: scale(3) translate(-15%, 0%);
  }
  25% {
    transform: scale(3) translate(-75%, -15%);
  }
  50% {
    transform: scale(3) translate(-55%, -75%);
  }
  75% {
    transform: scale(3) translate(0%, -55%);
  }
  100% {
    transform: scale(3) translate(-15%, 0%);
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

.travelMap {
  animation-name: moveMap;
  animation-duration: 45s;
  animation-iteration-count: infinite;
  animation-timing-function: linear;
  transform-origin: 0 0;
  opacity: 0.2;
  position: absolute;
  top: 0;
  left: 0;
  height: 200%;
  width: 200%;
}

</style>
<style type="scss">
.travelMap {
  path, line {
    fill: none;
    stroke: black;
    stroke-linecap: round;
  }

  .field {
    stroke: none;
    fill-opacity: 1.0;
  }

  .slope {
    stroke-width: 1;
  }

  .river {
    stroke-width: 2;
  }

  .coast {
    stroke-width: 4;
  }

  .border {
    stroke-width: 5;
    stroke-dasharray: 4,4;
    stroke-linecap: butt;
  }

  text {
    font-family: "Palatino Linotype", "Book Antiqua", Palatino, serif;
    color: black;
    stroke: white;
    stroke-width: 5;
    stroke-linejoin: round;
    paint-order: stroke;
  }

  text.region {
    stroke-width:10;
    font-variant: small-caps;
  }
}

</style>
