<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useClient } from '../../api';
import type { AuthState, SignInParams } from '../../api_models';
import ErrorMessage from '../App/Messages/ErrorMessage.vue';
import LoadingMessage from '../App/Messages/LoadingMessage.vue';
import DashboardSigninForm from './DashboardSigninForm.vue';
import DashboardSessionState from './DashboardSessionState.vue';

const emit = defineEmits<{
	(e: 'stateUpdate', state: AuthState): void;
}>();

enum Stage {
	Idle,
	Signin,
	Authorized,
};

const state = reactive({
	data: null as AuthState | null,
	ready: false,
	busy: false,
	stage: Stage.Idle,
	error: null as string | null,
});

const client = useClient();

const checkAuthState = async () => {

	const { data, error } = await client.auth.whoami();
	if (!data || error) {
		state.error = error?.message || 'Unable to check auth status';
		return;
	}

	if (data.session) {
		state.stage = Stage.Authorized
	}

	state.data = data;

	emit('stateUpdate', state.data);
};

const retryCheckAuthState = async () => {

	await checkAuthState();

	if (state.error) {
		setTimeout(retryCheckAuthState, 15_000);
	}

	state.ready = true;
};

onMounted(retryCheckAuthState);

const attemptSignin = async (params: SignInParams) => {

	state.error = null;
	state.data = null;

	const { data, error } = await client.auth.signin(params);
	if (!data || error) {
		state.error = error?.message || 'Invalid credentials';
		return;
	}

	if (data.session) {
		state.stage = Stage.Authorized
	}
	
	state.data = data;

	emit('stateUpdate', state.data);
}

const handleSignin = async (params: SignInParams) => {

	state.busy = true;
	state.stage = Stage.Signin;

	await attemptSignin(params);

	state.busy = false;
};

const handleSignout = async () => {

	if (!confirm('You sure want to quit?')) {
		return;
	}

	const { data, error } = await client.auth.signout()
	if (!data || error) {
		state.error = error?.message || 'Unable to sign out';
		return;
	}

	state.data = data;
	state.stage = Stage.Idle;
	state.error = null;

	emit('stateUpdate', state.data);
};

</script>

<template>
	<div class="dashboard-session" :class="{ expanded: state.stage !== Stage.Authorized }">

		<ErrorMessage v-if="state.error && state.stage === Stage.Idle">
			<template v-slot:message>
				Auth error
			</template>
			<template v-slot:details>
				{{ state.error }}
			</template>
		</ErrorMessage>

		<LoadingMessage v-else-if="!state.ready || state.busy">
			One second...
		</LoadingMessage>

		<DashboardSessionState v-else-if="state.stage === Stage.Authorized && state.data?.actor"
			:error="state.error"
			:actor="state.data.actor"
			@signout="handleSignout" />

		<DashboardSigninForm v-else :locked="state.busy" :error="state.error" @signin="handleSignin" />

	</div>
</template>

<style lang="scss" scoped>
	.dashboard-session {
		display: flex;
		flex-direction: column;
		
		&.expanded {
			align-items: center;
			justify-content: center;
			flex-grow: 1;
		}
	}
</style>
