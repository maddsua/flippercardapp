<script setup lang="ts">

import { computed } from 'vue';
import type { ResourceVisibility } from '@/api_models';

interface DeckMeta {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

const props = defineProps<{
	meta: DeckMeta;
	changed?: boolean;
	changesSaved?: boolean;
}>();

const nameInvalid = computed(() => !props.meta.name.trim().length);

//	todo: set icons

</script>

<template>
	<div class="deck-meta-info" >

		<div class="row">

			<div class="visibility-icon" :class="[ meta.visibility.toLowerCase() ]"></div>

			<input type="text"
				class="name"
				:class="{ invalid: nameInvalid }"
				v-model="props.meta.name"
				placeholder="Deck name (required)" />

			<div class="autosave-state" :class="{ changed, saved: changesSaved }">
				<template v-if="changesSaved">
					Changes saved locally
				</template>
				<template v-else-if="changed">
					Saving changes...
				</template>
				<template v-else>
					No local changes
				</template>
			</div>
		</div>

		<div class="description">
			<input type="text"
				class="description"
				v-model="props.meta.description"
				placeholder="[No description]" />
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.deck-meta-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		width: 25rem;
		max-width: 100%;
		min-width: 0;
		padding: 0.25rem 0.5rem;
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
		
		.autosave-state {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.25rem;
			align-items: center;
			flex-shrink: 0;
			font-size: 0.5rem;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;

			&::before {
				content: "";
				display: block;
				width: 0.65rem;
				height: 0.65rem;
				background-color: white;
				mask-type: alpha;
				mask-size: contain;
				mask-repeat: no-repeat;
				mask-position: center;
				mask-image: url(/src/assets/icons/broom-mask.svg);
			}

			&.changed::before {
				mask-image: url(/src/assets/icons/refresh-mask.svg);
			}

			&.saved::before {
				mask-image: url(/src/assets/icons/check-mask.svg);
			}
		}
	}
</style>
