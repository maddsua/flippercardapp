<script setup lang="ts">
import CardPoll from './CardPoll.vue';
import CardTextBox from './CardTextBox.vue';
import CardTextNode from './CardTextNode.vue';
import CardTitle from './CardTitle.vue';
import type { CardContentFace } from '../../content';
import { computed, type CSSProperties } from 'vue';
import CardImage from './CardImage.vue';

const props = defineProps<{
	entry: CardContentFace;
	decoration?: 'question-mark';
	is3dBackface?: boolean;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
	(e: 'score', score: number): void;
}>();

const canvasStyle = computed((): CSSProperties => ({
	color: props.entry.theme?.card?.mask_color
}));

const canvasClasses = computed(() => ({
	[`decoration-${props.decoration}`]: !!props.decoration,
	backface: props.is3dBackface,
}));

const contentStyle = computed((): CSSProperties => {

	const { card } = props.entry.theme || {};

	return {
		backgroundColor: card?.fill_color,
		borderColor: card?.outline_color || card?.fill_color,
		color: card?.mask_color,
	};
});

</script>

<template>
	<div class="card-canvas" :style="canvasStyle" :class="canvasClasses">

		<div class="card-content" :style="contentStyle">

			<template v-for="node of entry.content">

				<CardTitle v-if="node.type === 'title'">
					{{ node.content || '[Title]' }}
				</CardTitle>

				<CardImage v-else-if="node.type === 'image'" :entry="node" />

				<CardTextBox v-else-if="node.type === 'textbox'">
					<template v-for="txtnode of node.content">
						<CardTextNode v-if="txtnode.type === 'text' && txtnode.content.length" :theme="txtnode.theme">
							{{ txtnode.content }}
						</CardTextNode>
						<br v-else-if="txtnode.type === 'newline'" />
					</template>
				</CardTextBox>

				<CardPoll v-else-if="node.type === 'poll'"
					:entry="node"
					:theme="entry.theme?.interactives"
					@score="(score) => emit('score', score)"
					@flip="emit('flip')"
					@next="emit('next')" />

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
		box-shadow: 0.5rem 0.5rem 2rem rgba(0, 0, 0, 0.25);

		&.backface {
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

		// helps prevent container corners from poking out on some mobile browsers
		will-change: transform;
	}

</style>
