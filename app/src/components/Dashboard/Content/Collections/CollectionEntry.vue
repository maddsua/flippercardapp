<script setup lang="ts">
import { computed } from 'vue';
import type { CollectionMetadata } from '../../../../api_models';
import GenericButton from '../../../App/GenericButton.vue';

const props = defineProps<{
	entry: CollectionMetadata;
}>();

const emit = defineEmits<{
	(e: 'manage'): void;
}>();

const date = computed(() => new Date(props.entry.updated).toLocaleDateString('en-UK', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
}));

</script>

<template>
	<div class="collection">

		<div class="header">

			<div class="row">
				<div class="name">
					{{ entry.name }}
				</div>
				<div class="date">
					{{ date }}
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
				<div class="size">
					{{ entry.size }} deck(s)
				</div>
			</div>

		</div>

		<div class="actions">
			<GenericButton variant="thin" @click="emit('manage')">
				Manage
			</GenericButton>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.collection {
		display: flex;
		flex-flow: row nowrap;
		justify-content: space-between;
		gap: 2rem;
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

			.date, .size {
				font-size: 0.75rem;
				color: var(--app-theme-mysterious-white);
				flex-shrink: 0;
			}
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 0.5rem;
			flex-shrink: 0;
		}
	}
</style>
