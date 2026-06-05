<script setup lang="ts">
import { computed } from 'vue';
import { UploadReadyState, type CardContentNode } from '@/content';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import CardTitleNodeEditor from './CardTitleNodeEditor.vue';
import CardTextNodeEditor from './CardTextNodeEditor.vue';
import CardPollEditor from './CardPollEditor.vue';
import CardImageNodeEditor from './CardImageNodeEditor.vue';

const props = defineProps<{
	isFront?: boolean;
}>();

const model = defineModel<CardContentNode[]>();

type NodeType = CardContentNode['type'];

const nodeSortOrder: NodeType[] = ['title', 'image', 'textbox', 'poll'];

const addNode = (type: NodeType) => {

	if (!model.value) {
		throw new Error('Editor model binding may not be omitted');
	}

	const sortIndex = new Map<NodeType, number>(nodeSortOrder.map((item, idx) => ([item, idx])));

	switch (type) {
		case 'title':
			model.value?.push({ type: 'title', content: '' });
			break;
		case 'image':
			model.value?.push({ type: 'image', state: UploadReadyState.Idle });
			break;
		case 'textbox':
			model.value?.push({ type: 'textbox', content: [] });
			break;
		case 'poll':
			model.value?.push({ type: 'poll', content: [] });
			break;
	}

	model.value.sort((a, b) => (sortIndex.get(a.type) ?? sortIndex.size) - (sortIndex.get(b.type) ?? sortIndex.size));
};

const removeNode = (idx: number) => model.value?.splice(idx, 1);

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
	<div class="card-face-content-editor">
		<div class="title">
			Face content
		</div>
		<div class="tree">

			<template v-for="(item, idx) of model">

				<CardTitleNodeEditor v-if="item.type === 'title'"
					v-model="item.content"
					@remove="removeNode(idx)" />

				<CardImageNodeEditor v-else-if="item.type === 'image'"
					v-model="item.media_id"
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
