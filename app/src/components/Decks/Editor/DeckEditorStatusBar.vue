<script setup lang="ts">
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import type { ResourceVisibility } from '@/api_models';
import DeckMetaInfo from './Toolbar/DeckMetaInfo.vue';

interface DeckMeta {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

const props = defineProps<{
	meta: DeckMeta;
	edited?: boolean;
	published?: boolean;
	valid?: boolean;
}>();

const emit = defineEmits<{
	(e: 'disacard'): void;
	(e: 'publish'): void;
	(e: 'import'): void;
	(e: 'export'): void;
	(e: 'versions'): void;
}>();

</script>

<template>
	<div class="status-bar">

		<div class="wrapper">

			<DeckMetaInfo :meta="props.meta" />

			<div class="publish-actions">
				<GenericButton variant="thin" theme="orange" @click="emit('disacard')">
					<template v-if="edited">
						Discard
					</template>
					<template v-else>
						Exit
					</template>
				</GenericButton>
				<GenericButton variant="thin" theme="blue" @click="emit('versions')">
					Versions
				</GenericButton>
				<GenericButton variant="thin" theme="blue" @click="emit('import')">
					Import
				</GenericButton>
				<GenericButton variant="thin" theme="blue" :disabled="!valid" @click="emit('export')">
					Export
				</GenericButton>
				<GenericButton variant="thin" theme="blue" :disabled="!valid || !edited" @click="emit('publish')">
					Publish
				</GenericButton>
			</div>

		</div>
	</div>
</template>

<style lang="scss" scoped>
	.status-bar {
		position: relative;
		display: flex;
		justify-content: center;
		z-index: 100;

		.wrapper {
			position: relative;
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 1rem;
			width: 100%;
			max-width: 50rem;
			padding: 0.5rem 1rem;
			background-color: var(--app-theme-carbon);
			border-bottom-left-radius: 0.75rem;
			border-bottom-right-radius: 0.75rem;
		}

		.publish-actions {
			display: flex;
			align-items: center;
			justify-content: end;
			gap: 1rem;
		}
	}
</style>
