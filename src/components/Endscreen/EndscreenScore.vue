<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	score: number;
}>();

const score = computed(() => props.score > 1 ? 100 : props.score < 0 ? 0 : props.score * 100);

</script>

<template>
	<div class="endscreen-score">
		<div class="labels">
			<div class="title">
				Your score
			</div>
			<div class="score">
				{{ score.toFixed(0) }}%
			</div>
		</div>
		<div class="progress">
			<div class="fill" :class="{ red: score < 30, orange: score < 60 }" :style="{ width: `${score}%` }"></div>
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
