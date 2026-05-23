<script setup lang="ts">
import type { CardTextboxElementTheme } from '../../content';

const props = defineProps<{
	theme?: CardTextboxElementTheme;
}>();
</script>

<template>
	<span class="card-text-node" :class="{ italic: theme?.italic, bold: theme?.bold, [`${theme?.decoration}`]: !!theme?.decoration }" :style="{ backgroundColor: theme?.highlight?.fill_color, color: theme?.highlight?.text_color }">
		<slot>
			[Text]
		</slot>
	</span>
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
	}
</style>
