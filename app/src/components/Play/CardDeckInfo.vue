<script setup lang="ts">
import { useLanguage, intl } from '@/intl';

const props = defineProps<{
	labels: string[]
	size: number;
	index: number;
	isMarked?: boolean;
	shareable?: boolean;
}>();

const emit = defineEmits<{
	(e: 'exit'): void;
	(e: 'toggleMarked'): void;
	(e: 'share'): void;
}>();

const lang = useLanguage();

</script>

<template>
	<div class="deck-info">

		<div class="deck-progress">
			<div v-for="idx of size" class="marker" :class="{ filled: idx < index + 1 }"></div>
		</div>

		<div class="details-row">

			<div class="deck-meta-actions">
				<button type="button" class="icon exit" title="Exit game" @click="emit('exit')"></button>
			</div>

			<div class="deck-labels">
				<template v-for="(item, idx) of labels">
					<template v-if="idx > 0">
						<hr />
					</template>
					<span>
						{{ item }}
					</span>
				</template>
			</div>

			<div class="deck-meta-actions">

				<button v-if="shareable" type="button" class="labeled with-icon share" @click="emit('share')">
					{{ intl(lang, {
						en: 'Share',
						de: 'Teilen',
						uk: 'Поширити'
					}) }}
				</button>

				<button type="button" class="labeled with-icon mark" :class="{ marked: isMarked }" @click="emit('toggleMarked')">
					<template v-if="isMarked">
						{{ intl(lang, {
							en: 'Unmark',
							de: 'Markiert',
							uk: 'Збережено'
						}) }}
					</template>
					<template v-else>
						{{ intl(lang, {
							en: 'Mark',
							de: 'Markieren',
							uk: 'Зберігти'
						}) }}
					</template>
				</button>

			</div>

		</div>

	</div>
</template>

<style lang="scss" scoped>
	.deck-info {
		position: absolute;
		left: 0.5rem;
		top: 0.5rem;
		right: 0.5rem;
		z-index: 10;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		user-select: none;

		.deck-progress {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.2rem;

			.marker {
				flex-grow: 1;
				height: 2px;
				background-color: rgba(255, 255, 255, 0.35);

				&.filled {
					background-color: white;
				}
			}
		}

		.details-row {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.75rem;
			padding: 0 0.5rem;
			align-items: center;
			justify-content: space-between;
		}

		.deck-labels {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 0.5rem;
			min-width: 0;
			overflow: hidden;

			span {
				display: block;
				font-size: 0.75rem;
				white-space: nowrap;
			}

			hr {
				display: block;
				width: 3px;
				height: 3px;
				border-radius: 100%;
				background-color: white;
				border: none;
				outline: none;
				margin: 0;
				padding: 0;
				flex-shrink: 0;
			}
		}

		.deck-meta-actions {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			justify-content: end;
			flex-shrink: 0;
			gap: 0.5rem;

			button.labeled {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 0.25rem;
				color: var(--app-theme-snow-white);
				background-color: var(--app-theme-powder-trail);
				font-weight: 600;
				border: none;
				outline: none;
				border-radius: 0.25rem;
				padding: 0.25rem 0.5rem;
				font-size: 0.7rem;

				&:hover {
					cursor: pointer;
				}

				&.with-icon::before {
					content: "";
					display: block;
					width: 0.85rem;
					height: 0.85rem;
					mask-type: alpha;
					mask-size: contain;
					mask-position: center;
					mask-repeat: no-repeat;
					background-color: var(--app-theme-snow-white);
				}

				&.mark {
					background-color: var(--app-theme-irish-green);

					&::before {
						mask-image: url(/src/assets/icons/star-mask.svg);
					}

					&.marked {
						background-color: var(--app-theme-spooky-orange);

						&::before {
							mask-image: url(/src/assets/icons/star-filled-mask.svg);
						}
					}
				}

				&.share::before {
					mask-image: url(/src/assets/icons/share-arrow-mask.svg);
				}
			}

			button.icon {
				display: block;
				width: 1.25rem;
				height: 1.25rem;
				outline: none;
				border: none;
				background: none;
				background-size: contain;
				background-repeat: no-repeat;
				background-position: center;
				background-image: url(/src/assets//icons/cross-cut-mask.svg);

				&:hover {
					cursor: pointer;
				}
			}
		}
	}
</style>
