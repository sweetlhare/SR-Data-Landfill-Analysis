<script setup>
import { useRoute, useRouter } from 'vue-router';
import { reactive, ref, unref, watchEffect } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { API_BASE_URL } from '../config';
import BaseModal from '../components/BaseModal.vue';
import { useGetSessionId } from '../composables/useGetSessionId';
import { useCreateRules } from '../composables/useCreateRules';

const router = useRouter();
const route = useRoute();

const isModalShown = ref(false);
const modalFormState = reactive({
  name: '',
  city: '',
  address: '',
  latitude: '',
  longitude: '',
  coordinates: [],
  illegal: false,
});
const rules = useCreateRules();
const v$ = useVuelidate(rules, modalFormState);

const fetching = reactive({
  landfills: false,
  newLandfill: false,
});
const errors = reactive({
  landfills: false,
  newLandfill: false,
});

watchEffect(() => {
  modalFormState.coordinates = modalFormState.latitude + ', ' + modalFormState.longitude;
});

const landfills = ref([]);
const cleanLandfills = ref([]);
const dirtyLandfills = ref([]);
const illegalLandfills = ref([]);
watchEffect(() => {
  illegalLandfills.value = [];
  cleanLandfills.value = [];
  dirtyLandfills.value = [];

  if (Array.isArray(landfills.value))
    landfills.value.forEach((landfill) => {
      if (landfill.illegal) illegalLandfills.value.push(landfill);
      else if (landfill.violations_count) dirtyLandfills.value.push(landfill);
      else cleanLandfills.value.push(landfill);
    });
});

getLandfills();

