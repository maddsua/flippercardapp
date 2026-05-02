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
} | null>(null);

const transformStyle = computed(() => ({
	rotate: dragState.value ? `${((dragState.value.x - dragState.value.initX) / 100).toFixed(1)}deg` : `${randomRotate.value.toFixed(1)}deg`,
	transform: dragState.value ? `translateX(${dragState.value.x - dragState.value.initX}px) translateY(${dragState.value.y - dragState.value.initY}px)` : undefined,
}));

const handleDragStart = (event: PointerEvent) => {

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
		return;
	}

	const { clientX, clientY } = event;

	dragState.value.x = clientX;
	dragState.value.y = clientY;
};

const handleDragDone = () => {

	if (!dragState.value) {
		return;
	}

	const delta = dragState.value.initY - dragState.value.y;

	if (Math.abs(delta) > 50) {
		if (delta > 0) {
			emit('next');
		} else {
			emit('prev');
		}
	}

	dragState.value = null;
};

//	todo: fix clicking

</script>

<template>
	<div class="card-container" :class="{ flipped, dragging: !!dragState }" :style="transformStyle" @pointerdown="handleDragStart" @pointermove="handleDragUpdate" @pointerup="handleDragDone" @pointercancel="handleDragDone" @pointerout="handleDragDone" @flip="flip">
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
			transition: none;
		}
	}
	
</style>
