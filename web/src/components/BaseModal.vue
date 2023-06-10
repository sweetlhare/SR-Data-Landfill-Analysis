<script setup>
defineProps({ isModalShown: { type: Boolean, default: false } });
defineEmits(['hideModal']);
</script>

<template>
  <transition name="modal">
    <div v-show="isModalShown" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header section-title">
            <slot name="header"> default header </slot>
            <button class="modal-close" @click="$emit('hideModal')">
              <img src="/icons/sharp.svg" alt="Скрыть боковую панель" />
            </button>
          </div>
          <div class="modal-body">
            <slot name="body"> default body </slot>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<style lang="scss">
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: table;
}

.modal-wrapper {
  display: table-cell;
  vertical-align: middle;
}

.modal-container {
  width: 600px;
  margin: 0px auto;
  padding: 40px 60px;
  background-color: #fff;
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  font-family: Helvetica, Arial, sans-serif;
}

.modal-header {
  position: relative;
}

.modal-close {
  position: absolute;
  top: 0;
  right: 0;
}

.modal-body {
  margin-bottom: 20px;
}

.modal-default-button {
  display: block;
  margin-top: 1rem;
}

/*
 * The following styles are auto-applied to elements with
 * transition="modal" when their visibility is toggled
 * by Vue.js.
 *
 * You can easily play with the modal transition by editing
 * these styles.
 */

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.5s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
