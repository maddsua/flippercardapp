<script setup lang="ts">

import { computed } from 'vue';
import type { ResourceVisibility } from '@/api_models';
import { blurInteractive } from '@/dom';
import type { ContentSummary } from '@/content';

interface DeckMeta {
	summary: ContentSummary;
	visibility: ResourceVisibility;
};

const props = defineProps<{
	meta: DeckMeta;
	showDescription?: boolean;
}>();

const nameInvalid = computed(() => !props.meta.summary.name.trim().length);

</script>

<template>
	<div class="deck-summary" >

		<div class="row">

			<div class="visibility-icon" :class="[ meta.visibility.toLowerCase() ]"></div>

			<input type="text"
				class="name"
				:class="{ invalid: nameInvalid }"
				v-model="props.meta.summary.name"
				placeholder="Deck name"
				@keydown.enter.stop="blurInteractive" />

		</div>

		<input v-if="showDescription" type="text"
			class="description"
			v-model="props.meta.summary.description"
			placeholder="[No description]"
			@keydown.enter.stop="blurInteractive" />

	</div>
</template>

<style lang="scss" scoped>
	.deck-summary {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		width: 15rem;
		max-width: 100%;
		min-width: 0;
		border-radius: 0.5rem;
		user-select: none;

		.row {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			width: 100%;
			gap: 0.5rem;
		}

		input {
			display: block;
			width: 100%;
			border: 1px solid transparent;
			outline: none;
			border-radius: 0.25rem;
			background: unset;
			padding: 0.125rem;

			&:focus {
				border-color: var(--app-theme-powder-trail);
			}

			&.invalid {
				border-color: red;
			}

			&.name {
				font-size: 0.75rem;
				font-weight: 600;
			}

			&.description {
				font-size: 0.65rem;
				font-weight: 400;
			}
		}

		.visibility-icon {
			display: block;
			width: 0.75rem;
			height: 0.75rem;
			flex-shrink: 0;
			background-position: center;
			background-repeat: no-repeat;
			background-size: contain;

			&.public {
				background-image: url(/src/assets/icons/world-mask.svg);
			}

			&.hidden {
				background-image: url(/src/assets/icons/link-mask.svg);
			}

			&.private {
				background-image: url(/src/assets/icons/lock-mask.svg);
			}
		}
	}
</style>
