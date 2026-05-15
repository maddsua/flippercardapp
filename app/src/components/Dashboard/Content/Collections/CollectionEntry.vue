<script setup lang="ts">
import { computed, ref } from 'vue';
import type { CardDeckMetadata, CollectionMetadata } from '../../../../api_models';
import LoadingMessage from '../../../App/LoadingMessage.vue';
import DeckEntry from './DeckEntry.vue';

interface Entry extends CollectionMetadata {
	decks?: CardDeckMetadata[] | null;
};

const props = defineProps<{
	entry: Entry;
}>();

const decksShown = ref(false);

const emit = defineEmits<{
	(e: 'edit'): void;
	(e: 'showDecks'): void;
	(e: 'addDeck'): void;
	(e: 'editDeck', id: string): void;
	(e: 'deleteDeck', id: string): void;
}>();

const showDecks = () => {
	decksShown.value = true;
	emit('showDecks');
};

const date = computed(() => new Date(props.entry.updated).toLocaleDateString('en-UK', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
}));

//	todo: make this ui not suck

</script>

<template>
	<div class="collection">

		<div class="header">

			<div class="row">
				<div class="name">
					{{ entry.name }}
				</div>
				<div class="stats">
					{{ date }}
				</div>
				<div class="actions">

					<button type="button" @click="emit('addDeck')">
						+ Add deck
					</button>

					<button type="button" @click="emit('edit')">
						Edit
					</button>

				</div>
			</div>

			<div class="row">
				<div class="description">
					<template v-if="entry.description">
						{{ entry.description }}
					</template>
					<template v-else>
						[No description provided]
					</template>
				</div>
			</div>

		</div>

		<template v-if="entry.size > 0">

			<div v-if="!decksShown" class="decks-trigger">
				<button type="button" @click="showDecks">
					Show {{ entry.size }} deck(s)
				</button>
			</div>

			<LoadingMessage v-else-if="decksShown && !entry.decks?.length">
				Loading decks...
			</LoadingMessage>

			<template v-else>

				<div v-if="entry.decks?.length" class="deck-list">
					<DeckEntry v-for="deck of entry.decks" :key="deck.id" :entry="deck"
						@edit="emit('editDeck', deck.id)"
						@delete="emit('deleteDeck', deck.id)" />
				</div>
				
			</template>
			
		</template>

	</div>
</template>

<style lang="scss" scoped>
	.collection {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 0.5rem 0.75rem;
		border-radius: 0.25rem;
		background-color: var(--app-theme-ghostly-glow);

		.row {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 1rem;
		}

		.header {
			flex-grow: 1;

			.name {
				font-size: 0.95rem;
				font-weight: 600;
				flex-grow: 1;
			}

			.description {
				font-size: 0.75rem;
				font-weight: 400;
				color: var(--app-theme-mysterious-white);
			}

			.name, .description {
				flex-grow: 1;
				overflow: hidden;
				white-space: nowrap;
				text-overflow: ellipsis;
			}

			.stats {
				font-size: 0.75rem;
				color: var(--app-theme-mysterious-white);
				flex-shrink: 0;
			}

			.actions {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 0.5rem;
				flex-shrink: 0;
	
				button {
					display: block;
					background-color: unset;
					border: none;
					color: var(--app-theme-sky-blue);
					font-size: 0.85rem;
	
					&:hover {
						cursor: pointer;
						color: var(--app-theme-deep-lavender);
					}
				}
			}
		}

		.decks-trigger {
			display: flex;
			justify-content: center;

			button {
				display: block;
				background-color: unset;
				border: none;
				color: var(--app-theme-spooky-orange);
				font-size: 0.75rem;

				&:hover {
					cursor: pointer;
					color: var(--app-theme-snow-white);
				}
			}
		}

		.deck-list {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
			padding: 0.5rem 1rem;
		}
	}
</style>
