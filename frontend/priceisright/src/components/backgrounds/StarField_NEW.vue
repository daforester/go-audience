<template>
  <div></div>
</template>

<script>

class Star {
  constructor(x, y, z) {
    this.x = x; this.y = y; this.z = z;

    this.size = 5;
    return this;
  }
}

export default {
  name: 'StarField',
  data() {
    return {
      canvas: undefined,
      context: undefined,
      maxDepth: 2500,
      prevTime: undefined,
      stars: undefined,
    };
  },
  mounted() {
    this.canvas = document.createElement('canvas');
    // this.canvas.classList.add('m-0', 'p-0', 'overflow-hidden', 'vh-100', 'vw-100');
    this.canvas.style.top = 0;
    this.canvas.style.left = 0;
    this.canvas.style.zIndex = -100;
    this.canvas.style.position = 'fixed';

    document.body.prepend(this.canvas);

    this.context = this.canvas.getContext('2d');

    this.setCanvasSize();
    this.stars = this.makeStars(200);

    window.requestAnimationFrame(this.init);

    window.addEventListener('resize', () => {
      this.setCanvasSize();
      this.stars = this.makeStars(200);
    });
  },
  beforeDestroy() {
    this.canvas.remove();
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
      for (let index = 0; index < count; index += 1) {
        out.push(new Star(
          Math.random() * this.canvas.width,
          Math.random() * this.canvas.height,
          index * (this.maxDepth / 200),
        ));
      }

      return out;
    },
    setCanvasSize() {
      this.canvas.height = window.innerHeight;
      this.canvas.width = window.innerWidth;
    },
    tick() {
      this.clear();

      for (let index = this.stars.length - 1; index > -1; index -= 1) {
        const star = this.stars[index];

        star.z -= 5;

        if (star.z < 0) {
          // this.stars.push(this.stars.splice(index, 1)[0]);
          star.z = this.maxDepth;
        } else {
          const translateX = this.canvas.width * 0.5;
          const translateY = this.canvas.height * 0.5;

          const fieldOfView = (this.canvas.height + this.canvas.width) * 0.25;

          const starX = (star.x - translateX) / (star.z / fieldOfView) + translateX;
          const starY = (star.y - translateY) / (star.z / fieldOfView) + translateY;

          const scale = fieldOfView / (fieldOfView + star.z);

          const color = Math.floor(scale * 512);

          this.context.fillStyle = `rgb(${color},${color},${color})`;
          if (starX < 0 || starY < 0 || starX > this.canvas.width || starY > this.canvas.height) {
            star.z = this.maxDepth;
          } else {
            this.context.fillRect(starX, starY, star.size * scale, star.size * scale);
          }
        }
      }

      window.requestAnimationFrame(this.tick);
    },
  },
};

</script>

<style scoped>

</style>
