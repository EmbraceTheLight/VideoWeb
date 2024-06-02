<template>
    <div>
        <canvas id="cav">你的浏览器不支持canvas</canvas>
    </div>
</template>

<script>
    export default {
        mounted() {
            console.log('test')
        }
    }
</script>

<style scoped>
      canvas {
        width: 1280px;
        height: 720px;
        background-color: rgba(0, 0, 0, 0.2);
      }
</style>