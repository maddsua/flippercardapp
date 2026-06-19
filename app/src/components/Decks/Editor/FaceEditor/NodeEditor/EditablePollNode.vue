<script setup lang="ts">
import type { CardPollNodeOption } from '@/content';
import EditableNodeHarness from './EditableNodeHarness.vue';
import EditablePollNodeOption from './EditablePollNodeOption.vue';

const model = defineModel<CardPollNodeOption[]>();

const emit = defineEmits<{
	(e: 'setQuizFlag', flag: boolean): void;
}>();

const setOptionCheck = (markIdx: number) => {

	if (!model.value) {
		return;
	}

	model.value?.forEach((val, idx) => val.is_answer = idx === markIdx);

	emit('setQuizFlag', true);
};

const addOption = () => {
	if (!model.value) {
		model.value = [];
	}
	model.value.push({ value: '' });
};

const removeOption = (idx: number) => {
	const removed = model.value?.splice(idx, 1);
	if (removed?.some(item => item.is_answer)) {
		emit('setQuizFlag', false);
	}
};

</script>

<template>
	<EditableNodeHarness>
		<template v-slot:title>
			Poll / Quiz
		</template>

		<template v-slot:content>

			<div class="option-list">
				<template v-if="model?.length">
					<EditablePollNodeOption v-for="(option, idx) of model"
						:checked="!!option.is_answer"
						v-model="option.value"
						@check="setOptionCheck(idx)"
						@remove="removeOption(idx)" />
				</template>
				<template v-else>
					<div class="no-option-message">
						No options added yet
					</div>
				</template>
			</div>

			<div v-if="!model?.length || model.length < 4" class="list-actions">
				<button type="button" @click="addOption">+ Add option</button>
			</div>

		</template>
	</EditableNodeHarness>
</template>

<style lang="scss" scoped>
	.option-list {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		width: 20rem;
		max-width: 100%;

		.no-option-message {
			font-size: 0.75rem;
			color: var(--app-theme-mysterious-white);
			text-align: center;
		}
	}

	.list-actions {
		display: flex;
		justify-content: center;

		button {
			display: block;
			background-color: unset;
			border: none;
			color: var(--app-theme-deep-lavender);
			font-weight: 600;
			font-size: 0.75rem;

			&:hover {
				cursor: pointer;
				color: var(--app-theme-snow-white);
			}
		}
	}
</style>
