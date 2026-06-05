<script setup lang="ts">
import type { AuthActor } from '../../api_models';
import GenericButton from '../App/Inputs/GenericButton.vue';

const props = defineProps<{
	actor: AuthActor;
	error?: string | null;
}>();

const emits = defineEmits<{
	(e: 'signout'): void;
}>();

</script>

<template>
	<div class="session-state">

		<div class="state-row">
			<div class="actor-name">
				<span class="label">
					Signed in as
				</span>
				<span class="value">
					{{ actor.name }}
				</span>
			</div>
			<GenericButton theme="orange" variant="thin" @click="emits('signout')">
				Sign out
			</GenericButton>
		</div>

		<div v-if="error" class="state-row">
			<div class="error-message">
				{{ error }}
			</div>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.session-state {
		display: flex;
		flex-direction: column;
		gap: 1rem;

		.state-row {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			justify-content: center;
		}

		.actor-name {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			font-size: 0.95rem;
			gap: 0.5em;

			.label {
				font-size: 0.75em;
				font-weight: 400;
			}

			.value {
				font-weight: 600;
			}
		}

		.error-message {
			color: var(--app-theme-blood-red);
			font-weight: 600;
			font-size: 0.75rem;
			width: 100%;
			max-width: 25rem;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}
	}
</style>
