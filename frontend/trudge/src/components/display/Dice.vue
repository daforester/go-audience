<template>
  <div>
    <div id="dice-box" />
  </div>
</template>
<script>

import DiceBox from "@3d-dice/dice-box";

export default {
  name: "DiceComponent",
  data() {
    return {
      Box: null,
      label: '',
      params: {},
    };
  },
  created() {
    this.emitter.on("rollDice", this.rollDice);
  },
  unmounted() {
    this.emitter.off("rollDice", this.rollDice);
  },
  methods: {
    rollDice() {
      document.getElementById("dice-box").replaceChildren();
      const Box = new DiceBox("#dice-box", {
        assetPath: "/assets/",
        theme: "default",
        offscreen: true,
        scale: 6,
      });

      Box.init().then(async () => {
        Box.roll(["1d20"]);
      });
    },
  },
};

</script>
<style lang="scss">

#dice-box {
  position: fixed;
  box-sizing: border-box;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

#dice-box canvas {
  width: 100%;
  height: 100%;
}

#roll {
  position: absolute;
  top: 5px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1;
  cursor: pointer;
}

</style>
