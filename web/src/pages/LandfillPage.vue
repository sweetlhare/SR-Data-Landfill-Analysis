<script setup>
import { ref, computed, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import * as XLSX from 'xlsx';

import BaseDictionary from '../components/BaseDictionary.vue';
import BaseMap from '../components/BaseMap.vue';
import BaseTable from '../components/BaseTable.vue';
import UploadSidebar from '../components/UploadSidebar.vue';

import { API_BASE_URL } from '../config';
import { useFormatDate } from '../composables/useFormatDate';
import { useGetSessionId } from '../composables/useGetSessionId';
import BaseBreadcrumbs from '../components/BaseBreadcrumbs.vue';

const router = useRouter();
const route = useRoute();

const landfillData = ref(null);
const landfillId = computed(() => landfillData.value?.id ?? null);
const fetching = reactive({
  landfillData: false,
  uploadImages: false,
});
const errors = reactive({
  landfillData: false,
  uploadImages: false,
  excel: false,
});

getLandfillData();
async function getLandfillData() {
  fetching.landfillData = true;
  errors.landfillData = false;
  const sessionId = useGetSessionId();

  try {
    const newData = await fetch(`${API_BASE_URL}/landfills/${route.params.id}`, {
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
    // const newData = (await import('../fixtures/landFillData')).default;

    landfillData.value = newData;
  } catch (error) {
    errors.landfillData = true;
    console.error('Error: ', error);
  }

  fetching.landfillData = false;
}

const surveys = computed(() =>
  landfillData.value ? landfillData.value.surveys?.toReversed() : []
);
// todo убрать ?? когда будет гарантированный массив

const criteria = computed(() => (landfillData.value ? landfillData.value.violations : []));
const violationsCount = computed(() => landfillData.value.violations_count);

const landfillDictionaryItems = computed(() =>
  landfillData.value
    ? [
        {
          key: 'Наименование полигона',
          value: landfillData.value.name,
        },
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
          value: landfillData.value.coordinates,
        },
        // todo
        //  {
        //   key: 'ФИО ответственного за полигон',
        //   value: landfillData.value.manager?.name ?? 'не указаны',
        // },
        // {
        //   key: 'Должность ответственного за полигон',
        //   value: landfillData.value.manager?.position ?? 'не указана',
        // },
        // {
        //   key: 'Контакты ответственного за полигон',
        //   value: landfillData.value.manager
        //     ? landfillData.value.manager.phone + ', ' + landfillData.value.manager.email
        //     : 'не указаны',
        // },
      ]
    : []
);

const breadcrumbItems = reactive([
  {
    title: 'Карточки полигонов',
    routeName: 'region',
    params: computed(() => ({
      id: landfillData.value?.region.id,
    })),
  },
  {
    title: computed(() => 'Полигон №' + landfillData.value?.id),
    routeName: '',
  },
]);

const audits = ref([]);
const lastAudit = computed(() => audits.value.at(-1));
async function getAudits() {
  const sessionId = useGetSessionId();
  const lastSurvey = await fetch(`${API_BASE_URL}/surveys/${surveys.value[0].id}`, {
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
  audits.value = lastSurvey.audits;
}

// Скачать отчет в excel
async function generateExcelFile() {
  await getAudits();
  if (!lastAudit.value) {
    errors.excel = true;
    setTimeout(() => {
      errors.excel = false;
    }, 2000);
    return;
  }

  const dataToExport = [
    ['Идентификатор полигона', String(landfillData.value.id)],
    [
      'Фактический адрес',
      // landfillData.value.city + ', ' +
      landfillData.value.address,
    ],
    ['Дата проверки', useFormatDate(lastAudit.value.date)],
    ['ФИО проверившего', lastAudit.value.author.name],
    ['Должность проверившего', lastAudit.value.author.position],
    ['Нарушения:', ''],
  ];

  for (const violation of lastAudit.value.violations) {
    dataToExport.push([violation.status ? 'Обнаружено' : 'Не обнаружено', violation.title]);
  }

  const ws = XLSX.utils.json_to_sheet(dataToExport);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, 'Отчет по полигону');
  const reportDate = useFormatDate(lastAudit.value.date);
  XLSX.writeFile(wb, `Отчет по полигону №${landfillData.value.id} ${reportDate}.xlsx`);
}

// Данные для таблицы съёмок
const surveysTableColumns = reactive([
  { title: '№', field: 'id' },
  { title: 'Дата съёмки', field: 'date', template: 'date' },
  { title: 'Загрузил снимки', field: 'author', template: 'author' },
  { title: 'Нарушения', field: 'violations_count', template: 'violations' },
  { title: 'Результаты', field: 'id', template: 'results' },
]);

// Новая съёмка
const isSidebarShown = ref(false);
</script>

<template>
  <main class="layout landfill">
    <div v-if="fetching.landfillData" class="info-message">Подождите, идёт загрузка...</div>
    <div v-else-if="errors.landfillData" class="error-message">
      Ошибка при загрузке. Перезагрузите страницу или попробуйте позднее
    </div>
    <template v-else>
      <BaseBreadcrumbs class="landfill__breadcrumbs" :items="breadcrumbItems"></BaseBreadcrumbs>
      <article class="landfill__info">
        <h1 class="page-title">Сводка по полигону №{{ landfillData.id }}</h1>
        <!-- <h2 class="page-subtitle">{{ landfillData.name }}</h2> -->
        <BaseDictionary :items="landfillDictionaryItems"></BaseDictionary>
        <!-- <BaseDictionary :items="auditDictionaryItems"></BaseDictionary> -->
        <section class="landfill__map">
          <h3 class="landfill__map-title">Полигон на карте</h3>
          <BaseMap
            style="width: 100%; height: 300px"
            :map-center="landfillData.coordinates.split(', ')"
            class="landfill__map-canvas"
          ></BaseMap>
        </section>
        <section class="landfill__criteria criteria">
          <h2 class="criteria__title section-title">
            {{
              'Статус: Полигон ' +
              (violationsCount ? `с ${violationsCount} нарушениями` : 'без нарушений')
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
            </li>
          </ul>
          <button class="btn btn_primary" @click="generateExcelFile()">
            Скачать отчет по полигону
          </button>
          <div v-if="errors.excel" class="error-message">
            Отсутствуют проверки для составления отчёта
          </div>
        </section>
      </article>
      <section class="landfill__images images">
        <h3 class="images__title section-title">Спутниковый снимок полигона</h3>
        <img
          v-if="landfillData.preview_image_path"
          class="images__main"
          :src="'backend' + landfillData.preview_image_path"
        />
        <!-- fixture -->
        <!-- :src="landfillData.preview_image_path" -->
      </section>
      <div class="landfill__surveys surveys">
        <h3 class="section-title surveys__title">История съёмок полигона</h3>
        <button class="btn btn_primary btn_small" @click="isSidebarShown = true">
          Загрузить новые снимки
        </button>
        <BaseTable
          v-if="surveys.length"
          class="surveys__table"
          :rows="surveys"
          :columns="surveysTableColumns"
        >
          <template #date="dateProps">
            {{ useFormatDate(dateProps.value) }}
          </template>
          <template #violations="violationsProps">
            {{ violationsProps.value }}
          </template>
          <template #results="resultsProps">
            <RouterLink
              class="btn btn_primary-outlined btn_small"
              :to="{ name: 'survey', params: { id: resultsProps.value } }"
            >
              Просмотреть
            </RouterLink>
          </template>
          <template #author="authorProps">
            {{ authorProps.value.name }}
          </template>
        </BaseTable>
      </div>
    </template>
  </main>
  <UploadSidebar
    :is-sidebar-shown="isSidebarShown"
    :url="`${API_BASE_URL}/landfills/${landfillId}/survey`"
    @hide-sidebar="isSidebarShown = false"
  />
</template>

<style lang="scss" scoped>
@import '../assets/scss/landfill';
@import '../assets/scss/images';
@import '../assets/scss/criteria';

.surveys {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.surveys__title {
  margin-bottom: 0;
}

.surveys__table {
  width: 100%;
}
</style>
