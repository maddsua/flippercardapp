<script setup lang="ts">
import { computed, ref } from 'vue';
import type { CardNode } from './content';
import CardFace from './CardFace.vue';

const props = defineProps<{
	card: CardNode;
}>();

const emit = defineEmits<{
	(e: 'prev'): void;
	(e: 'next'): void;
}>();

const flipped = ref(false);
const randomRotate = computed(() => (Math.random() - 0.5) * 6);

const flip = () => flipped.value = !flipped.value;

const dragState = ref<{
	initX: number;
	initY: number;
	x: number;
	y: number;
	targetInteractive: boolean;
} | null>(null);

const dragDelta = computed(() => dragState.value ? ({ x: dragState.value.x - dragState.value.initX, y: dragState.value.y - dragState.value.initY }) : null);
const dragging = computed(() => dragDelta.value ? Math.abs(dragDelta.value.x) + Math.abs(dragDelta.value.y) > 1 : false);

const transformStyle = computed(() => ({
	rotate: dragDelta.value && dragging.value ? `${((dragDelta.value.x) / 25).toFixed(1)}deg` : `${randomRotate.value.toFixed(1)}deg`,
	transform: dragDelta.value ? `translateX(${dragDelta.value.x}px) translateY(${dragDelta.value.y}px) rotateY(${flipped.value ? 180 : 0}deg)` : undefined,
}));

const handleDragStart = (event: PointerEvent) => {

	const { clientX, clientY } = event;

	dragState.value = {
		initX: clientX,
		initY: clientY,
		x: clientX,
		y: clientY,
		targetInteractive: !!(event.target as HTMLElement | undefined)?.closest('button, a, input, textarea'),
	};
};

const handleDragUpdate = (event: PointerEvent) => {

	if (!dragState.value) {
		return;
	}

	const { clientX, clientY } = event;

	dragState.value.x = clientX;
	dragState.value.y = clientY;
};

const relativeSwipeThreshold = 0.25;

const handleDragDone = () => {

	if (dragDelta.value) {

		const { x, y } = dragDelta.value;

		const delta = Math.abs(x) + Math.abs(y)
	
		if (Math.abs(y) > window.innerHeight * relativeSwipeThreshold) {
			y > 1 ? emit('prev') : emit('next');
		} else if (Math.abs(x) > window.innerWidth * relativeSwipeThreshold) {
			flip();
		} else if (delta < 1 && !dragState.value?.targetInteractive) {
			flip();
		}
	}

	dragState.value = null;
};

</script>

<template>
	<div class="card-container" :class="{ flipped, dragging }" :style="transformStyle" @pointerdown="handleDragStart" @pointermove="handleDragUpdate" @pointerup="handleDragDone" @pointercancel="handleDragDone" @pointerout="handleDragDone" @flip="flip">
		<CardFace :entry="card.front" @flip="flip" @next="emit('next')" />
		<CardFace :entry="card.back" @flip="flip" @next="emit('next')" />
	</div>
</template>

<style lang="scss" scoped>
	.card-container {
		//	todo: research
		//	font-size: clamp(2rem, 5vw, 5rem);
		position: relative;
		width: 100%;
		padding-bottom: 150%;

		user-select: none;

		transition: transform 400ms ease;
		transform-style: preserve-3d;
		perspective: 50cm;
		scale: 0.95;

		&.flipped {
			transform: rotateY(180deg);
		}

		&.dragging {
			scale: 0.9;
			transition: none;
		}
	}
	
</style>
