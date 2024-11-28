<script setup lang="ts">
import { Photo } from '../types/photo'

const props = defineProps<{
  photo: Photo
}>()

const emit = defineEmits(['preview'])

const handleClick = () => {
  emit('preview', props.photo)
}
</script>

<template>
  <div class="photo-card" @click="handleClick">
    <img :src="photo.netPath" :alt="photo.name">
    <div class="photo-info">
      <h3>{{ photo.name }}</h3>
    </div>
  </div>
</template>

<style scoped>
.photo-card {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.photo-card:hover {
  transform: translateY(-4px);
}

.photo-card img {
  width: 100%;
  height: auto;
  display: block;
}

.photo-info {
  padding: 12px;
}
</style>