<script setup lang="ts">
import { computed, ref } from 'vue';
import type { CardNode } from './content';
import Card from './Card.vue';

const props = defineProps<{
	entries: CardNode[];
}>();

const activeIdx = ref(0);
const activeCard = computed(() => props.entries[activeIdx.value]);

const nextCard = () => {
	if (activeIdx.value < props.entries.length - 1) {
		activeIdx.value++
		return;
	}

	alert('End of the line!');

	//	todo: show end screen
};

const prevCard = () => {
	if (activeIdx.value > 0) {
		activeIdx.value--
	}
};

//	todo: animate A/B

</script>

<template>
	<Card :key="activeCard.id" :card="activeCard" @next="nextCard" @prev="prevCard" />
</template>

<style lang="scss" scoped>
	.card-view {
		display: block;
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
	}
</style>
