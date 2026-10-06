<template>
  <div>
    <svg ref="map"
         class="map"
         preserveAspectRatio="none"
    />
  </div>
</template>

<script>

import * as d3 from 'd3';
import { doMap, generateCoast } from '../libs/terrain-master';

export default {
  name: 'MapMaker',
  data() {
    return {
      mapParams: {
        extent: {
          width: 1,
          height: 1,
        },
        generator: generateCoast,
        npts: 16384,
        ncities: 30,
        nterrs: 60,
        fontsizes: {
          region: 10,
          city: 7,
          town: 5,
        },
      },
    };
  },
  created() {
    window.addEventListener('resize', this.adjustViewBox);
  },
  unmounted() {
    window.removeEventListener('resize', this.adjustViewBox);
  },
  mounted() {
    this.adjustViewBox();
    setTimeout(this.resetBoard, 1);
  },
  methods: {
    adjustViewBox() {
      const l = -1000 * (this.mapParams.extent.width / 2);
      const t = -1000 * (this.mapParams.extent.height / 2);
      const w = 1000 * this.mapParams.extent.width;
      const h = 1000 * this.mapParams.extent.height;
      this.$refs.map.setAttribute('viewBox',`${l} ${t} ${w} ${h}`);
    },
    resetBoard() {
      return new Promise((r) => {
        const svg = d3.select(this.$refs.map);
        doMap(svg, this.mapParams);
        r();
      });
    },
  },
};

</script>

<style scoped lang="scss">
.map {
  position: absolute;
  top: 0;
  left: 0;
  transform-origin: top left;
  transform: scale(1);
}

</style>
<style lang="scss">
.map {
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
