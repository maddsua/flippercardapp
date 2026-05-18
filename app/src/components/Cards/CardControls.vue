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
		<button type="button" class="prev" aria-label="Go to the previous card" :class="{ initial: !has_prev }" @click="emit('prev')"></button>
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
		z-index: 10;
	}

	button {
		display: block;
		border-radius: 0.5rem;
		border: 1px solid transparent;
		width: 3.5rem;
		height: 2.5rem;
		background-color: var(--app-theme-ghostly-glow);
		background-size: 1.75rem;
		background-position: center;
		background-repeat: no-repeat;
		cursor: pointer;
		transition: border-color 0.25s;

		&.prev {
			background-image: url(/src/assets/icons/arrow-right-mask.svg);
			rotate: 180deg;

			&.initial {
				background-image: url(/src/assets//icons/cancel-mask.svg);
				background-size: 1rem;
			}
		}

		&.next {
			background-image: url(/src/assets/icons/arrow-right-mask.svg);

			&.final {
				background-image: url(/src/assets/icons/check-mask.svg);
				background-size: 1.5rem;
			}
		}

		&:disabled {
			cursor: pointer;
			pointer-events: none;
			opacity: 0.5;
		}

		&:hover {
			cursor: pointer;
			border-color: var(--app-theme-sky-blue);
		}

		&:focus, &:focus-visible {
			outline: 4px auto -webkit-focus-ring-color;
		}
	}
</style>
