<script setup lang="ts">
import type { ResourceVisibility } from '../../api_models';
import GenericButton from '../App/Inputs/GenericButton.vue';
import ContentEntryBadge from './ContentEntryBadge.vue';

const props = defineProps<{
	title: string;
	summary?: string | null;
	visibility?: ResourceVisibility;
	starrable?: boolean | null;
	starred?: boolean | null;
	cardCount?: number | null;
	deckCount?: number | null;
	score?: number | null;
	editable?: boolean;
	playable?: boolean;
}>();

const emit = defineEmits<{
	(e: 'click'): void;
	(e: 'edit'): void;
	(e: 'delete'): void;
}>();

</script>

<template>

	<div class="content-list-entry" :class="{ editable }">

		<div v-if="editable" class="editor-actions">
			<div class="actions-row">
				<GenericButton variant="thin" theme="green" @click="emit('click')">
					<template v-if="playable">
						Play
					</template>
					<template v-else>
						Inspect
					</template>
				</GenericButton>
				<GenericButton variant="thin" @click="emit('edit')">
					Edit
				</GenericButton>
				<GenericButton variant="thin" theme="red" @click="emit('delete')">
					Delete
				</GenericButton>
			</div>
		</div>

		<button type="button" class="primary-action" @click="emit('click')">
	
			<div class="row-group">
	
				<div class="title">
					{{ title }}
				</div>
	
				<div v-if="starred || starrable || visibility || score || cardCount || deckCount" class="entry-badges">

					<ContentEntryBadge v-if="score" icon="score">
						{{ score.toFixed(0) }}%
					</ContentEntryBadge>

					<ContentEntryBadge v-if="cardCount" icon="cards">
						{{ cardCount.toFixed(0) }}
					</ContentEntryBadge>

					<ContentEntryBadge v-if="deckCount" icon="decks">
						{{ deckCount.toFixed(0) }}
					</ContentEntryBadge>

					<template v-if="visibility">
						<ContentEntryBadge v-if="visibility === 'PRIVATE'" icon="lock" />
						<ContentEntryBadge v-else-if="visibility === 'HIDDEN'" icon="link" />
					</template>

					<template v-if="starred || starrable">
						<ContentEntryBadge v-if="starred" icon="star-filled" />
						<ContentEntryBadge v-else icon="star" />
					</template>

				</div>

			</div>
	
			<div v-if="summary" class="summary">
				{{ summary }}
			</div>
	
		</button>

	</div>

</template>

<style lang="scss" scoped>

	@use '@/media.scss';

	.content-list-entry {
		position: relative;
		border-radius: 1rem;
		overflow: hidden;
		user-select: none;

		.editor-actions {
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
			z-index: 2;
			padding: 1rem;
			background-color: rgba(0, 0, 0, 0.4);
			opacity: 0;
			
			.actions-row {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				justify-content: end;
				gap: 1rem;
				transition: all 150ms ease;
				width: 100%;
				height: 100%;
			}
		}

		.primary-action {
			display: flex;
			flex-direction: column;
			align-items: start;
			gap: 0.25rem;
			width: 100%;
			padding: 1rem 1.5rem;
			border: none;
			outline: none;
			transition: all 150ms ease;
			text-align: start;
			text-align: unset;

			background-color: var(--app-theme-irish-green);
			color: var(--app-theme-snow-white);

			@media (orientation: portrait) {
				padding: 1rem 1.25rem;
			}

			@include media.non-sticky-hover {
				background-color: var(--app-theme-spooky-orange);
			}
		}

		&.editable {

			.primary-action {
				pointer-events: none;
			}

			@media (orientation: portrait) {
				.actions-row {
					visibility: collapse;
				}

				&:hover .actions-row {
					visibility: unset;
				}
			}
			
			&:hover {
				.primary-action {
					pointer-events: none;
					filter: blur(2px);
					opacity: 0.75;
				}
				.editor-actions {
					opacity: 1;
				}
			}
		}
	}

	.row-group {
		display: flex;
		flex-flow: row nowrap;
		align-items: start;
		width: 100%;
		gap: 1rem;
	}

	.title {
		font-size: 1rem;
		font-weight: 600;
		width: 100%;
	}

	.summary {
		font-size: 0.85rem;
		font-weight: 400;
		width: 100%;
	}

	.entry-badges {
		display: flex;
		flex-flow: row nowrap;
		flex-shrink: 0;
		gap: 0.5rem;
	}

</style>
