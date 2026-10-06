<template>
  <div>
    <button class="btn btn-primary d-none position-absolute" style="top: 0; left: 0;" @click="toggleAnimation">
      Play / Pause
    </button>
    <div ref="bglayer" class="background vw-100 vh-100 position-absolute overflow-hidden">
      <div class="travel">
        <div class="items">
          <div v-for="(icon, k) in icons"
               :key="`icon.${k}`"
               ref="icon"
               class="icon"
          >
            <FontAwesomeIcon :icon="['fal', icon]" />
          </div>
        </div>
        <div class="lines">
          <div v-for="l in travelLineCount"
               :key="`line.${l}`"
               ref="line"
               class="line"
          >
            <!-- -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

const icons = [
  'warehouse-alt',
  'torii-gate',
  'synagogue',
  'store',
  'store-alt',
  'school',
  'place-of-worship',
  'mosque',
  'monument',
  'landmark',
  'landmark-alt',
  'house',
  'gopuram',
  'university',
  'dungeon',
  'vihara',
  'igloo',
  'farm',
  'campground',
  'archway',
  'trees',
  'tree-palm',
  'tree-large',
  'tree',
  'tree-alt',
  'squirrel',
  'cloud-sleet',
  'cloud-hail-mixed',
  'clouds',
  'cloud-snow',
  'cloud-showers-heavy',
  'cloud-showers',
  'cloud-rainbow',
  'cloud-rain',
  'cloud-drizzle',
  'cloud-hail',
  'thunderstorm',
  'wind',
  'shovel',
  'treasure-chest',
  'scarecrow',
  'dove',
  'duck',
  'kiwi-bird',
  'crow',
  'water',
  'fire-smoke',
  'smoke',
  'dragon',
  'ufo',
  'campfire',
  'mountains',
  'acorn',
  'frog',
  'cauldron',
  'bat',
  'tombstone',
  'wagon-covered',
  'horse-saddle',
];

export default {
  name: 'TravelOriginal',
  components: {
    FontAwesomeIcon,
  },
  data() {
    return {
      angle: 0,
      destination: {
        x: 0,
        y: 0,
      },
      distance: 0,
      drawTimer: null,
      fps: 60,
      iconCount: 60,
      icons: [],
      maxDuration: 30,
      resetTimer: null,
      speed: 3,
      stepsNeeded: 0,
      stepsRemaining: 0,
      travelLineCount: 0,
    };
  },
  created() {
    this.resetTimer = setInterval(this.resetBoard, this.maxDuration * 1000);
    this.drawTimer = setInterval(this.drawStep, 1000 / this.fps);
    for (let n = 0; n < this.iconCount; n += 1) {
      this.icons.push(this.randIcon());
    }
  },
  unmounted() {
    if (this.drawTimer) {
      clearInterval(this.drawTimer);
    }
    if (this.resetTimer) {
      clearInterval(this.resetTimer);
    }
  },
  mounted() {
    this.resetBoard();
  },
  methods: {
    drawStep() {
      if (this.stepsRemaining === 0) {
        this.travelLineCount += 1;
        this.distance = Math.floor((Math.random()
          * Math.sqrt((window.innerWidth ** 2) + (window.innerHeight ** 2))) + 200);
        this.angle = Math.random() * Math.PI * 2;
        // const icon = this.$refs.icon[Math.floor(Math.random() * this.$refs.icon.length)];
        // this.distance = Math.sqrt((((window.innerWidth / 2) - icon.data_left) ** 2)
        //   + ((icon.data_top - (window.innerHeight / 2)) ** 2)) / 2;
        // this.angle = Math.atan2((window.innerWidth / 2) - icon.data_left,
        // icon.data_top - (window.innerHeight / 2)) + (Math.PI / 2);
        this.stepsNeeded = Math.ceil(this.distance / this.speed);
        this.stepsRemaining = this.stepsNeeded;
        return;
      }
      const line = this.$refs.line[this.$refs.line.length - 1];
      if (this.stepsRemaining === this.stepsNeeded) {
        // First frame - set line position
        line.data_width = 0;
        line.data_left = window.innerWidth / 2;
        line.data_top = window.innerHeight / 2;
        line.style.left = `${line.data_left}px`;
        line.style.top = `${line.data_top}px`;
        line.style.transform = `rotate(${this.angle + Math.PI}rad)`;
        this.destination.x = this.distance * Math.cos(this.angle - (Math.PI));
        this.destination.y = 0 - (this.distance * Math.sin(this.angle));
      }

      const sx = this.destination.x / this.stepsNeeded;
      const sy = this.destination.y / this.stepsNeeded;

      line.data_width += Math.sqrt((sx ** 2) + (sy ** 2));
      line.style.width = `${line.data_width}px`;

      // Move everything
      this.destination.x += sx;
      this.destination.y += sy;

      this.$refs.icon.forEach((i) => {
        i.data_left += sx;
        i.data_top += sy;
        i.style.left = `${i.data_left}px`;
        i.style.top = `${i.data_top}px`;
      });
      this.$refs.line.forEach((i, k) => {
        if (k === this.$refs.line.length - 1) {
          return;
        }
        i.data_left += sx;
        i.data_top += sy;
        i.style.left = `${i.data_left}px`;
        i.style.top = `${i.data_top}px`;
      });

      this.stepsRemaining -= 1;
    },
    randIcon() {
      return icons[Math.floor(Math.random() * icons.length)];
    },
    resetBoard() {
      this.$refs.icon.forEach((i) => {
        i.data_left = Math.floor(
          Math.random() * window.innerWidth * 6,
        ) - (window.innerWidth * 3);
        i.style.left = `${i.data_left}px`;
        i.data_top = Math.floor(
          Math.random() * window.innerHeight * 6,
        ) - (window.innerHeight * 3);
        i.style.top = `${i.data_top}px`;
      });
      this.travelLineCount = 0;
      this.stepsRemaining = 0;
      this.stepsNeeded = 0;
      this.distance = 0;
      this.angle = 0;
      this.destination.x = 0;
      this.destination.y = 0;
    },
    toggleAnimation() {
      if (this.drawTimer) {
        clearInterval(this.drawTimer);
        this.drawTimer = null;
      } else {
        this.drawTimer = setInterval(this.drawStep, 1000 / this.fps);
      }
      if (this.resetTimer) {
        clearInterval(this.resetTimer);
        this.resetTimer = null;
      } else {
        this.resetTimer = setInterval(this.resetBoard, this.maxDuration * 1000);
      }
    }
  },
};

</script>

<style scoped lang="scss">

button {
  z-index: 99;
}

.background {
  top: 0;
  left: 0;
  z-index: -99;
}

.icon {
  color: darkgray;
  font-size: 4rem;
  position: absolute;
}

.line {
  background: #e74c3c;
  border-radius: 2rem;
  height: 2rem;
  margin-top: 8rem;
  position: absolute;
  transform-origin: top left;
}

</style>
