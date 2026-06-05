<script lang="ts" setup>

const props = defineProps<{
	required?: boolean;
	disabled?: boolean;
	label?: string;
}>();

const model = defineModel<boolean>();

</script>

<template>
	<label class="checkbox" :class="{ disabled }">
		<input type="checkbox" :required="props.required" v-model="model" :disabled="disabled" />
		<span v-if="label" class="label-text">
			{{ label }}
		</span>
	</label>
</template>

<style lang="scss" scoped>

label {
	display: flex;
	flex-flow: row nowrap;
	gap: 0.25rem;
	align-items: center;
	padding: 0.25rem;

	input {
		display: block;
		width: 0.75rem;
		height: 0.75rem;
		border-radius: 0.25rem;
		border: 1px solid var(--grey-200);
		position: relative;

		&:checked {
			background-color: var(--purple-200);
			border: none;

			&::after {
				position: absolute;
				left: 0;
				right: 0;
				width: 100%;
				height: 100%;
				display: block;
				content: '';
				mask-image: url(/src/assets/icons/check-mask.svg);
				mask-size: 0.75rem;
				mask-position: center;
				mask-repeat: no-repeat;
				background-color: white;
			}
		}

		&:hover {
			cursor: pointer;
		}
	}

	.label-text {
		font-size: 0.8rem;
		line-height: 1.5em;
		user-select: none;
		font-weight: 500;

		&:hover {
			cursor: pointer;
		}
	}

	&.disabled {
		opacity: 0.8;
		filter: saturate(0);
		pointer-events: none;
	}
}

</style>
