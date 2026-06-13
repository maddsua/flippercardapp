<script setup lang="ts">
import { computed } from 'vue';
import { UploadReadyState, type CardContentNode } from '@/content';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import EditableTitleNode from './EditableTitleNode.vue';
import EditableTextNode from './EditableTextNode.vue';
import EditablePollNode from './EditablePollNode.vue';
import EditableImageNode from './EditableImageNode.vue';
import EditorGroup from '../EditorGroup.vue';

const props = defineProps<{
	isFront?: boolean;
}>();

const model = defineModel<CardContentNode[]>();

type NodeType = CardContentNode['type'];

const nodeSortOrder: NodeType[] = ['title', 'image', 'textbox', 'poll'];

const nodeOrder = new Map<NodeType, number>(nodeSortOrder.map((item, idx) => ([item, idx])));

const createNode = (type: NodeType): CardContentNode | null => {
	switch (type) {
		case 'title':
			return { type: 'title', content: '' };
		case 'image':
			return { type: 'image', state: UploadReadyState.Idle };
		case 'textbox':
			return { type: 'textbox', content: [] };
		case 'poll':
			return { type: 'poll', content: [] };
		default:
			return null;
	}
};

const addNode = (type: NodeType) => {

	if (!model.value) {
		throw new Error('Editor model binding may not be omitted');
	}

	const next = createNode(type);
	if (!next) {
		throw new Error(`Unknown node type: ${type}`);
	}

	const nextPosition = nodeOrder.get(type) ?? nodeOrder.size;

	for (let idx = 0; idx < model.value.length; idx++) {
		const node = model.value[idx];
		const nodePosition = nodeOrder.get(node.type) ?? 0;
		if (nodePosition > nextPosition) {
			model.value.splice(idx, 0, next);
			return;
		}
	}

	model.value.push(next);
};

const removeNode = (idx: number) => model.value?.splice(idx, 1);

const reorderNode = (idx: number, delta: number) => {

	const newIdx = idx + delta;
	if (!model.value || newIdx < 0 || newIdx >= (model.value?.length ?? 0)) {
		return;
	}

	const node = model.value[idx];

	model.value[idx] = model.value[newIdx];
	model.value[newIdx] = node;
};

const availableNodes = computed((): Record<NodeType, boolean> => {
	const typeSet = new Set(model.value?.map(item => item.type));
	return {
		title: !typeSet.has('title'),
		image: !typeSet.has('image'),
		textbox: !typeSet.has('textbox'),
		poll: props.isFront && !typeSet.has('poll'),
	};
});

</script>

<template>

	<EditorGroup title="Content">

		<div class="content-tree">
			<template v-for="(item, idx) of model">

				<EditableTitleNode v-if="item.type === 'title'"
					v-model="item.content"
					@up="reorderNode(idx, -1)"
					@down="reorderNode(idx, 1)"
					@remove="removeNode(idx)" />

				<EditableImageNode v-else-if="item.type === 'image'"
					v-model="item.media_id"
					@up="reorderNode(idx, -1)"
					@down="reorderNode(idx, 1)"
					@remove="removeNode(idx)" />

				<EditableTextNode v-else-if="item.type === 'textbox'"
					v-model="item.content"
					@up="reorderNode(idx, -1)"
					@down="reorderNode(idx, 1)"
					@remove="removeNode(idx)" />

				<EditablePollNode v-else-if="item.type === 'poll'"
					v-model="item.content"
					@up="reorderNode(idx, -1)"
					@down="reorderNode(idx, 1)"
					@setQuizFlag="flag => item.is_quiz = flag"
					@remove="removeNode(idx)" />
			</template>
		</div>

		<div class="add-actions">

			<GenericButton v-if="availableNodes.title" variant="thin" @click="addNode('title')">
				+ Add title
			</GenericButton>

			<GenericButton v-if="availableNodes.image" variant="thin" @click="addNode('image')">
				+ Add image
			</GenericButton>

			<GenericButton v-if="availableNodes.textbox" variant="thin" @click="addNode('textbox')">
				+ Add textbox
			</GenericButton>

			<GenericButton v-if="availableNodes.poll" variant="thin" @click="addNode('poll')">
				+ Add poll
			</GenericButton>

		</div>

	</EditorGroup>

</template>

<style lang="scss" scoped>

	.content-tree {
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

</style>
