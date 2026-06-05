<script setup lang="ts">
import { unwrapErrorMessage, useClient } from '@/api';
import type { AuthActor } from '@/api_models';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import LoadingMessage from '@/components/App/Messages/LoadingMessage.vue';
import { computed, onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';

const client = useClient();
const router = useRouter();

const state = reactive({
	actor: null as AuthActor | null,
	ready: false,
	error: null as string | null,
});

onMounted(async () => {
	const { data, error } = await client.auth.whoami();
	if (!data || error) {
		state.error = unwrapErrorMessage(error);
		state.ready = true;
		return;
	}
	state.actor = data.actor || null;
	state.ready = true;
});

const permissions = computed(() => Array.from(Object.entries(state.actor?.permissions || {})).filter(([_, set]) => set).map(([key]) => key));

const signin = () => {
	router.push('/dashboard/signin');
};

const signout = async () => {

	if (!confirm('You sure want to quit?')) {
		return;
	}

	state.ready = false;

	const { data, error } = await client.auth.signout()
	if (!data || error) {
		state.ready = true;
		state.error = unwrapErrorMessage(error);
		return;
	}

	state.actor = null;
	state.error = null;
	state.ready = true;
};

</script>

<template>

	<div class="account-settings">

		<LoadingMessage v-if="!state.ready" />

		<InlineErrorMessage v-else-if="state.error">

			<template v-slot:title>
				Unable to check auth status
			</template>

			{{ state.error }}

		</InlineErrorMessage>

		<div v-else class="session-state">

			<div v-if="state.actor" class="session-info">
				<div class="actor-id">
					<div class="summary">
						Authorized as <span class="username">{{ state.actor.name }}</span>
					</div>
					<GenericButton theme="orange" variant="thin" @click="signout">
						Sign out?
					</GenericButton>
				</div>
				<div class="actor-permissions">
					<span v-for="item of permissions">
						{{ item }}
					</span>
				</div>
			</div>

			<div v-else class="signin-prompt">
				
				<div class="status-message">
					Not signed in
				</div>

				<GenericButton variant="thin" @click="signin">
					Sign in
				</GenericButton>
			</div>

		</div>
	
	</div>

</template>

<style lang="scss" scoped>

	.account-settings {
		display: flex;
		flex-direction: column;
		gap: 2.5rem;
	}

	.session-state {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		padding: 0.75rem 1.5rem;
		background-color: var(--app-theme-ghostly-glow);
		border-radius: 0.75rem;

		.session-info {
			display: flex;
			flex-direction: column;
			gap: 1rem;

			.actor-id {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				justify-content: space-between;
				gap: 1rem;

				.summary {
					display: flex;
					flex-flow: row nowrap;
					align-items: center;
					font-size: 0.75rem;
					gap: 0.5rem;

					.username {
						display: block;
						padding: 0.125rem 0.25rem;
						border-radius: 0.25rem;
						line-height: 1em;
						font-weight: 600;
						background-color: var(--app-theme-kinda-white);
						color: var(--app-theme-carbon);
					}
				}
			}

			.actor-permissions {
				display: flex;
				flex-flow: row wrap;
				gap: 0.5rem;

				span {
					display: block;
					padding: 0.125rem 0.5rem;
					border-radius: 0.5rem;
					font-size: 0.65rem;
					font-weight: 600;
					color: var(--app-theme-snow-white);
					background-color: var(--app-theme-rich-mint);
				}
			}
		}
	}

	.signin-prompt {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;

		.status-message {
			font-size: 0.75rem;
		}
	}

</style>
