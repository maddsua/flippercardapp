<script setup lang="ts">

const props = defineProps<{
	has_next?: boolean;
	has_prev?: boolean;
}>();

const emit = defineEmits<{
	(e: 'next'): void;
	(e: 'prev'): void;
}>();

</script>

<template>
	<div class="card-navigation">
		<button type="button" class="prev" aria-label="Go to the previous card" :disabled="!has_prev" @click="emit('prev')"></button>
		<button type="button" class="next" aria-label="Go to the next card" :class="{ final: !has_next }" @click="emit('next')"></button>
	</div>
</template>

<style lang="scss" scoped>
	.card-navigation {
		position: absolute;
		bottom: 1rem;
		left: 0;
		width: 100%;
		display: flex;
		justify-content: space-around;
		z-index: 20;
	}

	button {
		width: 5rem;
		height: 3rem;
		background-size: 2rem;
		background-position: center;
		background-repeat: no-repeat;

		&.prev {
			background-image: url(/src/assets/icons/arrow-right-mask.svg);
			rotate: 180deg;
		}

		&.next {
			background-image: url(/src/assets/icons/arrow-right-mask.svg);
		}

		&.next.final {
			background-image: url(/src/assets/icons/check-mask.svg);
			background-size: 1.75rem;
		}

		&:disabled {
			cursor: pointer;
			pointer-events: none;
			opacity: 0.5;
		}
	}
</style>
