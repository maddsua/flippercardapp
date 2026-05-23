<script setup lang="ts">
import { computed, reactive, watch } from 'vue';
import type { CardImageElement } from '../../content';
import LoadingMessage from '../App/LoadingMessage.vue';

const props = defineProps<{
	entry: CardImageElement;
}>();

const mediaURL = computed(() => props.entry.media_id ? `/media/images/${props.entry.media_id}` : null);

enum ReadyState {
	Idle,
	Failed,
	Ready
};

const state = reactive({
	ready: ReadyState.Idle,
	rotation: `${((Math.random() - 0.5) * 10).toFixed(0)}deg`,
});

watch(() => props.entry.media_id, () => state.ready = ReadyState.Idle);

</script>

<template>
	<div class="card-image" :class="{ placeholder: state.ready !== ReadyState.Ready }" :style="{ rotate: state.rotation }">

		<img v-if="mediaURL" :src="mediaURL" loading="lazy" @load="state.ready = ReadyState.Ready" @error="state.ready = ReadyState.Failed" />

		<div v-if="state.ready !== ReadyState.Ready" class="placeholder error">

			<template v-if="state.ready === ReadyState.Failed">
				Unable to load image
			</template>

			<template v-else-if="state.ready === ReadyState.Idle">
				<LoadingMessage />
			</template>

			<template v-else>
				No image selected
			</template>

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
