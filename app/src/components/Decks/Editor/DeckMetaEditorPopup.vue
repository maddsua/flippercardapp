<script setup lang="ts">

import type { ResourceVisibility } from '@/api_models';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import InputLabel from '@/components/App/Inputs/InputLabel.vue';
import { isContainerOutsideClick } from '@/dom';
import { resourceVisibilityOptions } from '@/inputs';
import { computed, onMounted, reactive, ref } from 'vue';

interface DeckMeta {
	name: string;
	description: string | null;
	visibility: ResourceVisibility;
};

const props = defineProps<{
	meta: DeckMeta;
}>();

const emit = defineEmits<{
	(e: 'edit', meta: DeckMeta): void;
	(e: 'done'): void;
}>();

const state = reactive({
	shown: false,
	inputs: {
		name: props.meta.name,
		description: props.meta.description || '',
		visibilty: props.meta.visibility,
	}
});

const isEdited = computed(() =>
	state.inputs.name !== props.meta.name ||
	state.inputs.description !== (props.meta.description || '') ||
	state.inputs.visibilty !== props.meta.visibility);

const openPrompt = () => {
	state.shown = true;
	window.addEventListener('click', handleOutsideClicks);
};

const closePrompt = () => {
	state.shown = false;
	setTimeout(() => emit('done'), 200);
	window.removeEventListener('click', closePrompt);
};

const rootElement = ref<HTMLElement | null>(null);

const handleOutsideClicks = (event: Event) =>
	isContainerOutsideClick(rootElement.value, event.target) ? closePrompt() : null;

onMounted(() => setTimeout(openPrompt, 10));

const edit = () => {
	emit('edit', {
		name: state.inputs.name,
		description: state.inputs.description || null,
		visibility: state.inputs.visibilty,
	});
	closePrompt();
};

</script>

<template>
	<div class="meta-editor">
		<div class="wrapper" :class="{ hidden: !state.shown }" ref="rootElement">

			<div class="header">
				<div class="title">
					Edit deck
				</div>
				<div class="actions">
					<GenericButton variant="thin" theme="orange" @click="closePrompt">
						Discard
					</GenericButton>
					<GenericButton variant="thin" theme="blue" :disabled="!isEdited" @click="edit">
						Save
					</GenericButton>
				</div>
			</div>

			<div class="form">

				<InputLabel variant="slick">
					<template v-slot:label>
						Deck name
					</template>
					<GenericInput type="text" placeholder="Name" v-model="state.inputs.name" />
				</InputLabel>

				<InputLabel variant="slick">
					<template v-slot:label>
						Deck name
					</template>
					<GenericInput type="text" placeholder="Description" :multiline="true" v-model="state.inputs.description" />
				</InputLabel>

				<InputLabel variant="slick">
					<template v-slot:label>
						Visibility
					</template>
					<GenericDropdown :options="resourceVisibilityOptions" v-model="state.inputs.visibilty" />
				</InputLabel>

			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.meta-editor {
		position: absolute;
		left: 0;
		bottom: 0;
		width: 100%;
		width: 20rem;

		& > .wrapper {
			position: absolute;
			left: 0;
			top: 1rem;
			width: 100%;
			display: flex;
			flex-direction: column;
			gap: 1rem;
			background-color: var(--app-theme-carbon);
			border: 1px solid var(--app-theme-powder-trail);
			padding: 1rem;
			border-radius: 0.5rem;
			transition: all 150ms ease;

			&.hidden {
				transform: translateY(-5rem);
				opacity: 0;
			}
		}

		.header {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 1rem;
			justify-content: space-between;

			.title {
				font-size: 0.85rem;
				color: var(--app-theme-mysterious-white);
			}

			.actions {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 0.5rem;
			}
		}

		.form {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
		}
	}
</style>
