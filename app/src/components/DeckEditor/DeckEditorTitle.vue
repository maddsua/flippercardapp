<script setup lang="ts">
import type { ResourceVisibility } from '../../api_models';


interface MetaEntry {
	name: string;
	description?: string | null;
	visibility: ResourceVisibility;
}

const props = defineProps<{
	meta: MetaEntry;
}>();

</script>

<template>
	<div class="deck-editor-title" >
		<div class="group">
			<div class="title">
				<div class="name">
					{{ meta.name }}
				</div>
				<div class="visibility" :class="[ meta.visibility.toLowerCase() ]"></div>
			</div>
			<div class="description">
				<template v-if="meta.description">
					{{ meta.description }}
				</template>
				<template v-else>
					[No description]
				</template>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.deck-editor-title {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.75rem;
		align-items: center;
		width: 100%;
		min-width: 0;
		padding: 0.25rem 0.5rem;
		border-radius: 0.5rem;
		user-select: none;

		.group {
			display: flex;
			flex-direction: column;
			gap: 0.25rem;
			max-width: 100%;

			.title {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				justify-content: center;
				width: 100%;
				gap: 1rem;
			}

			.name, .description {
				min-width: 0;
				overflow: hidden;
				text-overflow: ellipsis;
				white-space: nowrap;
			}

			.name {
				font-size: 0.75rem;
				font-weight: 600;
			}

			.visibility {
				display: block;
				width: 1rem;
				height: 1rem;
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
			
			.description {
				font-size: 0.65rem;
				font-weight: 400;
			}
		}

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-ghostly-glow);
		}
	}
</style>
