<script setup lang="ts">
import { ref, watch } from 'vue';
import type { CardContentFace } from '@/content';
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

const listRef = ref<Array<Element | null>>([]);

let scrollTimeout: NodeJS.Timeout | null = null;

watch(() => props.list.length, () => {

	if (scrollTimeout) {
		clearTimeout(scrollTimeout);
	}

	scrollTimeout = setTimeout(() => {
		listRef.value[props.pointer]?.scrollIntoView({ behavior: 'smooth', block: 'center' });
		scrollTimeout = null;
	}, 100);

});

</script>

<template>
	<div class="editor-deck-navigation">
		<div class="scroll-wrapper" :class="{ shaded: !!list.length }">
			<ul v-if="list.length" class="card-list">
				<li v-for="(card, idx) of list" :ref="elem => listRef[idx] = elem as Element | null">
					<button type="button" class="card" @click="emit('select', idx)">
						<CardThumbnail :card="card" :label="idx + 1" :active="idx === pointer" :showControls="true" @duplicate="emit('duplicate', idx)" @remove="emit('remove', idx)"/>
					</button>
				</li>
			</ul>
		</div>
		<button type="button" class="new-card" title="Add card" @click="emit('add')">+ Add card</button>
	</div>
</template>

<style lang="scss" scoped>
	.editor-deck-navigation {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1rem;
		user-select: none;
		min-height: 0;
		max-height: 100%;

		.scroll-wrapper {
			position: relative;
			min-height: 0;
			max-height: 100%;

			&.shaded {

				&::after, &::before {
					content: "";
					position: absolute;
					left: 0;
					z-index: 10;
					display: block;
					width: 100%;
					height: 2rem;
					background: linear-gradient(0deg,rgba(0, 0, 0, 0) 0%, #2f2f2f 100%);
				}

				&::before {
					top: 0;
				}

				&::after {
					bottom: 0;
					transform: rotate(180deg);
				}
			}
		}

		.card-list {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
			list-style: none;
			padding: 2rem 0;
			margin: 0;
			overflow: hidden auto;
			scrollbar-width: thin;
			padding-right: 0.25rem;
			scroll-behavior: smooth;
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
