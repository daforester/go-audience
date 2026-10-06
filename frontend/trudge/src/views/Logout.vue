<template>
  <div />
</template>
<script>

export default {
  name: 'LogoutView',
  created() {
    this.emitter.emit('removeTokensByService', 'Trudge');
    this.$router.push({ path: '/' });
  },
};

</script>
<style scoped>
</style>
