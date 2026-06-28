<script setup lang="ts">
import { computed, type CSSProperties } from 'vue';
import type { CardTextboxElementTheme } from '../../content';

const props = defineProps<{
	theme?: CardTextboxElementTheme;
}>();

const elementClass = computed((): Record<string, boolean> => ({
	italic: !!props.theme?.italic,
	bold: !!props.theme?.bold,
	[`${props.theme?.decoration}`]: !!props.theme?.decoration,
	[`size-${props.theme?.size}`]: !!props.theme?.size,
}));

const elementStyle = computed((): CSSProperties => ({
	backgroundColor: props.theme?.highlight?.fill_color || undefined,
	color: props.theme?.highlight?.text_color || undefined,
}));

</script>

<template>
	<div class="card-text-node" :class="elementClass" :style="elementStyle">
		<slot>
			[Text]
		</slot>
	</div>
</template>

<style lang="scss" scoped>

	.card-text-node {
		display: inline;
		font-weight: 500;
		font-size: 1.8em;
		line-height: 1.25em;

		&.italic {
			font-style: italic;
		}

		&.bold {
			font-weight: 700;
		}

		&.underline {
			text-decoration: underline;
			text-decoration-thickness: 0.125em;
		}

		&.strikethrough {
			text-decoration: line-through;
			text-decoration-thickness: 0.25em;
		}

		&.size-xs {
			font-size: 1.25em;
		}

		&.size-s {
			font-size: 1.5em;
		}

		&.size-m {
			font-size: 1.75em;
		}

		&.size-l {
			font-size: 2em;
		}

		&.size-xl {
			font-size: 2.25em;
		}
	}

</style>
