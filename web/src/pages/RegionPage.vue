<script setup>
import { useRoute } from 'vue-router';
import { reactive, ref, watchEffect } from 'vue';
import { API_BASE_URL } from '../config';
import BaseModal from '../components/BaseModal.vue';

const route = useRoute();
const showModal = ref(false);
const modalFormState = reactive({
  type: '',
  city: '',
  address: '',
  latitude: '',
  longitude: '',
  coordinates: [],
  illegal: false,
});
const fetching = reactive({
  landfills: false,
  newLandfill: false,
});
const errors = reactive({
  landfills: false,
  newLandfill: false,
});

watchEffect(() => {
  modalFormState.coordinates = [modalFormState.latitude, modalFormState.longitude];
});

const landfills = ref([]);
const cleanLandfills = ref([]);
const dirtyLandfills = ref([]);
const illegalLandfills = ref([]);
watchEffect(() => {
  illegalLandfills.value = [];
  cleanLandfills.value = [];
  dirtyLandfills.value = [];

  landfills.value.forEach((landfill) => {
    if (landfill.illegal) illegalLandfills.value.push(landfill);
    else if (landfill.violations_count) dirtyLandfills.value.push(landfill);
    else cleanLandfills.value.push(landfill);
  });
});

getLandfills();

async function getLandfills() {
  const parsedRes = await fetch(`${API_BASE_URL}/landfills/?regionID=${route.params.id}`).then(
    (response) => response.json()
  );

  landfills.value = parsedRes;
}

async function sendNewLandfill() {
  fetching.newLandfill = true;
  errors.newLandfill = false;

  const newLandfillData = (({ type, city, address, coordinates, illegal }) => ({
    type,
    city,
    address,
    coordinates,
    illegal,
  }))(modalFormState);

  try {
    await fetch(`${API_BASE_URL}/landfills/`, {
      method: 'POST',
      body: JSON.stringify(newLandfillData),
      headers: {
        'Content-Type': 'application/json',
      },
    }).then((res) => {
      if (!res.ok) throw new Error('HTTP error ' + res.status);
      return res.json();
    });

    showModal.value = false;
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
    <!-- <del>
      <ul class="region__search search">
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по регионам</label>
            <input id="name" class="input-block__input" placeholder="Введите регион" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по городам</label>
            <input id="name" class="input-block__input" placeholder="Введите город" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по адресу</label>
            <input id="name" class="input-block__input" placeholder="Введите адрес" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по координатам</label>
            <input id="name" class="input-block__input" placeholder="55°45′07″ N, 37°36′56″  E" />
          </div>
        </li>
        <li class="search__item">
          <div class="search__input-block modal__input-block input-block">
            <label class="input-block__label" for="name">Поиск по названию полигона</label>
            <input id="name" class="input-block__input" placeholder="Введите название полигона" />
          </div>
        </li>
      </ul>
    </del> -->

    <div class="region__nav-buttons">
      <button disabled class="btn btn_primary">Показать полигоны на карте</button>
      <button disabled class="btn btn_primary-outlined">CRM для полигонов РФ</button>
      <button id="show-modal" class="btn btn_primary" @click="showModal = true">
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
              <h4 class="dashboard-column__item-title h4">{{ landfill.type }}</h4>
              <span><strong>Адрес:</strong> {{ landfill.city + ', ' + landfill.address }}</span>
              <span><strong>Координаты:</strong> {{ landfill.coordinates.join(', ') }}</span>
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
              <h4 class="dashboard-column__item-title h4">{{ landfill.type }}</h4>
              <span><strong>Адрес:</strong> {{ landfill.city + ', ' + landfill.address }}</span>
              <span><strong>Координаты:</strong> {{ landfill.coordinates.join(', ') }}</span>
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
          <li v-if="landfill.isIllegal" class="dashboard-column__item">
            <RouterLink
              class="dashboard-column__link red"
              :to="{ name: 'landfill', params: { id: landfill.id } }"
            >
              <h4 class="dashboard-column__item-title h4">{{ landfill.type }}</h4>
              <span><strong>Адрес:</strong> {{ landfill.city + ', ' + landfill.address }}</span>
              <span><strong>Координаты:</strong> {{ landfill.coordinates.join(', ') }}</span>
            </RouterLink>
          </li>
        </ul>
      </div>
    </section>
  </main>
  <transition name="modal">
    <BaseModal v-if="showModal" class="modal" @close="showModal = false">
      <template #header>
        <h2 class="h2">Добавить полигон</h2>
      </template>
      <template #body>
        <form class="modal__form">
          <div class="modal__input-block input-block">
            <label class="input-block__label" for="type">Наименование полигона</label>
            <input
              id="type"
              v-model="modalFormState.type"
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
          <div class="modal__input-block input-block">
            <label class="input-block__label">Координаты</label>
            <input
              v-model="modalFormState.latitude"
              class="input-block__input"
              placeholder="Широта в формате XX.XXXXX"
            />
            <input
              v-model="modalFormState.longitude"
              class="input-block__input"
              placeholder="Долгота в формате XX.XXXXX"
            />
          </div>
          <button type="submit" class="btn btn_primary" @click="sendNewLandfill">Добавить</button>
          <button type="button" class="btn btn_danger" @click="showModal = false">Отмена</button>
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
  justify-content: center;

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

.modal__input-block {
  margin-bottom: 18px;
}
</style>
