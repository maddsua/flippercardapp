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

interface DragState {
	initX: number;
	initY: number;
	x: number;
	y: number;
	targetInteractive: boolean;
};

const dragState = ref<DragState | null>(null);

interface DragDelta {
	x: number;
	y: number;
};

const dragDelta = computed((): DragDelta | null => dragState.value ? ({ x: dragState.value.x - dragState.value.initX, y: dragState.value.y - dragState.value.initY }) : null);
const dragging = computed(() => dragDelta.value ? Math.abs(dragDelta.value.x) + Math.abs(dragDelta.value.y) > 1 : false);

const transformStyle = computed(() => ({
	rotate: dragDelta.value && dragging.value ? `${((dragDelta.value.x) / 25).toFixed(1)}deg` : `${randomRotate.value.toFixed(1)}deg`,
	transform: dragDelta.value ? `translateX(${dragDelta.value.x}px) translateY(${dragDelta.value.y}px) rotateY(${flipped.value ? 180 : 0}deg)` : undefined,
}));

const pointerID = (event: PointerEvent): number | null => {
	if (event.pointerId > 0) {
		return event.pointerId
	}
	return null;
};

const capturePointer = (event: PointerEvent) => {
	const target = event.target as HTMLElement;
	const ptrID = pointerID(event);
	if (ptrID) {
		target.setPointerCapture(ptrID);
	}
};

const releasePointerCapture = (event: PointerEvent) => {
	const target = event.target as HTMLElement;
	const ptrID = pointerID(event);
	if (ptrID) {
		target.releasePointerCapture(ptrID);
	}
};

const handleDragStart = (event: PointerEvent) => {

	// don't process right and auxilary button events
	if (event.button < 0 || event.button > 1) {
		return;
	}

	// allegedly, this prevents chrome from breaking the drag logic
	capturePointer(event);

	const { clientX, clientY } = event;
	const target = event.target as HTMLElement;

	dragState.value = {
		initX: clientX,
		initY: clientY,
		x: clientX,
		y: clientY,
		targetInteractive: !!target.closest('button, a, input, textarea, [data-interactive]'),
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

const handleSwipeGesture = (state: DragState) => {

	const dx = dragDelta.value?.x ?? 0;
	const dy = dragDelta.value?.y ?? 0;

	const thresholdY = window.innerHeight * 0.25;
	if (Math.abs(dy) > thresholdY) {
		dy > 1 ? emit('prev') : emit('next');
		return;
	}

	const thresholdX = window.innerWidth * 0.4;
	if (Math.abs(dx) > thresholdX) {
		flip();
		return;
	}

	const delta = Math.abs(dx) + Math.abs(dy);
	if (delta < 1 && !state.targetInteractive) {
		flip();
		return;
	}

};

const handleDragDone = (event?: PointerEvent) => {

	if (event) {
		releasePointerCapture(event);
	}

	if (dragState.value) {
		handleSwipeGesture(dragState.value);
	}

	dragState.value = null;
};

</script>

<template>
	<div class="card-container" :class="{ flipped, dragging }" :style="transformStyle" @pointerdown="handleDragStart" @pointermove="handleDragUpdate" @pointerup="handleDragDone" @pointercancel="handleDragDone">
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
