<script setup lang="ts">
import { nextTick, reactive, ref } from 'vue';
import type { CardNode } from '../../content';
import Card from './Card.vue';
import CardControls from './CardControls.vue';
import CardDeckInfo from './CardDeckInfo.vue';

const props = defineProps<{
	labels: string[];
	entries: CardNode[];
}>();

const emit = defineEmits<{
	(e: 'finish'): void;
	(e: 'exit'): void;
	(e: 'score', score: number): void;
}>();

const activeIdx = ref(0);

interface CardState {
	card: CardNode;
	flags: {
		active: boolean;
		animate?: boolean;
		fade?: boolean;
		swipe_up?: boolean;
		swipe_down?: boolean;
	}
};

type SlideInDirection = 'from-top' | 'from-bottom';
type SlideOutDirection = 'to-top' | 'to-bottom';

const slideOut = (direction?: SlideInDirection): SlideOutDirection | undefined =>
	direction === 'from-bottom' ? 'to-top' :
		direction === 'from-top' ? 'to-bottom' : undefined;

const newAnimatedState = (direction?: SlideInDirection): CardState => {

	if (!props.entries.length) {
		throw new Error('Mounting card view without any cards is undefined behavior');
	}

	const nextCard = props.entries[activeIdx.value];
	const state: CardState = {
		card: nextCard,
		flags: {
			active: true,
		}
	};

	switch (direction) {
		case 'from-top':

			state.flags.swipe_up = true;
			state.flags.fade = true;

			nextTick(() => {
				state.flags.animate = true;
				nextTick(() => {
					state.flags.swipe_up = false;
					state.flags.fade = false;
				});
			});

			break;
	
		case 'from-bottom':

			state.flags.swipe_down = true;
			state.flags.fade = true;

			nextTick(() => {
				state.flags.animate = true;
				nextTick(() => {
					state.flags.swipe_down = false;
					state.flags.fade = false;
				});
			});

			break;
	}

	return state;
}

const pairState = reactive({
	a: newAnimatedState(),
	b: null as CardState | null,
});

const ejectAnimatedState = async (state: CardState, direction?: SlideOutDirection) => {

	state.flags.animate = true;
	
	nextTick(() => {

		state.flags.active = false;

		if (!direction) {
			state.flags.fade = true;
			return;
		}

		switch (direction) {
			case 'to-bottom':
				state.flags.swipe_down = true;
				break;
		
			case 'to-top':
				state.flags.swipe_up = true;
				break;
		}

		setTimeout(() => state.flags.fade = true, 200);
	});
};

const switchCards = (direction?: SlideInDirection) => {

	if (!pairState.b) {
		ejectAnimatedState(pairState.a, slideOut(direction));
		pairState.b = newAnimatedState(direction);
		return
	}

	if (pairState.a.flags.active) {
		ejectAnimatedState(pairState.a, slideOut(direction));
		pairState.b = newAnimatedState(direction);
	} else {
		ejectAnimatedState(pairState.b, slideOut(direction));
		pairState.a = newAnimatedState(direction);
	}
};

const nextCard = (): boolean => {

	if (activeIdx.value < props.entries.length - 1) {
		activeIdx.value++
		switchCards('from-bottom');
		return true;
	}

	emit('finish');
	return false;
};

const prevCard = (): boolean => {

	if (activeIdx.value > 0) {
		activeIdx.value--
		switchCards('from-top');
		return true;
	}

	return false
};

const scoreSet = new Set<string>();
const countScore = (score: number) => {

	const key = props.entries[activeIdx.value].id;

	if (scoreSet.has(key)) {
		return;
	}

	scoreSet.add(key);
	emit('score', score);
};

const handleCtrlBack = () => {

	if (!prevCard()) {
		//	todo: add a prompt
		emit('exit');
	}

};

</script>

<template>
	<div class="card-widget">
		<CardDeckInfo :labels="labels" :size="entries.length" :index="activeIdx" />
		<div class="card-screen-container">
			<div class="card-transition-slot" v-for="(item,idx) of [pairState.a, pairState.b]" :key="`${idx}:${item?.card.id || 'null'}`" :class="item?.flags">
				<Card v-if="item" :key="item.card.id" :card="item.card" @score="countScore" @next="nextCard" @prev="prevCard" />
			</div>
		</div>
		<CardControls :has_prev="activeIdx > 0" :has_next="activeIdx < entries.length - 1" @prev="handleCtrlBack" @next="nextCard" />
	</div>
</template>

<style lang="scss" scoped>
	.card-widget {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: visible;
		display: flex;
		flex-direction: column;
		align-items: center;

		@media (orientation: landscape) {
			max-width: 50rem;
		}

		.card-screen-container {
			position: relative;
			width: 55vh;
			height: 100%;
			overflow: visible;

			@media (max-aspect-ratio: 1/1.75) {
				width: 100vw;
			}
		}

		.card-transition-slot {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			justify-content: center;
			overflow: visible;
			
			&.animate {
				transition: all 200ms ease;
			}

			&.active {
				z-index: 1;
			}

			&.swipe_up {
				transform: translateX(25vw) translateY(-100vh);
			}

			&.swipe_down {
				transform: translateX(-25vw) translateY(100vh);
			}

			&.fade {
				opacity: 0;
			}
		}
	}
</style>
