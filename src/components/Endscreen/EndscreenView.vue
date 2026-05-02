<script setup lang="ts">
import { computed } from 'vue';
import EndscreenScore from './EndscreenScore.vue';
import EndscreenStats from './EndscreenStats.vue';
import EndscreenStatTile from './EndscreenStatTile.vue';
import EndscreenActions from './EndscreenActions.vue';
import EndscreenButton from './EndscreenButton.vue';
import EndscreenHeader from './EndscreenHeader.vue';

const props = defineProps<{
	collectionName: string;
	tagName?: string;
	size: number;
	correct: number;
	time: number;
}>();

//	todo: calc
const time = computed(() => '4:32');

//	todo: style buttons

//	todo: add events

//	todo: style stats

//	todo: wrap the whole thing int a card

//	todo: make reactive

</script>

<template>
	<div class="endscreen-view">
		<div class="endscreen-card">

			<EndscreenHeader>
				<template v-slot:summary>
					Good job!
				</template>
			</EndscreenHeader>

			<EndscreenScore :score="correct/size" />
			<EndscreenStats>
				<EndscreenStatTile icon="target">
					<template v-slot:title>
						Correct
					</template>
					{{ correct }}
					<template v-slot:after>
						/{{ size }}
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

					{{ (correct/size)*100 }}%
				</EndscreenStatTile>
			</EndscreenStats>
			<EndscreenActions>
				<EndscreenButton icon="retry" :filled="true">
					Try again
				</EndscreenButton>
				<EndscreenButton icon="home">
					Home
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
