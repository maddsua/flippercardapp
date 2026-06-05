<script setup lang="ts">
import { reactive } from 'vue';
import GenericButton from '../App/Inputs/GenericButton.vue';
import type { SignInParams } from '@/api_models';
import GenericInput from '../App/Inputs/GenericInput.vue';
import InputGroup from '../App/Inputs/InputGroup.vue';
import InlineErorrMessage from '../App/Messages/InlineErorrMessage.vue';

const props = defineProps<{
	locked?: boolean;
	error?: string | null;
}>();

const emit = defineEmits<{
	(e: 'signin', params: SignInParams): void;
}>();

const state = reactive<SignInParams>({
	username: '',
	password: ''
});

</script>

<template>
	<div class="signin-form" :class="{ locked }">

		<div class="header">
			<div class="title">
				Account sign-in
			</div>
		</div>

		<InputGroup >

			<GenericInput type="text" placeholder="Username" v-model="state.username" />
			<GenericInput type="password" placeholder="Password" v-model="state.password" />

			<InlineErorrMessage v-if="error">
				{{ error }}
			</InlineErorrMessage>

			<GenericButton variant="wide" @click="emit('signin', state)">Sign in</GenericButton>

		</InputGroup>

	</div>
</template>

<style lang="scss" scoped>

	.signin-form {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		width: 100%;
		max-width: 25rem;

		&.locked {
			pointer-events: none;
			filter: saturate(0);
		}

		.header {
			display: flex;
			flex-flow: row nowrap;
			justify-content: center;
			flex-grow: 1;

			.title {
				font-size: 1.125rem;
				font-weight: 600;
			}
		}
	}
	
</style>
