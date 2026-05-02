<script setup lang="ts">
import CardBody from './CardBody.vue';
import CardPoll from './CardPoll.vue';
import CardTextBox from './CardTextBox.vue';
import CardTextNode from './CardTextNode.vue';
import CardTitle from './CardTitle.vue';
import type { CardSideNode } from './content';

const props = defineProps<{
	entry: CardSideNode;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
}>();

</script>

<template>
	<CardBody>
		<template v-for="node of entry.content">
			<CardTitle v-if="node.type === 'title'">
				{{ node.content }}
			</CardTitle>
			<CardTextBox v-else-if="node.type === 'textbox'">
				<template v-for="txtnode of node.content">
					<CardTextNode v-if="txtnode.type === 'text'">
						{{ txtnode.content }}
					</CardTextNode>
					<br v-else-if="txtnode.type === 'newline'" />
				</template>
			</CardTextBox>
			<CardPoll v-else-if="node.type === 'poll'" :entry="node" @flip="emit('flip')" @next="emit('next')" />
		</template>
	</CardBody>
</template>
