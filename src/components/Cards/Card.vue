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

//	todo: handle swiping

const flipped = ref(false);
const randomRotate = computed(() => (Math.random() - 0.5) * 6);

const flip = () => flipped.value = !flipped.value;

</script>

<template>
	<div class="card-container" :class="{ flipped }" :style="{ rotate: `${randomRotate.toFixed(1)}deg` }" @flip="flip">
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
	}
	
</style>