async function getLandfills() {
  const sessionId = useGetSessionId();
  try {
    const parsedRes = await fetch(`${API_BASE_URL}/landfills?regionID=${route.params.id}`, {
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

    landfills.value = parsedRes;
  } catch (error) {
    errors.landfills = true;
    console.error('Error: ', error);
  }
}

async function sendNewLandfill() {
  const isFormCorrect = await unref(v$).$validate();
  if (!isFormCorrect) return;

  fetching.newLandfill = true;
  errors.newLandfill = false;

  const newLandfillData = (({ name, city, address, coordinates, illegal }) => ({
    region_id: +route.params.id,
    name,
    city,
    address,
    coordinates,
    illegal,
  }))(modalFormState);
  const sessionId = useGetSessionId();

  try {
    await fetch(`${API_BASE_URL}/landfills`, {
      method: 'POST',
      body: JSON.stringify(newLandfillData),
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

    isModalShown.value = false;
    getLandfills();
  } catch (error) {
    errors.newLandfill = true;
    console.error('Error: ', error);
  }

  fetching.newLandfill = false;
}
</script>

<template>
  <main class="layout region">
    <h2 class="h2">Список полигонов региона</h2>
    <div v-if="fetching.landfills" class="info-message">Подождите, идёт загрузка...</div>
    <div v-else-if="errors.landfills" class="error-message">
      Ошибка при загрузке. Перезагрузите страницу или попробуйте позднее
    </div>
    <template v-else>
      <ul class="region__search search">
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="search-city">Поиск по населенным пунктам</label>
            <input id="search-city" class="input-block__input" placeholder="Введите название" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по широте</label>
            <input id="name" class="input-block__input" placeholder="55°45′07″ N, 37°36′56″  E" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по долготе</label>
            <input id="name" class="input-block__input" placeholder="55°45′07″ N, 37°36′56″  E" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по наименованию полигона</label>
            <input
              id="name"
              class="input-block__input"
              placeholder="Введите наименование полигона"
            />
          </div>
        </li>
      </ul>
      <div class="region__nav-buttons">
        <button disabled class="btn btn_primary">Показать полигоны на карте</button>
        <button disabled class="btn btn_primary-outlined">CRM для полигонов РФ</button>
        <button id="show-modal" class="btn btn_primary" @click="isModalShown = true">
          Добавить полигон
        </button>
      </div>
      <section class="region__dashboard dashboard">
        <div class="dashboard__column dashboard-column">
          <h3 class="dashboard-column__title h3">Полигоны без нарушений</h3>
          <ul
            v-for="landfill of cleanLandfills"
            :key="landfill.id"
            class="dashboard-column__list dashboard-column__list_normal"
          >
            <li class="dashboard-column__item">
              <RouterLink
                class="dashboard-column__link green"
                :to="{ name: 'landfill', params: { id: landfill.id } }"
              >
                <h4 class="dashboard-column__item-title h4">
                  {{ '№' + landfill.id + ' ' + landfill.name }}
                </h4>
                <span
                  ><strong>Адрес:</strong>
                  {{
                    // landfill.city + ', ' +
                    landfill.address || 'не указан'
                  }}</span
                >
                <span><strong>Координаты:</strong> {{ landfill.coordinates }}</span>
              </RouterLink>
            </li>
          </ul>
        </div>
        <div class="dashboard__column dashboard-column">
          <h3 class="dashboard-column__title h3">Полигоны с нарушениями</h3>
          <ul
            v-for="landfill of dirtyLandfills"
            :key="landfill.id"
            class="dashboard-column__list dashboard-column__list_violated"
          >
            <li class="dashboard-column__item">
              <RouterLink
                class="dashboard-column__link yellow"
                :to="{ name: 'landfill', params: { id: landfill.id } }"
              >
                <h4 class="dashboard-column__item-title h4">
                  {{ '№' + landfill.id + ' ' + landfill.name }}
                </h4>
                <span
                  ><strong>Адрес:</strong>
                  {{
                    // landfill.city + ', ' +
                    landfill.address
                  }}</span
                >
                <span><strong>Координаты:</strong> {{ landfill.coordinates }}</span>
                <span><strong>Нарушения:</strong> {{ landfill.violations_count }}</span>
              </RouterLink>
            </li>
          </ul>
        </div>
        <div class="dashboard__column dashboard-column">
          <h3 class="dashboard-column__title h3">Незаконные свалки</h3>
          <ul
            v-for="landfill of illegalLandfills"
            :key="landfill.id"
            class="dashboard-column__list dashboard-column__list_illegal"
          >
            <li class="dashboard-column__item">
              <RouterLink
                class="dashboard-column__link red"
                :to="{ name: 'landfill', params: { id: landfill.id } }"
              >
                <h4 class="dashboard-column__item-title h4">
                  {{ '№' + landfill.id + ' ' + landfill.name }}
                </h4>
                <span
                  ><strong>Адрес:</strong>
                  {{
                    // landfill.city + ', ' +
                    landfill.address
                  }}</span
                >
                <span><strong>Координаты:</strong> {{ landfill.coordinates }}</span>
              </RouterLink>
            </li>
          </ul>
        </div>
      </section>
    </template>
  </main>
  <transition name="modal">
    <BaseModal :is-modal-shown="isModalShown" class="modal" @hide-modal="isModalShown = false">
      <template #header>
        <h2 class="h2">Добавить полигон</h2>
      </template>
      <template #body>
        <form class="modal__form" @submit.prevent="sendNewLandfill">
          <div class="modal__input-block input-block">
            <label class="input-block__label" for="name">Наименование полигона</label>
            <input
              id="name"
              v-model="modalFormState.name"
              class="input-block__input"
              placeholder="Введите наименование полигона"
            />
          </div>
          <div class="modal__input-block input-block">
            <label class="input-block__label" for="city">Город</label>
            <input
              id="city"
              v-model="modalFormState.city"
              class="input-block__input"
              placeholder="Введите город"
            />
          </div>
          <div class="modal__input-block input-block">
            <label class="input-block__label" for="address">Фактический адрес</label>
            <input
              id="address"
              v-model="modalFormState.address"
              class="input-block__input"
              placeholder="Введите фактический адрес"
            />
          </div>
          <div class="modal__coordinates">
            <div class="modal__input-block input-block">
              <label class="input-block__label" for="latitude">Широта</label>
              <input
                id="latitude"
                v-model="modalFormState.latitude"
                class="input-block__input"
                placeholder="Широта в формате XX.XXXXX"
              />
            </div>
            <div class="modal__input-block input-block">
              <label for="longitude" class="input-block__label">Долгота</label>
              <input
                id="longitude"
                v-model="modalFormState.longitude"
                class="input-block__input"
                placeholder="Долгота в формате XX.XXXXX"
              />
            </div>
          </div>
          <div class="modal__input-block input-block">
            <input
              id="illegal"
              v-model="modalFormState.illegal"
              type="checkbox"
              class="input-block__checkbox"
              placeholder="Введите фактический адрес"
            /><label class="input-block__label input-block__label_checkbox" for="illegal"
              >Незаконная свалка</label
            >
          </div>
          <button type="submit" class="modal__button btn btn_primary">Добавить</button>
          <button type="button" class="modal__button btn btn_danger" @click="isModalShown = false">
            Отмена
          </button>
          <div v-if="v$.$errors.length" class="error-message">
            Все поля обязательны для заполнения
          </div>
          <div v-if="fetching.newLandfill" class="info-message">Подождите, идёт добавление...</div>
          <div v-else-if="errors.newLandfill" class="error-message">Ошибка при добавлении.</div>
        </form>
      </template>
    </BaseModal>
  </transition>
</template>

<style lang="scss" scoped>
.region__nav-buttons {
  margin-bottom: 36px;
  display: flex;
  align-items: center;

  button:not(:last-child) {
    margin-right: 48px;
  }
}

.search {
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  .input-block {
    width: 240px;

    &__input {
      border-radius: 0;
    }
  }

  &__item:not(:last-child) {
    margin-right: 24px;
  }

  &__item:first-child .input-block__input {
    border-top-left-radius: 12px;
    border-bottom-left-radius: 12px;
  }
  &__item:last-child .input-block__input {
    border-top-right-radius: 12px;
    border-bottom-right-radius: 12px;
  }
}
.region__dashboard {
  margin-bottom: 56px;
}
.dashboard {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 24px;
  color: var(--black);
}

.dashboard-column__title {
  text-align: center;
}

.dashboard-column__list {
  margin-top: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;

  li:not(:last-child) {
    margin-bottom: 11px;
  }
}

.dashboard-column__item {
  width: 360px;
  transition: transform 0.3s ease-in-out;

  &:hover {
    transform: scale(1.1);
  }
}

.dashboard-column__item-title {
  font-weight: 700;
}

.dashboard-column__link {
  padding: 28px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;

  span {
    margin-bottom: 8px;
  }
}
.h4 {
  margin-bottom: 8px;
}
.green {
  background: rgba(111, 204, 134, 0.3);
}
.red {
  background: rgba(255, 65, 65, 0.3);
}
.yellow {
  background: rgba(255, 190, 65, 0.3);
}

.modal__form {
  display: flex;
  flex-wrap: wrap;
  gap: 18px;
}

.modal__input-block {
  width: 100%;
}

.modal__coordinates {
  display: flex;
  gap: 18px;
}

.modal__button {
  margin-right: 36px;
}
</style>
