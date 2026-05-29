<template>
  <label class="platform-card" :class="{ selected: checked }">
    <input :checked="checked" type="checkbox" @change="$emit('toggle', platform.id)" />
    <div class="platform-card__header">
      <strong>{{ platform.name }}</strong>
      <span>{{ platform.id }}</span>
    </div>
    <p>{{ platform.description }}</p>
    <div class="platform-card__chips">
      <span v-for="hint in platform.styleHints" :key="hint">{{ hint }}</span>
    </div>
  </label>
</template>

<script setup>
defineProps({
  platform: {
    type: Object,
    required: true,
  },
  checked: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['toggle'])
</script>
