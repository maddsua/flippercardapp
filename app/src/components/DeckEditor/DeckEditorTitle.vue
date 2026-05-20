<script setup lang="ts">

interface MetaEntry {
	name: string;
	description?: string | null;
}

const props = defineProps<{
	meta: MetaEntry;
}>();

</script>

<template>
	<div class="deck-editor-title" >
		<div class="group">
			<div class="name">
				{{ meta.name }}
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

			.name, .description {
				max-width: 100%;
				overflow: hidden;
				text-overflow: ellipsis;
				white-space: nowrap;
			}

			.name {
				font-size: 0.75rem;
				font-weight: 600;
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
