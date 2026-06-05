<script setup lang="ts">
import { reactive } from 'vue';
import DashboardSessionWidget from './DashboardSessionWidget.vue';
import type { AuthState } from '../../api_models';
import AppUiHeader from '../App/Layout/AppUiHeader.vue';

const state = reactive({
	auth: null as  AuthState | null,
});

</script>

<template>

	<AppUiHeader>
		<template v-slot:title>
			Dashboard
		</template>
	</AppUiHeader>

	<div class="options-menu" v-if="state.auth?.actor">

		<template v-if="state.auth.actor.permissions.team_member">

			Congrats, you're an admin!

		</template>

		<div v-else class="noop-filler">
			No settings are available as of now
		</div>

	</div>

	<DashboardSessionWidget @stateUpdate="val => state.auth = val" />

	<p v-if="!state.auth?.actor" class="service-notice">
		Designed by maddsua. Provided by MWS via Railway Corp.
	</p>

</template>

<style lang="scss" scoped>
	.options-menu {
		display: flex;
		flex-direction: column;
		flex-grow: 1;
		
		.noop-filler {
			display: flex;
			flex-grow: 1;
			align-items: center;
			justify-content: center;
			font-size: 0.75rem;
		}
	}

	.service-notice {
		width: 100%;
		text-align: center;
		font-size: 0.65rem;
	}
</style>
