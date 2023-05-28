<script setup>
import { ref, computed, reactive } from 'vue';
import { useRoute } from 'vue-router';
import * as XLSX from 'xlsx';

import BaseDictionary from '../components/BaseDictionary.vue';
import FileDrop from '../components/FileDrop.vue';
import BaseMap from '../components/BaseMap.vue';

import { API_BASE_URL } from '../config';
import { useFormatDate } from '../composables/useFormatDate';
import { useGetCookie } from '../composables/useGetCookie.js';

const route = useRoute();

const landfillData = ref(null);
const landfillId = computed(() => landfillData.value?.id ?? null);
const mode = ref('view');
const fetching = reactive({
  upload: false,
  audit: false,
});
const errors = reactive({
  upload: false,
  audit: false,
});

getLandfillData();
async function getLandfillData() {
  const response = await fetch(`${API_BASE_URL}/landfills/${route.params.id}`).then((response) =>
    response.json()
  );

  landfillData.value = response;
}

const surveys = computed(() => (landfillData.value ? landfillData.value.surveys.toReversed() : []));
const viewedSurveyIndex = ref(0);
const audits = computed(() =>
  surveys.value.length
    ? surveys.value[viewedSurveyIndex.value].audits.filter((audit) => !audit.ai_generated_status)
    : []
);
function setViewedSurvey(index) {
  viewedSurveyIndex.value = index;
}
const images = computed(() =>
  surveys.value.length ? surveys.value[viewedSurveyIndex.value].ai_images : []
);
const mainImageIndex = ref(0);
function setMainImage(index) {
  mainImageIndex.value = index;
}

const lastAudit = computed(() =>
  surveys.value.length ? surveys.value.at(-1).audits.at(-1) : null
);
const lastAuditViolations = computed(() =>
  lastAudit.value ? lastAudit.value.violations.filter((criterion) => criterion.status) : []
);

const landfillDictionaryItems = computed(() =>
  landfillData.value
    ? [
        {
          key: 'Населенный пункт',
          value: landfillData.value.city,
        },
        {
          key: 'Адрес',
          value: landfillData.value.address,
        },
        {
          key: 'Координаты',
          value: landfillData.value.coordinates.join(', '),
        },
        {
          key: 'ФИО ответственного за полигон',
          value: landfillData.value.manager.name,
        },
        {
          key: 'Должность ответственного за полигон',
          value: landfillData.value.manager.position,
        },
        {
          key: 'Контакты ответственного за полигон',
          value: landfillData.value.manager.phone + ', ' + landfillData.value.manager.email,
        },
      ]
    : []
);

const auditDictionaryItems = computed(() =>
  lastAudit.value
    ? [
        {
          key: 'Дата последней проверки',
          value: useFormatDate(lastAudit.value.date),
        },
        {
          key: 'ФИО проверившего',
          value: lastAudit.value.auditor.name,
        },
        {
          key: 'Должность проверившего',
          value: lastAudit.value.auditor.position,
        },
      ]
    : []
);

// Скачать отчет в excel
function generateExcelFile() {
  const dataToExport = [
    ['Идентификатор полигона', String(landfillData.value.id)],
    ['Фактический адрес', landfillData.value.city + ', ' + landfillData.value.address],
    ['Дата проверки', useFormatDate(lastAudit.value.date)],
    ['ФИО проверившего', lastAudit.value.auditor.name],
    ['Должность проверившего', lastAudit.value.auditor.position],
    ['Нарушения:', ''],
  ];

  for (const violation of lastAudit.value.violations) {
    dataToExport.push([violation.status ? 'Обнаружено' : 'Не обнаружено', violation.title]);
  }

  const ws = XLSX.utils.json_to_sheet(dataToExport);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, 'Отчет по полигону');
  const reportDate = useFormatDate(
    landfillData.value.surveys.toReversed()[0].audits.toReversed()[0].date
  );
  XLSX.writeFile(wb, `Отчет по полигону №${landfillData.value.id} ${reportDate}.xlsx`);
}

// Новая съемка
const newSurveyDate = ref(null);
const newSurveyImages = ref([]);
const isAuditComplete = computed(
  () => criteria.value.filter((element) => element.isAudited === false).length === 0
);
function updateImages(newArray) {
  newSurveyImages.value = newArray;
}
function switchToUploadMode() {
  if (!newSurveyDate.value || new Date(newSurveyDate.value) > new Date()) return;
  mode.value = 'upload';
}
async function uploadImages() {
  fetching.upload = true;
  errors.upload = false;

  const formData = new FormData();
  formData.append('date', newSurveyDate.value);
  for (const image of newSurveyImages.value) {
    formData.append('files', image);
  }

  try {
    const newSurvey = await fetch(`${API_BASE_URL}/landfills/${landfillId.value}/survey`, {
      method: 'POST',
      body: formData,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }).then((res) => {
      if (!res.ok) throw new Error('HTTP error ' + res.status);
      mode.value = 'audit';
      return res.json();
    });
    landfillData.value.surveys.push(newSurvey);
    switchToAuditMode();
  } catch (error) {
    errors.upload = true;
    console.error('Error: ', error);
  }

  fetching.upload = false;
}

