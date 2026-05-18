<script setup lang="ts">
import { computed } from 'vue';
import type { CardContentElement } from '../../../content';
import GenericButton from '../../App/GenericButton.vue';
import CardTitleNodeEditor from './CardTitleNodeEditor.vue';
import CardTextNodeEditor from './CardTextNodeEditor.vue';
import CardPollEditor from './CardPollEditor.vue';

const props = defineProps<{
	isFront?: boolean;
}>();

const model = defineModel<CardContentElement[]>();

type NodeType = 'title' | 'textbox' | 'poll';

const addNode = (type: NodeType) => {

	if (!model.value) {
		throw new Error('Editor model binding may not be omitted');
	}

	const entries = model.value || [];

	const lastTitle = entries.findLastIndex(item => item.type === 'title');
	const lastTextbox = entries.findLastIndex(item => item.type === 'textbox');
	const lastPoll = entries.findLastIndex(item => item.type === 'poll');

	const insert = (idx: number, entry: CardContentElement) => {

		if (!model.value?.length) {
			model.value?.push(entry);
			return
		}

		if (idx < 0) {
			model.value?.unshift(entry);
			return
		} else if (idx >= model.value.length) {
			model.value?.push(entry);
			return
		}

		model.value?.splice(idx + 1, 0, entry);
	};

	switch (type) {
		case 'title':
			insert(lastTitle, { type: 'title', content: '' });
			break;
		case 'textbox':
			insert(Math.max(lastTitle, lastTextbox), { type: 'textbox', content: [] });
			break;
		case 'poll':
			insert(Math.max(lastTitle, lastTextbox, lastPoll), { type: 'poll', content: [] });
			break;
	}
};

const removeNode = (idx: number) => model.value?.splice(idx, 1);

const availableNodes = computed((): Record<NodeType, boolean> => {
	const typeSet = new Set(model.value?.map(item => item.type));
	return {
		title: !typeSet.has('title'),
		textbox: !typeSet.has('textbox'),
		poll: props.isFront && !typeSet.has('poll'),
	};
});

</script>

<template>
	<div class="card-face-content-editor">
		<div class="title">
			Face content
		</div>
		<div class="tree">

			<template v-for="(item, idx) of model">

				<CardTitleNodeEditor v-if="item.type === 'title'"
					v-model="item.content"
					@remove="removeNode(idx)" />

				<CardTextNodeEditor v-else-if="item.type === 'textbox'"
					v-model="item.content"
					@remove="removeNode(idx)" />

				<CardPollEditor v-else-if="item.type === 'poll'"
					v-model="item.content"
					@setQuizFlag="flag => item.is_quiz = flag"
					@remove="removeNode(idx)" />
			</template>
			
			<div class="add-actions">

				<GenericButton v-if="availableNodes.title" variant="thin" @click="addNode('title')">
					+ Add title
				</GenericButton>

				<GenericButton v-if="availableNodes.textbox" variant="thin" @click="addNode('textbox')">
					+ Add textbox
				</GenericButton>

				<GenericButton v-if="availableNodes.poll" variant="thin" @click="addNode('poll')">
					+ Add poll
				</GenericButton>

			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.card-face-content-editor {
		display: flex;
		flex-direction: column;
		gap: 2rem;

		.title {
			font-size: 0.75rem;
			font-weight: 600;
		}

		.tree {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			height: 100%;
		}

		.add-actions {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 1.25rem;
		}
	}
</style>
