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
	(e: 'score', score: number): void;
}>();

const flipped = ref(false);
const randomRotate = computed(() => (Math.random() - 0.5) * 6);

const flip = () => flipped.value = !flipped.value;

const dragState = ref<{
	initX: number;
	initY: number;
	x: number;
	y: number;
} | null>(null);

const dragDelta = computed(() => dragState.value ? ({ x: dragState.value.x - dragState.value.initX, y: dragState.value.y - dragState.value.initY }) : null);
const dragging = computed(() => dragDelta.value ? Math.abs(dragDelta.value.x) + Math.abs(dragDelta.value.y) > 1 : false);

const transformStyle = computed(() => ({
	rotate: dragDelta.value && dragging.value ? `${((dragDelta.value.x) / 25).toFixed(1)}deg` : `${randomRotate.value.toFixed(1)}deg`,
	transform: dragDelta.value ? `translateX(${dragDelta.value.x}px) translateY(${dragDelta.value.y}px) rotateY(${flipped.value ? 180 : 0}deg)` : undefined,
}));

const capturePointer = (event: PointerEvent) => {
	const target = event.target as HTMLElement;
	target.setPointerCapture(event.pointerId);
};

const releasePointerCapture = (event: PointerEvent) => {
	const target = event.target as HTMLElement;
	target.releasePointerCapture(event.pointerId);
};

const targetDraggable = (event: PointerEvent): boolean => {
	const target = event.target as HTMLElement;
	return !target.closest('button, a, input, textarea');
};

const handleDragStart = (event: PointerEvent) => {

	if (!targetDraggable(event)) {
		return;
	}

	capturePointer(event);

	const { clientX, clientY } = event;

	dragState.value = {
		initX: clientX,
		initY: clientY,
		x: clientX,
		y: clientY,
	};
};

const handleDragUpdate = (event: PointerEvent) => {

	if (!dragState.value) {
		releasePointerCapture(event);
		return;
	}

	const { clientX, clientY } = event;

	dragState.value.x = clientX;
	dragState.value.y = clientY;
};

const relativeSwipeThresholdY = 0.25;
const relativeSwipeThresholdX = 0.45;

const handleDragDone = (event?: PointerEvent) => {

	if (event) {
		releasePointerCapture(event);
	}

	if (dragDelta.value) {

		const { x, y } = dragDelta.value;

		const delta = Math.abs(x) + Math.abs(y)
	
		if (Math.abs(y) > window.innerHeight * relativeSwipeThresholdY) {
			y > 1 ? emit('prev') : emit('next');
		} else if (Math.abs(x) > (window.innerWidth * relativeSwipeThresholdX) || delta < 1) {
			flip();
		}
	}

	dragState.value = null;
};

</script>

<template>
	<div class="card-container" :class="{ flipped, dragging }" :style="transformStyle" @pointerdown="handleDragStart" @pointermove="handleDragUpdate" @pointerup="handleDragDone" @pointercancel="handleDragDone" @pointerout="handleDragDone" @flip="flip">
		<CardFace :entry="card.front" @flip="flip" @score="(score) => emit('score', score)" @next="emit('next')" />
		<CardFace :entry="card.back" @flip="flip" @score="(score) => emit('score', score)" @next="emit('next')" />
	</div>
</template>

<style lang="scss" scoped>
	.card-container {
		position: relative;
		width: 100%;
		padding-bottom: 150%;

		user-select: none;
		touch-action: none;

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