const criteria = ref([]);
function switchToAuditMode() {
  const aiGeneratedAudits = landfillData.value.surveys
    .at(-1)
    .audits.filter((audit) => audit.ai_generated_status);
  criteria.value = aiGeneratedAudits.violations;
  criteria.value.map((criterion) => (criterion.isAudited = false));
  mode.value = 'audit';
}
function switchAudited(criterion) {
  criterion.isAudited = !criterion.isAudited;
}
function switchViolation(criterion) {
  criterion.isViolated = !criterion.isViolated;
}
async function sendAuditData() {
  fetching.audit = true;
  errors.audit = false;

  const auditData = criteria.value.map(({ id, status }) => {
    id, status;
  });
  const sessionId = useGetCookie('reoLandfillsSessionId');

  try {
    const newSurvey = await fetch(`${API_BASE_URL}/landfills/audit/`, {
      method: 'POST',
      body: JSON.stringify(auditData),
      headers: {
        'Content-Type': 'application/json',
        'X-Session-ID': sessionId,
      },
    }).then((res) => {
      if (!res.ok) throw new Error('HTTP error ' + res.status);
      mode.value = 'audit';
      return res.json();
    });
    landfillData.value.surveys.push(newSurvey);
    switchToAuditMode();
  } catch (error) {
    errors.audit = true;
    console.error('Error: ', error);
  }

  fetching.upload = false;
  mode.value = 'view';
}
</script>

<template>
  <main class="layout landfill">
    <template v-if="landfillData">
      <article class="landfill__info">
        <h1>Сводка по полигону №{{ landfillData.id }}</h1>
        <h2>{{ landfillData.type }}</h2>
        <BaseDictionary :items="landfillDictionaryItems"></BaseDictionary>
        <BaseDictionary :items="auditDictionaryItems"></BaseDictionary>
        <section class="landfill__preview" style="height: 250px">
          <h3 class="landfill__preview-title">Предпросмотр на карте и со спутника</h3>
          <BaseMap
            style="width: 48%"
            :map-center="landfillData.coordinates"
            class="landfill__ymap-preview"
          ></BaseMap>
          <img
            v-if="landfillData.preview_image_path"
            class="landfill__satellite-preview"
            :src="landfillData.preview_image_path"
          />
        </section>
      </article>
      <section class="landfill__images images">
        <template v-if="mode === 'upload' || !surveys.length">
          <div v-if="fetching.upload" class="info-message">Подождите, идёт проверка...</div>
          <FileDrop v-else class="images__main" @update-files="updateImages"></FileDrop>
          <button class="btn btn_primary" @click="uploadImages">Найти нарушения</button>
          <div v-if="errors.upload" class="error-message">Ошибка при отправке</div>
          <button class="btn btn_danger" @click="mode = 'view'">Отменить загрузку</button>
        </template>
        <template v-if="mode !== 'upload' && surveys.length">
          <img class="images__main" :src="images[mainImageIndex]" alt="" />
          <ul class="images__list">
            <li v-for="(image, index) of images" :key="image" class="images__list-item">
              <a href="#" @click.prevent="setMainImage(index)">
                <img class="images__img" :src="image" alt="landfill photo" />
              </a>
            </li>
          </ul>
        </template>
      </section>
      <section v-if="mode === 'audit'" class="landfill__criteria criteria">
        <h2 class="criteria__title">Проверка по критериям</h2>
        <ul class="criteria__list">
          <li v-for="criterion of criteria" :key="criterion.id" class="criteria__item">
            <div>
              <img :src="criterion.isViolated ? '/icons/alert-circle.svg' : '/icons/audit.svg'" />
              <p class="criteria__text" :class="{ criteria__text_audited: criterion.isAudited }">
                {{ criterion.text }}
              </p>
            </div>
            <div class="criteria__buttons criteria-buttons">
              <button
                class="criteria-buttons__approve btn btn_primary btn_small"
                @click="switchAudited(criterion)"
              >
                {{ !criterion.isAudited ? 'Подтвердить' : 'Отменить' }}
              </button>
              <button
                v-if="!criterion.isAudited"
                class="criteria-buttons__reject btn btn_danger btn_small"
                @click="switchViolation(criterion)"
              >
                Изменить
              </button>
            </div>
          </li>
        </ul>
        <button v-if="isAuditComplete" class="btn btn_primary-outlined" @click="sendAuditData">
          Завершить проверку
        </button>
      </section>
      <section class="landfill__violations violations">
        <template v-if="lastAuditViolations.length > 0">
          <h2 class="violations__title">
            Статус: Полигон с {{ lastAuditViolations.length }} нарушениями
          </h2>
          <ul class="violations__list">
            <li
              v-for="violation of lastAuditViolations"
              :key="violation.id"
              class="violations__item"
            >
              {{ violation.title }}
            </li>
          </ul>
        </template>
        <h2 v-else class="violations__title">Статус: Полигон без нарушений</h2>
        <button class="btn btn_primary" @click="generateExcelFile()">
          Скачать отчет по полигону
        </button>
      </section>
      <div class="landfill__surveys surveys">
        <h2 class="surveys__title">Съёмки по датам</h2>
        <ul class="surveys__list">
          <li v-for="(survey, index) of surveys" :key="survey.id" class="surveys__item">
            <button
              class="surveys__item-btn btn btn_small btn_info-outlined"
              @click="setViewedSurvey(index)"
            >
              {{ useFormatDate(survey.date) }}
            </button>
          </li>
        </ul>
        <div class="input-group">
          <input
            v-model="newSurveyDate"
            :disabled="mode === 'upload'"
            class="form-input"
            type="date"
          />
          <button class="surveys__load btn btn_primary btn_small" @click="switchToUploadMode">
            Загрузить новые снимки
          </button>
        </div>
      </div>
      <section v-if="surveys.length" class="landfill__audits">
        <h2 class="landfill__audits-title">История проверок по данной съёмке</h2>
        <ul class="landfill__audits-list">
          <li v-for="audit of audits" :key="audit.id" class="landfill__audits-item">
            <h4>Проверка</h4>
            <span><strong>Дата: </strong>{{ useFormatDate(audit.date) }}</span>
            <span><strong>ФИО: </strong>{{ audit.auditor.name }}</span>
            <span><strong>Должность: </strong>{{ audit.auditor.position }}</span>
            <span>
              <strong>Статус: </strong> обнаружено
              {{ audit.violations.filter((violation) => violation.status).length }} нарушений
            </span>
          </li>
        </ul>
      </section>
    </template>
  </main>
