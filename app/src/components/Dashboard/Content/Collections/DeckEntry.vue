<script setup lang="ts">
import { computed } from 'vue';
import type { CardDeckMetadata } from '../../../../api_models';

const props = defineProps<{
	entry: CardDeckMetadata;
}>();

const emit = defineEmits<{
	(e: 'edit'): void;
	(e: 'delete'): void;
}>();

const date = computed(() => new Date(props.entry.updated).toLocaleDateString('en-UK', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric'
}));

</script>

<template>
	<div class="deck">

		<div class="row">
			<div class="name">
				{{ entry.name }}
			</div>
			<div class="stats">
				{{ date }}
			</div>
			<div class="actions">
				<button type="button" @click="emit('edit')">
					Open editor
				</button>
				<button type="button" class="danger" @click="emit('delete')">
					Delete
				</button>
			</div>
		</div>

		<div class="description">
			<template v-if="entry.description">
				{{ entry.description }}
			</template>
			<template v-else>
				[No description provided]
			</template>
		</div>
	
	</div>
</template>

<style lang="scss" scoped>
	.deck {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		background-color: var(--app-theme-ghostly-glow);

		.row {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 1rem;
		}

		.name {
			font-size: 0.85rem;
			font-weight: 600;
			flex-grow: 1;
		}

		.description {
			font-size: 0.65rem;
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

			button {
				display: block;
				background-color: unset;
				border: none;
				color: var(--app-theme-sky-blue);
				font-size: 0.75rem;

				&:hover {
					cursor: pointer;
					color: var(--app-theme-deep-lavender);
				}

				&.danger {
					color: var(--app-theme-blood-red);

					&:hover {
						color: var(--app-theme-spooky-orange);
					}
				}
			}
		}
	}

</style>
