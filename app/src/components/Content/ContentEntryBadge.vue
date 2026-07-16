<script setup lang="ts">
const props = defineProps<{
	icon?: 'cards' | 'decks' | 'score' | 'star' | 'star-filled' | 'world' | 'link' | 'lock' | 'clock' | 'completion' | 'card-stars';
}>();
</script>

<template>
	<div class="badge" :class="{ [`icon-${icon}`]: !!icon, icon: icon }">
		<span v-if="$slots.default">
			<slot />
		</span>
	</div>
</template>

<style lang="scss" scoped>
	.badge {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		gap: 0.25em;
		flex-shrink: 0;

		span {
			white-space: nowrap;
			font-size: 1em;
			font-weight: 600;
		}

		&.icon::before {
			content: "";
			display: block;
			width: 1.35em;
			height: 1.35em;
			background-size: contain;
			background-repeat: no-repeat;
			background-position: center;
		}

		&.icon-cards::before {
			background-image: url(/src/assets/icons/quiz-card-mask.svg);
		}
		&.icon-decks::before {
			background-image: url(/src/assets/icons/card-deck-mask.svg);
		}
		&.icon-score::before {
			background-image: url(/src/assets/icons/target-mask.svg);
		}
		&.icon-star::before {
			background-image: url(/src/assets/icons/star-mask.svg);
		}
		&.icon-star-filled::before {
			background-image: url(/src/assets/icons/star-filled-mask.svg);
		}
		&.icon-world::before {
			background-image: url(/src/assets/icons/world-mask.svg);
		}
		&.icon-link::before {
			background-image: url(/src/assets/icons/link-mask.svg);
		}
		&.icon-lock::before {
			background-image: url(/src/assets/icons/lock-mask.svg);
		}
		&.icon-clock::before {
			background-image: url(/src/assets/icons/clock-mask.svg);
		}
		&.icon-completion::before {
			background-image: url(/src/assets/icons/check-mask.svg);
		}
		&.icon-card-stars::before {
			background-image: url(/src/assets/icons/card-stars-mask.svg);
		}
	}
</style>