</template>

<style lang="scss" scoped>
.landfill {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 64px;

  &__preview {
    padding-top: 18px;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
  }

  &__preview-title {
    width: 100%;
    margin-bottom: 18px;
  }

  &__ymap-preview {
    border-radius: 20px;
    overflow: hidden;
  }

  &__satellite-preview {
    object-fit: fill;
    width: 48%;
    height: 240px;
    border-radius: 20px;

    transition: transform 0.3s ease-in-out;

    &:hover {
      transform: scale(3);
    }
  }

  &__audits {
    grid-column: 1 / span 2;
    margin-top: 30px;
  }

  &__audits-list {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    margin-top: 18px;
  }

  &__audits-item {
    display: flex;
    flex-direction: column;
    padding: 15px;
    background: #f7f7f7;
    border-radius: 8px;

    &:not(:last-child) {
      margin-right: 20px;
    }

    span {
      margin-top: 10px;
    }
  }

  &__surveys {
    grid-column: 1 / span 2;
  }

  &__images {
    grid-row: span 2;
  }

  &__violations {
    grid-column: 1;
  }
}

.images {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 18px;

  &__main {
    height: 600px;
    width: 100%;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    margin-top: 64px;
    border-radius: 20px;
  }

  &__list {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    justify-content: space-between;

    .images__list-item img {
      width: 120px;
      height: 120px;
    }

    .images__list-item:not(:last-child) {
      margin-right: 20px;
    }
  }

  &__img {
    transition: transform 0.3s ease-in-out;
    border-radius: 10px;

    &:hover {
      transform: scale(1.2);
    }
  }
}

.surveys {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 18px;

  &__title {
    width: 100%;
    margin-bottom: 18px;
  }

  &__list {
    display: flex;
    align-items: center;
    gap: inherit;
  }
}

.violations {
  &__title {
    margin-bottom: 18px;
  }

  &__list {
    margin-bottom: 30px;
  }

  &__item {
    font-weight: 400;
    font-size: 14px;
    line-height: 18px;
    color: #ffffff;
    padding: 4px 10px;
    background: #ff4141;
    border-radius: 50px;
    margin-bottom: 10px;
    width: fit-content;
  }
}

.criteria {
  &__list {
    margin-bottom: 30px;
  }

  &__item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 10px;
    font-weight: 400;
    font-size: 14px;
    line-height: 18px;

    img {
      margin-right: 10px;
    }

    div:first-child {
      display: flex;
    }
  }

  &__text {
    &_audited {
      font-weight: bold;
    }
  }
}

.criteria-buttons {
  display: flex;
  align-items: center;
  gap: 24px;
}
</style>
