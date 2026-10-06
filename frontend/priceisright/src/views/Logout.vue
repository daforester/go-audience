<template>
  <div />
</template>
<script>

export default {
  name: "LogoutView",
  created() {
    this.emitter.emit("removeTokensByService", "PriceIsRight");
    this.$router.push({ path: "/" });
  },
};

</script>
<style scoped>
</style>
