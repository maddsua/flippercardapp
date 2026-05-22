<script setup lang="ts">
const props = defineProps<{
	title: string;
	summary?: string;
	starrable?: boolean;
	starred?: boolean;
	cardCount?: number;
	deckCount?: number;
	score?: number;
}>();
</script>

<template>
	<button type="button" class="content-list-entry">

		<div class="flex-group">

			<div class="title">
				{{ title }}
			</div>

			<template v-if="starred || starrable || score || cardCount || deckCount">

				<div class="stats">
					<div v-if="score" class="item score">
						{{ score.toFixed(0) }}%
					</div>
					<div v-if="cardCount" class="item cards">
						{{ cardCount.toFixed(0) }}
					</div>
					<div v-if="deckCount" class="item decks">
						{{ deckCount.toFixed(0) }}
					</div>
					<div class="item star" :class="{ starred }"></div>
				</div>

			</template>

		</div>

		<div v-if="summary" class="summary">
			{{ summary }}
		</div>

	</button>
</template>

<style lang="scss" scoped>
	.content-list-entry {
		display: flex;
		display: flex;
		flex-direction: column;
		align-items: start;
		gap: 0.25rem;
		padding: 1rem 2rem;
		border-radius: 1rem;
		border: none;
		outline: none;
		transition: all 150ms ease;
		text-align: start;
		text-align: unset;

		background-color: var(--app-theme-irish-green);
		color: var(--app-theme-snow-white);

		@media (orientation: portrait) {
			padding: 1rem 1.5rem;
		}

		.flex-group {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			justify-content: space-between;
			width: 100%;
		}

		.title {
			font-size: 1rem;
			font-weight: 600;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			width: 100%;
		}

		.summary {
			font-size: 0.85rem;
			font-weight: 400;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
			width: 100%;
		}

		hr {
			display: block;
			width: 100%;
			height: 1px;
			background-color: var(--app-theme-powder-trail);
			border: none;
			outline: none;
		}

		.stats {
			display: flex;
			flex-flow: row nowrap;
			flex-shrink: 0;
			gap: 0.5rem;

			.item {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				white-space: nowrap;
				gap: 0.25em;
				font-size: 0.75rem;
				flex-shrink: 0;
				font-weight: 600;

				&::before {
					content: "";
					display: block;
					width: 1.35em;
					height: 1.35em;
					background-size: contain;
					background-repeat: no-repeat;
					background-position: center;
				}

				&.cards::before {
					background-image: url(/src/assets/icons/quiz-card-mask.svg);
				}

				&.decks::before {
					background-image: url(/src/assets/icons/card-deck-mask.svg);
				}

				&.score::before {
					background-image: url(/src/assets/icons/target-mask.svg);
				}

				&.star {
					
					&::before {
						background-image: url(/src/assets/icons/star-mask.svg);
					}

					&.starred::before {
						background-image: url(/src/assets/icons/star-filled-mask.svg);
					}
				}
			}
		}

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-spooky-orange);
		}
	}
</style>
