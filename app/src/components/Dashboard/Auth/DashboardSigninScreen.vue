<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue';
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import { unwrapErrorMessage, useClient } from '@/api';
import { useRouter } from 'vue-router';
import InlineErrorMessage from '@/components/App/Messages/InlineErrorMessage.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';
import InputGroup from '@/components/App/Inputs/InputGroup.vue';
import GenericInput from '@/components/App/Inputs/GenericInput.vue';
import type { AuthSession } from '@/api_models';
import { useLanguage, intl } from '@/intl';

const client = useClient();
const router = useRouter();
const lang = useLanguage();

const state = reactive({
	inputs: {
		username: '',
		password: '',
	},
	ready: false,
	busy: false,
	session: null as AuthSession | null,
	error: null as string | null,
});

const valid = computed(() => state.inputs.username.length > 0 && state.inputs.password.length > 0)

const backHref = '/dashboard';
const exitScreen = () => router.push(backHref);

onMounted(async () => {

	const { data, error } = await client.auth.whoami();
	if (!data || error) {
		state.error = unwrapErrorMessage(error);
		return;
	}

	if (data.session?.id) {
		state.session = data.session || null;
		exitScreen();
		return;
	}

	state.ready = true;
});

const signin = async () => {

	state.error = null;
	state.session = null;
	state.busy = true;

	await checkCredentials();

	state.busy = false;

	if (state.session) {
		exitScreen();
	}
};

const checkCredentials = async () => {

	const { data, error } = await client.auth.signin(state.inputs);
	if (!data || error) {
		state.error = error?.message || 'Invalid credentials';
		return;
	}

	state.session = data.session || null;
};

</script>

<template>

	<AppUiHeader :backHref="backHref">
		<template v-slot:title>
			{{ intl(lang, {
				en: 'Account',
				de: 'Konto',
				uk: 'Обліковий запис'
			}) }}
		</template>
	</AppUiHeader>

	<div class="signin-screen">

		<div class="signin-form" :class="{ locked: state.busy }">

			<div class="header">
				<div class="message-title">
					{{ intl(lang, {
						en: 'Account sign-in',
						de: 'Anmeldung am Konto',
						uk: 'Вхід в обліковий запис'
					}) }}
				</div>
			</div>

			<InputGroup>
				<GenericInput type="text" placeholder="Username" v-model="state.inputs.username" />
				<GenericInput type="password" placeholder="Password" v-model="state.inputs.password" />
				<GenericButton variant="wide" :disabled="state.busy || !valid" :spinner="state.busy" @click="signin">
					{{ intl(lang, {
						en: 'Sign in',
						de: 'Anmelden',
						uk: 'Увійти'
					}) }}
				</GenericButton>
			</InputGroup>

			<div v-if="state.session || state.error" class="status-messages">

				<InlineErrorMessage v-if="state.error">

					<template v-slot:title>
						{{ intl(lang, {
							en: 'Sign-in error',
							de: 'Anmeldefehler',
							uk: 'Помилка входу'
						}) }}
					</template>

					{{ state.error }}

				</InlineErrorMessage>

				<div v-if="state.session" class="suceess-message">
					{{ intl(lang, {
						en: 'Signed in successfully!',
						de: 'Erfolgreich angemeldet!',
						uk: 'Вхід виконано!'
					}) }}
				</div>

			</div>

		</div>

	</div>

</template>

<style lang="scss" scoped>

	.signin-screen {
		display: flex;
		justify-content: center;
		align-items: center;
		flex-grow: 1;
	}

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

			.message-title {
				font-size: 1.125rem;
				font-weight: 600;
			}
		}
	}

	.status-messages {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.suceess-message {
		font-size: 0.85rem;
		font-weight: 600;
		text-align: center;
		padding: 0.5rem 1rem;
		background-color: var(--app-theme-irish-green);
		border-radius: 0.5rem;
		border: 1px solid var(--app-theme-powder-trail);
	}

</style>
