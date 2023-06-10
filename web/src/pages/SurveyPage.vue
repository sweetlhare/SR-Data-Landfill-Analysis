<script setup>
import { ref, computed, reactive, watchEffect, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { View360, CylindricalProjection } from '@egjs/vue3-view360';
import '@egjs/view360/css/view360.min.css';

import BaseDictionary from '../components/BaseDictionary.vue';

import { API_BASE_URL } from '../config';
import { useFormatDate } from '../composables/useFormatDate';
import { useGetSessionId } from '../composables/useGetSessionId';
import BaseTable from '../components/BaseTable.vue';
import BaseBreadcrumbs from '../components/BaseBreadcrumbs.vue';
import BaseModal from '../components/BaseModal.vue';

const router = useRouter();
const route = useRoute();

const surveyData = ref(null);
const mode = ref('view');

const fetching = reactive({
  get: false,
  audit: false,
});
const errors = reactive({
  get: false,
  audit: false,
});

getSurveyData();
async function getSurveyData() {
  fetching.get = true;
  errors.get = false;
  const sessionId = useGetSessionId();

  try {
    const newData = await fetch(`${API_BASE_URL}/surveys/${route.params.id}`, {
      headers: {
        'X-Session-ID': sessionId,
      },
    }).then((res) => {
      if (!res.ok) {
        if (res.status === 401) router.push({ name: 'sign-in' });
        throw new Error('HTTP error ' + res.status);
      }
      return res.json();
    });

    // fixture
    // const newData = (await import('../fixtures/surveyData')).default;

    surveyData.value = newData;
  } catch (error) {
    errors.get = true;
    console.error('Error: ', error);
  }

  fetching.get = false;
}
const images = computed(
  () => (surveyData.value ? surveyData.value.ai_images.map((path) => 'backend' + path) : [])

  // fixture
  // surveyData.value ? surveyData.value.ai_images : []
);
const mainImageIndex = ref(0);
function setMainImage(index) {
  mainImageIndex.value = index;
}
const projection = ref(null);
const rawImages = computed(
  () =>
    // todo уточнить расположение
    surveyData.value ? surveyData.value.ai_images.map((path) => 'backend' + path) : []
  // surveyData.value ? surveyData.value.raw_images.map((path) => 'backend' + path) : []
);
const rawProjection = ref(null);
const isRawImageShown = ref(false);
const isModalOpen = ref(false);

// todo убрать присвоение rawProjection
watchEffect(() => {
  projection.value = new CylindricalProjection({
    src: images.value[mainImageIndex.value],
    partial: true,
  });
});
watch(isRawImageShown, () => {
  nextTick(() => {
    isModalOpen.value = isRawImageShown.value;
  });
  rawProjection.value = new CylindricalProjection({
    src: rawImages.value[mainImageIndex.value],
    partial: true,
  });
});

const audits = computed(() => (surveyData.value ? surveyData.value.audits.toReversed() : []));
const viewedAuditIndex = ref(0);
const viewedAudit = computed(() => audits.value.at(viewedAuditIndex.value));
function setViewedAudit(index) {
  viewedAuditIndex.value = index;
}
const criteria = ref([]);
watchEffect(() => {
  criteria.value = viewedAudit.value?.violations ?? [];
});
const violatedCriteria = computed(() => criteria.value.filter((criterion) => criterion.status));

function resetCriteria() {
  const aiGeneratedAudits = surveyData.value.audits.filter((audit) => audit.ai_generated_status);
  criteria.value = aiGeneratedAudits.at(-1).violations;
  criteria.value.map((criterion) => (criterion.isAudited = false));
}
function switchToAuditMode() {
  resetCriteria();
  mode.value = 'audit';
}
function switchAudited(criterion) {
  criterion.isAudited = !criterion.isAudited;
}
function switchViolation(criterion) {
  criterion.status = !criterion.status;
}
const isAuditComplete = computed(
  () => criteria.value.filter((element) => element.isAudited === false)?.length === 0
);
async function sendAuditData() {
  fetching.audit = true;
  errors.audit = false;

  const violationsData = criteria.value.map(({ id, status }) => ({
    id,
    status,
  }));
  const auditData = {
    survey_id: surveyData.value.id,
    violations: violationsData,
  };
  const sessionId = useGetSessionId();

  try {
    const newAudit = await fetch(`${API_BASE_URL}/audits`, {
      method: 'POST',
      body: JSON.stringify(auditData),
      headers: {
        'Content-Type': 'application/json',
        'X-Session-ID': sessionId,
      },
    }).then((res) => {
      if (!res.ok) {
        if (res.status === 401) router.push({ name: 'sign-in' });
        throw new Error('HTTP error ' + res.status);
      }
      return res.json();
    });
    surveyData.value.audits.push(newAudit);
    mode.value = 'view';
  } catch (error) {
    errors.audit = true;
    console.error('Error: ', error);
  }

  fetching.upload = false;
}
function cancelAudit() {
  resetCriteria();
  mode.value = 'view';
}

const landfillDictionaryItems = computed(() =>
  surveyData.value
    ? [
        {
          key: 'Загрузивший съёмку',
          value: surveyData.value.author.name,
        },
        {
          key: 'Должность загрузившего съёмку',
          value: surveyData.value.author.position,
        },
      ]
    : []
);
const auditDictionaryItems = computed(() =>
  viewedAudit.value
    ? [
        {
          key: 'ФИО проверившего',
          value: viewedAudit.value.author.name,
        },
        {
          key: 'Должность проверившего',
          value: viewedAudit.value.author.position,
        },
      ]
    : []
);

const breadcrumbItems = reactive([
  {
    title: 'Карточки полигонов',
    routeName: 'region',
    params: computed(() => ({ id: surveyData.value?.region.id })),
  },
  {
    title: computed(() => 'Полигон №' + surveyData.value?.landfill_id),
    routeName: 'landfill',
    params: computed(() => ({ id: surveyData.value?.landfill_id })),
  },
  {
    title: computed(() => 'Съёмка №' + surveyData.value?.id),
    routeName: '',
  },
]);

// Данные для таблицы съёмок
const auditsTableColumns = reactive([
  { title: '№', field: 'id' },
  { title: 'Дата и время проверки', field: 'date', template: 'date' },
  { title: 'Проверяющий', field: 'author', template: 'author' },
  { title: 'Нарушения', field: 'true_violations_count' },
  { title: 'Результаты', field: 'index', template: 'results' },
]);
</script>

<template>
  <main class="layout landfill">
    <div v-if="fetching.get" class="info-message">Подождите, идёт загрузка...</div>
    <div v-else-if="errors.get" class="error-message">
      Ошибка при загрузке. Перезагрузите страницу или попробуйте позднее
    </div>
    <template v-else>
      <BaseBreadcrumbs class="landfill__breadcrumbs" :items="breadcrumbItems"></BaseBreadcrumbs>
      <article class="landfill__info">
        <h1 class="page-title">
          Съёмка полигона №{{ surveyData.landfill_id }} от
          {{ useFormatDate(new Date(surveyData.date)) }}
        </h1>
        <BaseDictionary :items="landfillDictionaryItems"></BaseDictionary>
        <h2 class="page-subtitle">
          Результаты {{ viewedAuditIndex === 0 ? 'последней' : 'выбранной' }} проверки<br />(№{{
            viewedAudit.id
          }}
          от {{ useFormatDate(new Date(viewedAudit.date), true, true) }})
        </h2>
        <BaseDictionary :items="auditDictionaryItems"></BaseDictionary>
        <section class="landfill__criteria criteria">
          <h2 v-if="mode === 'audit'" class="criteria__title section-title">Ручная проверка</h2>
          <h2 v-else class="criteria__title section-title">
            {{
              'Статус: Полигон ' +
              (violatedCriteria.length
                ? `с ${violatedCriteria.length} нарушениями`
                : 'без нарушений')
            }}
          </h2>
          <ul class="criteria__list">
            <li v-for="criterion of criteria" :key="criterion.id" class="criteria__item">
              <div>
                <img :src="criterion.status ? '/icons/alert-circle.svg' : '/icons/check.svg'" />
                <p class="criteria__text" :class="{ criteria__text_audited: criterion.status }">
                  {{ (criterion.status ? 'Обнаружено: ' : 'Не обнаружено: ') + criterion.title }}
                </p>
              </div>
              <div v-if="mode === 'audit'" class="criteria__buttons">
                <button
                  class="criteria__button btn btn_small"
                  :class="!criterion.isAudited ? 'btn_primary' : 'btn_danger'"
                  @click="switchAudited(criterion)"
                >
                  {{ !criterion.isAudited ? 'Подтвердить' : 'Отменить' }}
                </button>
                <button
                  v-if="!criterion.isAudited"
                  class="criteria__button btn btn_info-outlined btn_small"
                  @click="switchViolation(criterion)"
                >
                  Изменить
                </button>
              </div>
            </li>
          </ul>
          <template v-if="mode === 'audit'">
            <p class="criteria__tip tip">
              <em>
                Подтвердите финальные результаты по каждому критерию, чтобы завершить проверку
              </em>
            </p>
            <button :disabled="!isAuditComplete" class="btn btn_primary" @click="sendAuditData()">
              Завершить проверку
            </button>
            <button class="btn btn_danger" @click="cancelAudit()">Отменить проверку</button>
          </template>
        </section>
      </article>
      <section v-if="images.length" class="landfill__images images">
        <View360 class="images__main" :projection="projection" />
        <p class="images__tip tip">
          <em> Управляйте изображением, перетаскивая его и прокручивая колёсико мыши </em>
        </p>
        <button class="btn btn_primary-outlined btn_small" @click="isRawImageShown = true">
          Показать несжатую версию снимка
        </button>
        <h4 class="images__subtitle section-subtitle">Все снимки</h4>
        <ul class="images__list">
          <li v-for="(image, index) of images" :key="image" class="images__list-item">
            <a href="#" @click.prevent="setMainImage(index)">
              <img class="images__img" :src="image" alt="landfill photo" />
            </a>
          </li>
        </ul>
      </section>
      <section class="landfill__audits audits">
        <h3 class="section-title audits__title">История проверок снимков за данную дату</h3>
        <button class="audits__start-audit btn btn_primary btn_small" @click="switchToAuditMode">
          Новая проверка
        </button>
        <BaseTable class="audits__table" :rows="audits" :columns="auditsTableColumns">
          <template #date="dateProps">
            {{ useFormatDate(dateProps.value, true, true) }}
          </template>
          <template #results="resultsProps">
            <button
              class="btn btn_primary-outlined btn_small"
              @click="setViewedAudit(resultsProps.value)"
            >
              Просмотреть
            </button>
          </template>
          <template #author="authorProps">
            {{ authorProps.value.name }}
          </template>
        </BaseTable>
      </section>
    </template>
  </main>
  <BaseModal
    class="raw-image-modal"
    :is-modal-shown="isRawImageShown"
    @hide-modal="isRawImageShown = false"
  >
    <template #header>Просмотр снимка</template>
    <template #body>
      <View360 v-if="isModalOpen" class="raw-image" :projection="rawProjection" />
      <p class="images__tip tip">
        <em>Управляйте изображением, перетаскивая его и прокручивая колёсико мыши</em>
      </p>
    </template>
  </BaseModal>
</template>

<style lang="scss">
@import '../assets/scss/landfill';
@import '../assets/scss/images';
@import '../assets/scss/criteria';

.raw-image-modal {
  .modal-container {
    width: 80vw;
    min-height: 80vh;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .modal-body {
    display: flex;
    flex-direction: column;
  }
}

.raw-image {
  flex-grow: 1;
  flex-basis: 70vh;
  margin-bottom: 10px;
}

.audits {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.audits__title.section-title {
  margin-bottom: 0;
}

.audits__table {
  width: 100%;
}
</style>
