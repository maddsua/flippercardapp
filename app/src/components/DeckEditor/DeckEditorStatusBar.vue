<script setup lang="ts">
import { ref } from 'vue';
import GenericButton from '../App/GenericButton.vue';
import DeckMetaEditorPopup from './DeckMetaEditorPopup.vue';
import DeckEditorTitle from './DeckEditorTitle.vue';

interface DeckDetails {
	name: string;
	description: string | null;
};

const props = defineProps<{
	meta: DeckDetails;
	edited?: boolean;
	published?: boolean;
	valid?: boolean;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'updateDetails', details: DeckDetails): void;
	(e: 'disacard'): void;
	(e: 'publish'): void;
	(e: 'import'): void;
	(e: 'export'): void;
}>();

const metaEditorOpen = ref(false);

</script>

<template>
	<div class="status-bar">

		<div class="wrapper">

			<DeckEditorTitle :meta="meta" @click="metaEditorOpen = true" />

			<div class="view-actions">
				<button type="button" class="icon flip" title="Flip view" @click="emit('flip')"></button>
			</div>

			<div class="publish-actions">
				<GenericButton variant="thin" theme="orange" @click="emit('disacard')">
					<template v-if="edited">
						Discard
					</template>
					<template v-else>
						Exit
					</template>
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

			<DeckMetaEditorPopup v-if="metaEditorOpen" :deck="meta"
				@done="metaEditorOpen = false"
				@edit="val => emit('updateDetails', val)" />

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
			grid-template-columns: 1fr 1fr 1fr;
			gap: 1rem;
			width: 100%;
			max-width: 50rem;
			padding: 0.5rem 1rem;
			background-color: var(--app-theme-carbon);
			border-bottom-left-radius: 0.75rem;
			border-bottom-right-radius: 0.75rem;
		}

		.view-actions {
			display: flex;
			align-items: center;
			justify-content: center;
		}

		.publish-actions {
			display: flex;
			align-items: center;
			justify-content: end;
			gap: 1rem;
		}

		button.icon {
			display: block;
			width: 1.25rem;
			height: 1.25rem;
			flex-shrink: 0;
			mask-type: alpha;
			mask-size: contain;
			mask-position: center;
			mask-repeat: no-repeat;
			background-color: rgba(255, 255, 255, 0.8);
			
			&:hover {
				cursor: pointer;
				background-color: white;
			}
			
			&.flip {
				width: 2rem;
				height: 2rem;
				mask-image: url(/src/assets/icons/flip-mask.svg);
			}
		}
	}
</style>
