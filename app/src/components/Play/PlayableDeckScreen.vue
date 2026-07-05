<script setup lang="ts">
import { computed, nextTick, reactive, ref } from 'vue';
import type { CardNode } from '@/content';
import Card from '../Cards/Card.vue';
import CardDeckControls from './CardDeckControls.vue';
import CardDeckInfo from './CardDeckInfo.vue';
import { appCanShareData, appShareData } from '@/share';
import { useStorage } from '@/storage/storage';

const props = defineProps<{
	labels: string[];
	entries: CardNode[];
	isMarked?: boolean;
	showNavigation?: boolean;
	disableRotation?: boolean;
}>();

const emit = defineEmits<{
	(e: 'finish'): void;
	(e: 'exit'): void;
	(e: 'score', score: number): void;
	(e: 'toggleMarked'): void;
}>();

const store = useStorage();

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
};

const createShareState = (): ShareData | null => {

	const shareLinkOnly = store.preferences.sharing.linkOnly.load();

	return appCanShareData() ? {
		title: (!shareLinkOnly && props.labels.length) ? `Play: ${props.labels[0]}` : undefined,
		url: window.location.href,
	} : null;
};

const state = reactive({
	slots: {
		a: newAnimatedState(),
		b: null as CardState | null,
	},
	animating: false,
	shareable: createShareState(),
})

const cardSlots = computed(() => ([state.slots.a, state.slots.b]))

const cardSlotKey = (slot: CardState | null, idx: number) => `${idx}:${slot?.card.id || 'null'}`;

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

		setTimeout(() => state.flags.fade = true, 50);
	});
};

const switchCards = (direction?: SlideInDirection) => {

	if (!state.slots.b) {
		ejectAnimatedState(state.slots.a, slideOut(direction));
		state.slots.b = newAnimatedState(direction);
		return
	}

	if (state.slots.a.flags.active) {
		ejectAnimatedState(state.slots.a, slideOut(direction));
		state.slots.b = newAnimatedState(direction);
	} else {
		ejectAnimatedState(state.slots.b, slideOut(direction));
		state.slots.a = newAnimatedState(direction);
	}
};

const withAnimationLock = (): boolean => {
	if (state.animating) {
		return false;
	}
	state.animating = true;
	setTimeout(() => state.animating = false, 500);
	return true;
};

const navigateForward = () => {

	if (!withAnimationLock()) {
		return;
	}

	if (activeIdx.value < props.entries.length - 1) {
		activeIdx.value++
		switchCards('from-bottom');
		return;
	}

	emit('finish');
};

const prevCard = () => {

	if (!withAnimationLock()) {
		return;
	}

	if (activeIdx.value > 0) {
		activeIdx.value--
		switchCards('from-top');
		return;
	}
};

const navigateBack = () => {
	if (activeIdx.value > 0) {
		prevCard();
	} else {
		emit('exit');
	}
};

const scoreState = reactive({
	totalAnswers: 0,
	cardAnswerSet: new Set<string>(),
});

const countScore = (score: number) => {

	const { id: cardID } = props.entries[activeIdx.value];

	if (!scoreState.cardAnswerSet.has(cardID)) {
		scoreState.cardAnswerSet.add(cardID);
		emit('score', score);
	}

	scoreState.totalAnswers++;
};

</script>

<template>
	<div class="deck-screen">

		<CardDeckInfo
			:labels="labels"
			:size="entries.length"
			:index="activeIdx"
			:isMarked="isMarked"
			:shareable="!!state.shareable"
			@toggleMarked="emit('toggleMarked')"
			@share="appShareData(state.shareable)"
			@exit="emit('exit')" />

		<div class="card-viewport" :class="{ noinput: state.animating }">
			<div class="card-transition" v-for="(item, idx) of cardSlots" :key="cardSlotKey(item, idx)" :class="item?.flags">
				<Card v-if="item"
					:key="item.card.id"
					:card="item.card"
					:disableRotation="disableRotation"
					@score="countScore"
					@next="navigateForward"
					@prev="prevCard" />
			</div>
		</div>

		<CardDeckControls v-if="showNavigation"
			:has_prev="activeIdx > 0"
			:has_next="activeIdx < entries.length - 1"
			@prev="navigateBack"
			@next="navigateForward" />

	</div>
</template>

<style lang="scss" scoped>
	.deck-screen {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: visible;
		display: flex;
		flex-direction: column;
		align-items: center;
		user-select: none;

		@media (orientation: landscape) {
			max-width: 50rem;
		}

		.card-viewport {
			position: relative;
			width: 55vh;
			height: 100%;
			overflow: visible;

			&.noinput {
				pointer-events: none;
			}

			@media (orientation: portrait) {
				width: 45vh;
			}

			@media (max-aspect-ratio: 0.55) {
				width: 100%;
			}
		}

		.card-transition {
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
