<script setup lang="ts">
import { computed, reactive, ref, type CSSProperties } from 'vue';
import type { CardNode } from '../../content';
import CardFace from './CardFace.vue';

const props = defineProps<{
	card: CardNode;
}>();

const emit = defineEmits<{
	(e: 'prev'): void;
	(e: 'next'): void;
	(e: 'score', score: number): void;
}>();

const state = reactive({
	flipped: false,
	rotation: (Math.random() - 0.5) * 6,
	shake: false,
	animating: false,
});

const containerRef = ref<HTMLElement | null>(null);

const hasBackface = computed(() => props.card.back.content.length > 0);

const flip = () => {

	if (state.animating) {
		return;
	}

	if (!hasBackface.value) {

		state.animating = true;
		setTimeout(() => state.animating = false, 500);

		state.shake = true;
		setTimeout(() => state.shake = false, 450);

		return;
	}

	state.animating = true;
	setTimeout(() => state.animating = false, 350);

	state.flipped = !state.flipped
};

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

const containerClasses = computed((): Record<string, boolean> => ({
	flipped: state.flipped,
	noinput: state.animating,
	shake: state.shake,
	dragging: dragging.value,
}));

const containterTransforms = computed((): CSSProperties => {

	const baseR = state.rotation;

	if (!dragDelta.value) {
		return { rotate: `${baseR.toFixed(1)}deg` };
	}

	const { x: deltaX, y: deltaY } = dragDelta.value;

	const deltaR = (deltaX / 25);

	return {
		rotate: `${(deltaR + baseR).toFixed(1)}deg`,
		transform: `translateX(${deltaX}px) translateY(${deltaY}px) rotateY(${state.flipped ? 180 : 0}deg)`,
	};
});

interface BoxSize {
	width: number;
	height: number;
};

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

const viewportEdgeDistance = (viewport: number, point: number): number => {
	const center = viewport / 2;
	return point > center ? viewport - point : point;
};

const handleDragStart = (event: PointerEvent) => {

	// don't process right and auxilary button events
	if (event.button < 0 || event.button > 1) {
		return;
	}

	// only handle one pointer at the time
	if (dragState.value) {
		return;
	}

	const { clientX, clientY } = event;

	//	don't process events right next to the screen edge
	// as those likely are  os navigation gestures
	const osGestureMargin = window.innerWidth * 0.1;
	if (viewportEdgeDistance(window.innerWidth, clientX) < osGestureMargin) {
		return;
	}

	const target = event.target as HTMLElement;
	const interactiveTarget = target.closest<HTMLElement>('button, a, input, textarea, [data-interactive]');

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

const handlePollScore = (score: number, final?: boolean) => {

	if (final) {
		if (score === 0 && hasBackface.value && !state.flipped) {
			setTimeout(flip, 300);
		} else {
			setTimeout(() => emit('next'), 500);
		}
	}

	emit('score', score);
};

</script>

<template>
	<div class="card-container"
		:class="containerClasses"
		:style="containterTransforms"
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
			@pollScore="handlePollScore"
			@next="emit('next')" />

		<CardFace
			:entry="card.back"
			:is3dBackface="true"
			@flip="flip"
			@pollScore="handlePollScore"
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
		transition: transform 300ms ease;

		&.noinput {
			pointer-events: none;
		}

		&.flipped {
			transform: rotateY(180deg);
		}

		&.dragging {
			scale: 0.9;
			transition: none;
		}

		&.shake {
			animation: horizontal-flip-shake 300ms;
		}
	}

	@keyframes horizontal-flip-shake {
		0% { transform: rotateY(0deg) }
		25% { transform: rotateY(20deg) }
		66% { transform: rotateY(-10deg) }
		100% { transform: rotateY(0deg) }
	}
	
</style>
