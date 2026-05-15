<script setup lang="ts">
import CardPoll from './CardPoll.vue';
import CardTextBox from './CardTextBox.vue';
import CardTextNode from './CardTextNode.vue';
import CardTitle from './CardTitle.vue';
import type { CardFace } from '../../content';
import { computed } from 'vue';

const props = defineProps<{
	entry: CardFace;
	decoration?: 'question-mark';
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
	(e: 'score', score: number): void;
}>();

const cardTheme = computed(() => props.entry.theme?.card);

</script>

<template>
	<div class="card-canvas" :style="{ color: cardTheme?.mask_color }" :class="{ [`decoration-${decoration}`]: !!decoration }">
		<div class="card-content" :style="{ backgroundColor: cardTheme?.fill_color, borderColor: cardTheme?.outline_color || cardTheme?.fill_color, color: cardTheme?.mask_color }">
			<template v-for="node of entry.content">
				<CardTitle v-if="node.type === 'title'">
					{{ node.content || '[Title]' }}
				</CardTitle>
				<CardTextBox v-else-if="node.type === 'textbox'">
					<template v-for="txtnode of node.content">
						<CardTextNode v-if="txtnode.type === 'text'" :theme="txtnode.theme">
							{{ txtnode.content || '[Content]' }}
						</CardTextNode>
						<br v-else-if="txtnode.type === 'newline'" />
					</template>
				</CardTextBox>
				<CardPoll v-else-if="node.type === 'poll'" :entry="node" :theme="entry.theme?.interactives" @score="(score) => emit('score', score)" @flip="emit('flip')" @next="emit('next')" />
			</template>
		</div>
	</div>
</template>

<style lang="scss" scoped>

	.card-canvas {
		position: absolute;
		top: 0;
		bottom: 0;
		overflow: hidden;

		// sets adaptive font size for the card
		font-size: 4cqw;

		display: flex;
		flex-direction: column;
		width: 100%;
		height: 100%;
		color: var(--app-theme-spooky-orange);
		background-color: var(--app-theme-snow-white);
		border-radius: 2em;
		padding: 1em;
		
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;

		// these two somehow prevent jagged edges
		outline: 1px solid transparent;
		box-shadow: 1em 0.25rem 0.25rem rgba(0, 0, 0, 0.25);

		&:nth-child(2n) {
			transform: rotateY(180deg);
		}

		&.decoration-question-mark::before {
			content: "?";
			display: block;
			font-size: 3.25em;
			line-height: 1em;
			position: absolute;
			top: 0.5em;
			left: 0.75em;
			font-weight: 600;
			color: inherit;
		}
	}

	.card-content {
		display: flex;
		flex-direction: column;
		gap: 1.5em;
		align-items: center;
		flex-grow: 1;
		border: 3px solid var(--app-theme-spooky-orange);
		border-radius: 1.5em;
		padding: 2em;
		color: var(--app-theme-midnight);
		overflow: hidden;
	}

</style>
