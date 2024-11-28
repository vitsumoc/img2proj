<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Photo } from '../types/photo'
import PhotoCard from './PhotoCard.vue'

const showPreview = ref(false)
const currentPhoto = ref<Photo | null>(null)

const handlePreview = (photo: Photo) => {
  currentPhoto.value = photo
  showPreview.value = true
}

const closePreview = () => {
  showPreview.value = false
  currentPhoto.value = null
}

let photos = ref<Photo[]>([])

const fetchImages = async () => {
  try {
    photos.value = []
    const response = await fetch('/api/imgs')
    const res = await response.json()
    let id = 1
    res.data.forEach((item: any) => {
      photos.value.push({
        id: id++,
        dirName: item.dirName,
        dirPath: item.dirPath,
        imgPath: item.imgPath,
        netPath: '\\img\\' + item.dirName + '\\' + item.imgPath,
        name: item.name,
        updateAt: item.updateAt
      })
    })
    console.log(photos.value)
  } catch (error) {
    console.error('获取图片数据失败:', error)
  }
}

onMounted(() => {
  fetchImages()
})
</script>

<template>
  <div class="photo-gallery">
    <PhotoCard 
      v-for="photo in photos"
      :key="photo.id"
      :photo="photo"
      @preview="handlePreview"
    />
    
    <!-- 预览模态框 -->
    <div v-if="showPreview" class="preview-modal" @click.self="closePreview">
      <div class="preview-content">
        <img :src="currentPhoto?.netPath" :alt="currentPhoto?.name">
        <div class="preview-info">
          <h3>{{ currentPhoto?.name }}</h3>
          <p>图片文件: {{ currentPhoto?.name }}</p>
          <p>项目文件夹: {{ currentPhoto?.dirName }}</p>
          <p>项目路径: {{ currentPhoto?.dirPath }}</p>
        </div>
        <button class="close-button" @click="closePreview">&times;</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.photo-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 15px;
  padding: 20px;
}

.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.preview-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  max-width: 90%;
  max-height: 90vh;
  position: relative;
}

.preview-content img {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.preview-info {
  margin-top: 15px;
  color: #333;
}

.preview-info p {
  margin: 5px 0;
  font-size: 14px;
}

.close-button {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.close-button:hover {
  color: #000;
}
</style>