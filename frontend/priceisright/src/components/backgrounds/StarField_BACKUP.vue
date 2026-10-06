<template>
  <canvas id="canvas"
          ref="canvas"
          class="background m-0 p-0 position-fixed overflow-hidden vh-100 vw-100"
  />
</template>

<script>

export default {
  name: 'StarField',
  data() {
    return {
      canvas: undefined,
      context: undefined,
      prevTime: undefined,
      stars: undefined,
    };
  },
  mounted() {
    this.canvas = this.$refs.canvas;
    this.context = this.canvas.getContext('2d');

    this.stars = this.makeStars(200);
    this.setCanvasSize();

    window.requestAnimationFrame(this.init);

    window.addEventListener('resize', this.setCanvasSize);
  },
  methods: {
    clear() {
      this.context.fillStyle = 'black';
      this.context.fillRect(0, 0, this.canvas.width, this.canvas.height);
    },
    init(time) {
      this.prevTime = time;
      window.requestAnimationFrame(this.tick);
    },
    makeStars(count) {
      const out = [];
      for (let i = 0; i < count; i += 1) {
        const s = {
          x: Math.random() * 1600 - 800,
          y: Math.random() * 900 - 450,
          z: Math.random() * 1000,
        };
        out.push(s);
      }
      return out;
    },
    moveStars(distance) {
      const count = this.stars.length;
      for (let i = 0; i < count; i += 1) {
        const s = this.stars[i];
        s.z -= distance;
        while (s.z <= 1) {
          s.z += 1000;
        }
      }
    },
    putPixel(x, y, brightness) {
      const intensity = brightness * 255;
      this.context.fillStyle = `rgb(${intensity},${intensity},${intensity})`;
      this.context.fillRect(x, y, 2, 2);
      /*
      this.context.beginPath();
      this.context.arc(x, y, 2, 0, 2 * Math.PI, false);
      this.context.fill();
      this.context.stroke();
      */
    },
    setCanvasSize() {
      this.canvas.height = window.innerHeight;
      this.canvas.width = window.innerWidth;
    },
    tick(time) {
      const h = this.canvas.height;
      const w = this.canvas.width;

      const elapsed = time - this.prevTime;
      this.prevTime = time;

      this.moveStars(elapsed * 0.25);

      this.clear();

      const cx = w / 2;
      const cy = h / 2;

      const count = this.stars.length;
      for (let i = 0; i < count; i += 1) {
        const star = this.stars[i];

        const x = cx + star.x / (star.z * 0.001);
        const y = cy + star.y / (star.z * 0.001);

        if (!(x < 0 || x >= w || y < 0 || y >= h)) {
          const d = star.z / 1000.0;
          const b = 1 - d * d;

          this.putPixel(x, y, b * 2);
        }
      }

      window.requestAnimationFrame(this.tick);
    },
  },
};

</script>

<style scoped>
.background {
  top: 0;
  left: 0;
  z-index: -100;
}
</style>
