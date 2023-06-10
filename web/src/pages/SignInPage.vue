<script setup>
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { API_BASE_URL } from '../config';
import { useSetCookie } from '../composables/useSetCookie';

const router = useRouter();

const authError = ref(false);
const formData = reactive({
  email: '',
  password: '',
});

async function onSubmit() {
  try {
    const parsedRes = await fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      body: JSON.stringify(formData),
      headers: {
        'Content-Type': 'application/json',
      },
    }).then((res) => {
      if (!res.ok) {
        if (res.status === 401) router.push({ name: 'sign-in' });
        throw new Error('HTTP error ' + res.status);
      }
      return res.json();
    });
    const sessionId = parsedRes.session_id;
    useSetCookie('reoLandfillsSessionId', sessionId, { path: '/' });
    router.push({ name: 'region', params: { id: 30 } });
  } catch (error) {
    authError.value = true;
    console.error('Error: ', error);
  }
}
</script>

<template>
  <div class="page__layout">
    <div class="layout__flex">
      <img src="/icons/logo.svg" alt="" />
      <form class="sign-in-form" @submit.prevent="onSubmit">
        <h1 class="heading1">Вход</h1>
        <div class="input-block">
          <label class="input-block__label" for="email">Email</label>
          <input
            id="email"
            v-model="formData.email"
            class="input-block__input"
            type="email"
            placeholder="Введите фамилию, имя и отчество"
          />
        </div>
        <div class="input-block">
          <label class="input-block__label" for="password">Пароль</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            class="input-block__input"
            placeholder="Введите пароль"
          />
        </div>
        <button type="submit" class="btn btn_primary">Войти</button>
        <div v-if="authError" class="error-message">Ошибка при авторизации</div>
      </form>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.page__layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 200px;

  .layout__banner {
    img {
      height: 100vh;
    }
  }
}
.layout__flex {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  margin-top: 25px;
  margin-left: 55px;
  padding-left: 15px;

  h1 {
    margin-top: 102px;
    margin-bottom: 22px;
  }
}

input {
  width: 359px;
  padding-right: 20px;
}

.input-block:not(:last-child) {
  margin-bottom: 20px;
}
</style>
