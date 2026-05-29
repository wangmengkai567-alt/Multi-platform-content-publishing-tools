<template>
  <article class="preview-panel">
    <header class="preview-panel__header">
      <div>
        <h3>{{ preview.platformName }}</h3>
        <span>{{ preview.platformId }}</span>
      </div>
      <div class="preview-panel__tags">
        <span v-for="tag in preview.recommendedTags" :key="tag">#{{ tag }}</span>
      </div>
    </header>

    <h4>{{ preview.title }}</h4>
    <pre>{{ preview.formattedBody }}</pre>

    <section v-if="preview.notes?.length">
      <strong>适配建议</strong>
      <ul>
        <li v-for="note in preview.notes" :key="note">{{ note }}</li>
      </ul>
    </section>

    <section v-if="preview.warnings?.length" class="preview-panel__warnings">
      <strong>注意</strong>
      <ul>
        <li v-for="warning in preview.warnings" :key="warning">{{ warning }}</li>
      </ul>
    </section>
  </article>
</template>

<script setup>
defineProps({
  preview: {
    type: Object,
    required: true,
  },
})
</script>
