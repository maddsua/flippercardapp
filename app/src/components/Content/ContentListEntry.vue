<script setup lang="ts">
import { computed } from 'vue';
import type { ResourceVisibility } from '@/api_models';
import ContentEntryBadge from './ContentEntryBadge.vue';

const props = defineProps<{
	title: string;
	summary?: string | null;
	date?: string | null;
	visibility?: ResourceVisibility;
	starrable?: boolean | null;
	starred?: boolean | null;
	cardCount?: number | null;
	deckCount?: number | null;
	score?: number | null;
	completion?: number | null;
}>();

const newBadgeThreshold = 7 * 24 * 60 * 60 * 1000;

const showNewBadge = computed(() => {

	if (!props.date) {
		return false;
	}

	try {
		return new Date().getTime() - new Date(props.date).getTime() < newBadgeThreshold;
	} catch (_) {
		return false;
	}
});

</script>

<template>

	<button type="button" class="content-list-entry" :class="{ 'badge-new': showNewBadge }">

		<div class="row-group">

			<div class="entry-title">
				{{ title }}
			</div>

			<div v-if="starred || starrable || visibility || score || completion || cardCount || deckCount" class="entry-badges">

				<ContentEntryBadge v-if="completion" icon="completion">
					{{ completion.toFixed(0) }}%
				</ContentEntryBadge>

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

		<div v-if="summary" class="entry-summary">
			{{ summary }}
		</div>

	</button>

</template>

<style lang="scss" scoped>

	@use '@/media.scss';

	.content-list-entry {
		position: relative;
		user-select: none;
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
		border-radius: 1rem;

		background-color: var(--app-theme-irish-green);
		color: var(--app-theme-snow-white);

		@media (orientation: portrait) {
			padding: 1rem 1.25rem;
		}

		@include media.non-sticky-hover {
			background-color: var(--app-theme-spooky-orange);
		}

		&.badge-new::before {
			content: "NEW!";
			position: absolute;
			left: 0;
			top: 0;
			font-size: 0.55rem;
			font-weight: 600;
			padding: 0.25em;
			background-color: var(--app-theme-spooky-orange);
			color: var(--app-theme-snow-white);
			border-radius: 0.5em;
			transform: rotate(-25deg);
			z-index: 10;
			box-shadow: 0 0 0.5rem rgba(0, 0, 0, 0.25);
		}

		.row-group {
			display: flex;
			flex-flow: row nowrap;
			align-items: start;
			width: 100%;
			gap: 1rem;
		}

		.entry-title {
			font-size: 1rem;
			font-weight: 600;
			width: 100%;
		}

		.entry-summary {
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
	}

</style>
