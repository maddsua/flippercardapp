<script setup lang="ts">
import { computed, ref } from 'vue';
import type { CardContentNode } from '../../content';
import CardFace from './CardFace.vue';

const props = defineProps<{
	card: CardContentNode;
}>();

const emit = defineEmits<{
	(e: 'prev'): void;
	(e: 'next'): void;
	(e: 'score', score: number): void;
}>();

const flipped = ref(false);
const randomRotate = computed(() => (Math.random() - 0.5) * 6);
const containerRef = ref<HTMLElement | null>(null);

const flip = () => flipped.value = !flipped.value;

interface DragState {
	initX: number;
	initY: number;
	x: number;
	y: number;
	pointerID: number | null;
	interactiveTarget: HTMLElement | null;
};

const dragState = ref<DragState | null>(null);

interface Delta2D {
	x: number;
	y: number;
};

const dragDelta = computed((): Delta2D | null => {

	if (!dragState.value) {
		return null;
	}

	return {
		x: dragState.value.x - dragState.value.initX,
		y: dragState.value.y - dragState.value.initY,
	};
});

const dragging = computed(() => dragDelta.value ? Math.abs(dragDelta.value.x) + Math.abs(dragDelta.value.y) > 1 : false);

const transformStyle = computed(() => ({
	rotate: dragDelta.value && dragging.value ? `${((dragDelta.value.x) / 25).toFixed(1)}deg` : `${randomRotate.value.toFixed(1)}deg`,
	transform: dragDelta.value ? `translateX(${dragDelta.value.x}px) translateY(${dragDelta.value.y}px) rotateY(${flipped.value ? 180 : 0}deg)` : undefined,
}));

interface BoxSize {
	width: number;
	height: number;
}

const dragBoundBoxSize = (): BoxSize => {

	if (!containerRef.value) {
		return { width: window.innerWidth, height: window.innerHeight };
	}

	const rect = containerRef.value.getBoundingClientRect();
	return { width: rect.width, height: rect.height };
};

const capturePointer = ({ pointerId }: PointerEvent) => {

	if (typeof pointerId !== 'number') {
		return null;
	}

	containerRef.value?.setPointerCapture(pointerId);

	return pointerId;
};

const releasePointerCapture = ({ pointerID: capturedPointerID }: DragState) => {

	if (!capturedPointerID) {
		return;
	}

	containerRef.value?.releasePointerCapture(capturedPointerID);
};

const isCapturedPointer = (state: DragState, event: PointerEvent): boolean => {
	return state.pointerID === null || state.pointerID === event.pointerId;
};

const handleDragStart = (event: PointerEvent) => {

	// don't process right and auxilary button events
	if (event.button < 0 || event.button > 1) {
		return;
	}

	// ignore any other pointers
	if (dragState.value) {
		return;
	}

	const target = event.target as HTMLElement;
	const interactiveTarget = target.closest<HTMLElement>('button, a, input, textarea, [data-interactive]');

	const { clientX, clientY } = event;

	dragState.value = {
		initX: clientX,
		initY: clientY,
		x: clientX,
		y: clientY,
		interactiveTarget,
	
		// allegedly, this prevents chrome from breaking the drag logic
		pointerID: capturePointer(event),
	};
};

const handleDragUpdate = (event: PointerEvent) => {

	if (!dragState.value || !isCapturedPointer(dragState.value, event)) {
		return;
	}

	dragState.value.x = event.clientX;
	dragState.value.y = event.clientY;
};

const handleSwipeGesture = (state: DragState) => {

	const dx = dragDelta.value?.x ?? 0;
	const dy = dragDelta.value?.y ?? 0;

	const { width, height } = dragBoundBoxSize();

	// if swiped vertically more than the threshold
	if (Math.abs(dy) > (height * 0.25)) {
		console.debug('swiped')
		dy > 1 ? emit('prev') : emit('next');
		return;
	}

	// if swiped horizontally more than the threshold
	if (Math.abs(dx) > (width * 0.4)) {
		flip();
		return;
	}

	// if we consider this to be a click/tap
	const totalDragged = Math.abs(dx) + Math.abs(dy);
	if (totalDragged < 3) {
		if (state.interactiveTarget) {
			state.interactiveTarget.click();
		} else {
			flip();
		}
	}
};

const handleDragDone = (event: PointerEvent) => {

	if (!dragState.value || !isCapturedPointer(dragState.value, event)) {
		return;
	}

	handleSwipeGesture(dragState.value);
	releasePointerCapture(dragState.value);

	dragState.value = null;
};

</script>

<template>
	<div class="card-container" :class="{ flipped, dragging }" :style="transformStyle"
		@pointerdown="handleDragStart"
		@pointermove="handleDragUpdate"
		@pointerup="handleDragDone"
		@pointercancel="handleDragDone"
		@pointerleave.self="handleDragDone"
		ref="containerRef">

		<CardFace
			:entry="card.front"
			decoration="question-mark"
			@flip="flip"
			@score="(score) => emit('score', score)"
			@next="emit('next')" />

		<CardFace
			:entry="card.back"
			:is3dBackface="true"
			@flip="flip"
			@score="(score) => emit('score', score)"
			@next="emit('next')" />

	</div>
</template>

<style lang="scss" scoped>
	.card-container {

		// set fixed-ish size
		position: relative;
		width: 100%;
		padding-bottom: 150%;

		// use container width as the font reference
		container-type: inline-size;

		// disable unwanted default interactions
		user-select: none;
		touch-action: none;

		// make 3d go brrrr
		transform-style: preserve-3d;
		perspective: 50cm;

		scale: 0.95;
		transition: transform 400ms ease;

		&.flipped {
			transform: rotateY(180deg);
		}

		&.dragging {
			scale: 0.9;
			transition: none;
		}
	}
	
</style>
