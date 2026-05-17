<script setup lang="ts">
import { computed } from 'vue';
import type { CardDeckMetadata } from '../../../../api_models';
import GenericButton from '../../../App/GenericButton.vue';

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
			<div class="details">
				<div class="row">
					<div class="name">
						{{ entry.name }}
					</div>
					<div class="stats">
						{{ date }}
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
		</div>

		<div class="actions">
			<GenericButton variant="thin" @click="emit('edit')">
				Open editor
			</GenericButton>
			<GenericButton variant="thin" theme="red" @click="emit('delete')">
				Delete
			</GenericButton>
		</div>
	
	</div>
</template>

<style lang="scss" scoped>
	.deck {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		gap: 2rem;
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		background-color: var(--app-theme-ghostly-glow);

		.details {
			display: flex;
			flex-direction: column;
			gap: 0.25rem;
			flex-grow: 1;
		}

		.row {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 1rem;
			flex-grow: 1;
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
		}
	}

</style>
