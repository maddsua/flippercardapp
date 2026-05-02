<script setup lang="ts">
import CardBody from './CardBody.vue';
import CardPoll from './CardPoll.vue';
import CardTextBox from './CardTextBox.vue';
import CardTextNode from './CardTextNode.vue';
import CardTitle from './CardTitle.vue';
import type { CardSide, PollOption } from './content';

const props = defineProps<{
	entry: CardSide;
}>();

const emit = defineEmits<{
	(e: 'flip'): void;
	(e: 'next'): void;
}>();

const handlePollOptionSelect = (opt: PollOption) => {
	switch (opt.action) {
		case 'fail':
			//	todo: show as failed
			break;
		case 'fail-show':
			//	todo: show as failed and flip
			setTimeout(() => emit('flip'), 500);
			break;
		default:
			emit('next');
			break;
	}
};

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
			<CardPoll v-else-if="node.type === 'poll'" :entry="node" @select="handlePollOptionSelect" />
		</template>
	</CardBody>
</template>
