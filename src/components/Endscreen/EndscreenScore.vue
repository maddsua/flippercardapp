<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	rate: number;
}>();

const percentage = computed(() => props.rate > 1 ? 100 : props.rate < 0 ? 0 : props.rate * 100);

</script>

<template>
	<div class="endscreen-score">
		<div class="labels">
			<div class="title">
				Your score
			</div>
			<div class="score">
				{{ percentage.toFixed(0) }}%
			</div>
		</div>
		<div class="progress">
			<div class="fill" :class="{ red: rate < 0.3, orange: rate < 0.6 }" :style="{ width: `${percentage}%` }"></div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.endscreen-score {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;

		.labels {
			display: flex;
			flex-flow: row nowrap;
			justify-content: space-between;
			gap: 1rem;

			.title {
				color: var(--app-color-light-grey);
			}

			.score {
				font-weight: 600;
			}
		}

		.progress {
			width: 100%;
			height: 0.75rem;
			overflow: hidden;
			border-radius: 1rem;
			background-color: var(--app-color-blueish-dark);

			.fill {
				background-color: var(--app-accent-blue);
				height: 100%;

				&.red {
					background-color: var(--app-accent-red);
				}

				&.orange {
					background-color: var(--app-accent-orange);
				}
			}
		}
	}
</style>
