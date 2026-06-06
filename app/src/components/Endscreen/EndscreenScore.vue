<script setup lang="ts">
import { computed } from 'vue';
import { useLanguage, intl } from '@/intl';

const props = defineProps<{
	rate: number;
}>();

const lang = useLanguage();

const percentage = computed(() => props.rate > 1 ? 100 : props.rate < 0 ? 0 : props.rate * 100);

</script>

<template>
	<div class="endscreen-score">
		<div class="labels">
			<div class="title">
				{{ intl(lang, {
					en: 'Your score',
					de: 'Deine Punktzahl',
					uk: 'Результат'
				}) }}
			</div>
			<div class="score">
				{{ Math.floor(percentage) }}%
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
				color: var(--app-theme-mysterious-white);
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
			background-color: var(--app-theme-midnight-glow);

			.fill {
				background-color: var(--app-theme-sky-blue);
				height: 100%;

				&.red {
					background-color: var(--app-theme-blood-red);
				}

				&.orange {
					background-color: var(--app-theme-spooky-orange);
				}
			}
		}
	}
</style>
