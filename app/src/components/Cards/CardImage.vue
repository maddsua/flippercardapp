<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { CardImageElement } from '../../content';
import LoadingMessage from '../App/LoadingMessage.vue';

const props = defineProps<{
	entry: CardImageElement;
}>();

const url = computed(() => props.entry.media_id ? `/media/images/${props.entry.media_id}` : null);

const randomRotation = computed(() => `${((Math.random() - 0.5) * 10).toFixed(0)}deg`);

const ready = ref(false);

watch(() => props.entry.media_id, () => ready.value = false);

</script>

<template>
	<div class="card-image" :style="{ rotate: randomRotation }">

		<img v-if="url" :src="url" @load="ready = true" />

		<div v-if="!url" class="placeholder empty">
			No image selected
		</div>

		<div v-else-if="!ready" class="placeholder loading">
			<LoadingMessage />
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.card-image {
		position: relative;
		background-color: var(--app-theme-snow-white);
		border-radius: 0.5rem;
		overflow: hidden;
		border: 0.25em solid var(--app-theme-snow-white);
		box-shadow: 0 0 1.5em rgba(0, 0, 0, 0.25);

		.placeholder {
			position: absolute;
			z-index: 1;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: rgba(0, 0, 0, 0.75);
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 0.85rem;
		}

		img {
			display: block;
			width: 100%;
			height: 100%;
			object-fit: cover;
			object-position: center;
		}
	}
</style>
