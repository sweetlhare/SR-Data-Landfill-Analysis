<script setup>
import { reactive, ref } from 'vue';
import { useCreateURL } from '../composables/useCreateUrl';

const emit = defineEmits('update-files');

const isDragging = ref(false);
const fileInput = ref(null);
const files = reactive([]);

function onChange() {
  for (const file of fileInput.value.files) {
    if (files.find((element) => element.name === file.name && element.size === file.size)) continue;
    files.push(file);
  }
  emit('update-files', files);
}

function onDragover(e) {
  e.preventDefault();
  isDragging.value = true;
}

function onDragleave() {
  isDragging.value = false;
}

function onDrop(e) {
  e.preventDefault();
  fileInput.value.files = e.dataTransfer.files;
  onChange();
  isDragging.value = false;
}

function remove(i) {
  files.splice(i, 1);
}
</script>

<template>
  <div class="dropzone-container" @dragover="onDragover" @dragleave="onDragleave" @drop="onDrop">
    <input
      id="fileInput"
      ref="fileInput"
      type="file"
      multiple
      name="file"
      class="hidden-input"
      accept=".pdf,.jpg,.jpeg,.png"
      @change="onChange"
    />
    <label for="fileInput" class="file-label">
      <div v-if="isDragging">Отпустите изображения, чтобы загрузить их.</div>
      <div v-else>Перетащите изображения или <u>нажмите сюда</u> для загрузки.</div>
    </label>
    <div v-if="files.length" class="preview-container mt-4">
      <div v-for="file in files" :key="file.name" class="preview-card">
        <div>
          <!-- <img class="preview-img" :src="useCreateURL(file)" /> -->
          <p class="file-name">
            {{ file.name }}
          </p>
        </div>
        <div>
          <button
            class="ml-2"
            type="button"
            title="Remove file"
            @click="remove(files.indexOf(file))"
          >
            <b>×</b>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dropzone-container {
  padding: 4rem;
  border: 1px dashed var(--green);
  border-radius: 8px;

  font-size: 14px;
  line-height: 18px;
  color: var(--green);
  text-align: center;
}

.hidden-input {
  opacity: 0;
  overflow: hidden;
  position: absolute;
  width: 1px;
  height: 1px;
}

.file-label {
  font-size: 20px;
  display: block;
  cursor: pointer;
}

.preview-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  width: 100%;
  margin-top: 2rem;
}

.preview-card {
  display: flex;
  justify-content: space-between;
  border: 1px solid #a2a2a2;
  padding: 5px;
  margin-left: 5px;
}

.file-name {
  width: 150px;
  overflow-wrap: break-word;
}

.preview-img {
  width: 150px;
  height: 150px;
  border-radius: 5px;
}
</style>
