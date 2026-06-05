<script setup lang="ts">
import type { CardNode } from '@/content';
import CardThumbnail from './CardThumbnail.vue';

const props = defineProps<{
	cards: CardNode[];
}>();

const model = defineModel<Set<string>>();

const toggleId = (id: string) => model.value?.has(id) ? model.value.delete(id) : model.value?.add(id);

</script>

<template>
	<ul class="card-selector-grid">
		<li v-for="(card, idx) of cards">
			<button type="button" class="card" >
				<CardThumbnail :card="card.front" :active="modelValue?.has(card.id)" :label="idx + 1" @click="toggleId(card.id)" />
			</button>
		</li>
	</ul>
</template>

<style lang="scss" scoped>
	.card-selector-grid {
		display: flex;
		flex-flow: row wrap;
		gap: 0.5rem;
		list-style: none;
		padding: 0;
		margin: 0;
		min-height: 0;
		overflow: hidden auto;
		scrollbar-width: thin;
		height: 100%;

		li {
			display: block;
			padding: 0;
			margin: 0;

			button {
				display: block;
				padding: 0;
				font-size: inherit;
				outline: none;
				border: none;
				background: none;
			}
		}
	}
</style>
