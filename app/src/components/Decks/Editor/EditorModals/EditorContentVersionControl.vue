<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { unwrapErrorMessage, useClient } from '@/api';
import type { CardDeckVersionMetadata } from '@/api_models';
import { genericPageState, pageControls } from '@/dataloader';
import ErrorMessage from '@/components/App/Messages/ErrorMessage.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import EditorModal from '../EditorModal.vue';

const client = useClient();

const props = defineProps<{
	deckID: string;
}>();

const emit = defineEmits<{
	(e: 'done'): void;
	(e: 'rollback', version: CardDeckVersionMetadata): void;
}>();

const state = reactive({
	page: genericPageState<CardDeckVersionMetadata>(),
	rollback: {
		busy: false,
		error: null as string | null,
	}
});

const { reload } = pageControls(state.page, (page) => client.decks.versions(props.deckID, page));

onMounted(reload);

const fmtDate = (date: string) => new Date(date).toLocaleDateString('en-UK', {
	year: 'numeric',
	month: 'short',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric',
});

const rollbackVersion = async (versionID: string) => {

	if (!confirm('Roll back to this version?')) {
		return;
	}

	state.rollback.busy = true;

	const { data, error } = await client.decks.rollbackVersion(props.deckID, versionID);
	if (!data || error) {
		state.rollback.error = unwrapErrorMessage(error);
		return;
	}

	state.rollback.busy = false;

	emit('rollback', data);
	emit('done');
};

</script>

<template>
	<EditorModal title="Content versions" variant="narrow" @close="emit('done')">

		<div class="version-control">

			<ErrorMessage v-if="state.page.error">
				<template v-slot:message>
					Unable to load versions
				</template>
				<template v-slot:details>
					{{ state.page.error }}
				</template>
			</ErrorMessage>

			<LoadingMessage v-else-if="!state.page.ready" />

			<div v-else-if="!state.page.entries.length" class="list-placeholder">
				No versions available
			</div>

			<ErrorMessage v-else-if="state.rollback.error">
				<template v-slot:message>
					Unable to roll back a version
				</template>
				<template v-slot:details>
					{{ state.rollback.error }}
				</template>
			</ErrorMessage>

			<LoadingMessage v-else-if="state.rollback.busy">
				Rolling back...
			</LoadingMessage>

			<div v-else class="version-list">

				<div v-for="entry of state.page.entries" class="entry">

					<div class="summary">
						<div class="date">
							{{ fmtDate(entry.created) }}
						</div>
						<div class="label">
							{{ entry.label || 'No label' }}
						</div>
					</div>

					<div class="actions">
						<div v-if="entry.is_latest" class="flag-latest">Latest</div>
						<GenericButton v-else theme="orange" variant="thin" @click="rollbackVersion(entry.id)">
							Rollback
						</GenericButton>
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
