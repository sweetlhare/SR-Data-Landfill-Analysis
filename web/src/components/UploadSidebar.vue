<script setup>
import { reactive, ref } from 'vue';
import FileDrop from '../components/FileDrop.vue';
import { useRouter } from 'vue-router';
import { useGetSessionId } from '../composables/useGetSessionId';

const props = defineProps({
  url: { type: String, required: true },
  isSidebarShown: { type: Boolean, default: false },
});
defineEmits(['hideSidebar']);
const router = useRouter();

const fetching = reactive({
  upload: false,
});
const errors = reactive({
  upload: false,
  date: false,
  dateNotUnique: false,
  uploadTimeout: false,
});

const newSurveyDate = ref(null);
const newSurveyImages = ref([]);
function updateImages(newArray) {
  newSurveyImages.value = newArray;
}
function checkDate() {
  let result = true;
  if (!newSurveyDate.value || new Date(newSurveyDate.value) > new Date()) {
    errors.date = true;
    setTimeout(() => {
      errors.date = false;
    }, 2000);
    result = false;
  }
  return result;
}
async function uploadImages() {
  if (!checkDate()) return;

  fetching.upload = true;
  errors.dateNotUnique = false;
  errors.uploadTimeout = false;
  errors.upload = false;

  const formData = new FormData();
  formData.append('date', new Date(newSurveyDate.value).toISOString());
  // formData.append('max_size', '4096');
  for (const image of newSurveyImages.value) {
    formData.append('files', image);
  }
  const sessionId = useGetSessionId();

  try {
    const newSurvey = await fetch(props.url, {
      method: 'POST',
      body: formData,
      headers: {
        // 'Content-Type': 'multipart/form-data',
        'X-Session-ID': sessionId,
      },
    }).then((res) => {
      if (!res.ok) {
        if (res.status === 401) router.push({ name: 'sign-in' });
        else if (res.status === 400) errors.dateNotUnique = true;
        else if (res.status === 504) {
          errors.uploadTimeout = true;
        } else errors.upload = true;
        throw new Error('HTTP error ' + res.status);
      }
      return res.json();
    });
    router.push({ name: 'survey', params: { id: newSurvey.id } });
  } catch (error) {
    console.error('Error: ', error);
  }

  fetching.upload = false;
}
</script>

<template>
  <transition name="fade">
    <div v-if="isSidebarShown" class="overlay"></div>
  </transition>
  <transition name="slide">
    <div v-if="isSidebarShown" class="sidebar">
      <div class="sidebar__header">
        <h3 class="section-title sidebar__title">Добавление съемки</h3>
        <button @click="$emit('hideSidebar')">
          <img src="/icons/sharp.svg" alt="Скрыть боковую панель" />
        </button>
      </div>

      <form class="sidebar__form form" @submit.prevent="uploadImages">
        <div class="input-block">
          <label class="input-block__label" for="date">Дата съемки</label>
          <input
            id="date"
            v-model="newSurveyDate"
            class="input-block__input"
            type="date"
            placeholder="Введите дату"
          />
        </div>

        <div v-if="fetching.upload" class="info-message">Подождите, идёт проверка...</div>
        <FileDrop v-else class="images__main" @update-files="updateImages"></FileDrop>
        <div v-if="errors.dateNotUnique" class="error-message">
          Ошибка при отправке. Проверьте уникальность даты
        </div>
        <div v-else-if="errors.uploadTimeout" class="error-message">
          Превышено время ожидания от сервера. Попробуйте перезагрузить страницу.
        </div>
        <div v-else-if="errors.upload" class="error-message">
          Ошибка при отправке. Повторите позднее
        </div>
        <button v-if="newSurveyImages.length" type="submit" class="btn btn_primary">
          Найти нарушения
        </button>
        <button
          v-if="newSurveyImages.length"
          type="button"
          class="btn btn_danger"
          @click="newSurveyImages = []"
        >
          Очистить
        </button>
      </form>
    </div>
  </transition>
</template>

<style scoped lang="scss">
.overlay {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #000000da;
}
.sidebar {
  position: fixed;
  top: 0;
  right: 0;

  height: 100vh;
  width: 40%;
  min-width: 600px;
  padding: 42px;
  overflow-y: scroll;

  border-top-left-radius: 8px;
  background-color: white;
}

.sidebar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar__title {
  margin-bottom: 0;
}

.sidebar__form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.input-block__label {
  margin-bottom: 6px;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
  transition: all 0.3s ease-in 0s;
}

.fade-enter-active,
.fade-leave-active {
  transition: transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transition: all 0.3s ease-in 0s;
}
</style>
