<script setup lang="ts">
import { ref, watch } from 'vue';
import type { CardContentFace } from '../../content';
import CardThumbnail from './CardThumbnail.vue';

const props = defineProps<{
	list: CardContentFace[];
	pointer: number;
}>();

const emit = defineEmits<{
	(e: 'select', idx: number): void;
	(e: 'remove', idx: number): void;
	(e: 'duplicate', idx: number): void;
	(e: 'add'): void;
}>();

const scrollableRef = ref<HTMLElement | null>(null);

watch(() => props.list.length, (length, oldLength) => {
	if (length > oldLength && scrollableRef.value) {
		setTimeout(() => scrollableRef.value!.scrollTop = scrollableRef.value!.scrollHeight, 100);
	}
});

</script>

<template>
	<div class="editor-deck-navigation" ref="scrollableRef">
		<ul class="card-list">
			<li v-for="(card, idx) of list">
				<button type="button" class="card" @click="emit('select', idx)">
					<CardThumbnail :card="card" :label="idx + 1" :active="idx === pointer" :showControls="true" @duplicate="emit('duplicate', idx)" @remove="emit('remove', idx)"/>
				</button>
			</li>
		</ul>
		<button type="button" class="new-card" title="Add card" @click="emit('add')">+ Add card</button>
	</div>
</template>

<style lang="scss" scoped>
	.editor-deck-navigation {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		user-select: none;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding-right: 0.5rem;
		scroll-behavior: smooth;

		.card-list {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
			list-style: none;
			padding: 0;
			margin: 0;

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

		button.new-card {
			display: block;
			text-align: center;
			padding: 0.5rem;
			font-size: 0.65rem;
			font-weight: 600;
			background-color: var(--app-theme-sky-blue);
			color: var(--app-theme-snow-white);
			border: none;
			outline: none;
			border-radius: 0.25rem;

			&:hover {
				cursor: pointer;
				background-color: var(--app-theme-deep-lavender);
			}
		}
	}
</style>
