<script setup lang="ts">
import CardPollOption from './CardPollOption.vue';
import type { PollNode, PollOption } from './content';

const props = defineProps<{
	entry: PollNode;
}>();

const emit = defineEmits<{
	(e: 'select', opt: PollOption): void;
}>();

</script>

<template>
	<div class="card-poll">
		<CardPollOption v-for="option of entry.content" :entry="option" @select="emit('select', option)" />
	</div>
</template>

<style lang="scss" scoped>
	.card-poll {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		width: 100%;

		&:last-child {
			margin-top: auto;
		}
	}
</style>
