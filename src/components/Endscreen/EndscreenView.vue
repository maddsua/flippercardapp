<script setup lang="ts">
import { computed } from 'vue';
import EndscreenScore from './EndscreenScore.vue';
import EndscreenStats from './EndscreenStats.vue';
import EndscreenStatTile from './EndscreenStatTile.vue';
import EndscreenActions from './EndscreenActions.vue';
import EndscreenButton from './EndscreenButton.vue';
import EndscreenHeader from './EndscreenHeader.vue';

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

</script>

<template>
	<div class="endscreen-view">
		<div class="endscreen-card">

			<EndscreenHeader>
				<template v-slot:summary>

					<template v-if="scoreRate > 0.95">
						You've nailed it!
					</template>

					<template v-else-if="scoreRate > 0.75">
						Good job!
					</template>

					<template v-else-if="scoreRate > 0.6">
						Hey, not bad!
					</template>

					<template v-else-if="scoreRate > 0.3">
						You could do better!
					</template>

					<template v-else>
						Better luck next time!
					</template>

				</template>
			</EndscreenHeader>

			<EndscreenScore :rate="scoreRate" />
			<EndscreenStats>
				<EndscreenStatTile icon="target">
					<template v-slot:title>
						Correct
					</template>
					{{ stats.score }}
					<template v-slot:after>
						/{{ stats.questions }}
					</template>
				</EndscreenStatTile>
				<EndscreenStatTile icon="time">
					<template v-slot:title>
						Time
					</template>
					{{ time }}
				</EndscreenStatTile>
				<EndscreenStatTile icon="prize">
					<template v-slot:title>
						Accuracy
					</template>

					{{ Math.floor(scoreRate*100) }}%
				</EndscreenStatTile>
			</EndscreenStats>
			<EndscreenActions>
				<EndscreenButton icon="retry" :filled="true" @click="emit('reset')">
					Try again
				</EndscreenButton>
				<EndscreenButton icon="finish" @click="emit('finish')">
					Finish
				</EndscreenButton>
			</EndscreenActions>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.endscreen-view {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		justify-content: center;
	
		.endscreen-card {
			display: flex;
			flex-direction: column;
			gap: 1.5rem;
			padding: 2rem;
		}
	}
</style>
