<script setup lang="ts">
import type { ResourceVisibility } from '@/api_models';
import { resourceVisibilityOptions } from '@/inputs';
import type { CardNode } from '@/content';
import { blurInteractive } from '@/dom';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';
import EditorModal from '../EditorModal.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';

interface Content {
	summary: ContentSummary;
	cards: CardNode[]
};

interface ContentSummary {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

interface OriginMetadata {
	created: string | null;
	updated: string | null;
}

const props = defineProps<{
	content: Content;
	origin: OriginMetadata;
}>();

const emit = defineEmits<{
	(e: 'done'): void;
}>();

const fmtDate = (date: string | null) => {

	if (!date) {
		return 'N/A';
	}

	return new Date(date).toLocaleString('en-UK', {
		year: 'numeric',
		month: 'numeric',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
		second: 'numeric',
	});
};

</script>

<template>
	<EditorModal title="Deck details" variant="compact" @close="emit('done')">

		<div class="deck-details">

			<div class="summary-form">

				<InputLabel variant="slick">

					<template v-slot:label>
						Deck name
					</template>

					<GenericInput placeholder="Deck name"
						variant="borderless"
						v-model="content.summary.name"
						@keydown.enter.stop="blurInteractive"
						@keydown.escape.stop="blurInteractive" />

				</InputLabel>

				<InputLabel variant="slick">

					<template v-slot:label>
						Deck summary
					</template>

					<GenericInput placeholder="Deck summary"
						variant="borderless"
						:multiline="true"
						v-model="content.summary.description"
						@keydown.enter.stop="blurInteractive"
						@keydown.escape.stop="blurInteractive" />

				</InputLabel>

				<InputLabel variant="slick">

					<template v-slot:label>
						Deck visibility
					</template>

					<GenericDropdown :options="resourceVisibilityOptions" v-model="content.summary.visibility" />

				</InputLabel>

			</div>

			<div class="deck-metadata">
				<div class="property">
					<span class="label">
						Created:
					</span>
					{{ fmtDate(origin.created) }}
				</div>
				<div class="property">
					<span class="label">
						Updated:
					</span>
					{{ fmtDate(origin.updated) }}
				</div>
			</div>

			<GenericButton variant="thin-wide" @click="emit('done')">
				OK
			</GenericButton>

		</div>

	</EditorModal>
</template>

<style lang="scss" scoped>
	.deck-details {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		max-height: 100%;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding-right: 0.5rem;

		.summary-form {
			display: flex;
			flex-direction: column;
			gap: 1.5rem;
		}

		.deck-metadata {
			display: flex;
			flex-flow: row wrap;
			gap: 0.75rem;
			align-items: center;

			.property {
				font-size: 0.75rem;
				white-space: nowrap;
				color: var(--app-theme-kinda-white);

				.label {
					font-weight: 300;
					color: var(--app-theme-mysterious-white);
				}
			}
		}
	}
</style>
