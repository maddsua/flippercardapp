<script setup lang="ts">
import { useRouter } from 'vue-router';


const props = defineProps<{
	href: string;
	icon?: 'home' | 'star' | 'search' | 'menu';
	active?: boolean;
	label?: string;
}>();

const router = useRouter();

</script>

<template>
	<button type="button" :aria-label="label" :title="label" :class="{ active, [`icon-${icon}`]: !!icon }" @click="router.push(href)"></button>
</template>

<style lang="scss" scoped>
	button {
		display: block;
		padding: 0.5rem 0.75rem;
		opacity: 0.75;
		border: none;
		outline: none;
		background: none;
		background-color: unset;
		transition: all 150ms ease;

		&::after {
			content: "";
			display: block;
			width: 1.5rem;
			height: 1.5rem;
			mask-type: alpha;
			mask-size: contain;
			mask-repeat: no-repeat;
			mask-position: center;
			background-color: var(--app-theme-snow-white);
			transition: all 150ms ease;
		}

		@media (orientation: landscape) {

			padding: 0.35rem;

			&::after {
				width: 1.5rem;
				height: 1.5rem;
			}
		}

		&:hover {
			cursor: pointer;
			opacity: 1;
		}

		&.icon-home::after {
			mask-image: url(/src/assets/icons/home-mask.svg);
		}

		&.icon-star::after {
			mask-image: url(/src/assets/icons/star-mask.svg);
		}

		&.icon-search::after {
			mask-image: url(/src/assets/icons/search-mask.svg);
		}

		&.icon-menu::after {
			mask-image: url(/src/assets/icons/menu-mask.svg);
		}

		&.active {
			opacity: 1;

			&::after {
				background-color: var(--app-theme-midnight);
			}

			background-color: var(--app-theme-mysterious-white);
			border-radius: 0.5rem;
		}
	}
</style>
