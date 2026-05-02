<script setup lang="ts">
import { nextTick, reactive, ref } from 'vue';
import type { CardNode } from './content';
import Card from './Card.vue';
import CardNavigation from './CardNavigation.vue';
import CardDeckInfo from './CardDeckInfo.vue';

const props = defineProps<{
	entries: CardNode[];
}>();

const emit = defineEmits<{
	(e: 'end'): void;
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

const nextCard = () => {
	if (activeIdx.value < props.entries.length - 1) {
		activeIdx.value++
		switchCards('from-bottom');
		return;
	}

	emit('end');
};

const prevCard = () => {
	if (activeIdx.value > 0) {
		activeIdx.value--
		switchCards('from-top');
	}
};

</script>

<template>
	<div class="card-view">
		<CardDeckInfo collectionName="Test collection" tagName="All" :size="entries.length" :index="activeIdx" />
		<template v-for="(item,idx) of [pairState.a, pairState.b]" :key="`${idx}:${item?.card.id}`">
			<div class="card-slot" :class="item?.flags">
				<Card v-if="item" :key="item.card.id" :card="item.card" @next="nextCard" @prev="prevCard" />
			</div>
		</template>
		<CardNavigation :has_prev="activeIdx > 0" :has_next="activeIdx < entries.length - 1" @prev="prevCard" @next="nextCard" />
	</div>
</template>

<style lang="scss" scoped>
	.card-view {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;

		.card-slot {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			justify-content: center;
			overflow: hidden;
			
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
