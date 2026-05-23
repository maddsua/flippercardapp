<script setup lang="ts">
import { computed, reactive, watch } from 'vue';
import type { CardImageElement } from '../../content';
import LoadingMessage from '../App/LoadingMessage.vue';

const props = defineProps<{
	entry: CardImageElement;
}>();

const mediaURL = computed(() => props.entry.media_id ? `/media/images/${props.entry.media_id}` : null);

const state = reactive({
	ready: false,
	failed: false,
	rotation: `${((Math.random() - 0.5) * 10).toFixed(0)}deg`,
});

const resetState = () => {
	state.ready = false;
	state.failed = false;
};

watch(() => props.entry.media_id, resetState);

</script>

<template>
	<div class="card-image" :class="{ failed: state.failed }" :style="{ rotate: state.rotation }">

		<img v-if="mediaURL && !state.failed" :src="mediaURL" @load="state.ready = true" @error="state.failed = true" />

		<div v-else-if="state.failed" class="placeholder error">
			Unable to load image
		</div>

		<div v-if="!mediaURL" class="placeholder empty">
			No image selected
		</div>

		<div v-else-if="!state.ready" class="placeholder loading">
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

		&.failed {
			width: 15em;
			height: 6em;
		}

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
			font-size: 0.65rem;
			color: var(--app-theme-snow-white);
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
