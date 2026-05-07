<script setup lang="ts">
import { computed } from 'vue';
import EndscreenScore from './EndscreenScore.vue';
import EndscreenStats from './EndscreenStats.vue';
import EndscreenStatTile from './EndscreenStatTile.vue';
import EndscreenActions from './EndscreenActions.vue';
import EndscreenButton from './EndscreenButton.vue';
import EndscreenHeader from './EndscreenHeader.vue';
import { intl, useLanguage } from '../../intl';

interface DeckStats {
	questions: number;
	score: number;
	time: number;
}

const props = defineProps<{
	stats: DeckStats;
}>();

const emit = defineEmits<{
	(e: 'reset'): void;
	(e: 'finish'): void;
}>();

const scoreRate = computed(() => props.stats.score / props.stats.questions);

const time = computed(() => {

	const { time } = props.stats;

	const mins = Math.floor(time / 60);
	if (mins > 60) {
		return "1hr+";
	}

	const secs = time - (mins * 60);

	return `${mins}:${secs.toString().padStart(2, '0')}`;
});

const lang = useLanguage();

</script>

<template>
	<div class="endscreen">
		<div class="content">

			<EndscreenHeader>
				<template v-slot:summary>

					<template v-if="scoreRate > 0.95">
						{{ intl(lang, {
							en: `You've nailed it!`,
							de: 'Perfekt hingekriegt!',
							uk: 'Вау, потужно!'
						}) }}
					</template>

					<template v-else-if="scoreRate > 0.75">
						{{ intl(lang, {
							en: 'Good job!',
							de: 'Gute Arbeit!',
							uk: 'Гарна робота!'
						}) }}
					</template>

					<template v-else-if="scoreRate > 0.6">
						{{ intl(lang, {
							en: 'Hey, not bad!',
							de: 'Hey, nicht schlecht!',
							uk: 'Скоріше так ніж ні!'
						}) }}
					</template>

					<template v-else-if="scoreRate > 0.3">
						{{ intl(lang, {
							en: 'You could do better!',
							de: 'Es könnte besser sein.',
							uk: 'Могло бути краще.'
						}) }}
					</template>

					<template v-else>
						{{ intl(lang, {
							en: 'Better luck next time!',
							de: 'Viel Glück beim nächsten Mal!',
							uk: 'Спробуємо знов?'
						}) }}
					</template>

				</template>
			</EndscreenHeader>

			<EndscreenScore :rate="scoreRate" />
			<EndscreenStats>
				<EndscreenStatTile icon="target">
					<template v-slot:title>
						{{ intl(lang, {
							en: 'Correct',
							de: 'Richtig',
							uk: 'Вірно'
						}) }}
					</template>
					{{ stats.score }}
					<template v-slot:after>
						/{{ stats.questions }}
					</template>
				</EndscreenStatTile>
				<EndscreenStatTile icon="time">
					<template v-slot:title>
						{{ intl(lang, {
							en: 'Time',
							de: 'Zeit',
							uk: 'Час'
						}) }}
					</template>
					{{ time }}
				</EndscreenStatTile>
				<EndscreenStatTile icon="prize">
					<template v-slot:title>
						{{ intl(lang, {
							en: 'Accuracy',
							de: 'Präzision',
							uk: 'Точність'
						}) }}
					</template>

					{{ Math.floor(scoreRate*100) }}%
				</EndscreenStatTile>
			</EndscreenStats>
			<EndscreenActions>
				<EndscreenButton icon="retry" :filled="true" @click="emit('reset')">
					{{ intl(lang, {
						en: 'Try again',
						de: 'Noch einmal',
						uk: 'Ще раз'
					}) }}
				</EndscreenButton>
				<EndscreenButton icon="finish" @click="emit('finish')">
					{{ intl(lang, {
						en: 'Finish',
						de: 'Beended',
						uk: 'Завершити'
					}) }}
				</EndscreenButton>
			</EndscreenActions>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.endscreen {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 2rem;
		width: 100%;
		max-width: 30rem;
		height: 100%;

		.content {
			display: flex;
			flex-direction: column;
			gap: 2rem;
			width: 100%;
		}
	}
</style>
