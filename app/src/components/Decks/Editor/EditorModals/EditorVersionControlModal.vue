<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeckVersion, CardDeckVersionMetadata } from '@/api_models';
import { genericPageState, pageControls } from '@/dataloader';
import { fmtTimeString } from '@/date';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import EditorModal from '../EditorModal.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';

const client = useClient();

const props = defineProps<{
	deckID: string;
}>();

const emit = defineEmits<{
	(e: 'done'): void;
	(e: 'pull', version: CardDeckVersion): void;
}>();

const state = reactive({
	page: genericPageState<CardDeckVersionMetadata>(),
	actions: {
		busy: false,
		error: null as string | null,
	},
});

const { reload } = pageControls(state.page, (page) => client.decks.versions.list(props.deckID, page));

onMounted(reload);

const deleteVersion = async (versionID: string) => {

	if (!confirm('Delete this version?')) {
		return;
	}

	state.actions.busy = true;

	const { error } = await client.decks.versions.remove(props.deckID, versionID);
	if ( error) {
		state.actions.busy = false;
		state.actions.error = unwrapErrorMessage(error);
		return;
	}

	state.page.entries = state.page.entries.filter(item => item.id != versionID);
	state.actions.busy = false;
};

const loadVersion = async (versionID: string) => {

	if (!confirm('Replace current content with this version?')) {
		return;
	}

	state.actions.busy = true;

	const { data, error } = await client.decks.versions.load(props.deckID, versionID);
	if (!data || error) {
		state.actions.busy = false;
		state.actions.error = unwrapErrorMessage(error);
		return;
	}

	state.actions.busy = false;

	emit('pull', data);
	emit('done');
};

</script>

<template>
	<EditorModal title="Content versions" variant="narrow" @close="emit('done')">

		<div class="version-control">

			<InlineErrorMessage v-if="state.page.error">
				<template v-slot:title>
					Unable to load versions
				</template>
				{{ state.page.error }}
			</InlineErrorMessage>

			<LoadingMessage v-else-if="!state.page.ready" />

			<div v-else-if="!state.page.entries.length" class="list-placeholder">
				No versions available
			</div>

			<InlineErrorMessage v-else-if="state.actions.error">
				{{ state.actions.error }}
			</InlineErrorMessage>

			<LoadingMessage v-else-if="state.actions.busy">
				Updating versions...
			</LoadingMessage>

			<div v-else class="version-list">

				<div v-for="entry of state.page.entries" class="entry">

					<div class="summary">
						<div class="date">
							{{ fmtTimeString(entry.created) }}
						</div>
						<div class="label">
							{{ entry.label || 'No label' }}
						</div>
					</div>

					<div class="actions">

						<div v-if="entry.is_latest" class="flag-latest">Latest</div>

						<template v-else>
							<GenericButton variant="thin" :disabled="state.actions.busy" @click="loadVersion(entry.id)">
								Load
							</GenericButton>
							<GenericButton theme="red" variant="thin" :disabled="state.actions.busy" @click="deleteVersion(entry.id)">
								Delete
							</GenericButton>
						</template>

					</div>

				</div>
			</div>

		</div>

	</EditorModal>
</template>

<style lang="scss" scoped>
	.version-control {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		min-height: 0;

		.list-placeholder {
			font-size: 0.8rem;
			color: var(--app-theme-mysterious-white);
		}

		.version-list {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: start;
			align-items: start;
			gap: 0.5rem;
			height: 100%;
			width: 100%;
			min-height: 0;
			overflow: hidden auto;
			scrollbar-width: thin;

			.entry {
				display: flex;
				flex-flow: row nowrap;
				justify-content: space-between;
				align-items: center;
				width: 100%;
				gap: 2rem;
				padding: 0.75rem 1rem;
				border-radius: 0.25rem;
				background-color: var(--app-theme-ghostly-glow);

				.summary {
					display: flex;
					flex-direction: column;
					gap: 0.25rem;

					.date {
						font-size: 0.7rem;
					}

					.label {
						font-size: 0.65rem;
						min-width: 0;
						white-space: nowrap;
						overflow: hidden;
						text-overflow: ellipsis;
					}
				}

				.actions {
					display: flex;
					flex-flow: row nowrap;
					gap: 0.5rem;
					flex-shrink: 0;

					.flag-latest {
						padding: 0.25rem 0.5rem;
						border-radius: 1rem;
						background-color: var(--app-theme-irish-green);
						color: var(--app-theme-snow-white);
						font-weight: 600;
						font-size: 0.65rem;
						pointer-events: none;
					}
				}
			}
		}
	}
</style>
