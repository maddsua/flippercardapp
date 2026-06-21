<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import AppViewButton from './AppViewButton.vue';
import { detectVirtualKeyboard } from '@/keyboard';

const route = useRoute();
const appView = computed((): string | null => typeof route.meta.app_view === 'string' ? route.meta.app_view : null);

const vkbdState = window.navigator.maxTouchPoints > 0 ? detectVirtualKeyboard() : null;;

</script>

<template>
	<div class="app-viewport">

		<div class="app-content-viewport">
			<div class="app-content">
				<slot>
					[App content]
				</slot>
			</div>
		</div>

		<div class="app-toolbar" :class="{ hidden: vkbdState?.isOpen }">
			<div class="app-navigation">
				<AppViewButton href="/" :active="appView === 'home'" icon="home" />
				<AppViewButton href="/starred" :active="appView === 'starred'" icon="star" />
				<AppViewButton href="/collections/discover" :active="appView === 'discover'" icon="search" />
				<AppViewButton href="/dashboard" :active="appView === 'menu'" icon="menu" />
			</div>
		</div>

	</div>
</template>

<style lang="scss" scoped>

	@use '@/media.scss';

	.app-viewport {
		position: relative;
		width: 100%;
		height: 100%;
		height: 100dvh;
		min-height: 0;
		overflow: hidden auto;
		scrollbar-width: thin;

		@include media.pwa {
			height: 100vh;
		}

		.app-content-viewport {
			position: relative;
			display: flex;
			justify-content: center;
			width: 100%;
			min-height: 100%;
			padding-bottom: 6rem;

			.app-content {
				position: relative;
				display: flex;
				flex-direction: column;
				gap: 2.5rem;
				width: 100%;
				padding: 2rem 1rem;

				@include media.pwa {
					padding: 1rem;
				}

				@include media.desktop {
					max-width: 40rem;
				}
			}
		}

		.app-toolbar {
			position: fixed;
			bottom: 0;
			left: 0;
			z-index: 100;
			display: flex;
			align-items: center;
			justify-content: center;
			width: 100%;
			background-color: var(--app-theme-carbon);
			transition: all 200ms ease;

			&.hidden {
				transform: translateY(25vh);
			}

			.app-navigation {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				justify-content: space-between;
				padding: 0.5rem 1rem;
				width: 100%;
				flex-shrink: 0;
				border-radius: 1rem;
				background-color: var(--app-theme-ghostly-glow);

				@include media.desktop {
					max-width: 40rem;
				}
			}
		}
	}
</style>
